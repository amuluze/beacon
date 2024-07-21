// Package service
// Date: 2024/06/10 18:56:26
// Author: Amu
// Description:
package service

import (
	"fmt"
	"github.com/amuluze/amutool/logger"
)

func NewLogger(config *Config) *logger.Logger {
	logx := logger.NewJsonFileLogger(
		logger.SetLogFile(config.Log.Output),
		logger.SetLogLevel(config.Log.Level),
		logger.SetLogFileRotationTime(config.Log.Rotation),
		logger.SetLogFileMaxAge(config.Log.MaxAge),
	)
	fmt.Println("define logx")
	return logx
}
