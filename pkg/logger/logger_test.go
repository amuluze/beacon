// Package logger
// Date: 2024/03/19 17:57:41
// Author: Amu
// Description:
package logger

import (
	"fmt"
	"log/slog"
	"testing"
	"time"
)

func TestTextLogger(t *testing.T) {
	logx := NewTextLogger()
	slog.SetDefault(logx.Logger)
	logx.SetDebugLevel()
	slog.Debug("this is a debug message")
	slog.Info("this is a info message")
	slog.Error("this is a error message")
	fmt.Println("--------------------------")
	logx.SetErrorLevel()
	slog.Debug("this is a debug message")
	slog.Info("this is a info message")
	slog.Error("this is a error message")
}

func TestJsonLogger(t *testing.T) {
	logx := NewJsonLogger()
	slog.SetDefault(logx.Logger)
	logx.SetDebugLevel()
	slog.Debug("this is a debug message")
	slog.Info("this is a info message")
	slog.Error("this is a error message")
	fmt.Println("--------------------------")
	logx.SetErrorLevel()
	slog.Debug("this is a debug message")
	slog.Info("this is a info message")
	slog.Error("this is a error message")
}

func TestJsonFileLogger(t *testing.T) {
	logx := NewJsonFileLogger(
		SetName("load"),
		SetLogFile("load.log"),
		SetLogLevel("info"),
		SetLogFileRotationTime(1),
		SetLogFileMaxAge(7),
		SetLogFileSuffix(".%Y%m%d"),
	)
	slog.SetDefault(logx.Logger)
	logx.SetDebugLevel()
	slog.Debug("this is a debug message")
	slog.Info("this is a info message")
	slog.Error("this is a error message")

	logx.SetErrorLevel()
	slog.Debug("this is a debug message")
	slog.Info("this is a info message")
	slog.Error("this is a error message")
	time.Sleep(5 * time.Second)
}
