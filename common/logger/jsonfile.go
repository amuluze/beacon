// Package logger
// Description: JSON 文件日志构造器，兼容旧版 amutool/logger 的 NewJsonFileLogger API。
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewJsonFileLogger 以函数式选项构造一个 JSON 格式的 Logger。
// 默认输出到 stdout；当通过 SetLogOutput("file") 指定时追加按天切割的文件 sink。
// 不依赖包级 defaultLogger，可安全多次调用。
func NewJsonFileLogger(options ...Option) *Logger {
	config := defaultConfig()
	for _, option := range options {
		option(config)
	}
	// 强制 JSON 编码，对齐 NewJsonFileLogger 的语义。
	config.LogFormat = "json"

	encoder := getEncoder(config)
	level, err := zapcore.ParseLevel(config.LogLevel)
	if err != nil {
		level = zapcore.InfoLevel
	}

	sinks := []zapcore.WriteSyncer{zapcore.AddSync(os.Stdout)}
	if config.LogOutput == "file" {
		if fw := getFileWriter(config); fw != nil {
			sinks = append(sinks, fw)
		}
	}
	writer := zapcore.NewMultiWriteSyncer(sinks...)

	return &Logger{
		Logger: zap.New(
			zapcore.NewCore(encoder, writer, level),
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		),
		name:    config.Name,
		loggers: make(map[string]*Logger),
	}
}
