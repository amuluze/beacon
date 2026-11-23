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

	// 初始化日志：NewLogger 已把包级 slog 默认实例接到 log.output 声明的
	// 按天切割文件 sink，业务日志随配置落盘；injector.Logger 保留 zap 句柄
	// 供需要结构化 logger 的组件使用。

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
