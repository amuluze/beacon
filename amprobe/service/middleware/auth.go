// Package middleware
// Date: 2024/3/27 16:20
// Author: Amu
// Description:
package middleware

import (
	"log/slog"

	"amprobe/pkg/auth"

	"amprobe/pkg/contextx"

	"amprobe/pkg/fiberx"

	"amprobe/service/schema"

	"amprobe/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

func wrapUserAuthContext(c *fiber.Ctx, userID string, username string) {
	ctx := contextx.NewUserID(c.UserContext(), userID)
	ctx = contextx.NewUsername(c.UserContext(), username)
	c.SetUserContext(ctx)
}

func UserAuthMiddleware(a auth.Auther, skippers ...SkipperFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if SkipHandler(c, skippers...) {
			var args schema.LoginArgs
			_ = c.BodyParser(&args)
			if c.Path() == "/v1/auth/login" {
				a.RecordAudit(args.Username, "登录")
			}
			return c.Next()
		}

		slog.Info("auth middleware", "token", fiberx.GetToken(c))
		var userID string
		var username string
		var isAdmin string
		var err error
		userID, username, isAdmin, err = a.ParseToken(fiberx.GetToken(c), "access_token")

		if errors.Is(err, auth.ErrInvalidToken) {
			slog.Error("invalid token", "err", err)
			return fiberx.Unauthorized(c)
		} else if err != nil {
			slog.Error("logout failed", "error", err)
			return fiberx.Unauthorized(c)
		}

		slog.Info("user id", "user_id", userID)
		wrapUserAuthContext(c, userID, username)
		if c.Method() == "POST" && isAdmin != "1" && c.Path() != "/v1/auth/logout" {
			return fiberx.Forbidden(c)
		}
		if err := c.Next(); err == nil {
			if (c.Method() == "POST" && isAdmin == "1") || c.Path() == "/v1/auth/logout" {
				a.RecordAudit(username, OperateEvent[c.Path()])
			}
			return nil
		} else {
			return fiberx.Failure(c, errors.New400Error(err.Error()))
		}
	}
}
