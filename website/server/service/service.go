// Package service
// Date: 2025/02/12 15:10:33
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

// InitHttpServer 启动 HTTP 监听并返回 listenErr 与 shutdown。
// listen 失败不再 panic（避免 goroutine panic 直接 crash 进程、跳过清理），
// 而是通过 listenErr 通道通知主流程优雅退出。
// shutdown 使用 ShutdownWithContext 让 ShutdownTimeout 真正生效。
func InitHttpServer(ctx context.Context, config *Config, app *fiber.App) (<-chan error, func()) {
	appConfig := config.Fiber
	addr := fmt.Sprintf("%s:%d", appConfig.Host, appConfig.Port)
	slog.Info("start http server", "addr", addr)

	listenErr := make(chan error, 1)
	go func() {
		if err := app.Listen(addr); err != nil {
			listenErr <- err
		}
	}()

	shutdown := func() {
		shutdownCtx, cancel := context.WithTimeout(ctx, time.Duration(appConfig.ShutdownTimeout)*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			slog.Warn("app shut down error", "err", err)
		}
	}

	return listenErr, shutdown
}

func Init(ctx context.Context, opts ...Option) (func(), <-chan error, error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	injector, cleanFunc, err := BuildInjector(o.ConfigFile)
	if err != nil {
		slog.Error("build injector failed", "err", err)
		return nil, nil, err
	}

	// 初始化日志
	slog.SetDefault(injector.Logger.Logger)

	listenErr, httpServerShutdown := InitHttpServer(ctx, injector.Config, injector.App)

	return func() {
		httpServerShutdown()
		cleanFunc()
	}, listenErr, nil
}

// Run 阻塞等待退出信号或监听错误，返回后由 defer 保证清理（关闭 HTTP、DB）。
// 不再调用 os.Exit，避免跳过 defer 与清理链。
func Run(ctx context.Context, opts ...Option) error {
	cleanFunc, listenErr, err := Init(ctx, opts...)
	if err != nil {
		return err
	}
	defer cleanFunc()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case sig := <-sc:
			switch sig {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				return nil
			case syscall.SIGHUP:
				// 忽略终端挂起，继续运行
				continue
			}
		case err := <-listenErr:
			// listen 失败（如端口占用）触发优雅退出而非裸 panic
			slog.Error("http server listen failed", "err", err)
			return fmt.Errorf("http server listen failed: %w", err)
		}
	}
}
