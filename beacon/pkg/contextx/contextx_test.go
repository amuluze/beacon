package contextx

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// TestResolveAgentID_Missing 验证 context 缺失 agentID 时返回 ErrMissingAgentID。
func TestResolveAgentID_Missing(t *testing.T) {
	id, err := ResolveAgentID(context.Background())
	if !errors.Is(err, ErrMissingAgentID) {
		t.Fatalf("err = %v, want ErrMissingAgentID", err)
	}
	if id != "" {
		t.Fatalf("id = %q, want empty", id)
	}
}

// TestResolveAgentID_Invalid 验证格式非法的 agentID 返回 ErrInvalidAgentID。
func TestResolveAgentID_Invalid(t *testing.T) {
	ctx := NewAgentID(context.Background(), "agent/evil")
	id, err := ResolveAgentID(ctx)
	if !errors.Is(err, ErrInvalidAgentID) {
		t.Fatalf("err = %v, want ErrInvalidAgentID", err)
	}
	if id != "" {
		t.Fatalf("id = %q, want empty", id)
	}
}

// TestResolveAgentID_Valid 验证合法 agentID 正确解析。
func TestResolveAgentID_Valid(t *testing.T) {
	ctx := NewAgentID(context.Background(), "agent-a.1_2")
	id, err := ResolveAgentID(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "agent-a.1_2" {
		t.Fatalf("id = %q, want agent-a.1_2", id)
	}
}

// TestIsValidAgentID 覆盖 agent_id 格式校验的边界。
func TestIsValidAgentID(t *testing.T) {
	tests := []struct {
		id   string
		want bool
	}{
		{"", false},
		{"agent-a", true},
		{"node_1.example", true},
		{"HOST", true},
		{"123", true},
		{"agent a", false},
		{"agent/../../etc", false},
		{"agent;drop", false},
		{"agent\x00null", false},
		{"中文", false},
		{strings.Repeat("a", maxAgentIDLen), true},
		{strings.Repeat("a", maxAgentIDLen+1), false},
	}
	for _, tc := range tests {
		if got := IsValidAgentID(tc.id); got != tc.want {
			t.Errorf("IsValidAgentID(%q) = %v, want %v", tc.id, got, tc.want)
		}
	}
}
