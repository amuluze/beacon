package fiberx

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// IsHTTPS reports whether the original client request used HTTPS. The first
// X-Forwarded-Proto value represents the client-facing hop when Beacon is
// behind the bundled reverse proxy.
func IsHTTPS(c *fiber.Ctx) bool {
	if forwarded := c.Get("X-Forwarded-Proto"); forwarded != "" {
		proto := strings.TrimSpace(strings.Split(forwarded, ",")[0])
		return strings.EqualFold(proto, "https")
	}
	return strings.EqualFold(c.Protocol(), "https")
}
