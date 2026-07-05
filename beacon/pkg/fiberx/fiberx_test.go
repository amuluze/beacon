package fiberx

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
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
