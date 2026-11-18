// Package middleware
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package middleware

import (
	"beacon/pkg/contextx"
	"beacon/pkg/errors"
	"beacon/pkg/fiberx"
	"log/slog"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

func CasbinMiddleware(enforcer *casbin.SyncedEnforcer, skippers ...SkipperFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if SkipHandler(c, skippers...) {
			return c.Next()
		}
		return enforceCasbin(c, enforcer)
	}
}

// OptionalCasbinMiddleware enforces policies only for authenticated requests.
// Anonymous visitors in a no-auth deployment are limited by the read-only
// middleware instead of being rejected by Casbin before they can browse.
func OptionalCasbinMiddleware(enforcer *casbin.SyncedEnforcer, skippers ...SkipperFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if SkipHandler(c, skippers...) {
			return c.Next()
		}
		if contextx.FromUserID(c.UserContext()) == "" {
			return c.Next()
		}
		return enforceCasbin(c, enforcer)
	}
}

func enforceCasbin(c *fiber.Ctx, enforcer *casbin.SyncedEnforcer) error {
	ctx := c.UserContext()

	path := c.Path()
	method := c.Method()
	userID := contextx.FromUserID(ctx)

	slog.Debug("casbin middleware", "userID", userID, "path", path, "method", method)
	if enforces, err := enforcer.Enforce(userID, path, method); err != nil {
		return fiberx.Failure(c, errors.New500Error(err.Error()))
	} else if !enforces {
		return fiberx.Forbidden(c)
	}
	return c.Next()
}
