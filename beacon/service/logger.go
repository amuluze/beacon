// Package service
// Date: 2024/03/19 17:33:57
// Author: Amu
// Description:
package service

import (
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
	return logx
}
