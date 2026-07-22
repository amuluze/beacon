package middleware

import (
	"errors"
	"net/http/httptest"
	"testing"

	"beacon/pkg/auth"
	"beacon/pkg/contextx"

	"github.com/gofiber/fiber/v2"
)

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

func TestWebSocketUserAuthMiddlewareAcceptsQueryTokenAndCarriesIdentity(t *testing.T) {
	auther := &websocketTestAuth{}
	app := fiber.New()
	app.Use(WebSocketUserAuthMiddleware(auther))
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

	resp, err := app.Test(httptest.NewRequest("GET", "/ws/terminal?token=query-token", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("status = %d, want 204", resp.StatusCode)
	}
	if auther.token != "query-token" {
		t.Fatalf("parsed token = %q, want query-token", auther.token)
	}
}

func TestWebSocketUserAuthMiddlewareRejectsMissingOrInvalidToken(t *testing.T) {
	for _, tc := range []struct {
		name string
		url  string
		err  error
	}{
		{name: "missing", url: "/ws/terminal"},
		{name: "invalid", url: "/ws/terminal?token=bad", err: auth.ErrInvalidToken},
	} {
		t.Run(tc.name, func(t *testing.T) {
			auther := &websocketTestAuth{err: tc.err}
			app := fiber.New()
			app.Use(WebSocketUserAuthMiddleware(auther))
			app.Get("/ws/terminal", func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusNoContent)
			})

			resp, err := app.Test(httptest.NewRequest("GET", tc.url, nil))
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
