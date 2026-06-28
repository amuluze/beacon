// Package service
// Date: 2025/02/12 15:12:05
// Author: Amu
// Description:
package service

import (
	"github.com/amuluze/amutool/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

var InjectorSet = wire.NewSet(NewInjector)

type Injector struct {
	App    *fiber.App
	Router *Router
	Config *Config
	Logger *logger.Logger
}

func NewInjector(app *fiber.App, router *Router, config *Config, logx *logger.Logger) (*Injector, error) {
	return &Injector{
		App:    app,
		Router: router,
		Config: config,
		Logger: logx,
	}, nil
}
