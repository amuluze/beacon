// Package middleware
// Date: 2024/3/27 16:20
// Author: Amu
// Description:
package middleware

import (
	"github.com/amuluze/amprobe/pkg/auth"
	"github.com/amuluze/amprobe/pkg/contextx"
	"github.com/amuluze/amprobe/pkg/fiberx"
	"github.com/amuluze/amprobe/service/schema"
	"github.com/amuluze/amutool/errors"
	"github.com/gofiber/fiber/v2"
	"log/slog"
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
			if c.Path() == "/api/v1/auth/login" {
				a.RecordAudit(args.Username, "登录")
			}
			return c.Next()
		}
		slog.Info("auth middleware", "token", fiberx.GetToken(c))
		userID, username, isAdmin, err := a.ParseToken(fiberx.GetToken(c), "access_token")
		if errors.Is(err, auth.ErrInvalidToken) {
			slog.Error("invalid token", "err", err)
			return fiberx.Unauthorized(c)
		} else if err != nil {
			slog.Error("logout failed", "error", err)
			return fiberx.Unauthorized(c)
		}

		slog.Info("user id", "user_id", userID)
		wrapUserAuthContext(c, userID, username)
		if c.Method() == "POST" && isAdmin != "1" {
			return fiberx.Forbidden(c)
		}
		if c.Method() == "POST" && isAdmin == "1" {
			a.RecordAudit(username, OperateEvent[c.Path()])
		}
		return c.Next()
	}
}
