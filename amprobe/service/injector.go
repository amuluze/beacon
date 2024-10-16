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
	Task    *TimedTask
	Logger  *logger.Logger
}

func NewInjector(app *fiber.App, router *Router, prepare *Prepare, config *Config, task *TimedTask, logx *logger.Logger) (*Injector, error) {
	return &Injector{
		App:     app,
		Router:  router,
		Config:  config,
		Prepare: prepare,
		Task:    task,
		Logger:  logx,
	}, nil
}
