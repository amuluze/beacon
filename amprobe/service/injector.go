// Package service
package service

import (
	"amprobe/service/report"

	"github.com/amuluze/amutool/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

var InjectorSet = wire.NewSet(NewInjector)

type Injector struct {
	App           *fiber.App
	Router        *Router
	Config        *Config
	Prepare       *Prepare
	Task          *TimedTask
	ReportService *report.Service
	Logger        *logger.Logger
}

func NewInjector(app *fiber.App, router *Router, prepare *Prepare, config *Config, task *TimedTask, reportService *report.Service, logx *logger.Logger) (*Injector, error) {
	return &Injector{
		App:           app,
		Router:        router,
		Config:        config,
		Prepare:       prepare,
		Task:          task,
		ReportService: reportService,
		Logger:        logx,
	}, nil
}
