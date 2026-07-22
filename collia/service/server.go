// Package service
// Date: 2024/06/10 18:34:36
// Author: Amu
// Description:
package service

import (
	"log/slog"
)

func Run(configFile string, prefix Prefix, version string) (func(), error) {
	v := NewVersion(version)
	injector, clearFunc, err := BuildInjector(configFile, prefix, v)
	if err != nil {
		slog.Error("build injector failed:", "err", err)
		return nil, err
	}

	// 初始化日志
	// 注：vendored 版 logger（common/logger）底层为 zap，不再桥接到标准库 slog；
	// 业务日志统一走 injector.Logger，slog 包级调用保持默认输出。

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

	return func() {
		timedTask.Stop()
		err := rpcServer.Stop()
		if err != nil {
			slog.Error("rpc server stop failed:", "err", err)
		}
		clearFunc()
	}, nil
}
