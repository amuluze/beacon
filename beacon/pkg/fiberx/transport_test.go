package fiberx

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestIsHTTPSUsesOriginalForwardedProtocol(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		if !IsHTTPS(c) {
			t.Fatal("request should be recognized as HTTPS")
		}
		return c.SendStatus(fiber.StatusNoContent)
	})
	req := httptest.NewRequest(fiber.MethodGet, "http://beacon.internal/", nil)
	req.Header.Set("X-Forwarded-Proto", "https, http")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("status = %d, want 204", resp.StatusCode)
	}
}
