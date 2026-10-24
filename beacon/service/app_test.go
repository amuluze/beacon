package service

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// noopRouter 是测试用的最小 IRouter，不注册任何路由，
// 仅用于验证 NewFiberApp 的中间件挂载（如 pprof gate）。
type noopRouter struct{}

func (noopRouter) Register(app *fiber.App) error { return nil }
func (noopRouter) Prefixes() []string            { return nil }

// TestPprofDisabledInProduction 验证生产模式下 /debug/pprof 不被挂载（404），
// 防止运行时内存/goroutine 信息无鉴权泄漏或被用于 DoS。
func TestPprofDisabledInProduction(t *testing.T) {
	app := NewFiberApp(&Config{App: App{Env: "production"}}, noopRouter{})
	defer func() { _ = app.Shutdown() }()

	req := httptest.NewRequest("GET", "/debug/pprof/", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("production: /debug/pprof/ status = %d, want 404 (pprof must be disabled)", resp.StatusCode)
	}
}

func TestProductionRejectsNonHTTPSRequests(t *testing.T) {
	app := NewFiberApp(&Config{App: App{Env: "production"}}, noopRouter{})
	defer func() { _ = app.Shutdown() }()

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusUpgradeRequired {
		t.Fatalf("HTTP status = %d, want 426", resp.StatusCode)
	}
}

func TestProductionAddsBrowserSecurityHeaders(t *testing.T) {
	app := NewFiberApp(&Config{App: App{Env: "production"}}, noopRouter{})
	defer func() { _ = app.Shutdown() }()

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	for header, want := range map[string]string{
		"Content-Security-Policy":   "default-src 'self'",
		"Strict-Transport-Security": "max-age=31536000",
		"X-Frame-Options":           "DENY",
		"X-Content-Type-Options":    "nosniff",
		"Referrer-Policy":           "no-referrer",
		"Permissions-Policy":        "camera=()",
	} {
		if got := resp.Header.Get(header); !strings.Contains(got, want) {
			t.Fatalf("%s = %q, want substring %q", header, got, want)
		}
	}
}

// TestPprofEnabledInDevelopment 验证非生产模式 /debug/pprof 被挂载（非 404）。
func TestPprofEnabledInDevelopment(t *testing.T) {
	app := NewFiberApp(&Config{App: App{Env: "development"}}, noopRouter{})
	defer func() { _ = app.Shutdown() }()

	resp, err := app.Test(httptest.NewRequest("GET", "/debug/pprof/", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	// pprof index 返回 200；只要不是 404 即说明已挂载。
	if resp.StatusCode == fiber.StatusNotFound {
		t.Fatal("development: /debug/pprof/ returned 404, expected pprof to be mounted")
	}
}

func TestPprofDisabledWhenEnvironmentIsUnset(t *testing.T) {
	app := NewFiberApp(&Config{}, noopRouter{})
	defer func() { _ = app.Shutdown() }()

	resp, err := app.Test(httptest.NewRequest("GET", "/debug/pprof/", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("unset env: /debug/pprof/ status = %d, want 404", resp.StatusCode)
	}
}
