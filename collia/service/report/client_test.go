package report

import (
	"io"
	"net/http"
	"testing"

	rpcSchema "common/rpc/schema"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestPushAcceptsNilContext(t *testing.T) {
	var gotToken string

	client := NewClient("http://example.test/report", "secret")
	client.httpClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		gotToken = r.Header.Get("X-Install-Token")
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(http.NoBody),
			Header:     make(http.Header),
		}, nil
	})
	defer client.Close()

	err := client.Push(nil, rpcSchema.MonitorReportArgs{AgentID: "agent-a"})
	if err != nil {
		t.Fatalf("Push returned error: %v", err)
	}
	if gotToken != "secret" {
		t.Fatalf("expected install token header, got %q", gotToken)
	}
}
