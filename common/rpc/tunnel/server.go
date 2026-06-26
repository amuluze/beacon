// Package tunnel implements the reverse tunnel transport layer.
// Server accepts Agent connections and dispatches RPC frames.
package tunnel

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
)

// ServerOption configures a ServerTunnel.
type ServerOption func(*ServerTunnel)

// WithServerTLS enables TLS on the tunnel listener.
func WithServerTLS(cfg *tls.Config) ServerOption {
	return func(s *ServerTunnel) {
		s.tlsConfig = cfg
	}
}

// WithJoinToken sets the expected join token for agent registration.
// If empty (default), token validation is skipped for backward compatibility.
func WithJoinToken(token string) ServerOption {
	return func(s *ServerTunnel) {
		s.joinToken = token
	}
}

// AgentInfo carries the agent identity and metadata from the registration frame.
type AgentInfo struct {
	AgentID string
	Version string
	OS      string
	Arch    string
}

// AgentLifecycle is called when an agent connects or disconnects.
type AgentLifecycle interface {
	OnAgentConnect(info AgentInfo)
	OnAgentDisconnect(agentID string)
	OnAgentHeartbeat(agentID string)
}

// pendingCall represents an in-flight RPC call waiting for a response.
type pendingCall struct {
	replyChan chan *Frame
}

// pendingStream represents an in-flight streaming RPC.
type pendingStream struct {
	chunkChan chan *Frame
}

// ServerTunnel is the Server-side tunnel manager.
// It accepts Agent connections and provides a Call interface for dispatching RPCs.
type ServerTunnel struct {
	UnimplementedReverseTunnelServer
	grpcServer *grpc.Server
	listener   net.Listener
	tlsConfig  *tls.Config
	joinToken  string
	agents     sync.Map // map[string]*agentStream
	pending    sync.Map // map[string]*pendingCall   (call_id -> pendingCall)
	streams    sync.Map // map[string]*pendingStream (stream_id -> pendingStream)
	callID     atomic.Uint64
	lifecycle  AgentLifecycle
}

type agentStream struct {
	agentID string
	send    func(*Frame) error
}

// NewServerTunnel creates a new Server-side tunnel manager.
func NewServerTunnel(opts ...ServerOption) *ServerTunnel {
	s := &ServerTunnel{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// SetAgentLifecycle registers lifecycle callbacks for agent connect/disconnect.
func (s *ServerTunnel) SetAgentLifecycle(l AgentLifecycle) {
	s.lifecycle = l
}

// Start starts the gRPC server on the given address.
func (s *ServerTunnel) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	if s.tlsConfig != nil {
		lis = tls.NewListener(lis, s.tlsConfig)
	}
	s.listener = lis

	s.grpcServer = grpc.NewServer()
	RegisterReverseTunnelServer(s.grpcServer, s)
	slog.Info("server tunnel: listening", "addr", addr, "tls", s.tlsConfig != nil, "auth", s.joinToken != "")
	return s.grpcServer.Serve(lis)
}

// Stop stops the gRPC server.
func (s *ServerTunnel) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}

// parseRegistrationPayload parses the registration frame payload as JSON.
// Break change: payload must be a valid RegistrationPayload JSON.
func parseRegistrationPayload(payload []byte) (RegistrationPayload, error) {
	var reg RegistrationPayload
	if err := json.Unmarshal(payload, &reg); err != nil {
		return reg, fmt.Errorf("invalid registration payload: %w", err)
	}
	return reg, nil
}

