// Package tunnel implements the reverse tunnel transport layer.
// Server accepts Agent connections and dispatches RPC frames.
package tunnel

import (
	"context"
	"crypto/subtle"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// AgentInfo carries metadata about a connected agent.
type AgentInfo struct {
	AgentID string
	Version string
	OS      string
	Arch    string
}

// RegistrationPayload is the JSON payload sent by an agent during registration.
type RegistrationPayload struct {
	AgentID  string `json:"agent_id"`
	Version  string `json:"version"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	JoinToken string `json:"join_token"`
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
	agents     sync.Map // map[string]*agentStream
	pending    sync.Map // map[string]*pendingCall   (call_id -> pendingCall)
	streams    sync.Map // map[string]*pendingStream (stream_id -> pendingStream)
	callID     atomic.Uint64
	lifecycle  AgentLifecycle
	joinToken  string
	tlsConfig  *tls.Config
}

type agentStream struct {
	agentID string
	send    func(*Frame) error
}

// ServerOption configures the ServerTunnel.
type ServerOption func(*ServerTunnel)

// WithJoinToken sets the join token for agent registration authentication.
func WithJoinToken(token string) ServerOption {
	return func(s *ServerTunnel) {
		s.joinToken = token
	}
}

// WithAgentLifecycle sets the lifecycle callback handler.
func WithAgentLifecycle(l AgentLifecycle) ServerOption {
	return func(s *ServerTunnel) {
		s.lifecycle = l
	}
}

// NewServerTunnel creates a new Server-side tunnel manager.
func NewServerTunnel(opts ...ServerOption) *ServerTunnel {
	s := &ServerTunnel{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// SetJoinToken configures the optional registration token required from Agents.
func (s *ServerTunnel) SetJoinToken(token string) {
	s.joinToken = token
}

// SetTLSConfig configures optional TLS credentials for the gRPC server.
func (s *ServerTunnel) SetTLSConfig(cfg *tls.Config) {
	s.tlsConfig = cfg
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
	s.listener = lis

	var opts []grpc.ServerOption
	if s.tlsConfig != nil {
		opts = append(opts, grpc.Creds(credentials.NewTLS(s.tlsConfig)))
	}
	s.grpcServer = grpc.NewServer(opts...)
	RegisterReverseTunnelServer(s.grpcServer, s)
	slog.Info("server tunnel: listening", "addr", addr)
	return s.grpcServer.Serve(lis)
}

// Stop stops the gRPC server.
func (s *ServerTunnel) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
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
	if agentID == "" {
		err := &InvalidAgentIDError{}
		slog.Warn("server tunnel: empty agent id rejected")
		return err
	}
	// Parse registration payload for agent metadata
	info := AgentInfo{AgentID: agentID}
	var regPayload RegistrationPayload
	if len(frame.Payload) > 0 {
		if err := json.Unmarshal(frame.Payload, &regPayload); err != nil {
			slog.Warn("server tunnel: invalid registration payload", "agent_id", agentID, "err", err)
			return fmt.Errorf("invalid registration payload: %w", err)
		}
		info.Version = regPayload.Version
		info.OS = regPayload.OS
		info.Arch = regPayload.Arch
		if regPayload.AgentID != "" {
			info.AgentID = regPayload.AgentID
			agentID = regPayload.AgentID
		}
	}

	// Validate join token
	joinToken := regPayload.JoinToken
	if joinToken == "" {
		joinToken = string(frame.Payload)
	}
	if !s.validJoinToken(joinToken) {
		err := &AgentUnauthorizedError{AgentID: agentID}
		_ = stream.Send(&Frame{
			FrameType: FrameType_FRAME_REGISTER_REJECTED,
			Error:     "invalid join token",
		})
		slog.Warn("server tunnel: agent registration rejected", "agent_id", agentID)
		return err
	}

	as := &agentStream{
		agentID: agentID,
		send: func(f *Frame) error {
			return stream.Send(f)
		},
	}
	if _, loaded := s.agents.LoadOrStore(agentID, as); loaded {
		err := &DuplicateAgentError{AgentID: agentID}
		slog.Warn("server tunnel: duplicate agent rejected", "agent_id", agentID)
		return err
	}
	slog.Info("server tunnel: agent registered", "agent_id", agentID)
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

		case FrameType_FRAME_REPLY:
			// Route single response to the waiting caller
			if v, ok := s.pending.LoadAndDelete(frame.Id); ok {
				pc := v.(*pendingCall)
				pc.replyChan <- frame
			}

		case FrameType_FRAME_STREAM_END:
			// Route stream end to the waiting stream consumer
			if v, ok := s.streams.LoadAndDelete(frame.Id); ok {
				ps := v.(*pendingStream)
				select {
				case ps.chunkChan <- frame:
				default:
				}
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

func (s *ServerTunnel) validJoinToken(token string) bool {
	if s.joinToken == "" {
		return true
	}
	if len(token) != len(s.joinToken) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(token), []byte(s.joinToken)) == 1
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
	// Clean up stream entry when context is cancelled
	go func() {
		<-ctx.Done()
		s.streams.Delete(id)
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

// AgentInfo returns info about all connected agents.
func (s *ServerTunnel) AgentCount() int {
	count := 0
	s.agents.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

func (s *ServerTunnel) nextID() string {
	return "rpc-" + itoa(int(s.callID.Add(1)))
}

// AgentOfflineError is returned when the target agent is not connected.
type AgentOfflineError struct {
	AgentID string
}

func (e *AgentOfflineError) Error() string {
	return "agent " + e.AgentID + " is offline"
}

// InvalidAgentIDError is returned when an Agent registers without an identity.
type InvalidAgentIDError struct{}

func (e *InvalidAgentIDError) Error() string {
	return "agent id is required"
}

// AgentUnauthorizedError is returned when an Agent registration token is invalid.
type AgentUnauthorizedError struct {
	AgentID string
}

func (e *AgentUnauthorizedError) Error() string {
	return "agent " + e.AgentID + " is unauthorized"
}

// DuplicateAgentError is returned when an Agent ID is already connected.
type DuplicateAgentError struct {
	AgentID string
}

func (e *DuplicateAgentError) Error() string {
	return "agent " + e.AgentID + " is already connected"
}

// RPCError is returned when the agent returns an RPC error.
type RPCError struct {
	Method string
	Err    string
}

func (e *RPCError) Error() string {
	return "rpc " + e.Method + " failed: " + e.Err
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
