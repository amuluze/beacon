// Package service
// Date: 2024/3/6 11:07
// Author: Amu
// Description:
package service

import (
	"github.com/amuluze/amprobe/service/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

func NewFiberApp(config *Config, r IRouter) *fiber.App {
	fiberConfig := fiber.Config{
		Prefork:      config.Fiber.Prefork,
		AppName:      config.Fiber.AppName,
		ServerHeader: config.Fiber.SeverHeader,
	}

	app := fiber.New(fiberConfig)

	// 添加中间件
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(pprof.New())
	app.Use(middleware.PanicMiddleware())
	app.Use(middleware.StackMiddleware)

	err := r.Register(app)
	if err != nil {
		panic(err)
	}

	return app
}
