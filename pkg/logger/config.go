// Package logger
// Date: 2024/03/20 12:11:50
// Author: Amu
// Description:
package logger

import (
	"log/slog"
)

type Config struct {
	Name                string     `load:"load"`     // 【默认】Logger 名称
	LogFile             string     `load:"load.log"` // 【默认】日志文件名称
	LogLevel            slog.Level `load:"info"`     // 【默认】日志打印级别
	LogFileRotationTime int        `load:"1"`        // 【默认】日志文件切割间隔，单位 D
	LogFileMaxAge       int        `load:"7"`        // 【默认】日志文件保留时间，单位 D
	LogFileSuffix       string     `load:".%Y%m%d"`  // 【默认】归档日志后缀
}

type Option func(*Config)

func SetName(name string) Option {
	return func(config *Config) {
		config.Name = name
	}
}

func SetLogFile(logFile string) Option {
	return func(config *Config) {
		config.LogFile = logFile
	}
}

func SetLogLevel(level string) Option {
	return func(config *Config) {
		switch level {
		case "debug":
			config.LogLevel = slog.LevelDebug
		case "info":
			config.LogLevel = slog.LevelInfo
		case "warn":
			config.LogLevel = slog.LevelWarn
		case "error":
			config.LogLevel = slog.LevelError
		}
	}
}

func SetLogFileRotationTime(duration int) Option {
	return func(config *Config) {
		config.LogFileRotationTime = duration
	}
}

func SetLogFileMaxAge(duration int) Option {
	return func(config *Config) {
		config.LogFileMaxAge = duration
	}
}

func SetLogFileSuffix(suffix string) Option {
	return func(config *Config) {
		config.LogFileSuffix = suffix
	}
}
