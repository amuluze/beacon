// Package service
// Date: 2024/06/10 18:56:26
// Author: Amu
// Description:
package service

import (
	"fmt"

	"github.com/amuluze/amutool/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func NewLogger(config *Config) *logger.Logger {
	logx := logger.NewJsonFileLogger(
		logger.SetLogFile(config.Logger.File),
		logger.SetLogLevel(config.Logger.Level),
		logger.SetLogFileRotationTime(config.Logger.RotationTime),
		logger.SetLogFileMaxAge(config.Logger.MaxAge),
	)

	switch config.Logger.Level {
	case "debug":
		logx.SetDebugLevel()
	case "info":
		logx.SetInfoLevel()
	case "warn":
		logx.SetWarnLevel()
	case "error":
		logx.SetErrorLevel()
	default:
		logx.Info("use default logger level")
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file changed:", in.Name)
		if err := viper.Unmarshal(config); err != nil {
			fmt.Printf("unmarshal config error when change, %v", err)
		}
		switch config.Logger.Level {
		case "debug":
			logx.SetDebugLevel()
		case "info":
			logx.SetInfoLevel()
		case "warn":
			logx.SetWarnLevel()
		case "error":
			logx.SetErrorLevel()
		default:
			logx.Warn("use default logger level")
		}
	})
	return logx
}
