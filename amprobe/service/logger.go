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
	)
	return logx
}
