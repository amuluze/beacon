// Package logger
// Description: 标准库 slog 桥接：把 slog 默认输出接到 stdout + 按天切割的文件 sink。
// 旧 amutool logger 的文件日志由 slog JSONHandler(AddSource) 写出；vendored 实现改为
// zap 后业务代码仍以包级 slog 调用为主，需要通过本构造恢复同等的文件日志语义。
package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

// NewJSONSlogHandler 构造输出到 stdout + 按天切割文件 sink 的 slog JSON Handler
// （含 AddSource 调用位置）。文件 sink 仅在 LogOutput=="file" 时挂载，与
// NewJsonFileLogger 的语义保持一致；rotator 初始化失败时返回仅 stdout 的
// handler 并附带错误信息供调用方告警。
func NewJSONSlogHandler(options ...Option) (slog.Handler, error) {
	config := defaultConfig()
	for _, option := range options {
		option(config)
	}

	writers := []io.Writer{os.Stdout}
	if config.LogOutput == "file" {
		if fw := getFileWriter(config); fw != nil {
			writers = append(writers, fw)
		} else {
			return slog.NewJSONHandler(os.Stdout, slogHandlerOptions(config)),
				fmt.Errorf("file sink unavailable, fallback to stdout only")
		}
	}

	var out io.Writer = os.Stdout
	if len(writers) > 1 {
		out = io.MultiWriter(writers...)
	}
	return slog.NewJSONHandler(out, slogHandlerOptions(config)), nil
}

func slogHandlerOptions(config *Config) *slog.HandlerOptions {
	level := slog.LevelInfo
	if parsed, err := parseSlogLevel(config.LogLevel); err == nil {
		level = parsed
	}
	return &slog.HandlerOptions{Level: level, AddSource: true}
}

func parseSlogLevel(value string) (slog.Level, error) {
	switch value {
	case "debug", "DEBUG":
		return slog.LevelDebug, nil
	case "info", "INFO", "":
		return slog.LevelInfo, nil
	case "warn", "WARN":
		return slog.LevelWarn, nil
	case "error", "ERROR":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unknown slog level %q", value)
	}
}
