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
	ctx = contextx.NewUsername(ctx, username)
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

		var userID string
		var username string
		var err error
		userID, username, err = a.ParseToken(fiberx.GetToken(c), "access_token")

		if errors.Is(err, auth.ErrInvalidToken) {
			slog.Error("invalid token", "err", err)
			return fiberx.Unauthorized(c)
		} else if err != nil {
			slog.Error("token parse failed", "error", err)
			return fiberx.Unauthorized(c)
		}

		slog.Info("user id", "user_id", userID)
		wrapUserAuthContext(c, userID, username)

		err = c.Next()
		if err == nil {
			if c.Method() == "POST" {
				a.RecordAudit(username, OperateEvent[c.Path()])
			}
			return nil
		}
		return err
	}
}
