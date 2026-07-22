// Package logger
// Date: 2023/4/10 17:39
// Author: Amu
// Description:
package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var defaultLogger *Logger

func NewLogger(config *Config) *Logger {
	fmt.Printf("logger config: %#v\n", config)
	defaultLogger.NewLogger(
		SetName(config.Name),
		SetLogFile(config.LogFile),
		SetLogLevel(config.LogLevel),
		SetLogFormat(config.LogFormat),
		SetLogOutput(config.LogOutput),
		SetLogFileRotationTime(config.LogFileRotationTime),
		SetLogFileMaxAge(config.LogFileMaxAge),
		SetLogFileSuffix(config.LogFileSuffix),
	)
	defaultLogger.loggers[config.Name].Info("logger initialized")
	return defaultLogger.loggers[config.Name]
}

func WithField(fields ...zap.Field) {
	defaultLogger.WithField(fields...)
}

func Debug(args ...interface{}) {
	defaultLogger.Logger.WithOptions(zap.AddCallerSkip(1)).Debug(fmt.Sprint(args...))
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Logger.WithOptions(zap.AddCallerSkip(1)).Debug(fmt.Sprintf(format, args...))
}

func Info(args ...interface{}) {
	defaultLogger.Logger.WithOptions(zap.AddCallerSkip(1)).Info(fmt.Sprint(args...))
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Logger.WithOptions(zap.AddCallerSkip(1)).Info(fmt.Sprintf(format, args...))
}

func Warn(args ...interface{}) {
	defaultLogger.Logger.WithOptions(zap.AddCallerSkip(1)).Warn(fmt.Sprint(args...))
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Logger.WithOptions(zap.AddCallerSkip(1)).Warn(fmt.Sprintf(format, args...))
}

func Error(args ...interface{}) {
	defaultLogger.Logger.WithOptions(zap.AddCallerSkip(1)).Error(fmt.Sprint(args...))
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Logger.WithOptions(zap.AddCallerSkip(1)).Error(fmt.Sprintf(format, args...))
}

func Fatal(args ...interface{}) {
	defaultLogger.Logger.WithOptions(zap.AddCallerSkip(1)).Fatal(fmt.Sprint(args...))
}

func Fatalf(format string, args ...interface{}) {
	defaultLogger.Logger.WithOptions(zap.AddCallerSkip(1)).Fatal(fmt.Sprintf(format, args...))
}

func Panic(args ...interface{}) {
	defaultLogger.Logger.WithOptions(zap.AddCallerSkip(1)).Panic(fmt.Sprint(args...))
}

func Panicf(format string, args ...interface{}) {
	defaultLogger.Logger.WithOptions(zap.AddCallerSkip(1)).Panic(fmt.Sprintf(format, args...))
}

// Default 返回默认的 Logger 实例
func Default() *Logger { return defaultLogger }

// SetDefault 设置默认的 Logger 实例
func SetDefault(l *Logger) { defaultLogger = l }

// InfoFields 打印带字段的信息日志（包级，使用默认 Logger）
func InfoFields(msg string, fields ...zap.Field) {
	defaultLogger.Logger.WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

// DebugFields 打印带字段的调试日志（包级，使用默认 Logger）
func DebugFields(msg string, fields ...zap.Field) {
	defaultLogger.Logger.WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

// WarnFields 打印带字段的警告日志（包级，使用默认 Logger）
func WarnFields(msg string, fields ...zap.Field) {
	defaultLogger.Logger.WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
}

// ErrorFields 打印带字段的错误日志（包级，使用默认 Logger）
func ErrorFields(msg string, fields ...zap.Field) {
	defaultLogger.Logger.WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}
