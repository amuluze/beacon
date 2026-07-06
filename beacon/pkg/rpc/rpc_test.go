package rpc

import (
	"context"
	"errors"
	"testing"

	"amprobe/pkg/contextx"
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
