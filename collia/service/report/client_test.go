package report

import (
	"net/http"
	"net/http/httptest"
	"testing"

	rpcSchema "common/rpc/schema"
)

func TestPushAcceptsNilContext(t *testing.T) {
	var gotToken string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotToken = r.Header.Get("X-Install-Token")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "secret")
	defer client.Close()

	err := client.Push(nil, rpcSchema.MonitorReportArgs{AgentID: "agent-a"})
	if err != nil {
		t.Fatalf("Push returned error: %v", err)
	}
	if gotToken != "secret" {
		t.Fatalf("expected install token header, got %q", gotToken)
	}
}
