// Package middleware
// Date: 2026/06/27
// Author: Amu
// Description: 分层速率限制中间件测试。
package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// fiberApp 构造一个挂载目标中间件的测试 app，返回其实例与 handler 的 fiber.Handler 包装。
func fiberApp(t *testing.T, mw fiber.Handler, method, path string) *fiber.App {
	t.Helper()
	app := fiber.New()
	app.Use(mw)
	app.Add(method, path, func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	return app
}

// TestGlobalRateLimitMiddleware_SkipsStaticAndAuth 验证全局限流跳过静态资源与登录路径。
func TestGlobalRateLimitMiddleware_SkipsStaticAndAuth(t *testing.T) {
	app := fiber.New()
	app.Use(GlobalRateLimitMiddleware(2))
	app.Get("/index.html", func(c *fiber.Ctx) error { return c.SendString("static") })
	app.Post("/api/v1/auth/login", func(c *fiber.Ctx) error { return c.SendString("login") })

	// 静态资源多次访问不应被限流。
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(fiber.MethodGet, "/index.html", nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("static request %d: %v", i, err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("static path should not be limited, got %d", resp.StatusCode)
		}
	}
	// 登录路径应由全局限流器跳过（交给登录限流器），多次访问仍非 429。
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/auth/login", nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("login request %d: %v", i, err)
		}
		if resp.StatusCode == fiber.StatusTooManyRequests {
			t.Fatalf("global limiter should skip auth path, got 429 on iter %d", i)
		}
	}
}
