// Package middleware
// Date: 2026/6/26
// Author: Amu
// Description: unit tests for panic recovery middleware
package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestPanicMiddleware_RecoversPanic(t *testing.T) {
	app := fiber.New()
	app.Use(PanicMiddleware())
	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 500 {
		t.Errorf("status = %d, want 500", resp.StatusCode)
	}
}

func TestPanicMiddleware_NormalRequest(t *testing.T) {
	app := fiber.New()
	app.Use(PanicMiddleware())
	app.Get("/ok", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/ok", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}
