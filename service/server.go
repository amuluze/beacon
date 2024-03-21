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
}

type Option func(*options)

func SetConfigFile(s string) Option {
	return func(o *options) {
		o.ConfigFile = s
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
	injector, cleanFunc, err := BuildInjector(o.ConfigFile)
	if err != nil {
		slog.Error("build injector failed", "err", err)
		return nil, err
	}

	slog.SetDefault(injector.Logger.Logger)
	httpServerCleanFunc := InitHttpServer(ctx, injector.Config, injector.App)

	timedTask := injector.Task
	go timedTask.Run()

	return func() {
		timedTask.Stop()
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
