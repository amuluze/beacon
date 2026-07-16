// Package middleware
// Date: 2026/07/16
// Author: Amu
// Description: tests for panic and security headers middleware
package middleware

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// panic 响应体只应包含固定文案，不能泄露 panic 值或 Go 栈。
func TestPanicMiddleware_NoStackLeak(t *testing.T) {
	app := fiber.New()
	app.Use(PanicMiddleware())
	app.Get("/boom", func(c *fiber.Ctx) error { panic("explode-secret") })

	resp, err := app.Test(httptest.NewRequest("GET", "/boom", nil), -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != 500 {
		t.Errorf("状态码 = %d, want 500", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if strings.Contains(string(body), "explode-secret") || strings.Contains(string(body), "runtime") {
		t.Errorf("响应泄露了 panic/栈信息: %s", body)
	}
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(SecurityHeadersMiddleware())
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("ok") })

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil), -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	cases := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
	}
	for k, v := range cases {
		if got := resp.Header.Get(k); got != v {
			t.Errorf("header %s = %q, want %q", k, got, v)
		}
	}
}
