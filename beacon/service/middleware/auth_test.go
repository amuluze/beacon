package middleware

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"beacon/pkg/auth"
	"beacon/pkg/authcookie"
	"beacon/pkg/contextx"
	"beacon/service/testutil"

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
func (a *websocketTestAuth) Release() error                                            { return nil }
func (a *websocketTestAuth) RecordAudit(context.Context, string, string, string) error { return nil }

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

// 审计只记录成功的具名操作：403 拒绝与未命名路由不得写入空操作记录。
func TestUserAuthMiddlewareAuditsOnlySuccessfulNamedPosts(t *testing.T) {
	auther := testutil.NewFakeAuther()
	manager := authcookie.NewManager(false, 60, 120)
	var recorded []string
	auther.RecordAuditFn = func(_ context.Context, _ string, operate string, _ string) error {
		recorded = append(recorded, operate)
		return nil
	}

	app := fiber.New()
	app.Use(UserAuthMiddleware(auther, manager, nil))
	app.Post("/api/v1/network_create", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	}).Name("创建网络")
	app.Post("/api/v1/network_delete", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusForbidden)
	}).Name("删除网络")
	app.Post("/api/v1/unnamed", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	call := func(path string) {
		req := httptest.NewRequest(fiber.MethodPost, path, nil)
		req.AddCookie(&http.Cookie{Name: authcookie.AccessCookieName, Value: "access"})
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request %s failed: %v", path, err)
		}
		defer resp.Body.Close()
	}
	call("/api/v1/network_create")
	call("/api/v1/network_delete")
	call("/api/v1/unnamed")

	if len(recorded) != 1 || recorded[0] != "创建网络" {
		t.Fatalf("recorded audits = %v, want only 创建网络", recorded)
	}
}

func TestOptionalUserAuthMiddlewareAnonymousAndAuthenticated(t *testing.T) {
	auther := &websocketTestAuth{}
	manager := authcookie.NewManager(false, 60, 120)
	app := fiber.New()
	app.Use(OptionalUserAuthMiddleware(auther, manager, nil))
	app.Get("/whoami", func(c *fiber.Ctx) error {
		return c.SendString(contextx.FromUserID(c.UserContext()))
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/whoami", nil))
	if err != nil {
		t.Fatalf("anonymous request failed: %v", err)
	}
	if body := readTestBody(t, resp); body != "" {
		t.Fatalf("anonymous user id = %q, want empty", body)
	}

	req := httptest.NewRequest(fiber.MethodGet, "/whoami", nil)
	req.AddCookie(&http.Cookie{Name: authcookie.AccessCookieName, Value: "access-token"})
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("authenticated request failed: %v", err)
	}
	if body := readTestBody(t, resp); body != "user-1" {
		t.Fatalf("authenticated user id = %q, want user-1", body)
	}

	req = httptest.NewRequest(fiber.MethodGet, "/whoami", nil)
	req.AddCookie(&http.Cookie{Name: authcookie.RefreshCookieName, Value: "refresh-token"})
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("refresh-only request failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("refresh-only status = %d, want 401", resp.StatusCode)
	}
}

func readTestBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return string(body)
}