// Tunnel implements the ReverseTunnelServer interface.
func (s *ServerTunnel) Tunnel(stream grpc.BidiStreamingServer[Frame, Frame]) error {
	// Wait for the first frame — must be REGISTER
	frame, err := stream.Recv()
	if err != nil {
		return err
	}

	if frame.FrameType != FrameType_FRAME_REGISTER {
		slog.Warn("server tunnel: first frame is not REGISTER")
		return nil
	}

	agentID := frame.Method
	reg, err := parseRegistrationPayload(frame.Payload)
	if err != nil {
		slog.Warn("server tunnel: registration payload parse failed", "agent_id", agentID, "err", err)
		rejFrame := &Frame{
			FrameType: FrameType_FRAME_REGISTER_REJECTED,
			Error:     err.Error(),
		}
		_ = stream.Send(rejFrame)
		return nil
	}

	// Validate join token if server has one configured.
	// Empty joinToken on server means skip validation (backward compatible).
	if s.joinToken != "" && reg.JoinToken != s.joinToken {
		slog.Warn("server tunnel: registration rejected", "agent_id", agentID)
		rejFrame := &Frame{
			FrameType: FrameType_FRAME_REGISTER_REJECTED,
			Error:     "invalid join token",
		}
		_ = stream.Send(rejFrame)
		return nil
	}

	info := AgentInfo{
		AgentID: agentID,
		Version: reg.Version,
		OS:      reg.OS,
		Arch:    reg.Arch,
	}
	slog.Info("server tunnel: agent registered", "agent_id", agentID, "version", info.Version, "os", info.OS, "arch", info.Arch)

	as := &agentStream{
		agentID: agentID,
		send: func(f *Frame) error {
			return stream.Send(f)
		},
	}
	s.agents.Store(agentID, as)
	defer func() {
		s.agents.Delete(agentID)
		if s.lifecycle != nil {
			s.lifecycle.OnAgentDisconnect(agentID)
		}
		slog.Info("server tunnel: agent disconnected", "agent_id", agentID)
	}()

	if s.lifecycle != nil {
		s.lifecycle.OnAgentConnect(info)
	}

	for {
		frame, err := stream.Recv()
		if err != nil {
			return err
		}

		switch frame.FrameType {
		case FrameType_FRAME_HEARTBEAT:
			if s.lifecycle != nil {
				s.lifecycle.OnAgentHeartbeat(agentID)
			}

		case FrameType_FRAME_REPLY, FrameType_FRAME_STREAM_END:
			// Route single response to the waiting caller
			if v, ok := s.pending.LoadAndDelete(frame.Id); ok {
				pc := v.(*pendingCall)
				pc.replyChan <- frame
			}

		case FrameType_FRAME_STREAM_DATA:
			// Route streaming chunk
			if v, ok := s.streams.Load(frame.Id); ok {
				ps := v.(*pendingStream)
				select {
				case ps.chunkChan <- frame:
				default:
				}
				if frame.Eos {
					s.streams.Delete(frame.Id)
				}
			}

		default:
			// Legacy behavior: treat as reply
			if v, ok := s.pending.LoadAndDelete(frame.Id); ok {
				pc := v.(*pendingCall)
				pc.replyChan <- frame
			}
		}
	}
}

// Call dispatches an RPC call to an agent and waits for the response.
func (s *ServerTunnel) Call(ctx context.Context, agentID string, method string, args interface{}, reply interface{}) error {
	payload, err := marshalArgs(args)
	if err != nil {
		return err
	}

	id := s.nextID()
	frame := &Frame{
		Id:        id,
		Method:    method,
		Payload:   payload,
		FrameType: FrameType_FRAME_REQUEST,
	}

	as, ok := s.agents.Load(agentID)
	if !ok {
		return &AgentOfflineError{AgentID: agentID}
	}

	pc := &pendingCall{replyChan: make(chan *Frame, 1)}
	s.pending.Store(id, pc)
	defer s.pending.Delete(id)

	if err := as.(*agentStream).send(frame); err != nil {
		return err
	}

	select {
	case resp := <-pc.replyChan:
		if resp.Error != "" {
			return &RPCError{Method: method, Err: resp.Error}
		}
		return unmarshalReply(resp.Payload, reply)
	case <-ctx.Done():
		return ctx.Err()
	}
}

// StreamCall sends an RPC call and returns a channel to receive streaming chunks.
func (s *ServerTunnel) StreamCall(ctx context.Context, agentID string, method string, args interface{}) (<-chan *Frame, error) {
	payload, err := marshalArgs(args)
	if err != nil {
		return nil, err
	}

	id := s.nextID()
	frame := &Frame{
		Id:        id,
		Method:    method,
		Payload:   payload,
		FrameType: FrameType_FRAME_REQUEST,
	}

	as, ok := s.agents.Load(agentID)
	if !ok {
		return nil, &AgentOfflineError{AgentID: agentID}
	}

	ps := &pendingStream{chunkChan: make(chan *Frame, 64)}
	s.streams.Store(id, ps)
	defer func() {
		// If ctx done before stream ends, clean up
		select {
		case <-ctx.Done():
			s.streams.Delete(id)
		default:
		}
	}()

	if err := as.(*agentStream).send(frame); err != nil {
		s.streams.Delete(id)
		return nil, err
	}

	return ps.chunkChan, nil
}

// AgentIDs returns the list of connected agent IDs.
func (s *ServerTunnel) AgentIDs() []string {
	var ids []string
	s.agents.Range(func(key, _ interface{}) bool {
		ids = append(ids, key.(string))
		return true
	})
	return ids
}

// IsOnline checks if an agent is currently connected.
func (s *ServerTunnel) IsOnline(agentID string) bool {
	_, ok := s.agents.Load(agentID)
	return ok
}

// AgentCount returns the number of connected agents.
func (s *ServerTunnel) AgentCount() int {
	count := 0
	s.agents.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

func (s *ServerTunnel) nextID() string {
	return "rpc-" + strconv.FormatUint(s.callID.Add(1), 10)
}

// AgentOfflineError is returned when the target agent is not connected.
type AgentOfflineError struct {
	AgentID string
}

func (e *AgentOfflineError) Error() string {
	return "agent " + e.AgentID + " is offline"
}

// RPCError is returned when the agent returns an RPC error.
type RPCError struct {
	Method string
	Err    string
}

func (e *RPCError) Error() string {
	return "rpc " + e.Method + " failed: " + e.Err
}
