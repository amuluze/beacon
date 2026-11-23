package service

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewLoggerAttachesFileSink 锁定修复：NewLogger 后包级 slog 输出必须
// 落到 config.Log.Output 声明的文件（按天切割后缀命名）。修复前 collia 不声明
// SetLogOutput("file") 也不接管 slog 默认实例，该文件永远不产生、业务日志
// 全部走 stderr/stdout。
func TestNewLoggerAttachesFileSink(t *testing.T) {
	previous := slog.Default()
	defer slog.SetDefault(previous)

	logPath := filepath.Join(t.TempDir(), "collia.log")
	config := &Config{Log: Log{Output: logPath, Level: "info", Rotation: 1, MaxAge: 7}}

	NewLogger(config)

	slog.Info("logger file sink probe", "probe", "file-sink")
	matches, err := filepath.Glob(logPath + ".2*")
	if err != nil {
		t.Fatalf("glob rotated log files: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("expected a rotated log file under the configured log output path")
	}
	data, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	if !strings.Contains(string(data), "logger file sink probe") {
		t.Fatalf("expected package-level slog output to reach the configured file, got:\n%s", data)
	}
}

// TestNewLoggerKeepsStdoutWhenOutputEmpty log.output 为空时不得挂文件 sink：
// vendored logger 会把空路径回退为可执行文件同级的 default.log，必须保持
// stdout 输出且不替换 slog 默认实例。
func TestNewLoggerKeepsStdoutWhenOutputEmpty(t *testing.T) {
	previous := slog.Default()
	defer slog.SetDefault(previous)

	NewLogger(&Config{Log: Log{Output: "", Level: "info"}})

	if slog.Default() != previous {
		t.Fatal("empty log.output must not replace the default slog handler")
	}
}
