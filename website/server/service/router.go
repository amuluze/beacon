// Package service
// Date: 2025/02/12 15:09:22
// Author: Amu
// Description:
package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"

	statisticsAPI "server/service/statistics/api"
)

var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

var _ IRouter = (*Router)(nil)

type IRouter interface {
	Register(app *fiber.App) error
	Prefixes() []string
}

type Router struct {
	config        *Config
	statisticsAPI *statisticsAPI.StatisticsAPI
}

func (a *Router) RegisterAPI(app *fiber.App) {
	app.Get("/download/install.sh", a.WebsiteInstallScript)
	app.Get("/download/compose.yaml", a.WebsiteCompose)

	api := app.Group("api")
	{
		v1 := api.Group("v1")
		{
			gStatistics := v1.Group("statistics")
			{
				gStatistics.Get("/query", a.statisticsAPI.StatisticsQuery)
				gStatistics.Post("/update", a.statisticsAPI.StatisticsUpdate)
			}
			gInstall := v1.Group("install")
			{
				gInstall.Post("/report", a.statisticsAPI.InstallationReport)
			}
		}
	}
}

func (a *Router) Register(app *fiber.App) error {
	a.RegisterAPI(app)
	return nil
}
func (a *Router) Prefixes() []string {
	return []string{"/api/"}
}
