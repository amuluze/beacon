// Package logger
// Date: 2023/4/10 17:18
// Author: Amu
// Description:
package logger

type Config struct {
	Name                string `load:"std"`         // 【默认】Logger 名称
	LogFile             string `load:"scanner.log"` // 【默认】日志文件名称
	LogLevel            string `load:"info"`        // 【默认】日志打印级别
	LogFormat           string `load:"text"`        // 【默认】日志打印样式，支持 text 和 json
	LogFileRotationTime int    `load:"1"`           // 【默认】日志文件切割间隔，单位 D
	LogFileMaxAge       int    `load:"7"`           // 【默认】日志文件保留时间，单位 D
	LogOutput           string `load:"stdout"`      // 【默认】日志输出位置，只会 stdout iohelper
	LogFileSuffix       string `load:".%Y%m%d"`     // 【默认】归档日志后缀
}

func defaultConfig() *Config {
	return &Config{
		Name:                "default",
		LogFile:             "default.log",
		LogLevel:            "info",
		LogFormat:           "text",
		LogFileRotationTime: 1,
		LogFileMaxAge:       7,
		LogOutput:           "stdout",
		LogFileSuffix:       ".%Y%m%d",
	}
}

type Option func(*Config)

func SetName(name string) Option {
	return func(config *Config) {
		config.Name = name
	}
}

func SetLogFile(logFile string) Option {
	return func(config *Config) {
		if logFile == "" {
			return
		}
		config.LogFile = logFile
	}
}

func SetLogLevel(level string) Option {
	return func(config *Config) {
		if level == "" {
			return
		}
		config.LogLevel = level
	}
}

func SetLogFormat(format string) Option {
	return func(config *Config) {
		if format == "" {
			return
		}
		config.LogFormat = format
	}
}

func SetLogOutput(output string) Option {
	return func(config *Config) {
		if output == "" {
			return
		}
		config.LogOutput = output
	}
}

func SetLogFileRotationTime(duration int) Option {
	return func(config *Config) {
		if duration <= 0 {
			return
		}
		config.LogFileRotationTime = duration
	}
}

func SetLogFileMaxAge(duration int) Option {
	return func(config *Config) {
		if duration <= 0 {
			return
		}
		config.LogFileMaxAge = duration
	}
}

func SetLogFileSuffix(suffix string) Option {
	return func(config *Config) {
		if suffix == "" {
			return
		}
		config.LogFileSuffix = suffix
	}
}
