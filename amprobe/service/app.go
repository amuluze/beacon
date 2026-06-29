// Package service
// Date: 2024/3/6 11:07
// Author: Amu
// Description:
package service

import (
	"amprobe/service/middleware"
	"amprobe/web"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

func NewFiberApp(config *Config, r IRouter) *fiber.App {
	fiberConfig := fiber.Config{
		Prefork:      config.Fiber.Prefork,
		AppName:      config.Fiber.AppName,
		ServerHeader: config.Fiber.SeverHeader,
		BodyLimit:    1000 * 1024 * 1024,
	}

	app := fiber.New(fiberConfig)

	// 添加中间件
	app.Use(cors.New(buildCORSConfig(config.CORS)))
	// 全局兜底限流：仅作用于 /api/ 路由，按 IP 计数；登录类由独立严格限流器接管。
	if config.RateLimit.Enable {
		app.Use(middleware.GlobalRateLimitMiddleware(config.RateLimit.GlobalMax))
	}
	app.Use(compress.New())
	// pprof 暴露运行时内存/goroutine/profile，可泄漏密钥或被用于 DoS。
	// 仅在非生产模式启用；生产模式禁用，避免 /debug/pprof/* 无鉴权暴露。
	if !strings.EqualFold(config.App.Env, productionEnv) {
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
// 凭证（Authorization）跨域依赖 AllowCredentials 与具体 Origin（不可用通配 *）。
func buildCORSConfig(cfg CORS) cors.Config {
	origins := cfg.AllowOrigins
	if len(origins) == 0 {
		origins = devAllowOrigins
	}
	return cors.Config{
		AllowOrigins:     joinComma(origins),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Agent-ID,X-Install-Token",
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
