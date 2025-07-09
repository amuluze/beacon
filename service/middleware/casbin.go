// Package middleware
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package middleware

import (
	"amprobe/pkg/contextx"
	"amprobe/pkg/errors"
	"amprobe/pkg/fiberx"
	"log/slog"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

func CasbinMiddleware(enforcer *casbin.SyncedEnforcer, skippers ...SkipperFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if SkipHandler(c, skippers...) {
			return c.Next()
		}
		ctx := c.UserContext()

		path := c.Path()
		method := c.Method()
		userID := contextx.FromUserID(ctx)

		slog.Debug("casbin middleware", "userID", userID, "path", path, "method", method)
		if enforces, err := enforcer.Enforce(userID, path, method); err != nil {
			return fiberx.Failure(c, errors.New400Error(err.Error()))
		} else if !enforces {
			return fiberx.Forbidden(c)
		}
		return c.Next()
	}
}
