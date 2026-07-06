package service

import (
	"os"
	"path/filepath"
	"testing"
)

// writeTempYAML writes a YAML config string to a temporary file and returns its path.
func writeTempYAML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return path
}

func TestNewConfigParsesMinimal(t *testing.T) {
	const yaml = `
control:
  server: 127.0.0.1:17000
  agent_id: agent-1
  join_token: secret
task:
  interval: 30
  disk:
    devices: []
  ethernet:
    names: []
  report:
    url: "http://localhost:8000/report"
    token: "rpt-token"
    agent_id: agent-1
db:
  dbtype: sqlite
  dbname: /tmp/collia-test
`
	cfg, err := NewConfig(writeTempYAML(t, yaml), Prefix("/test"))
	if err != nil {
		t.Fatalf("NewConfig: %v", err)
	}
	if cfg.Control.Server != "127.0.0.1:17000" {
		t.Fatalf("control.server = %q, want 127.0.0.1:17000", cfg.Control.Server)
	}
	if cfg.Control.AgentID != "agent-1" {
		t.Fatalf("control.agent_id = %q, want agent-1", cfg.Control.AgentID)
	}
	if cfg.Task.Interval != 30 {
		t.Fatalf("task.interval = %d, want 30", cfg.Task.Interval)
	}
	if cfg.Task.Report.URL != "http://localhost:8000/report" {
		t.Fatalf("task.report.url = %q", cfg.Task.Report.URL)
	}
	if cfg.prefix != Prefix("/test") {
		t.Fatalf("prefix = %q, want /test", cfg.prefix)
	}
}

func TestNewConfigMissingFile(t *testing.T) {
	_, err := NewConfig("/nonexistent/config-12345.yml", Prefix(""))
	if err == nil {
		t.Fatalf("expected error for missing config file, got nil")
	}
}

func TestNewConfigEmptyControlDefaults(t *testing.T) {
	const yaml = `
task:
  interval: 60
  disk:
    devices: []
  ethernet:
    names: []
  report:
    url: ""
    token: ""
    agent_id: ""
db:
  dbtype: sqlite
  dbname: /tmp/collia-empty
`
	cfg, err := NewConfig(writeTempYAML(t, yaml), Prefix(""))
	if err != nil {
		t.Fatalf("NewConfig: %v", err)
	}
	// Control should be zero-value since not specified
	if cfg.Control.AgentID != "" {
		t.Fatalf("expected empty agent_id, got %q", cfg.Control.AgentID)
	}
}
