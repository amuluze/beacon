package authcookie

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestSetSessionUsesProtectedCookiesAndNoTokenBody(t *testing.T) {
	manager := NewManager(true, 120, 3600)
	manager.generateCSRFToken = func() (string, error) { return "csrf-value", nil }
	app := fiber.New()
	app.Post("/login", func(c *fiber.Ctx) error {
		if err := manager.SetSession(c, "access-secret", "refresh-secret"); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"must_change_password": false})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodPost, "/login", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	if len(cookies) != 3 {
		t.Fatalf("cookie count = %d, want 3", len(cookies))
	}
	byName := make(map[string]string, len(cookies))
	for _, cookie := range cookies {
		byName[cookie.Name] = cookie.Raw
	}
	for _, name := range []string{AccessCookieName, RefreshCookieName} {
		raw := byName[name]
		if !strings.Contains(raw, "HttpOnly") || !strings.Contains(strings.ToLower(raw), "secure") || !strings.Contains(raw, "SameSite=Strict") {
			t.Fatalf("%s missing protected attributes: %s", name, raw)
		}
	}
	if strings.Contains(byName[CSRFCookieName], "HttpOnly") {
		t.Fatalf("CSRF cookie must remain readable by the browser client: %s", byName[CSRFCookieName])
	}
}

func TestValidateCSRFRequiresMatchingCookieAndHeader(t *testing.T) {
	manager := NewManager(false, 120, 3600)
	app := fiber.New()
	app.Post("/protected", func(c *fiber.Ctx) error {
		if err := manager.ValidateCSRF(c); err != nil {
			return c.SendStatus(fiber.StatusForbidden)
		}
		return c.SendStatus(fiber.StatusNoContent)
	})

	for _, tc := range []struct {
		name   string
		cookie string
		header string
		want   int
	}{
		{name: "matching", cookie: "csrf-value", header: "csrf-value", want: fiber.StatusNoContent},
		{name: "missing header", cookie: "csrf-value", want: fiber.StatusForbidden},
		{name: "mismatch", cookie: "csrf-value", header: "other", want: fiber.StatusForbidden},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(fiber.MethodPost, "/protected", nil)
			if tc.cookie != "" {
				req.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: tc.cookie})
			}
			if tc.header != "" {
				req.Header.Set(CSRFHeaderName, tc.header)
			}
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != tc.want {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tc.want)
			}
		})
	}
}
