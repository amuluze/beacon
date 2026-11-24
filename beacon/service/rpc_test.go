package service

import (
	"errors"
	"net/http/httptest"
	"path/filepath"
	"testing"

	beaconrpc "beacon/pkg/rpc"
	"common/database"

	"github.com/gofiber/fiber/v2"
)

func TestNewRPCClientHonorsDisabledControl(t *testing.T) {
	globalTunnel = nil
	client, err := NewRPCClient(&Config{Control: Control{Enable: false}}, nil, nil)
	if err != nil {
		t.Fatalf("NewRPCClient() error = %v", err)
	}
	if _, ok := client.(*beaconrpc.DisabledClient); !ok {
		t.Fatalf("NewRPCClient() type = %T, want *rpc.DisabledClient", client)
	}
	if globalTunnel != nil {
		t.Fatal("disabled control must not create a tunnel")
	}
}

func TestNewRPCClientRejectsWeakProductionJoinToken(t *testing.T) {
	_, err := NewRPCClient(&Config{
		App:     App{Env: productionEnv},
		Control: Control{Enable: true, JoinToken: "change-me"},
	}, nil, nil)
	if !errors.Is(err, ErrInsecureControlToken) {
		t.Fatalf("NewRPCClient() error = %v, want ErrInsecureControlToken", err)
	}
}

// TestNewHealthProbeWiresDBForReadiness 验证探针真实接线数据库：
// DB 存活时 /ready 为 200，失联后必须翻转为 503 而不是恒为 ready。
func TestNewHealthProbeWiresDBForReadiness(t *testing.T) {
	globalTunnel = nil
	client, err := NewRPCClient(&Config{Control: Control{Enable: false}}, nil, nil)
	if err != nil {
		t.Fatalf("NewRPCClient: %v", err)
	}
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "health")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}

	probe := NewHealthProbe(client, db)
	app := fiber.New()
	app.Get("/ready", probe.Readiness)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/ready", nil), 2000)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("ready with live db = %d, want 200", resp.StatusCode)
	}

	db.Close()
	resp, err = app.Test(httptest.NewRequest(fiber.MethodGet, "/ready", nil), 2000)
	if err != nil {
		t.Fatalf("app.Test after close: %v", err)
	}
	if resp.StatusCode != fiber.StatusServiceUnavailable {
		t.Fatalf("ready with closed db = %d, want 503", resp.StatusCode)
	}
}
