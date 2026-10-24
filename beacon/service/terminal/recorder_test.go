package terminal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRecorder_WriteOutput(t *testing.T) {
	dir := t.TempDir()
	stalePath := filepath.Join(dir, "test.cast")
	if err := os.WriteFile(stalePath, []byte("stale recording"), 0644); err != nil {
		t.Fatalf("seed recording: %v", err)
	}
	if err := os.Chmod(stalePath, 0644); err != nil {
		t.Fatalf("set insecure recording permissions: %v", err)
	}

	r, path, err := NewRecorder(dir, "test", 80, 24)
	if err != nil {
		t.Fatalf("new recorder failed: %v", err)
	}

	if err := r.WriteOutput([]byte("hello ")); err != nil {
		t.Fatalf("write output failed: %v", err)
	}
	if err := r.WriteOutput([]byte("world")); err != nil {
		t.Fatalf("write output failed: %v", err)
	}
	if err := r.Resize(120, 30); err != nil {
		t.Fatalf("resize failed: %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat recording: %v", err)
	}
	if got := info.Mode().Perm(); got != 0600 {
		t.Fatalf("recording permissions = %o, want 600", got)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read recording failed: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}

	var header map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &header); err != nil {
		t.Fatalf("header unmarshal failed: %v", err)
	}
	if header["version"].(float64) != 2 {
		t.Fatalf("expected version 2, got %v", header["version"])
	}
	if header["width"].(float64) != 80 {
		t.Fatalf("expected width 80, got %v", header["width"])
	}

	var line1 []interface{}
	if err := json.Unmarshal([]byte(lines[1]), &line1); err != nil {
		t.Fatalf("line1 unmarshal failed: %v", err)
	}
	if line1[1] != "o" {
		t.Fatalf("expected event 'o', got %v", line1[1])
	}
	if line1[2] != "hello " {
		t.Fatalf("expected 'hello ', got %v", line1[2])
	}
}

func TestCleanSessionPath(t *testing.T) {
	dir := t.TempDir()
	path, err := CleanSessionPath(dir, "sess-001")
	if err != nil {
		t.Fatalf("clean path failed: %v", err)
	}
	if !strings.HasSuffix(path, "sess-001.cast") {
		t.Fatalf("unexpected path: %s", path)
	}

	if _, err := CleanSessionPath(dir, "../etc/passwd"); err == nil {
		t.Fatal("expected error for path traversal session_id")
	}
}

func TestRecorder_ConcurrentWrites(t *testing.T) {
	dir := t.TempDir()
	r, path, err := NewRecorder(dir, "concurrent", 80, 24)
	if err != nil {
		t.Fatalf("new recorder failed: %v", err)
	}

	done := make(chan struct{}, 2)
	writer := func() {
		defer func() { done <- struct{}{} }()
		for i := 0; i < 100; i++ {
			_ = r.WriteOutput([]byte("x"))
			time.Sleep(time.Microsecond * 100)
		}
	}
	go writer()
	go writer()
	<-done
	<-done

	if err := r.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read recording failed: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) < 200 { // header + at least 200 outputs
		t.Fatalf("expected at least 202 lines, got %d", len(lines))
	}
	outputCount := 0
	for i, line := range lines {
		if i == 0 {
			continue
		}
		var event []interface{}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("unmarshal line failed: %v", err)
		}
		if event[1] == "o" {
			outputCount++
		}
	}
	if outputCount < 200 {
		t.Fatalf("expected at least 200 output events, got %d", outputCount)
	}
}

func TestNewRecorderRejectsSymlinkEscape(t *testing.T) {
	dir := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.cast")
	if err := os.WriteFile(outside, []byte("do not overwrite"), 0600); err != nil {
		t.Fatalf("write outside file: %v", err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, "linked.cast")); err != nil {
		t.Fatalf("create symlink: %v", err)
	}

	if _, _, err := NewRecorder(dir, "linked", 80, 24); err == nil {
		t.Fatal("expected symlink escape to be rejected")
	}
	data, err := os.ReadFile(outside)
	if err != nil {
		t.Fatalf("read outside file: %v", err)
	}
	if string(data) != "do not overwrite" {
		t.Fatalf("outside file was modified: %q", data)
	}
}
