package errors

import (
	"net/http"
	"testing"

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
