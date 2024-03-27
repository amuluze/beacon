// Package middleware
// Date: 2024/3/27 16:20
// Author: Amu
// Description:
package middleware

import (
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

		tokenUserID, err := a.ParseUserID(fiberx.GetToken(c))
		if errors.Is(err, auth.ErrInvalidToken) {
			return fiberx.Unauthorized(c)
		}
		if err != nil {
			slog.Error("logout failed", "error", err)
			return fiberx.Failure(c, err)
		}

		slog.Info("user id %v", tokenUserID)
		wrapUserAuthContext(c, tokenUserID)
		return c.Next()
	}
}
