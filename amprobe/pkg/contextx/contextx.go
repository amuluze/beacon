// Package contextx
// Date: 2024/3/27 16:28
// Author: Amu
// Description:
package contextx

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// ErrMissingAgentID 表示请求无法解析出目标 Agent 标识。
// 控制调用（rpc.Call）和监控查询（agentDB）在缺失 agentID 时都必须返回此错误，
// 而非静默回退为全表/默认节点，以避免跨 Agent 数据聚合或误操作。
var ErrMissingAgentID = errors.New("missing agent_id: unable to determine target agent")

// ErrInvalidAgentID 表示解析出的 agentID 格式非法（含不安全字符或超长）。
var ErrInvalidAgentID = errors.New("invalid agent_id: failed format validation")

// maxAgentIDLen 限制 agent_id 长度，避免超长值入库或用于查询。
const maxAgentIDLen = 64

// MaxAgentIDLen 返回 agent_id 允许的最大长度，供测试与边界断言复用。
func MaxAgentIDLen() int { return maxAgentIDLen }

// IsValidAgentID 校验 agent_id 仅含安全字符（字母、数字、-、_、.），且长度合理。
// 这是格式校验，不验证身份；身份绑定需 per-agent 凭证。
func IsValidAgentID(id string) bool {
	if len(id) == 0 || len(id) > maxAgentIDLen {
		return false
	}
	for _, r := range id {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '-', r == '_', r == '.':
		default:
			return false
		}
	}
	return true
}

// ResolveAgentID 从 context 解析并校验目标 Agent 标识。
// 返回值：
//   - agentID 非空且格式合法：返回 (agentID, nil)。
//   - agentID 缺失：返回 ("", ErrMissingAgentID)。
//   - agentID 格式非法：返回 ("", ErrInvalidAgentID)。
//
// 调用方（控制路径与监控查询路径）应统一使用本函数，保证空/非法 agentID
// 返回明确错误而非静默回退。
func ResolveAgentID(ctx context.Context) (string, error) {
	agentID := FromAgentID(ctx)
	if agentID == "" {
		return "", ErrMissingAgentID
	}
	if !IsValidAgentID(agentID) {
		return "", ErrInvalidAgentID
	}
	return agentID, nil
}

// AgentScopedDB 返回按 context 中 agent_id 过滤后的 DB 查询对象。
// 监控查询读路径应统一使用该 helper，保持与控制调用写路径一致：
// 缺失或格式非法的 agent_id 返回明确错误，不回退默认节点或全表查询。
func AgentScopedDB(ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	agentID, err := ResolveAgentID(ctx)
	if err != nil {
		return nil, err
	}
	return db.Where("agent_id = ?", agentID), nil
}

type (
	userIDCtx   struct{}
	usernameCtx struct{}
	agentIDCtx  struct{}
)

func NewUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDCtx{}, userID)
}

func FromUserID(ctx context.Context) string {
	v := ctx.Value(userIDCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameCtx{}, username)
}

func FromUsername(ctx context.Context) string {
	v := ctx.Value(usernameCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewAgentID(ctx context.Context, agentID string) context.Context {
	return context.WithValue(ctx, agentIDCtx{}, agentID)
}

func FromAgentID(ctx context.Context) string {
	v := ctx.Value(agentIDCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
