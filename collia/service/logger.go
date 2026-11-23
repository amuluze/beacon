// Package service
// Date: 2024/06/10 18:56:26
// Author: Amu
// Description:
package service

import (
	"log/slog"

	"github.com/amuluze/amutool/logger"
)

func NewLogger(config *Config) *logger.Logger {
	// log.output 为空时不挂文件 sink：vendored logger 会把空路径回退为
	// 可执行文件同级的 default.log（隐性副作用），保持 stdout 输出。
	if config.Log.Output == "" {
		return logger.NewJsonFileLogger(logger.SetLogLevel(config.Log.Level))
	}
	logx := logger.NewJsonFileLogger(
		logger.SetLogFile(config.Log.Output),
		logger.SetLogLevel(config.Log.Level),
		logger.SetLogFileRotationTime(config.Log.Rotation),
		logger.SetLogFileMaxAge(config.Log.MaxAge),
		// vendored logger 默认 LogOutput=stdout，不显式声明则文件 sink 不会挂载。
		logger.SetLogOutput("file"),
	)
	// 业务代码以包级 slog 调用为主（默认输出 stderr）：把 slog 默认实例
	// 接到同一按天切割文件 sink，使 log.output 配置对 slog 路径同样生效。
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
