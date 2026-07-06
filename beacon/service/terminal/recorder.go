// Package terminal provides WebSocket terminal bridge and asciinema session recording.
package terminal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Recorder writes terminal output in asciinema v2 format.
type Recorder struct {
	mu        sync.Mutex
	file      *os.File
	startedAt time.Time
	width     int
	height    int
	closed    bool
}

// NewRecorder creates a recorder that writes to the given file path.
// The directory is created if it does not exist.
func NewRecorder(path string, width, height int) (*Recorder, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create session directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
	if err != nil {
		return nil, fmt.Errorf("failed to create recording file: %w", err)
	}

	r := &Recorder{
		file:      file,
		startedAt: time.Now(),
		width:     width,
		height:    height,
	}

	if err := r.writeHeader(); err != nil {
		_ = file.Close()
		return nil, err
	}
	return r, nil
}

func (r *Recorder) writeHeader() error {
	header := map[string]interface{}{
		"version":   2,
		"width":     r.width,
		"height":    r.height,
		"timestamp": r.startedAt.Unix(),
		"env": map[string]string{
			"SHELL": "/bin/bash",
			"TERM":  "xterm-256color",
		},
	}
	data, err := json.Marshal(header)
	if err != nil {
		return fmt.Errorf("failed to marshal header: %w", err)
	}
	if _, err := r.file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	return nil
}

// WriteOutput writes a terminal output chunk to the recording.
func (r *Recorder) WriteOutput(data []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return nil
	}
	if len(data) == 0 {
		return nil
	}

	elapsed := time.Since(r.startedAt).Seconds()
	line := []interface{}{elapsed, "o", string(data)}
	encoded, err := json.Marshal(line)
	if err != nil {
		return fmt.Errorf("failed to marshal output line: %w", err)
	}
	if _, err := r.file.Write(append(encoded, '\n')); err != nil {
		return fmt.Errorf("failed to write output line: %w", err)
	}
	return nil
}

// Resize updates the recorded terminal size. It does not rewrite past frames.
func (r *Recorder) Resize(width, height int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return nil
	}
	r.width = width
	r.height = height
	return nil
}

// Close flushes and closes the recording file.
func (r *Recorder) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return nil
	}
	r.closed = true
	if err := r.file.Sync(); err != nil {
		_ = r.file.Close()
		return fmt.Errorf("failed to sync recording: %w", err)
	}
	if err := r.file.Close(); err != nil {
		return fmt.Errorf("failed to close recording: %w", err)
	}
	return nil
}

// CleanSessionPath joins the session directory with the session ID and ensures
// the resulting path stays within the base directory to prevent path traversal.
func CleanSessionPath(baseDir, sessionID string) (string, error) {
	if strings.ContainsAny(sessionID, "/\\") {
		return "", fmt.Errorf("invalid session_id")
	}
	baseDir = filepath.Clean(baseDir)
	path := filepath.Clean(filepath.Join(baseDir, sessionID+".cast"))
	if !strings.HasPrefix(path, baseDir+string(filepath.Separator)) && path != baseDir {
		return "", fmt.Errorf("session path escapes base directory")
	}
	return path, nil
}

// ReadRecording reads the full contents of a recording file. It is intended for tests.
func ReadRecording(path string) ([]byte, error) {
	return os.ReadFile(path)
}

var _ io.Closer = (*Recorder)(nil)
