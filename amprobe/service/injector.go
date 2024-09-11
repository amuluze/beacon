// Package service
// Date: 2024/3/6 11:08
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
	App     *fiber.App
	Router  *Router
	Config  *Config
	Prepare *Prepare
	Logger  *logger.Logger
}

func NewInjector(app *fiber.App, router *Router, prepare *Prepare, config *Config, logx *logger.Logger) (*Injector, error) {
	return &Injector{
		App:     app,
		Router:  router,
		Config:  config,
		Prepare: prepare,
		Logger:  logx,
	}, nil
}
