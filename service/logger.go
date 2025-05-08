// Package service
// Date: 2024/03/19 17:33:57
// Author: Amu
// Description:
package service

import (
	"os"
	"path/filepath"

	"amprobe/pkg/logger"
)

func NewLogger(config *Config) *logger.Logger {
	// 确保日志目录存在
	logDir := filepath.Dir(config.Log.Output)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			// 如果无法创建目录，打印错误并继续（使用默认日志配置）
			println("无法创建日志目录:", err.Error())
		}
	}

	// 创建日志记录器，添加更多配置选项
	logx := logger.NewJsonFileLogger(
		logger.SetLogFile(config.Log.Output),
		logger.SetLogLevel(config.Log.Level),
		logger.SetLogFileRotationTime(config.Log.Rotation),
		logger.SetLogFileMaxAge(config.Log.MaxAge),
	)
	return logx
}
