// Package service
// Date: 2024/3/6 11:00
// Author: Amu
// Description:
package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

type options struct {
	ConfigFile string
	ModelFile  ModeConf
}

type Option func(*options)

func SetConfigFile(s string) Option {
	return func(o *options) {
		o.ConfigFile = s
	}
}

func SetModelFile(s string) Option {
	return func(o *options) {
		o.ModelFile = ModeConf(s)
	}
}

func InitHttpServer(ctx context.Context, config *Config, app *fiber.App) func() {
	appConfig := config.Fiber
	addr := fmt.Sprintf("%s:%d", appConfig.Host, appConfig.Port)
	slog.Info("start http server", "addr", addr)
	go func() {
		err := app.Listen(addr)
		if err != nil {
			panic(err)
		}
	}()

	return func() {
		_, cancel := context.WithTimeout(ctx, time.Second*time.Duration(appConfig.ShutdownTimeout))
		defer cancel()
		if err := app.Shutdown(); err != nil {
			slog.Warn("app shut down error", "err", err)
		}
	}
}

func Init(ctx context.Context, opts ...Option) (func(), error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	injector, cleanFunc, err := BuildInjector(o.ConfigFile, o.ModelFile)
	if err != nil {
		slog.Error("build injector failed", "err", err)
		return nil, err
	}

	// 初始化日志
	// 注：vendored 版 logger（common/logger）底层为 zap，不再桥接到标准库 slog；
	// 业务日志统一走 injector.Logger，slog 包级调用保持默认输出。

	// 安装统计上报不参与启动成败判定。
	go ReportInstallation(ctx, injector.Config)

	// 初始化预设数据
	injector.Prepare.Init(injector.App)

	// 定时任务
	timedTask := injector.Task
	go timedTask.Run()

	// 版本检查（仅当配置启用时）。失败不影响启动。
	var versionChecker *VersionChecker
	if injector.Config.Update.Enable {
		versionChecker = NewVersionChecker(injector.Config)
		SetGlobalVersionChecker(versionChecker)
		go versionChecker.Run()
	}

	httpServerCleanFunc := InitHttpServer(ctx, injector.Config, injector.App)

	return func() {
		timedTask.Stop()
		if versionChecker != nil {
			versionChecker.Stop()
		}
		httpServerCleanFunc()
		cleanFunc()
	}, nil
}

func Run(ctx context.Context, opts ...Option) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := Init(ctx, opts...)
	if err != nil {
		return err
	}

EXIT:
	for {
		sig := <-sc
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	cleanFunc()
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}
