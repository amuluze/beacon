package service

import (
	"net/http/httptest"
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

	resp, err := app.Test(httptest.NewRequest("GET", "/debug/pprof/", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("production: /debug/pprof/ status = %d, want 404 (pprof must be disabled)", resp.StatusCode)
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
