// Package tunnel implements the reverse tunnel transport layer.
// Agent connects to Server via gRPC bidirectional stream, and
// RPC frames are multiplexed through the same connection.
package tunnel

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// ErrDataStreamNotReady is returned when a payload is pushed before the data
// stream is established (or after it was lost).
var ErrDataStreamNotReady = errors.New("tunnel data stream not ready")

// Handler is the callback type for processing incoming RPC calls.
// The handler receives JSON-encoded args and returns JSON-encoded reply.
// For streaming methods, the handler sends multiple frames via the provided sender.
type Handler func(ctx context.Context, method string, payload []byte, streamSender func(*Frame)) ([]byte, error)

// AgentTunnel is the Agent-side tunnel client.
// It connects to the Server and waits for incoming RPC frames.
type AgentTunnel struct {
	serverAddr string
	agentID    string
	joinToken  string
	version    string
	os         string
	arch       string
	conn       *grpc.ClientConn
	client     ReverseTunnelClient
	stream     grpc.BidiStreamingClient[Frame, Frame]
	dataStream grpc.BidiStreamingClient[Frame, Frame]
	tlsConfig  *tls.Config

	mu          sync.Mutex
	sendMu      sync.Mutex
	dataSendMu  sync.Mutex
	handler     Handler
	closed      bool
	reconnect   bool
	heartbeatCh chan struct{}
	dataReady   atomic.Bool
	dataSeq     atomic.Uint64
}

// NewAgentTunnel creates a new Agent-side tunnel connection.
// The agent will connect to serverAddr and identify as agentID.
func NewAgentTunnel(serverAddr string, agentID string) *AgentTunnel {
	return &AgentTunnel{
		serverAddr:  serverAddr,
		agentID:     agentID,
		reconnect:   true,
		heartbeatCh: make(chan struct{}, 1),
	}
}

// SetJoinToken sets the join token for agent registration.
func (a *AgentTunnel) SetJoinToken(token string) {
	a.joinToken = token
}

// SetMetadata configures the immutable build/runtime metadata sent in every
// registration frame. Reconnects reuse the same values for the process.
func (a *AgentTunnel) SetMetadata(version, osName, arch string) {
	a.version = version
	a.os = osName
	a.arch = arch
}

func (a *AgentTunnel) registrationPayload() RegistrationPayload {
	return RegistrationPayload{
		AgentID:   a.agentID,
		Version:   a.version,
		OS:        a.os,
		Arch:      a.arch,
		JoinToken: a.joinToken,
	}
}

// dataRegistrationPayload builds the registration payload for the data stream.
// The data stream shares the connection (and thus the client certificate
// identity) with the control stream but is registered with the data role so
// the Server does not register a second online session for this agent.
func (a *AgentTunnel) dataRegistrationPayload() RegistrationPayload {
	payload := a.registrationPayload()
	payload.Role = RoleData
	return payload
}

// DataReady reports whether the data stream is established and usable.
func (a *AgentTunnel) DataReady() bool {
	return a.dataReady.Load()
}

// PushData sends a payload to the Server over the data stream.
// It returns an error when the data stream is not available or the send
// failed; callers may fall back to another transport. Payloads are fire and
// forget: no queueing is performed, so a disconnected stream drops the
// payload instead of buffering it.
func (a *AgentTunnel) PushData(method string, payload []byte) error {
	if !a.dataReady.Load() {
		return ErrDataStreamNotReady
	}
	a.dataSendMu.Lock()
	defer a.dataSendMu.Unlock()
	if a.dataStream == nil {
		return ErrDataStreamNotReady
	}
	frame := &Frame{
		Id:        fmt.Sprintf("d%d", a.dataSeq.Add(1)),
		FrameType: FrameType_FRAME_STREAM_DATA,
		Method:    method,
		Payload:   payload,
	}
	if err := a.dataStream.Send(frame); err != nil {
		a.dataReady.Store(false)
		return fmt.Errorf("send data frame: %w", err)
	}
	return nil
}

