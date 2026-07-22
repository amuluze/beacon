package fiberx

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"beacon/pkg/contextx"
	pkgerrors "beacon/pkg/errors"
	tunnelpkg "common/rpc/tunnel"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TestGetTokenMissingAuthorization(t *testing.T) {
	app := fiber.New()
	var got string
	app.Get("/", func(c *fiber.Ctx) error {
		got = GetToken(c)
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest("GET", "/", nil)
	if _, err := app.Test(req); err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty token, got %q", got)
	}
}

func TestGetTokenBearerAuthorization(t *testing.T) {
	app := fiber.New()
	var got string
	app.Get("/", func(c *fiber.Ctx) error {
		got = GetToken(c)
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer access-token")
	if _, err := app.Test(req); err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if got != "access-token" {
		t.Fatalf("expected bearer token, got %q", got)
	}
}

func TestGetWebSocketTokenFromQuery(t *testing.T) {
	app := fiber.New()
	var got string
	app.Get("/ws/terminal", func(c *fiber.Ctx) error {
		got = GetWebSocketToken(c)
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest("GET", "/ws/terminal?token=query-token", nil)
	if _, err := app.Test(req); err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if got != "query-token" {
		t.Fatalf("expected query token, got %q", got)
	}
}

func TestGetWebSocketTokenPrefersAuthorizationHeader(t *testing.T) {
	app := fiber.New()
	var got string
	app.Get("/ws/terminal", func(c *fiber.Ctx) error {
		got = GetWebSocketToken(c)
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest("GET", "/ws/terminal?token=query-token", nil)
	req.Header.Set("Authorization", "Bearer header-token")
	if _, err := app.Test(req); err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if got != "header-token" {
		t.Fatalf("expected header token, got %q", got)
	}
}

// TestServiceError_Mapping 覆盖领域错误到 HTTP 状态码的映射，
// 验证 Domain R001（不可达/未实现必须返回可区分错误，禁止统一降级 500）。
func TestServiceError_Mapping(t *testing.T) {
	cases := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"agent offline", &tunnelpkg.AgentOfflineError{AgentID: "a-1"}, 503},
		{"wrapped agent offline", fmt.Errorf("call: %w", &tunnelpkg.AgentOfflineError{AgentID: "a-1"}), 503},
		{"missing agent id", contextx.ErrMissingAgentID, 400},
		{"wrapped missing agent id", fmt.Errorf("ctx: %w", contextx.ErrMissingAgentID), 400},
		{"invalid agent id", contextx.ErrInvalidAgentID, 400},
		{"record not found", gorm.ErrRecordNotFound, 404},
		{"wrapped record not found", fmt.Errorf("db: %w", gorm.ErrRecordNotFound), 404},
		{"deadline exceeded", context.DeadlineExceeded, 504},
		{"wrapped deadline", fmt.Errorf("rpc: %w", context.DeadlineExceeded), 504},
		{"generic error stays 500", fmt.Errorf("boom"), 500},
		{"predefined error passthrough", pkgerrors.New409Error("dup"), 409},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ServiceError(tc.err)
			if got.Status != tc.wantStatus {
				t.Fatalf("ServiceError(%v).Status = %d, want %d", tc.err, got.Status, tc.wantStatus)
			}
		})
	}
}
