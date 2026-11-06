// Package service
// Date: 2024/03/19 17:33:57
// Author: Amu
// Description:
package service

import (
	"log/slog"

	"github.com/amuluze/amutool/logger"
)

func NewLogger(config *Config) *logger.Logger {
	logx := logger.NewJsonFileLogger(
		logger.SetLogFile(config.Log.Output),
		logger.SetLogLevel(config.Log.Level),
		logger.SetLogFileRotationTime(config.Log.Rotation),
		logger.SetLogFileMaxAge(config.Log.MaxAge),
		// vendored logger 默认 LogOutput=stdout，不显式声明则文件 sink
		// 不会挂载（workspace 构建切换 vendored 实现后曾因此丢失文件日志）。
		logger.SetLogOutput("file"),
	)
	// 业务代码以包级 slog 调用为主（默认输出 stderr）：把 slog 默认实例
	// 接到 stdout + 同一按天切割文件 sink，恢复旧版文件日志语义。
	handler, err := logger.NewJSONSlogHandler(
		logger.SetLogFile(config.Log.Output),
		logger.SetLogLevel(config.Log.Level),
		logger.SetLogFileRotationTime(config.Log.Rotation),
		logger.SetLogFileMaxAge(config.Log.MaxAge),
		logger.SetLogOutput("file"),
	)
	if err != nil {
		slog.Error("build slog file handler failed, stderr only", "err", err)
	}
	slog.SetDefault(slog.New(handler))
	return logx
}
