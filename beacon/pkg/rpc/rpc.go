// Package rpc
// Date: 2024/6/12 10:30
// Author: Amu
// Description:
package rpc

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"beacon/pkg/contextx"
	tunnelpkg "common/rpc/tunnel"
)

const (
	DefaultNetwork = "tcp"
	// maxStreamDuration 是所有流式 RPC 的最长生命周期兜底。
	// 业务侧 ctx（如 WebSocket）取消会立刻释放资源；这里防止 agent PTY hang
	// 且客户端半开时 ctx 永不取消，导致 server streams map + 缓冲 + goroutine 永驻。
	// 超过后强制关闭流，触发 server.StreamCall 的清理路径。
	maxStreamDuration = 2 * time.Hour
)

// ErrMissingAgentID 在控制调用无法解析出目标 Agent 时返回。
// 权威定义在 contextx.ErrMissingAgentID，此处保留为兼容别名，
// 避免破坏已有的包外引用（如 repository、terminal 等）。
var ErrMissingAgentID = contextx.ErrMissingAgentID

// ErrControlDisabled 表示部署显式关闭了 Agent 控制通道。
var ErrControlDisabled = errors.New("agent control channel is disabled")

// Caller is the interface for making RPC calls to agents.
type Caller interface {
	Call(ctx context.Context, method string, args interface{}, reply interface{}) error
	StreamCall(ctx context.Context, method string, args interface{}) (<-chan []byte, error)
	Close() error
}

// DisabledClient 在 Control.Enable=false 时提供稳定、可诊断的 Caller。
type DisabledClient struct{}

func NewDisabledClient() *DisabledClient { return &DisabledClient{} }

func (*DisabledClient) Call(context.Context, string, interface{}, interface{}) error {
	return ErrControlDisabled
}

func (*DisabledClient) StreamCall(context.Context, string, interface{}) (<-chan []byte, error) {
	return nil, ErrControlDisabled
}

func (*DisabledClient) Close() error { return nil }

// TunnelClient wraps ServerTunnel to implement the Caller interface.
type TunnelClient struct {
	tunnel  *tunnelpkg.ServerTunnel
	mu      sync.Mutex
	started bool
}

// NewTunnelClient creates a new tunnel-based RPC caller.
func NewTunnelClient(tunnel *tunnelpkg.ServerTunnel) *TunnelClient {
	return &TunnelClient{tunnel: tunnel, started: true}
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
	// 兜底超时：即便调用方 ctx 永不取消（agent hang / 客户端半开），
	// maxStreamDuration 后也会取消，触发 server 侧 streams map 与缓冲的回收。
	streamCtx, cancel := context.WithTimeout(ctx, maxStreamDuration)
	chunkChan, err := tc.tunnel.StreamCall(streamCtx, agentID, method, args)
	if err != nil {
		cancel()
		return nil, err
	}
	// Convert tunnel Frame channel to []byte channel
	out := make(chan []byte, 64)
	go func() {
		defer close(out)
		defer cancel()
		for {
			select {
			case frame, ok := <-chunkChan:
				if !ok {
					return
				}
				out <- frame.Payload
				if frame.Eos || frame.FrameType == tunnelpkg.FrameType_FRAME_STREAM_END {
					if frame.Error != "" {
						slog.Warn("tunnel stream ended incomplete", "error", frame.Error)
					}
					return
				}
			case <-streamCtx.Done():
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
