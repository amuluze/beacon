// Package middleware
// Date: 2024/3/27 16:21
// Author: Amu
// Description:
package middleware

import (
	"github.com/gofiber/fiber/v2"
)

type SkipperFunc func(c *fiber.Ctx) bool

func AllowPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *fiber.Ctx) bool {
		path := c.Path()
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

func SkipHandler(c *fiber.Ctx, skippers ...SkipperFunc) bool {
	for _, skipper := range skippers {
		if skipper(c) {
			return true
		}
	}
	return false
}
