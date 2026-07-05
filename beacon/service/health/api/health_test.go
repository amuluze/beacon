// Package health tests for Probe Liveness / Readiness.
package health

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// fakeDBPinger satisfies DBPinger for unit tests.
type fakeDBPinger struct{ err error }

func (f *fakeDBPinger) Ping() error { return f.err }

// fakeTunnelMonitor satisfies TunnelMonitor; we only care about whether the call
// is reached, not the value returned.
type fakeTunnelMonitor struct{ called bool }

func (f *fakeTunnelMonitor) AgentCount() int { f.called = true; return 1 }

func doRequest(t *testing.T, p *Probe, path string) (*http.Response, string) {
	t.Helper()
	app := fiber.New()
	app.Get("/health", p.Liveness)
	app.Get("/ready", p.Readiness)

	req := httptest.NewRequest("GET", path, nil)
	resp, err := app.Test(req, 200)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	return resp, string(body)
}

func TestLiveness_Always200(t *testing.T) {
	p := NewProbe()
	resp, body := doRequest(t, p, "/health")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if !strings.Contains(body, `"status":"alive"`) {
		t.Fatalf("missing alive status in body: %s", body)
	}
}

func TestReadiness_NoDeps(t *testing.T) {
	p := NewProbe()
	resp, body := doRequest(t, p, "/ready")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 (degraded-less), got %d body=%s", resp.StatusCode, body)
	}
	if !strings.Contains(body, `"status":"ready"`) {
		t.Fatalf("expected status=ready, got %s", body)
	}
}

func TestReadiness_DBError_Returns503(t *testing.T) {
	p := NewProbe()
	p.SetDB(&fakeDBPinger{err: errors.New("db down")})
	resp, body := doRequest(t, p, "/ready")
	if resp.StatusCode != 503 {
		t.Fatalf("expected 503, got %d body=%s", resp.StatusCode, body)
	}
	if !strings.Contains(body, `"status":"not_ready"`) {
		t.Fatalf("expected status=not_ready, got %s", body)
	}
}

func TestReadiness_DBSuccess_TunnelQueried(t *testing.T) {
	p := NewProbe()
	tm := &fakeTunnelMonitor{}
	p.SetDB(&fakeDBPinger{err: nil})
	p.SetTunnel(tm)
	resp, _ := doRequest(t, p, "/ready")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if !tm.called {
		t.Fatal("expected TunnelMonitor.AgentCount to be invoked when DB is healthy")
	}
}
