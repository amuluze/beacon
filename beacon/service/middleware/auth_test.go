package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"beacon/pkg/auth"
	"beacon/pkg/authcookie"
	"beacon/pkg/contextx"

	"github.com/gofiber/fiber/v2"
)

type passwordChangeCheckerStub struct {
	mustChange bool
	err        error
}

func (s passwordChangeCheckerStub) MustChangePassword(context.Context, string) (bool, error) {
	return s.mustChange, s.err
}

type websocketTestAuth struct {
	token string
	err   error
}

func (a *websocketTestAuth) GenerateToken(string, string) (auth.TokenInfo, error) {
	return nil, errors.New("not implemented")
}
func (a *websocketTestAuth) DestroyToken(string) error { return nil }
func (a *websocketTestAuth) ParseToken(token, tokenType string) (string, string, error) {
	a.token = token
	if a.err != nil {
		return "", "", a.err
	}
	if token == "" || tokenType != "access_token" {
		return "", "", auth.ErrInvalidToken
	}
	return "user-1", "admin", nil
}
func (a *websocketTestAuth) Release() error             { return nil }
func (a *websocketTestAuth) RecordAudit(string, string) {}

func TestWebSocketUserAuthMiddlewareUsesCookieAndCarriesIdentity(t *testing.T) {
	auther := &websocketTestAuth{}
	manager := authcookie.NewManager(false, 60, 120)
	app := fiber.New()
	app.Use(WebSocketUserAuthMiddleware(auther, manager))
	app.Get("/ws/terminal", func(c *fiber.Ctx) error {
		if got := contextx.FromUserID(c.UserContext()); got != "user-1" {
			t.Fatalf("user id in context = %q, want user-1", got)
		}
		if got := contextx.FromUsername(c.UserContext()); got != "admin" {
			t.Fatalf("username in context = %q, want admin", got)
		}
		if got := c.Locals("user_id"); got != "user-1" {
			t.Fatalf("user id local = %v, want user-1", got)
		}
		if got := c.Locals("username"); got != "admin" {
			t.Fatalf("username local = %v, want admin", got)
		}
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest("GET", "http://example.com/ws/terminal", nil)
	req.Header.Set("Origin", "http://example.com")
	req.AddCookie(&http.Cookie{Name: authcookie.AccessCookieName, Value: "access-token"})
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("status = %d, want 204", resp.StatusCode)
	}
	if auther.token != "access-token" {
		t.Fatalf("parsed token = %q, want access-token", auther.token)
	}
}

func TestWebSocketUserAuthMiddlewareRejectsMissingCookieOrCrossOrigin(t *testing.T) {
	for _, tc := range []struct {
		name        string
		origin      string
		accessToken string
	}{
		{name: "missing cookie", origin: "http://example.com"},
		{name: "cross origin", origin: "http://evil.example", accessToken: "access-token"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			app := fiber.New()
			app.Use(WebSocketUserAuthMiddleware(&websocketTestAuth{}, authcookie.NewManager(false, 60, 120)))
			app.Get("/ws/terminal", func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusNoContent)
			})

			req := httptest.NewRequest("GET", "http://example.com/ws/terminal", nil)
			req.Header.Set("Origin", tc.origin)
			if tc.accessToken != "" {
				req.AddCookie(&http.Cookie{Name: authcookie.AccessCookieName, Value: tc.accessToken})
			}
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != fiber.StatusUnauthorized {
				t.Fatalf("status = %d, want 401", resp.StatusCode)
			}
		})
	}
}

var _ auth.Auther = (*websocketTestAuth)(nil)

func TestUserAuthMiddlewareRestrictsMandatoryPasswordChangeSession(t *testing.T) {
	auther := &websocketTestAuth{}
	manager := authcookie.NewManager(false, 60, 120)
	app := fiber.New()
	app.Use(UserAuthMiddleware(auther, manager, passwordChangeCheckerStub{mustChange: true}))
	app.Get("/api/v1/private", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })
	app.Post("/api/v1/auth/pass_update", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	request := httptest.NewRequest(fiber.MethodGet, "/api/v1/private", nil)
	request.AddCookie(&http.Cookie{Name: authcookie.AccessCookieName, Value: "access-token"})
	resp, err := app.Test(request)
	if err != nil {
		t.Fatalf("protected request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("protected request status = %d, want 403", resp.StatusCode)
	}

	request = httptest.NewRequest(fiber.MethodPost, "/api/v1/auth/pass_update", nil)
	request.AddCookie(&http.Cookie{Name: authcookie.AccessCookieName, Value: "access-token"})
	resp, err = app.Test(request)
	if err != nil {
		t.Fatalf("password update request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("password update status = %d, want 204", resp.StatusCode)
	}
}

func TestUserAuthMiddlewareFailsClosedWhenPasswordStateLookupFails(t *testing.T) {
	app := fiber.New()
	app.Use(UserAuthMiddleware(&websocketTestAuth{}, authcookie.NewManager(false, 60, 120), passwordChangeCheckerStub{err: errors.New("database secret detail")}))
	app.Get("/api/v1/private", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	request := httptest.NewRequest(fiber.MethodGet, "/api/v1/private", nil)
	request.AddCookie(&http.Cookie{Name: authcookie.AccessCookieName, Value: "access-token"})
	resp, err := app.Test(request)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", resp.StatusCode)
	}
}

func TestCSRFMiddlewareProtectsCookieAuthenticatedMutations(t *testing.T) {
	manager := authcookie.NewManager(false, 60, 120)
	app := fiber.New()
	app.Use(CSRFMiddleware(manager))
	app.Post("/api/v1/private", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	req := httptest.NewRequest(fiber.MethodPost, "/api/v1/private", nil)
	req.AddCookie(&http.Cookie{Name: authcookie.AccessCookieName, Value: "access-token"})
	req.AddCookie(&http.Cookie{Name: authcookie.CSRFCookieName, Value: "csrf-token"})
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("missing CSRF header status = %d, want 403", resp.StatusCode)
	}

	req = httptest.NewRequest(fiber.MethodPost, "/api/v1/private", nil)
	req.AddCookie(&http.Cookie{Name: authcookie.AccessCookieName, Value: "access-token"})
	req.AddCookie(&http.Cookie{Name: authcookie.CSRFCookieName, Value: "csrf-token"})
	req.Header.Set(authcookie.CSRFHeaderName, "csrf-token")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("matching CSRF token status = %d, want 204", resp.StatusCode)
	}
}

func TestCSRFMiddlewareProtectsLoginWithoutExistingSession(t *testing.T) {
	manager := authcookie.NewManager(false, 60, 120)
	app := fiber.New()
	app.Use(CSRFMiddleware(manager))
	app.Post("/api/v1/auth/login", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	resp, err := app.Test(httptest.NewRequest(fiber.MethodPost, "/api/v1/auth/login", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("login without CSRF proof status = %d, want 403", resp.StatusCode)
	}

	req := httptest.NewRequest(fiber.MethodPost, "/api/v1/auth/login", nil)
	req.AddCookie(&http.Cookie{Name: authcookie.CSRFCookieName, Value: "pre-auth-csrf"})
	req.Header.Set(authcookie.CSRFHeaderName, "pre-auth-csrf")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("login with CSRF proof status = %d, want 204", resp.StatusCode)
	}
}
