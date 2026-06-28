package tunnel

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

type fakeTunnelStream struct {
	ctx  context.Context
	recv chan *Frame
	send chan *Frame
}

func newFakeTunnelStream(ctx context.Context, frames ...*Frame) *fakeTunnelStream {
	s := &fakeTunnelStream{
		ctx:  ctx,
		recv: make(chan *Frame, len(frames)),
		send: make(chan *Frame, 16),
	}
	for _, frame := range frames {
		s.recv <- frame
	}
	return s
}

func (s *fakeTunnelStream) Recv() (*Frame, error) {
	select {
	case <-s.ctx.Done():
		return nil, s.ctx.Err()
	case frame, ok := <-s.recv:
		if !ok {
			return nil, io.EOF
		}
		return frame, nil
	}
}

func (s *fakeTunnelStream) Send(frame *Frame) error {
	s.send <- frame
	return nil
}

func (s *fakeTunnelStream) SetHeader(metadata.MD) error  { return nil }
func (s *fakeTunnelStream) SendHeader(metadata.MD) error { return nil }
func (s *fakeTunnelStream) SetTrailer(metadata.MD)       {}
func (s *fakeTunnelStream) Context() context.Context     { return s.ctx }
func (s *fakeTunnelStream) SendMsg(any) error            { return nil }
func (s *fakeTunnelStream) RecvMsg(any) error            { return nil }

type lifecycleRecorder struct {
	connected    chan string
	disconnected chan string
	heartbeat    chan string
}

func newLifecycleRecorder() *lifecycleRecorder {
	return &lifecycleRecorder{
		connected:    make(chan string, 4),
		disconnected: make(chan string, 4),
		heartbeat:    make(chan string, 4),
	}
}

func (l *lifecycleRecorder) OnAgentConnect(agentID string) {
	l.connected <- agentID
}

func (l *lifecycleRecorder) OnAgentDisconnect(agentID string) {
	l.disconnected <- agentID
}

func (l *lifecycleRecorder) OnAgentHeartbeat(agentID string) {
	l.heartbeat <- agentID
}

func TestTunnelRejectsInvalidRegistration(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		frame   *Frame
		wantErr any
	}{
		{
			name:    "empty agent id",
			frame:   &Frame{FrameType: FrameType_FRAME_REGISTER, Method: ""},
			wantErr: &InvalidAgentIDError{},
		},
		{
			name:    "invalid token",
			token:   "secret",
			frame:   &Frame{FrameType: FrameType_FRAME_REGISTER, Method: "agent-a", Payload: []byte("wrong")},
			wantErr: &AgentUnauthorizedError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tun := NewServerTunnel()
			tun.SetJoinToken(tt.token)
			stream := newFakeTunnelStream(context.Background(), tt.frame)

			err := tun.Tunnel(stream)
			if err == nil {
				t.Fatal("expected error")
			}
			switch tt.wantErr.(type) {
			case *InvalidAgentIDError:
				var target *InvalidAgentIDError
				if !errors.As(err, &target) {
					t.Fatalf("error = %T, want InvalidAgentIDError", err)
				}
			case *AgentUnauthorizedError:
				var target *AgentUnauthorizedError
				if !errors.As(err, &target) {
					t.Fatalf("error = %T, want AgentUnauthorizedError", err)
				}
			}
		})
	}
}

func TestTunnelRejectsDuplicateAgentID(t *testing.T) {
	tun := NewServerTunnel()
	lifecycle := newLifecycleRecorder()
	tun.SetAgentLifecycle(lifecycle)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	first := newFakeTunnelStream(ctx, &Frame{FrameType: FrameType_FRAME_REGISTER, Method: "agent-a"})
	errCh := make(chan error, 1)
	go func() {
		errCh <- tun.Tunnel(first)
	}()

	select {
	case got := <-lifecycle.connected:
		if got != "agent-a" {
			t.Fatalf("connected agent = %q, want agent-a", got)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for first agent connection")
	}

	second := newFakeTunnelStream(context.Background(), &Frame{FrameType: FrameType_FRAME_REGISTER, Method: "agent-a"})
	err := tun.Tunnel(second)
	var duplicate *DuplicateAgentError
	if !errors.As(err, &duplicate) {
		t.Fatalf("error = %T, want DuplicateAgentError", err)
	}

	cancel()
	<-errCh
}

func TestTunnelLifecycleHeartbeatAndDisconnect(t *testing.T) {
	tun := NewServerTunnel()
	lifecycle := newLifecycleRecorder()
	tun.SetAgentLifecycle(lifecycle)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream := newFakeTunnelStream(ctx, &Frame{FrameType: FrameType_FRAME_REGISTER, Method: "agent-a"})
	errCh := make(chan error, 1)
	go func() {
		errCh <- tun.Tunnel(stream)
	}()

	select {
	case got := <-lifecycle.connected:
		if got != "agent-a" {
			t.Fatalf("connected agent = %q, want agent-a", got)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for connect")
	}

	stream.recv <- &Frame{FrameType: FrameType_FRAME_HEARTBEAT}
	select {
	case got := <-lifecycle.heartbeat:
		if got != "agent-a" {
			t.Fatalf("heartbeat agent = %q, want agent-a", got)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for heartbeat")
	}

	close(stream.recv)
	<-errCh
	select {
	case got := <-lifecycle.disconnected:
		if got != "agent-a" {
			t.Fatalf("disconnected agent = %q, want agent-a", got)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for disconnect")
	}
}
