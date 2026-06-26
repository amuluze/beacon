// Package tunnel implements the reverse tunnel transport layer.
// Agent connects to Server via gRPC bidirectional stream, and
// RPC frames are multiplexed through the same connection.
package tunnel

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// RegistrationPayload is the JSON payload sent in the FRAME_REGISTER frame.
// It carries agent identity, version metadata, and the join token.
type RegistrationPayload struct {
	AgentID   string `json:"agent_id"`
	Version   string `json:"version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	JoinToken string `json:"join_token"`
}

// AgentOption configures an AgentTunnel.
type AgentOption func(*AgentTunnel)

// WithAgentTLS enables TLS credentials for the agent connection.
func WithAgentTLS(creds credentials.TransportCredentials) AgentOption {
	return func(a *AgentTunnel) {
		a.tlsCredentials = creds
	}
}

// Handler is the callback type for processing incoming RPC calls.
// The handler receives JSON-encoded args and returns JSON-encoded reply.
// For streaming methods, the handler sends multiple frames via the provided sender.
type Handler func(ctx context.Context, method string, payload []byte, streamSender func(*Frame)) ([]byte, error)

// AgentTunnel is the Agent-side tunnel client.
// It connects to the Server and waits for incoming RPC frames.
type AgentTunnel struct {
	serverAddr     string
	agentID        string
	joinToken      string
	version        string
	osName         string
	arch           string
	conn           *grpc.ClientConn
	client         ReverseTunnelClient
	stream         grpc.BidiStreamingClient[Frame, Frame]
	tlsCredentials credentials.TransportCredentials

	mu          sync.Mutex
	handler     Handler
	closed      bool
	reconnect   bool
	heartbeatCh chan struct{}
}

// NewAgentTunnel creates a new Agent-side tunnel connection.
// The agent will connect to serverAddr and identify as agentID.
func NewAgentTunnel(serverAddr string, agentID string, opts ...AgentOption) *AgentTunnel {
	a := &AgentTunnel{
		serverAddr:  serverAddr,
		agentID:     agentID,
		reconnect:   true,
		heartbeatCh: make(chan struct{}, 1),
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// SetJoinToken sets the join token for agent registration.
func (a *AgentTunnel) SetJoinToken(token string) {
	a.joinToken = token
}

// SetVersionInfo sets the agent version, OS, and architecture for registration.
func (a *AgentTunnel) SetVersionInfo(version, osName, arch string) {
	a.version = version
	a.osName = osName
	a.arch = arch
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
				slog.Debug("agent tunnel: close previous connection", "err", err)
			}
		}

		slog.Info("agent tunnel: connecting to server", "addr", a.serverAddr, "agent_id", a.agentID)

		creds := a.tlsCredentials
		if creds == nil {
			creds = insecure.NewCredentials()
		}

		conn, err := grpc.NewClient(a.serverAddr,
			grpc.WithTransportCredentials(creds),
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

		stream, err := a.client.Tunnel(ctx)
		if err != nil {
			slog.Error("agent tunnel: create stream failed", "err", err)
			if closeErr := conn.Close(); closeErr != nil {
				slog.Debug("agent tunnel: close connection failed", "err", closeErr)
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
		slog.Info("agent tunnel: connected to server", "tls", a.tlsCredentials != nil)

		// Send registration frame with agent metadata
		regPayload := RegistrationPayload{
			AgentID:   a.agentID,
			Version:   a.version,
			OS:        a.osName,
			Arch:      a.arch,
			JoinToken: a.joinToken,
		}
		payloadBytes, _ := json.Marshal(regPayload)
		regFrame := &Frame{
			FrameType: FrameType_FRAME_REGISTER,
			Method:    a.agentID,
			Payload:   payloadBytes,
		}
		if err := stream.Send(regFrame); err != nil {
			slog.Error("agent tunnel: send registration failed", "err", err)
			if closeErr := conn.Close(); closeErr != nil {
				slog.Debug("agent tunnel: close connection failed", "err", closeErr)
			}
			continue
		}
		slog.Info("agent tunnel: registered as", "agent_id", a.agentID)

		// Start heartbeat goroutine
		heartbeatCtx, heartbeatCancel := context.WithCancel(ctx)
		go a.heartbeatLoop(heartbeatCtx, stream)

		if err := a.processStream(heartbeatCtx, stream); err != nil {
			if rej, ok := err.(*RegistrationRejectedError); ok {
				slog.Error("agent tunnel: registration rejected", "reason", rej.Reason)
				heartbeatCancel()
				if closeErr := conn.Close(); closeErr != nil {
					slog.Debug("agent tunnel: close connection failed", "err", closeErr)
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
			slog.Warn("agent tunnel: stream ended", "err", err)
		}

		heartbeatCancel()

		if !a.reconnect {
			return err
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
			if err := stream.Send(frame); err != nil {
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

		switch frame.FrameType {
		case FrameType_FRAME_REGISTER_REJECTED:
			return &RegistrationRejectedError{Reason: frame.Error}
		case FrameType_FRAME_REQUEST:
			go a.dispatch(ctx, stream, frame)
		default:
			// ignore other frame types
		}
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
		_ = stream.Send(resp)
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
		_ = stream.Send(f)
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

	if err := stream.Send(resp); err != nil {
		slog.Error("agent tunnel: send response failed", "id", req.Id, "method", req.Method, "err", err)
	}
}

// Close stops the tunnel and disconnects.
func (a *AgentTunnel) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.reconnect = false
	a.closed = true
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

// RegistrationRejectedError is returned when the server rejects agent registration.
type RegistrationRejectedError struct {
	Reason string
}

func (e *RegistrationRejectedError) Error() string {
	return "registration rejected: " + e.Reason
}

// marshalArgs encodes the RPC arguments as JSON.
func marshalArgs(args interface{}) ([]byte, error) {
	if args == nil {
		return nil, nil
	}
	return json.Marshal(args)
}

// unmarshalReply decodes the JSON response into reply.
func unmarshalReply(data []byte, reply interface{}) error {
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
