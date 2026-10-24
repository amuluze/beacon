// Package service
// Date: 2024/3/6 11:07
// Author: Amu
// Description:
package service

import (
	"beacon/pkg/fiberx"
	"beacon/service/middleware"
	"beacon/web"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

// fiber 默认超时常量。配置缺值时兜底，确保慢速/半开连接不会无限占用内存。
// 单批监控上报远小于 32MB，1GB 的历史默认值会让单个异常请求耗尽容器内存。
const (
	defaultReadTimeout  = 30 * time.Second
	defaultWriteTimeout = 30 * time.Second
	defaultIdleTimeout  = 120 * time.Second
	defaultBodyLimit    = 32 * 1024 * 1024 // 32MB
)

func NewFiberApp(config *Config, r IRouter) *fiber.App {
	fiberConfig := fiber.Config{
		Prefork:      config.Fiber.Prefork,
		AppName:      config.Fiber.AppName,
		ServerHeader: config.Fiber.ServerHeader,
		ReadTimeout:  durationSeconds(config.Fiber.ReadTimeout, defaultReadTimeout),
		WriteTimeout: durationSeconds(config.Fiber.WriteTimeout, defaultWriteTimeout),
		IdleTimeout:  durationSeconds(config.Fiber.IdleTimeout, defaultIdleTimeout),
		BodyLimit:    defaultBodyLimit,
	}

	app := fiber.New(fiberConfig)

	// 添加中间件
	app.Use(securityHeadersMiddleware(config))
	app.Use(cors.New(buildCORSConfig(config.CORS)))
	if strings.EqualFold(config.App.Env, productionEnv) {
		app.Use(func(c *fiber.Ctx) error {
			if c.Path() == "/health" || c.Path() == "/ready" || fiberx.IsHTTPS(c) {
				return c.Next()
			}
			return c.Status(fiber.StatusUpgradeRequired).JSON(fiber.Map{"error": "HTTPS required"})
		})
	}
	// 全局兜底限流：仅作用于 /api/ 路由，按 IP 计数；登录类由独立严格限流器接管。
	if config.RateLimit.Enable {
		app.Use(middleware.GlobalRateLimitMiddleware(config.RateLimit.GlobalMax))
	}
	app.Use(compress.New())
	// pprof 暴露运行时内存/goroutine/profile，可泄漏密钥或被用于 DoS。
	// 仅在明确的 development 模式启用；空值和未知环境均按安全默认禁用。
	if strings.EqualFold(config.App.Env, "development") {
		app.Use(pprof.New())
	}
	app.Use(middleware.PanicMiddleware())
	app.Use(middleware.StackMiddleware)

	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(web.FS),
		PathPrefix: "/dist",
		Browse:     true,
	}))

	err := r.Register(app)
	if err != nil {
		panic(err)
	}

	return app
}

func securityHeadersMiddleware(config *Config) fiber.Handler {
	production := strings.EqualFold(config.App.Env, productionEnv)
	return func(c *fiber.Ctx) error {
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("Referrer-Policy", "no-referrer")
		c.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=(), usb=()")
		if production {
			c.Set("Content-Security-Policy", "default-src 'self'; base-uri 'self'; object-src 'none'; frame-ancestors 'none'; form-action 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: blob:; font-src 'self' data:; connect-src 'self' wss:")
			if fiberx.IsHTTPS(c) {
				c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}
		}
		return c.Next()
	}
}

// devAllowOrigins 是未配置 CORS 白名单时的本地开发回退域（Vite 常用端口）。
// 生产部署应显式配置 [CORS].AllowOrigins，避免回退到宽松的开发域。
var devAllowOrigins = []string{
	"http://localhost:5173",
	"http://127.0.0.1:5173",
	"http://localhost:3000",
	"http://127.0.0.1:3000",
}

// buildCORSConfig 根据配置生成 CORS 中间件配置。
//   - Enable=false：返回空配置（fiber cors 默认放行所有 Origin），保持同源可用。
//   - AllowOrigins 非空：严格限定到白名单。
//   - AllowOrigins 为空：回退到本地开发域（仅适用于开发模式）。
//
// HttpOnly Cookie 跨域依赖 AllowCredentials 与具体 Origin（不可用通配 *）。
func buildCORSConfig(cfg CORS) cors.Config {
	origins := cfg.AllowOrigins
	if len(origins) == 0 {
		origins = devAllowOrigins
	}
	return cors.Config{
		AllowOrigins:     joinComma(origins),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,X-Agent-ID,X-Install-Token,X-CSRF-Token",
		ExposeHeaders:    "Content-Disposition,X-Agent-ID",
		AllowCredentials: true,
		MaxAge:           300,
	}
}

// joinComma 将 origin 列表合并为 fiber cors 接受的逗号分隔字符串。
func joinComma(items []string) string {
	out := ""
	for i, s := range items {
		if i > 0 {
			out += ","
		}
		out += s
	}
	return out
}

// durationSeconds 把秒数转成 time.Duration；非正值回落到 fallback。
func durationSeconds(seconds int, fallback time.Duration) time.Duration {
	if seconds <= 0 {
		return fallback
	}
	return time.Duration(seconds) * time.Second
}
