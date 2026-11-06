package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"beacon/pkg/authcookie"
	authservice "beacon/service/auth/service"
	"beacon/service/middleware"
	"beacon/service/testutil"

	"github.com/gofiber/fiber/v2"
)

func newTestAPI() (*AuthAPI, *testutil.FakeAuther) {
	auther := testutil.NewFakeAuther()
	service := authservice.NewAuthService(auther, testutil.NewFakeAuthRepo())
	return NewLoginAPI(service, authcookie.NewManager(true, 120, 3600)), auther
}

func TestLoginSetsHttpOnlyCookiesWithoutReturningTokens(t *testing.T) {
	api, _ := newTestAPI()
	app := fiber.New()
	app.Post("/api/v1/auth/login", api.Login)

	req := httptest.NewRequest(fiber.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"username":"admin","password":"password"}`))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if strings.Contains(string(body), "access_token") || strings.Contains(string(body), "refresh_token") {
		t.Fatalf("login response leaked token material: %s", body)
	}

	byName := make(map[string]*http.Cookie)
	for _, cookie := range resp.Cookies() {
		byName[cookie.Name] = cookie
	}
	for _, name := range []string{authcookie.AccessCookieName, authcookie.RefreshCookieName} {
		cookie := byName[name]
		if cookie == nil {
			t.Fatalf("missing %s", name)
		}
		if !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteStrictMode {
			t.Fatalf("%s attributes: HttpOnly=%v Secure=%v SameSite=%v", name, cookie.HttpOnly, cookie.Secure, cookie.SameSite)
		}
	}
}

func TestTokenUpdateReadsRefreshCookieAndReturnsNoTokenBody(t *testing.T) {
	api, auther := newTestAPI()
	var parsedToken, parsedType string
	auther.ParseTokenFn = func(token, tokenType string) (string, string, error) {
		parsedToken, parsedType = token, tokenType
		return "user-1", "admin", nil
	}
	app := fiber.New()
	app.Post("/api/v1/auth/token_update", api.TokenUpdate)

	req := httptest.NewRequest(fiber.MethodPost, "/api/v1/auth/token_update", nil)
	req.AddCookie(&http.Cookie{Name: authcookie.RefreshCookieName, Value: "refresh-secret"})
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("status = %d, want 204", resp.StatusCode)
	}
	if parsedToken != "refresh-secret" || parsedType != "refresh_token" {
		t.Fatalf("parsed token/type = %q/%q", parsedToken, parsedType)
	}
	if resp.ContentLength > 0 {
		t.Fatalf("refresh response must not expose a body, length=%d", resp.ContentLength)
	}
}

func TestLogoutClearsCookiesAndInvalidatesBothTokens(t *testing.T) {
	api, auther := newTestAPI()
	var destroyed []string
	auther.DestroyTokenFn = func(token string) error {
		destroyed = append(destroyed, token)
		return nil
	}
	app := fiber.New()
	app.Post("/api/v1/auth/logout", api.Logout)

	req := httptest.NewRequest(fiber.MethodPost, "/api/v1/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: authcookie.AccessCookieName, Value: "access-secret"})
	req.AddCookie(&http.Cookie{Name: authcookie.RefreshCookieName, Value: "refresh-secret"})
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("status = %d, want 204", resp.StatusCode)
	}
	if strings.Join(destroyed, ",") != "access-secret,refresh-secret" {
		t.Fatalf("destroyed tokens = %v", destroyed)
	}
	for _, cookie := range resp.Cookies() {
		if cookie.Name == authcookie.AccessCookieName || cookie.Name == authcookie.RefreshCookieName || cookie.Name == authcookie.CSRFCookieName {
			if cookie.Value != "" || !cookie.Expires.Before(time.Now()) {
				t.Fatalf("cookie %s was not expired: value=%q expires=%v", cookie.Name, cookie.Value, cookie.Expires)
			}
		}
	}
}

func TestRefreshWithoutAccessCookieStillRequiresCSRF(t *testing.T) {
	for _, tc := range []struct {
		name, cookie, header string
		want                 int
	}{
		{"matching", "csrf-test", "csrf-test", http.StatusNoContent},
		{"missing header", "csrf-test", "", http.StatusForbidden},
		{"missing cookie", "", "csrf-test", http.StatusForbidden},
		{"mismatch", "csrf-test", "different", http.StatusForbidden},
	} {
		t.Run(tc.name, func(t *testing.T) {
			api, auther := newTestAPI()
			calls := 0
			auther.ParseTokenFn = func(token, typ string) (string, string, error) {
				calls++
				if token != "refresh-test" || typ != "refresh_token" {
					t.Errorf("unexpected credential type")
				}
				return "user-1", "admin", nil
			}
			app := fiber.New()
			app.Use(middleware.CSRFMiddleware(api.CookieManager))
			app.Post("/api/v1/auth/token_update", api.TokenUpdate)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token_update", nil)
			req.AddCookie(&http.Cookie{Name: authcookie.RefreshCookieName, Value: "refresh-test"})
			if tc.cookie != "" {
				req.AddCookie(&http.Cookie{Name: authcookie.CSRFCookieName, Value: tc.cookie})
			}
			req.Header.Set(authcookie.CSRFHeaderName, tc.header)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != tc.want {
				t.Fatalf("status=%d want=%d", resp.StatusCode, tc.want)
			}
			if (calls == 1) != (tc.want == http.StatusNoContent) {
				t.Fatalf("unexpected refresh calls: %d", calls)
			}
		})
	}
}
