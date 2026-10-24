// Package middleware
// Date: 2024/3/27 16:20
// Author: Amu
// Description:
package middleware

import (
	"context"
	"log/slog"
	"net/url"
	"strings"

	"beacon/pkg/auth"
	"beacon/pkg/authcookie"
	"beacon/pkg/contextx"
	"beacon/pkg/errors"
	"beacon/pkg/fiberx"
	"beacon/service/schema"

	"github.com/gofiber/fiber/v2"
)

type PasswordChangeChecker interface {
	MustChangePassword(ctx context.Context, userID string) (bool, error)
}

var passwordChangeAllowedPaths = map[string]struct{}{
	"/api/v1/auth/logout":      {},
	"/api/v1/auth/pass_update": {},
	"/api/v1/auth/user_info":   {},
}

func wrapUserAuthContext(c *fiber.Ctx, userID string, username string) {
	ctx := contextx.NewUserID(c.UserContext(), userID)
	ctx = contextx.NewUsername(ctx, username)
	c.SetUserContext(ctx)
}

// CSRFMiddleware applies double-submit CSRF protection to browser sessions.
// Public machine-to-machine endpoints remain unaffected because they do not
// carry Beacon session cookies. Login always requires the pre-authentication
// CSRF cookie/header pair to prevent login-CSRF session swapping.
func CSRFMiddleware(manager *authcookie.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := manager.EnsureCSRFToken(c); err != nil {
			slog.Error("create csrf token failed", "error", err)
			return fiberx.Failure(c, errors.Error{Status: fiber.StatusInternalServerError, Msg: errors.InternalServerError})
		}
		if isSafeMethod(c.Method()) {
			return c.Next()
		}
		if c.Path() != "/api/v1/auth/login" && manager.AccessToken(c) == "" && manager.RefreshToken(c) == "" {
			return c.Next()
		}
		if err := manager.ValidateCSRF(c); err != nil {
			slog.Warn("csrf token rejected", "path", c.Path())
			return fiberx.Forbidden(c)
		}
		return c.Next()
	}
}

func isSafeMethod(method string) bool {
	switch method {
	case fiber.MethodGet, fiber.MethodHead, fiber.MethodOptions:
		return true
	default:
		return false
	}
}

// WebSocketUserAuthMiddleware authenticates a same-origin WebSocket handshake
// with the HttpOnly access-token cookie before the connection is upgraded.
// Identity is stored in both UserContext and Fiber locals because websocket.Conn
// carries locals across the upgrade.
func WebSocketUserAuthMiddleware(a auth.Auther, manager *authcookie.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !sameWebSocketOrigin(c) {
			slog.Warn("cross-origin websocket rejected")
			return fiberx.Unauthorized(c)
		}
		userID, username, err := a.ParseToken(manager.AccessToken(c), "access_token")
		if err != nil {
			slog.Warn("websocket session rejected")
			return fiberx.Unauthorized(c)
		}

		wrapUserAuthContext(c, userID, username)
		c.Locals("user_id", userID)
		c.Locals("username", username)
		return c.Next()
	}
}

func sameWebSocketOrigin(c *fiber.Ctx) bool {
	origin, err := url.Parse(c.Get(fiber.HeaderOrigin))
	if err != nil || origin.Scheme == "" || origin.Host == "" {
		return false
	}
	wantScheme := "http"
	if fiberx.IsHTTPS(c) {
		wantScheme = "https"
	}
	return strings.EqualFold(origin.Scheme, wantScheme) && strings.EqualFold(origin.Host, string(c.Context().Host()))
}

func UserAuthMiddleware(a auth.Auther, manager *authcookie.Manager, checker PasswordChangeChecker, skippers ...SkipperFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if SkipHandler(c, skippers...) {
			var args schema.LoginArgs
			_ = fiberx.ParseBody(c, &args)
			if c.Path() == "/api/v1/auth/login" {
				a.RecordAudit(args.Username, "登录")
			}
			return c.Next()
		}

		var userID string
		var username string
		var err error
		userID, username, err = a.ParseToken(manager.AccessToken(c), "access_token")

		if errors.Is(err, auth.ErrInvalidToken) {
			slog.Error("invalid token", "err", err)
			return fiberx.Unauthorized(c)
		} else if err != nil {
			slog.Error("token parse failed", "error", err)
			return fiberx.Unauthorized(c)
		}

		wrapUserAuthContext(c, userID, username)
		if checker != nil {
			mustChange, checkErr := checker.MustChangePassword(c.UserContext(), userID)
			if checkErr != nil {
				slog.Error("password change state lookup failed", "error", checkErr)
				return fiberx.Failure(c, errors.Error{Status: fiber.StatusInternalServerError, Msg: errors.InternalServerError})
			}
			if _, allowed := passwordChangeAllowedPaths[c.Path()]; mustChange && !allowed {
				return fiberx.Failure(c, errors.Error{Status: fiber.StatusForbidden, Msg: "password change required"})
			}
		}

		err = c.Next()
		if err == nil {
			if c.Method() == "POST" {
				a.RecordAudit(username, c.Route().Name)
			}
			return nil
		}
		return err
	}
}
