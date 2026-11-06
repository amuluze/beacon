package errors

import (
	"fmt"
	"net/http"
	"testing"

	"beacon/pkg/validatex"
	tunnel "common/rpc/tunnel"
)

func TestFromErrorMapsAgentOfflineToServiceUnavailable(t *testing.T) {
	got := FromError(&tunnel.AgentOfflineError{AgentID: "agent-a"})
	if got.Status != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", got.Status, http.StatusServiceUnavailable)
	}
	if got.Err == "" {
		t.Fatal("expected service error detail")
	}
}

func TestFromErrorMapsDuplicateAgentToConflict(t *testing.T) {
	got := FromError(&tunnel.DuplicateAgentError{AgentID: "agent-a"})
	if got.Status != http.StatusConflict {
		t.Fatalf("status = %d, want %d", got.Status, http.StatusConflict)
	}
}

// 请求参数校验失败是调用方输入错误，不能当作内部故障返回 500。
func TestFromErrorMapsRequestValidationFailureToBadRequest(t *testing.T) {
	type queryArgs struct {
		Page int `validate:"required,gte=1"`
	}

	validationErr := validatex.ValidateStruct(&queryArgs{})
	if validationErr == nil {
		t.Fatal("expected validation failure for zero page")
	}

	got := FromError(validationErr)
	if got.Status != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", got.Status, http.StatusBadRequest)
	}
	if got.Err == "" {
		t.Fatal("expected validation error detail")
	}
}

// 受保护容器的控制操作是拒绝语义，包装后仍必须映射为 403 而不是 500。
func TestFromErrorMapsProtectedContainerToForbidden(t *testing.T) {
	got := FromError(fmt.Errorf("%w: beacon", ErrProtectedContainer))
	if got.Status != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", got.Status, http.StatusForbidden)
	}
	if got.Err == "" {
		t.Fatal("expected service error detail")
	}
}
