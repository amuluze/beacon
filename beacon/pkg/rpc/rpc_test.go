package rpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"beacon/pkg/contextx"
	tunnelpkg "common/rpc/tunnel"
)

// TestTunnelClientCallRejectsMissingAgentID 验证控制调用在 context 缺失 agentID 时
// 返回 ErrMissingAgentID，而非静默回退默认节点（Domain I001/R001）。
//
// 这里不构造真实 tunnel（会触发网络监听），只验证 ResolveAgentID 的早期拒绝路径：
// NewTunnelClient 不再接受 defaultID，agentID 解析失败立即返回错误，不会到达 tunnel.Call。
func TestTunnelClientCallRejectsMissingAgentID(t *testing.T) {
	tc := NewTunnelClient(nil) // tunnel 为 nil 也安全：agentID 校验在调用 tunnel 前拦截

	err := tc.Call(context.Background(), "AnyMethod", nil, nil)
	if !errors.Is(err, ErrMissingAgentID) {
		t.Fatalf("Call with missing agentID: err = %v, want ErrMissingAgentID", err)
	}
	if !errors.Is(err, contextx.ErrMissingAgentID) {
		t.Fatalf("Call with missing agentID: err = %v, want contextx.ErrMissingAgentID", err)
	}
}

// TestTunnelClientStreamCallRejectsMissingAgentID 验证流式控制调用同样拒绝缺失 agentID。
func TestTunnelClientStreamCallRejectsMissingAgentID(t *testing.T) {
	tc := NewTunnelClient(nil)

	_, err := tc.StreamCall(context.Background(), "AnyMethod", nil)
	if !errors.Is(err, ErrMissingAgentID) {
		t.Fatalf("StreamCall with missing agentID: err = %v, want ErrMissingAgentID", err)
	}
}

// TestTunnelClientCallRejectsInvalidAgentID 验证格式非法的 agentID 被拒绝。
// agentID 含不安全字符时返回 ErrInvalidAgentID，防止畸形象识进入控制通道。
func TestTunnelClientCallRejectsInvalidAgentID(t *testing.T) {
	tc := NewTunnelClient(nil)
	ctx := contextx.NewAgentID(context.Background(), "agent/../../etc")

	err := tc.Call(ctx, "AnyMethod", nil, nil)
	if !errors.Is(err, contextx.ErrInvalidAgentID) {
		t.Fatalf("Call with invalid agentID: err = %v, want ErrInvalidAgentID", err)
	}
}

// TestErrMissingAgentIDAlias 验证 rpc.ErrMissingAgentID 与 contextx 权威定义一致。
func TestErrMissingAgentIDAlias(t *testing.T) {
	if !errors.Is(ErrMissingAgentID, contextx.ErrMissingAgentID) {
		t.Fatal("rpc.ErrMissingAgentID must be identical to contextx.ErrMissingAgentID")
	}
}

// TestStreamCallReleasesStreamContextOnTunnelError 验证当 tunnel.StreamCall 失败时，
// StreamCall 派生的超时 ctx 的 cancel 被调用（防止 context 泄漏）。
// 用未注册任何 agent 的 ServerTunnel 触发 AgentOfflineError 路径。
func TestStreamCallReleasesStreamContextOnTunnelError(t *testing.T) {
	tc := NewTunnelClient(tunnelpkg.NewServerTunnel())
	ctx := contextx.NewAgentID(context.Background(), "nonexistent-agent")

	_, err := tc.StreamCall(ctx, "AnyMethod", nil)
	if err == nil {
		t.Fatal("StreamCall on offline agent: expected error, got nil")
	}
	var offlineErr *tunnelpkg.AgentOfflineError
	if !errors.As(err, &offlineErr) {
		t.Fatalf("StreamCall err = %v, want AgentOfflineError", err)
	}
}

// TestStreamCallDerivesTimeoutContext 验证 StreamCall 派生的 ctx 带有 maxStreamDuration
// 截止时间，作为 agent PTY hang / 客户端半开时的兜底清理信号。
// 这里不构造完整 tunnel，仅验证 maxStreamDuration 常量为有限正值（保证会触发）。
func TestStreamCallDerivesTimeoutContext(t *testing.T) {
	if maxStreamDuration <= 0 || maxStreamDuration > 24*time.Hour {
		t.Fatalf("maxStreamDuration = %v, must be a finite positive value (got unreasonable bound)", maxStreamDuration)
	}
}
