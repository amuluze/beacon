// Package logger
// Date: 2024/03/19 17:09:20
// Author: Amu
// Description:
package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	rotator "github.com/lestrrat-go/file-rotatelogs"
)

var defaultLevel slog.LevelVar

type Logger struct {
	*slog.Logger
}

func NewTextLogger() *Logger {
	defaultLevel.Set(slog.LevelInfo)
	log := &Logger{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     &defaultLevel,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					if t, ok := a.Value.Any().(time.Time); ok {
						a.Value = slog.StringValue(t.Format(time.DateTime))
					}
				}
				return a
			},
		})),
	}
	return log
}

func NewJsonLogger() *Logger {
	defaultLevel.Set(slog.LevelInfo)
	log := &Logger{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     &defaultLevel,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					if t, ok := a.Value.Any().(time.Time); ok {
						a.Value = slog.StringValue(t.Format(time.DateTime))
					}
				}
				return a
			},
		})),
	}
	return log
}

func NewJsonFileLogger(options ...Option) *Logger {
	config := &Config{
		Name:                "load",
		LogFile:             "load.log",
		LogLevel:            slog.LevelInfo,
		LogFileRotationTime: 1,
		LogFileMaxAge:       7,
		LogFileSuffix:       ".%Y%m%d",
	}
	for _, option := range options {
		option(config)
	}

	defaultLevel.Set(config.LogLevel)

	logFilePath := config.LogFile
	if !filepath.IsAbs(config.LogFile) {
		absPath, _ := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), config.LogFile))
		logFilePath = absPath
	}
	fmt.Printf("log file path: %s\n", logFilePath)
	_log, _ := rotator.New(
		filepath.Join(logFilePath+config.LogFileSuffix),
		rotator.WithLinkName(logFilePath),
		rotator.WithMaxAge(time.Duration(config.LogFileMaxAge)*24*time.Hour),
		rotator.WithRotationTime(time.Duration(config.LogFileRotationTime)*time.Hour*24),
	)
	defaultLogger := slog.New(slog.NewJSONHandler(_log, &slog.HandlerOptions{
		AddSource: true,
		Level:     &defaultLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.DateTime))
				}
			}
			return a
		},
	}))
	return &Logger{Logger: defaultLogger}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *Logger) SetDebugLevel() {
	defaultLevel.Set(slog.LevelDebug)
}

func (l *Logger) SetInfoLevel() {
	defaultLevel.Set(slog.LevelInfo)
}

func (l *Logger) SetWarnLevel() {
	defaultLevel.Set(slog.LevelWarn)
}

func (l *Logger) SetErrorLevel() {
	defaultLevel.Set(slog.LevelError)
}
