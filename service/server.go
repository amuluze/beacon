// Package service
// Date: 2024/3/6 11:00
// Author: Amu
// Description:
package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func InitHttpServer(config *Config, app *fiber.App) func() {
	appConfig := config.Fiber
	addr := fmt.Sprintf("%s:%d", appConfig.Host, appConfig.Port)
	slog.Info("start http server", "addr", addr)
	go func() {
		err := app.Listen(addr)
		if err != nil {
			slog.Warn("http server start error", "err", err)
			panic(err)
		}
	}()

	return func() {
		_, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(appConfig.ShutdownTimeout))
		defer cancel()
		if err := app.Shutdown(); err != nil {
			slog.Warn("app shut down error", "err", err)
		}
	}
}

func Run(configFile string, modelFile ModelConfig) (func(), error) {
	injector, cleanFunc, err := BuildInjector(configFile, modelFile)
	if err != nil {
		slog.Error("build injector failed", "err", err)
		return nil, err
	}

	// 初始化日志
	slog.SetDefault(injector.Logger.Logger)

	// 初始化预设数据
	injector.Prepare.Init(injector.App)

	// 定时任务
	timedTask := injector.Task
	go timedTask.Run()

	slog.Info("service starting")
	httpServerCleanFunc := InitHttpServer(injector.Config, injector.App)
	slog.Info("service start success")

	return func() {
		timedTask.Stop()
		httpServerCleanFunc()
		cleanFunc()
	}, nil
}
