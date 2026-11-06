package logger

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewJSONSlogHandlerWritesRotatedFile 验证 slog 桥接会写入按天切割的
// 文件 sink 并同步输出 stdout，恢复旧版文件日志语义。
func TestNewJSONSlogHandlerWritesRotatedFile(t *testing.T) {
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")

	handler, err := NewJSONSlogHandler(
		SetLogFile(logFile),
		SetLogLevel("info"),
		SetLogOutput("file"),
	)
	if err != nil {
		t.Fatalf("NewJSONSlogHandler: %v", err)
	}

	logger := slog.New(handler)
	logger.Info("hello file log", "key", "value")

	matches, err := filepath.Glob(logFile + ".2*")
	if err != nil {
		t.Fatalf("glob rotated files: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("rotated file count = %d, want 1", len(matches))
	}
	content, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	if !bytes.Contains(content, []byte(`"msg":"hello file log"`)) {
		t.Fatalf("log file missing entry: %s", content)
	}
	if !bytes.Contains(content, []byte(`"source"`)) {
		t.Fatalf("log file missing source info: %s", strings.TrimSpace(string(content)))
	}
}
