// Package middleware
// Date: 2026/07/16
// Author: Amu
// Description:
package middleware

import "github.com/gofiber/fiber/v2"

// SecurityHeadersMiddleware 注入基础安全响应头，收敛点击劫持、MIME 嗅探、协议降级等攻击面。
// CSP 的强约束需与前端资源策略对齐，暂不设置以避免误拦预渲染资源。
func SecurityHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		// HSTS 仅在 HTTPS 下被浏览器采纳，HTTP 场景无害
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		return c.Next()
	}
}
