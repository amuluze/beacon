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
		logger.SetLogFile(config.Logger.File),
		logger.SetLogLevel(config.Logger.Level),
		logger.SetLogFileRotationTime(config.Logger.RotationTime),
		logger.SetLogFileMaxAge(config.Logger.MaxAge),
	)
	return logx
}
