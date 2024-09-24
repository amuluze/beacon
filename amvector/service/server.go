// Package service
// Date: 2024/06/10 18:34:36
// Author: Amu
// Description:
package service

import (
	"log/slog"
)

func Run(configFile string, prefix Prefix) (func(), error) {
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

	return func() {
		timedTask.Stop()
		err := rpcServer.Stop()
		if err != nil {
			slog.Error("rpc server stop failed:", "err", err)
		}
		clearFunc()
	}, nil
}
