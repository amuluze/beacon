// Package service
// Date: 2024/06/10 18:34:36
// Author: Amu
// Description:
package service

import (
	"fmt"
	"log/slog"

	"github.com/amuluze/amprobe/amvector/pkg/profile"
	"github.com/amuluze/amprobe/amvector/service/container"
)

func Run(configFile string, prefix Prefix) (func(), error) {
	slog.Info("config file", "info", configFile)
	injector, clearFunc, err := BuildInjector(configFile, prefix)
	if err != nil {
		slog.Error("build injector failed:", "err", err)
		return nil, err
	}

	// 初始化日志
	slog.SetDefault(injector.Logger.Logger)

	// 定时任务
	timedTask := injector.Task
	go timedTask.Run()

	// rpc server
	rpcServer := injector.RPCServer
	go func() {
		err := rpcServer.Start()
		if err != nil {
			slog.Error("rpc server start failed:", "err", err)
		}
	}()

	if err := setupService(injector.Config.ServiceProfile); err != nil {
		slog.Error("setup service failed:", "err", err)
	}

	return func() {
		timedTask.Stop()
		err := rpcServer.Stop()
		if err != nil {
			slog.Error("rpc server stop failed:", "err", err)
		}
		clearFunc()
	}, nil
}

func setupService(serviceProfile string) error {
	containerManager := container.NewContainerManager()
	cfg, err := profile.ReadProfile(serviceProfile)
	if err != nil {
		return err
	}
	fmt.Printf("service profile: %#v\n", cfg)
	if err := containerManager.CreateNetwork(cfg.Services.Network); err != nil {
		return err
	}

	if err := containerManager.CreateAmprobe(cfg.Services.Amprobe); err != nil {
		return err
	}
	return nil
}
