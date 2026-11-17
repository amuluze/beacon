package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func newReadOnlyApp() *fiber.App {
	app := fiber.New()
	app.Use(ReadOnlyMiddleware(false))
	app.Get("/read", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Post("/write", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Delete("/write", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Post("/api/v1/auth/login", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Post("/api/v1/auth/pass_update", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Post("/api/v1/host/report", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Post("/api/v1/host/file_download", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	return app
}

func TestReadOnlyMiddlewareAllowsReadsAndIngress(t *testing.T) {
	app := newReadOnlyApp()
	for _, tc := range []struct {
		method string
		path   string
	}{
		{fiber.MethodGet, "/read"},
		{fiber.MethodPost, "/api/v1/auth/login"},
		{fiber.MethodPost, "/api/v1/auth/pass_update"},
		{fiber.MethodPost, "/api/v1/host/report"},
		{fiber.MethodPost, "/api/v1/host/file_download"},
	} {
		resp, err := app.Test(httptest.NewRequest(tc.method, tc.path, nil))
		if err != nil {
			t.Fatalf("%s %s: %v", tc.method, tc.path, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("%s %s status = %d, want 200", tc.method, tc.path, resp.StatusCode)
		}
	}
}

func TestReadOnlyMiddlewareRejectsMutations(t *testing.T) {
	app := newReadOnlyApp()
	for _, tc := range []struct {
		method string
		path   string
	}{
		{fiber.MethodPost, "/write"},
		{fiber.MethodDelete, "/write"},
		{fiber.MethodPost, "/api/v1/container/container_create"},
		{fiber.MethodPost, "/api/v1/user/user_create"},
	} {
		resp, err := app.Test(httptest.NewRequest(tc.method, tc.path, nil))
		if err != nil {
			t.Fatalf("%s %s: %v", tc.method, tc.path, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("%s %s status = %d, want 403", tc.method, tc.path, resp.StatusCode)
		}
	}
}

func TestReadOnlyMiddlewareNoAuthRejectsAuthenticatedWrites(t *testing.T) {
	app := fiber.New()
	app.Use(ReadOnlyMiddleware(true))
	app.Post("/api/v1/auth/login", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Post("/api/v1/auth/pass_update", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Post("/api/v1/auth/logout", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	for _, path := range []string{"/api/v1/auth/login", "/api/v1/auth/pass_update"} {
		resp, err := app.Test(httptest.NewRequest(fiber.MethodPost, path, nil))
		if err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("%s status = %d, want 403", path, resp.StatusCode)
		}
	}

	resp, err := app.Test(httptest.NewRequest(fiber.MethodPost, "/api/v1/auth/logout", nil))
	if err != nil {
		t.Fatalf("logout: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("logout status = %d, want 200", resp.StatusCode)
	}
}
