package authcookie

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"time"

	"beacon/pkg/fiberx"

	"github.com/gofiber/fiber/v2"
)

const (
	AccessCookieName  = "beacon_access_token"
	RefreshCookieName = "beacon_refresh_token"
	CSRFCookieName    = "beacon_csrf_token"
	CSRFHeaderName    = "X-CSRF-Token"
)

var ErrInvalidCSRFToken = errors.New("invalid CSRF token")

// Manager owns the browser-session cookie policy. Access and refresh tokens
// are never exposed to JavaScript; the CSRF cookie intentionally is, so the
// browser client can echo it in CSRFHeaderName (double-submit protection).
type Manager struct {
	forceSecure       bool
	accessMaxAge      int
	refreshMaxAge     int
	generateCSRFToken func() (string, error)
}

func NewManager(forceSecure bool, accessMaxAge, refreshMaxAge int) *Manager {
	return &Manager{
		forceSecure:       forceSecure,
		accessMaxAge:      accessMaxAge,
		refreshMaxAge:     refreshMaxAge,
		generateCSRFToken: newCSRFToken,
	}
}

func (m *Manager) AccessToken(c *fiber.Ctx) string {
	return c.Cookies(AccessCookieName)
}

func (m *Manager) RefreshToken(c *fiber.Ctx) string {
	return c.Cookies(RefreshCookieName)
}

func (m *Manager) SetSession(c *fiber.Ctx, accessToken, refreshToken string) error {
	csrfToken, err := m.generateCSRFToken()
	if err != nil {
		return err
	}

	secure := m.forceSecure || fiberx.IsHTTPS(c)
	m.setCookie(c, AccessCookieName, accessToken, "/", m.accessMaxAge, true, secure)
	m.setCookie(c, RefreshCookieName, refreshToken, "/api/v1/auth", m.refreshMaxAge, true, secure)
	m.setCookie(c, CSRFCookieName, csrfToken, "/", m.refreshMaxAge, false, secure)
	return nil
}

// EnsureCSRFToken provisions the readable CSRF cookie on safe requests so a
// fresh browser can include it with its first authenticated state change.
func (m *Manager) EnsureCSRFToken(c *fiber.Ctx) error {
	if c.Cookies(CSRFCookieName) != "" {
		return nil
	}
	csrfToken, err := m.generateCSRFToken()
	if err != nil {
		return err
	}
	m.setCookie(c, CSRFCookieName, csrfToken, "/", m.refreshMaxAge, false, m.forceSecure || fiberx.IsHTTPS(c))
	return nil
}

func (m *Manager) ValidateCSRF(c *fiber.Ctx) error {
	cookieToken := c.Cookies(CSRFCookieName)
	headerToken := c.Get(CSRFHeaderName)
	if cookieToken == "" || headerToken == "" || len(cookieToken) != len(headerToken) {
		return ErrInvalidCSRFToken
	}
	if subtle.ConstantTimeCompare([]byte(cookieToken), []byte(headerToken)) != 1 {
		return ErrInvalidCSRFToken
	}
	return nil
}

func (m *Manager) ClearSession(c *fiber.Ctx) {
	secure := m.forceSecure || fiberx.IsHTTPS(c)
	m.clearCookie(c, AccessCookieName, "/", true, secure)
	m.clearCookie(c, RefreshCookieName, "/api/v1/auth", true, secure)
	m.clearCookie(c, CSRFCookieName, "/", false, secure)
}

func (m *Manager) setCookie(c *fiber.Ctx, name, value, path string, maxAge int, httpOnly, secure bool) {
	expires := time.Now().Add(time.Duration(maxAge) * time.Second)
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		MaxAge:   maxAge,
		Expires:  expires,
		Secure:   secure,
		HTTPOnly: httpOnly,
		SameSite: fiber.CookieSameSiteStrictMode,
	})
}

func (m *Manager) clearCookie(c *fiber.Ctx, name, path string, httpOnly, secure bool) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    "",
		Path:     path,
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
		Secure:   secure,
		HTTPOnly: httpOnly,
		SameSite: fiber.CookieSameSiteStrictMode,
	})
}

func newCSRFToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}
