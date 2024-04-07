// Package middleware
// Date: 2024/3/27 16:20
// Author: Amu
// Description:
package middleware

import (
	"fmt"
	"github.com/amuluze/amprobe/pkg/auth"
	"github.com/amuluze/amprobe/pkg/contextx"
	"github.com/amuluze/amprobe/pkg/fiberx"
	"github.com/amuluze/amutool/errors"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func wrapUserAuthContext(c *fiber.Ctx, userID string) {
	ctx := contextx.NewUserID(c.UserContext(), userID)
	c.SetUserContext(ctx)
}

func UserAuthMiddleware(a auth.Auther, skippers ...SkipperFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if SkipHandler(c, skippers...) {
			return c.Next()
		}
		slog.Info("auth middleware", "token", fiberx.GetToken(c))
		tokenUserID, err := a.ParseUserID(fiberx.GetToken(c))
		if errors.Is(err, auth.ErrInvalidToken) {
			slog.Error("invalid token", "err", err)
			return fiberx.Unauthorized(c)
		}
		slog.Info("auth middleware qqq")
		if err != nil {
			slog.Error("logout failed", "error", err)
			return fiberx.Failure(c, err)
		}

		slog.Info("user id %v", tokenUserID)
		fmt.Println(c.Method(), c.Path(), c.Route().Name, tokenUserID)
		wrapUserAuthContext(c, tokenUserID)
		return c.Next()
	}
}