// openDataStream opens the push-only stream on the already established
// connection. The Server multiplexes both streams over the same HTTP/2
// connection, so the node certificate identity is shared and control calls
// are not blocked by monitoring payloads.
func (a *AgentTunnel) openDataStream(ctx context.Context, client ReverseTunnelClient) error {
	stream, err := client.Tunnel(ctx)
	if err != nil {
		return fmt.Errorf("create data stream: %w", err)
	}
	regPayload, _ := json.Marshal(a.dataRegistrationPayload())
	regFrame := &Frame{
		FrameType: FrameType_FRAME_REGISTER,
		Method:    a.agentID,
		Payload:   regPayload,
	}
	if err := stream.Send(regFrame); err != nil {
		return fmt.Errorf("send data registration: %w", err)
	}
	a.dataStream = stream
	a.dataReady.Store(true)

	// Drain server acknowledgements so the server's writes never stall.
	go func() {
		for {
			frame, err := stream.Recv()
			if err != nil {
				a.dataReady.Store(false)
				return
			}
			if frame.FrameType == FrameType_FRAME_REPLY && frame.Error != "" {
				slog.Error("agent tunnel: data payload rejected by server", "agent_id", a.agentID, "error", frame.Error)
			}
		}
	}()
	return nil
}

// SetTLSConfig configures optional TLS credentials for the gRPC client.
func (a *AgentTunnel) SetTLSConfig(cfg *tls.Config) {
	a.tlsConfig = cfg
}

// SetHandler registers the RPC handler for incoming requests.
func (a *AgentTunnel) SetHandler(h Handler) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.handler = h
}

// Start connects to the Server and starts processing frames.
// Blocks until the connection is closed or ctx is cancelled.
func (a *AgentTunnel) Start(ctx context.Context) error {
	for {
		if a.conn != nil {
			if err := a.conn.Close(); err != nil {
				slog.Debug("agent tunnel: close previous connection failed", "err", err)
			}
		}

		slog.Info("agent tunnel: connecting to server", "addr", a.serverAddr, "agent_id", a.agentID)
		transportCreds := insecure.NewCredentials()
		if a.tlsConfig != nil {
			transportCreds = credentials.NewTLS(a.tlsConfig)
		}
		conn, err := grpc.NewClient(a.serverAddr,
			grpc.WithTransportCredentials(transportCreds),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                30 * time.Second,
				Timeout:             10 * time.Second,
				PermitWithoutStream: true,
			}),
		)
		if err != nil {
			slog.Error("agent tunnel: connect failed", "err", err)
			if !a.reconnect {
				return err
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(5 * time.Second):
				continue
			}
		}
		a.conn = conn
		a.client = NewReverseTunnelClient(conn)
		a.dataReady.Store(false)

		stream, err := a.client.Tunnel(ctx)
		if err != nil {
			slog.Error("agent tunnel: create stream failed", "err", err)
			if closeErr := conn.Close(); closeErr != nil {
				slog.Debug("agent tunnel: close failed connection failed", "err", closeErr)
			}
			if !a.reconnect {
				return err
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(5 * time.Second):
				continue
			}
		}
		a.stream = stream
		slog.Info("agent tunnel: connected to server")

		// Send registration frame with JSON payload carrying agent metadata.
		regPayload, _ := json.Marshal(a.registrationPayload())
		regFrame := &Frame{
			FrameType: FrameType_FRAME_REGISTER,
			Method:    a.agentID,
			Payload:   regPayload,
		}
		if err := a.sendFrame(stream, regFrame); err != nil {
			slog.Error("agent tunnel: send registration failed", "err", err)
			if closeErr := conn.Close(); closeErr != nil {
				slog.Debug("agent tunnel: close unregistered connection failed", "err", closeErr)
			}
			continue
		}
		slog.Info("agent tunnel: registered as", "agent_id", a.agentID)

		// Open the push-only data stream on the same connection. Failures are
		// not fatal: the agent keeps working and falls back to the legacy
		// HTTP report channel until the data stream is available.
		if err := a.openDataStream(ctx, a.client); err != nil {
			slog.Warn("agent tunnel: data stream unavailable, falling back to http report", "agent_id", a.agentID, "err", err)
		}

		// Start heartbeat goroutine
		heartbeatCtx, heartbeatCancel := context.WithCancel(ctx)
		go a.heartbeatLoop(heartbeatCtx, stream)

		streamErr := a.processStream(heartbeatCtx, stream)
		if streamErr != nil {
			slog.Warn("agent tunnel: stream ended", "err", streamErr)
		}

		heartbeatCancel()

		if !a.reconnect {
			return streamErr
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(3 * time.Second):
			continue
		}
	}
}

