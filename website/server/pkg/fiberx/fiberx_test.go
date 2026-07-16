// Package fiberx
// Date: 2026/07/16
// Author: Amu
// Description: tests for fiberx helpers
package fiberx

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"server/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

type sampleBody struct {
	ID uint `json:"id" validate:"required"`
}

func TestGetToken(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString(GetToken(c)) })

	// 带 Bearer
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer abc")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	b, _ := io.ReadAll(resp.Body)
	if strings.TrimSpace(string(b)) != "abc" {
		t.Errorf("got %q, want abc", b)
	}

	// 无 header 不应 panic
	if _, err := app.Test(httptest.NewRequest("GET", "/", nil), -1); err != nil {
		t.Fatalf("无 Authorization 不应报错: %v", err)
	}
}

func TestParseBodyValidate_RejectsZeroID(t *testing.T) {
	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		var b sampleBody
		if err := ParseBodyValidate(c, &b); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		return c.SendString("ok")
	})

	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"id":0}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("id=0 应被校验拒绝，状态码 = %d, want 400", resp.StatusCode)
	}
}

func TestParseBodyValidate_Accepts(t *testing.T) {
	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		var b sampleBody
		if err := ParseBodyValidate(c, &b); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		return c.SendString("ok")
	})

	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"id":5}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("合法 id 状态码 = %d, want 200", resp.StatusCode)
	}
}

func TestSuccessAndFailure(t *testing.T) {
	app := fiber.New()
	app.Get("/ok", func(c *fiber.Ctx) error {
		return Success(c, fiber.Map{"hello": "world"})
	})
	app.Get("/fail", func(c *fiber.Ctx) error {
		return Failure(c, errors.New404Error("nope"))
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/ok", nil), -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Success 状态码 = %d, want 200", resp.StatusCode)
	}

	resp2, err := app.Test(httptest.NewRequest("GET", "/fail", nil), -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp2.StatusCode != 404 {
		t.Errorf("Failure 状态码 = %d, want 404", resp2.StatusCode)
	}
}
