// Package service
// Date: 2025/02/12 15:09:22
// Author: Amu
// Description:
package service

import (
	"time"

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
	// 健康检查：进程存活即 200，供 docker / 编排器探活，不依赖下游 DB 以避免抖动重启
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// 发布物静态托管：manager.sh / compose.yaml / version.json 等对外暴露在 /release/*
	a.registerReleaseStatic(app)

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
			gVersion := v1.Group("version")
			{
				gVersion.Get("/latest", a.VersionLatest)
			}
		}
	}
}

func (a *Router) Register(app *fiber.App) error {
	a.RegisterAPI(app)
	return nil
}

// registerReleaseStatic 将发布物目录以 /release/* 暴露为静态资源；
// 目录未配置时跳过注册，避免启动期因路径缺失 panic。
func (a *Router) registerReleaseStatic(app *fiber.App) {
	dir := a.config.Release.Dir
	if dir == "" {
		return
	}
	app.Static("/release", dir, fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        false, // 禁止目录列表，避免泄露发布物清单
		MaxAge:        300,
		CacheDuration: 10 * time.Second,
	})
}
func (a *Router) Prefixes() []string {
	return []string{"/api/"}
}