func (a *AgentTunnel) heartbeatLoop(ctx context.Context, stream grpc.BidiStreamingClient[Frame, Frame]) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			frame := &Frame{
				FrameType: FrameType_FRAME_HEARTBEAT,
				Method:    a.agentID,
			}
			if err := a.sendFrame(stream, frame); err != nil {
				slog.Debug("agent tunnel: heartbeat send failed", "err", err)
				return
			}
		}
	}
}

func (a *AgentTunnel) processStream(ctx context.Context, stream grpc.BidiStreamingClient[Frame, Frame]) error {
	for {
		frame, err := stream.Recv()
		if err != nil {
			return err
		}
		if frame.Eos {
			slog.Info("agent tunnel: received eos")
			return nil
		}
		if frame.FrameType != FrameType_FRAME_REQUEST {
			continue
		}

		// Dispatch the RPC call
		go a.dispatch(ctx, stream, frame)
	}
}

func (a *AgentTunnel) dispatch(ctx context.Context, stream grpc.BidiStreamingClient[Frame, Frame], req *Frame) {
	a.mu.Lock()
	handler := a.handler
	a.mu.Unlock()

	if handler == nil {
		resp := &Frame{
			Id:        req.Id,
			Error:     "no handler registered",
			FrameType: FrameType_FRAME_REPLY,
		}
		_ = a.sendFrame(stream, resp)
		return
	}

	// Stream sender: agent can send multiple frames for streaming methods
	var streamEnd atomicBool
	streamSender := func(f *Frame) {
		if streamEnd.get() {
			return
		}
		f.Id = req.Id
		f.FrameType = FrameType_FRAME_STREAM_DATA
		if f.Eos {
			f.FrameType = FrameType_FRAME_STREAM_END
			streamEnd.set(true)
		}
		_ = a.sendFrame(stream, f)
	}

	replyPayload, err := handler(ctx, req.Method, req.Payload, streamSender)

	// If the handler already sent stream frames, don't send a reply
	if streamEnd.get() {
		return
	}

	// Send single reply
	resp := &Frame{
		Id:        req.Id,
		FrameType: FrameType_FRAME_REPLY,
	}
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Payload = replyPayload
	}

	if err := a.sendFrame(stream, resp); err != nil {
		slog.Error("agent tunnel: send response failed", "id", req.Id, "method", req.Method, "err", err)
	}
}

func (a *AgentTunnel) sendFrame(stream grpc.BidiStreamingClient[Frame, Frame], frame *Frame) error {
	a.sendMu.Lock()
	defer a.sendMu.Unlock()
	return stream.Send(frame)
}

// Close stops the tunnel and disconnects.
func (a *AgentTunnel) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.reconnect = false
	a.closed = true
	a.dataReady.Store(false)
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

// marshalArgs encodes the RPC arguments as JSON.
func marshalArgs(args interface{}) ([]byte, error) {
	if args == nil {
		return nil, nil
	}
	return json.Marshal(args)
}

// unmarshalReply decodes the JSON response into reply.
// reply may be nil for fire-and-forget RPCs (e.g. ImagesPrune);
// in that case the payload is intentionally discarded.
func unmarshalReply(data []byte, reply interface{}) error {
	if reply == nil {
		return nil
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, reply)
}

// atomicBool is a simple atomic boolean.
type atomicBool struct {
	v int32
}

func (a *atomicBool) get() bool {
	return atomic.LoadInt32(&a.v) != 0
}

func (a *atomicBool) set(v bool) {
	if v {
		atomic.StoreInt32(&a.v, 1)
	} else {
		atomic.StoreInt32(&a.v, 0)
	}
}
