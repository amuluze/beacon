// Package service
// Date: 2025/02/12 15:01:26
// Author: Amu
// Description:
package service

import (
	"strings"
	"time"

	"server/pkg/errors"
	"server/pkg/fiberx"
	"server/service/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

// 写接口的访问频率上限，避免统计/上报被脚本刷量
const (
	writeLimiterMax    = 30
	writeLimiterExpiry = 1 * time.Minute
)

func NewFiberApp(config *Config, r IRouter) *fiber.App {
	fiberConfig := fiber.Config{
		Prefork:      config.Fiber.Prefork,
		AppName:      config.Fiber.AppName,
		ServerHeader: config.Fiber.SeverHeader,
		// 官网后端仅承载统计与脚本下发，收窄请求体上限以降低 DoS 面
		BodyLimit: 4 * 1024 * 1024,
	}

	app := fiber.New(fiberConfig)

	// 跨域：按白名单放行，未配置来源时仅允许同源
	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(config.App.CORSAllowOrigins, ","),
		AllowMethods: strings.Join([]string{"GET", "POST", "OPTIONS"}, ","),
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(compress.New())
	// 基础安全响应头：nosniff / 防点击劫持 / Referrer 策略 / HSTS
	app.Use(middleware.SecurityHeadersMiddleware())

	// pprof 仅在非生产环境暴露，避免公网读取运行时/堆信息
	if !config.App.IsProduction() {
		app.Use(pprof.New())
	}

	// PanicMiddleware 必须是包裹业务路由的最外层 recover，
	// 不能再叠加 fiber 官方 recover（其内层 defer 会先吞掉 panic，使本中间件失效）
	app.Use(middleware.PanicMiddleware())

	// 对写接口施加 IP 级速率限制
	writeLimiter := limiter.New(limiter.Config{
		Max:        writeLimiterMax,
		Expiration: writeLimiterExpiry,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			// 统一走 fiberx 信封，前端可用一致的 {err,msg} 结构处理
			return fiberx.Failure(c, errors.TooManyRequestsError)
		},
	})
	app.Use("/api/v1/statistics/update", writeLimiter)
	app.Use("/api/v1/install/report", writeLimiter)

	err := r.Register(app)
	if err != nil {
		panic(err)
	}

	return app
}
