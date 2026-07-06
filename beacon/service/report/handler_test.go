package report

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"beacon/pkg/contextx"
	"beacon/service/model"
	"common/database"
	rpcSchema "common/rpc/schema"

	"github.com/gofiber/fiber/v2"
)

// newTestApp 构造一个只挂载 HandleReport 的最小 Fiber app。
// token 为空时跳过 token 校验；设置为非空时要求 X-Install-Token 匹配。
func newTestApp(t *testing.T, token string) (*fiber.App, func()) {
	t.Helper()
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	if err := db.AutoMigrate(
		new(model.MonitorHost), new(model.MonitorCPU), new(model.MonitorMemory),
		new(model.MonitorDisk), new(model.MonitorNet), new(model.MonitorDocker),
		new(model.MonitorContainer), new(model.MonitorImage), new(model.MonitorNetwork),
	); err != nil {
		db.Close()
		t.Fatalf("auto migrate: %v", err)
	}

	svc := NewService(db, token)
	app := fiber.New()
	app.Post("/api/v1/host/report", svc.HandleReport)
	return app, func() { db.Close() }
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func TestHandleReport_RejectsMissingAgentID(t *testing.T) {
	app, cleanup := newTestApp(t, "")
	defer cleanup()

	body := mustMarshal(rpcSchema.MonitorReportArgs{
		Host: rpcSchema.HostReport{
			Timestamp: time.Now(),
			Hostname:  "test-host",
		},
	})
	req := httptest.NewRequest("POST", "/api/v1/host/report", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", resp.StatusCode, string(respBody))
	}
	if !strings.Contains(string(respBody), contextx.ErrMissingAgentID.Error()) {
		t.Fatalf("expected ErrMissingAgentID in body, got: %s", string(respBody))
	}
}

func TestHandleReport_RejectsInvalidAgentID(t *testing.T) {
	app, cleanup := newTestApp(t, "")
	defer cleanup()

	body := mustMarshal(rpcSchema.MonitorReportArgs{
		AgentID: "agent/evil",
		Host: rpcSchema.HostReport{
			Timestamp: time.Now(),
		},
	})
	req := httptest.NewRequest("POST", "/api/v1/host/report", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400 for invalid agent_id, got %d", resp.StatusCode)
	}
}

func TestHandleReport_AcceptsValidBatch(t *testing.T) {
	app, cleanup := newTestApp(t, "")
	defer cleanup()

	now := time.Now()
	body := mustMarshal(rpcSchema.MonitorReportArgs{
		AgentID: "agent-a",
		Host: rpcSchema.HostReport{
			Timestamp: now, Hostname: "host-a",
		},
		CPU: rpcSchema.CPUReport{
			Timestamp: now, CPUPercent: 42.0,
		},
		Memory: rpcSchema.MemoryReport{
			Timestamp: now, MemPercent: 55.0, MemTotal: 8192, MemUsed: 4096,
		},
	})
	req := httptest.NewRequest("POST", "/api/v1/host/report", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(respBody))
	}
	if !strings.Contains(string(respBody), `"ok":true`) {
		t.Fatalf("expected ok response, got: %s", string(respBody))
	}
}

func TestHandleReport_RejectsInvalidToken(t *testing.T) {
	app, cleanup := newTestApp(t, "secret-token")
	defer cleanup()

	body := mustMarshal(rpcSchema.MonitorReportArgs{
		AgentID: "agent-a",
		Host: rpcSchema.HostReport{
			Timestamp: time.Now(),
		},
	})
	req := httptest.NewRequest("POST", "/api/v1/host/report", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	// 不设置 token header

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected 401 for missing token, got %d", resp.StatusCode)
	}
}

func TestHandleReport_AcceptsWithValidToken(t *testing.T) {
	app, cleanup := newTestApp(t, "secret-token")
	defer cleanup()

	now := time.Now()
	body := mustMarshal(rpcSchema.MonitorReportArgs{
		AgentID: "agent-a",
		Host: rpcSchema.HostReport{
			Timestamp: now, Hostname: "host-a",
		},
	})
	req := httptest.NewRequest("POST", "/api/v1/host/report", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Install-Token", "secret-token")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
