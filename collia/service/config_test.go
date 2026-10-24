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
  tls:
    enable: true
task:
  interval: 30
  disk:
    devices: []
  ethernet:
    names: []
  report:
    url: "https://localhost:8000/report"
    token: "rpt-token"
    agent_id: agent-1
db:
  dbtype: sqlite
  dbname: /tmp/collia-test
`
	configPath := writeTempYAML(t, yaml)
	if err := os.Chmod(configPath, 0o644); err != nil {
		t.Fatalf("chmod config: %v", err)
	}
	cfg, err := NewConfig(configPath, Prefix("/test"))
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
	if cfg.Task.Report.URL != "https://localhost:8000/report" {
		t.Fatalf("task.report.url = %q", cfg.Task.Report.URL)
	}
	if cfg.prefix != Prefix("/test") {
		t.Fatalf("prefix = %q, want /test", cfg.prefix)
	}
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("stat config: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("config permissions = %o, want 600", got)
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

func TestNormalizeAgentIDReportInheritsControl(t *testing.T) {
	const yaml = `
control:
  agent_id: node-7
task:
  report:
    agent_id: ""
`
	cfg, err := NewConfig(writeTempYAML(t, yaml), Prefix(""))
	if err != nil {
		t.Fatalf("NewConfig: %v", err)
	}
	if cfg.Control.AgentID != "node-7" {
		t.Fatalf("control.agent_id = %q, want node-7", cfg.Control.AgentID)
	}
	if cfg.Task.Report.AgentID != "node-7" {
		t.Fatalf("report.agent_id = %q, want node-7 (inherited)", cfg.Task.Report.AgentID)
	}
}

func TestNormalizeAgentIDMismatchUsesControl(t *testing.T) {
	const yaml = `
control:
  agent_id: node-a
task:
  report:
    agent_id: node-b
`
	cfg, err := NewConfig(writeTempYAML(t, yaml), Prefix(""))
	if err != nil {
		t.Fatalf("NewConfig: %v", err)
	}
	if cfg.Control.AgentID != "node-a" {
		t.Fatalf("control.agent_id = %q, want node-a", cfg.Control.AgentID)
	}
	if cfg.Task.Report.AgentID != "node-a" {
		t.Fatalf("report.agent_id = %q, want node-a (aligned to control)", cfg.Task.Report.AgentID)
	}
}
