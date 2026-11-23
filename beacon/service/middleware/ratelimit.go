// Package middleware
// Date: 2026/06/27
// Author: Amu
// Description: 分层速率限制中间件，防止暴力破解与 DoS。
package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

const (
	defaultGlobalMax = 300 // 全局每 IP 每分钟兜底上限
	rateLimitWindow  = 1 * time.Minute
)

// limitReachedHandler 统一返回 429 响应，避免暴露内部细节。
func limitReachedHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "too many requests"})
}

// GlobalRateLimitMiddleware 对所有 /api/ 路由按 IP 做兜底限流。
// 跳过非 API 路径（前端静态资源）与登录类路径（由 router 内联严格限流器接管），
// 避免静态资源被计数、登录请求被双重计数。
func GlobalRateLimitMiddleware(max int) fiber.Handler {
	if max <= 0 {
		max = defaultGlobalMax
	}
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: rateLimitWindow,
		Next: func(c *fiber.Ctx) bool {
			path := c.Path()
			if !strings.HasPrefix(path, "/api/") {
				return true
			}
			// 登录类敏感端点交给严格限流器，避免双重计数。
			return hasAnyPrefix(path, authRateLimitedPrefixes...)
		},
		LimitReached: limitReachedHandler,
	})
}

// authRateLimitedPrefixes 是由登录限流器接管的敏感路径前缀。
var authRateLimitedPrefixes = []string{
	"/api/v1/auth/login",
	"/api/v1/auth/token_update",
}

// hasAnyPrefix 报告 path 是否以任一前缀开头。
func hasAnyPrefix(path string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
