// Package rpc
// Date: 2024/6/12 10:30
// Author: Amu
// Description:
package rpc

import (
	"context"
	"sync"

	"beacon/pkg/contextx"
	tunnelpkg "common/rpc/tunnel"
)

const (
	DefaultNetwork = "tcp"
)

// ErrMissingAgentID 在控制调用无法解析出目标 Agent 时返回。
// 权威定义在 contextx.ErrMissingAgentID，此处保留为兼容别名，
// 避免破坏已有的包外引用（如 repository、terminal 等）。
var ErrMissingAgentID = contextx.ErrMissingAgentID

// Caller is the interface for making RPC calls to agents.
type Caller interface {
	Call(ctx context.Context, method string, args interface{}, reply interface{}) error
	StreamCall(ctx context.Context, method string, args interface{}) (<-chan []byte, error)
	Close() error
}

// TunnelClient wraps ServerTunnel to implement the Caller interface.
type TunnelClient struct {
	tunnel  *tunnelpkg.ServerTunnel
	mu      sync.Mutex
	started bool
}

// NewTunnelClient creates a new tunnel-based RPC caller.
func NewTunnelClient(tunnel *tunnelpkg.ServerTunnel) *TunnelClient {
	return &TunnelClient{tunnel: tunnel}
}

func (tc *TunnelClient) Call(ctx context.Context, method string, args interface{}, reply interface{}) error {
	agentID, err := contextx.ResolveAgentID(ctx)
	if err != nil {
		return err
	}
	return tc.tunnel.Call(ctx, agentID, method, args, reply)
}

// StreamCall sends an RPC call and returns a channel of response chunks.
func (tc *TunnelClient) StreamCall(ctx context.Context, method string, args interface{}) (<-chan []byte, error) {
	agentID, err := contextx.ResolveAgentID(ctx)
	if err != nil {
		return nil, err
	}
	chunkChan, err := tc.tunnel.StreamCall(ctx, agentID, method, args)
	if err != nil {
		return nil, err
	}
	// Convert tunnel Frame channel to []byte channel
	out := make(chan []byte, 64)
	go func() {
		defer close(out)
		for {
			select {
			case frame, ok := <-chunkChan:
				if !ok {
					return
				}
				out <- frame.Payload
				if frame.Eos || frame.FrameType == tunnelpkg.FrameType_FRAME_STREAM_END {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, nil
}

func (tc *TunnelClient) Close() error {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	if tc.started {
		tc.tunnel.Stop()
		tc.started = false
	}
	return nil
}
