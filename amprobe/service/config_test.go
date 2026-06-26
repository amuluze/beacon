package service

import (
	"os"
	"path/filepath"
	"testing"
)

// writeConfig 写入一个最小可用的临时 toml 配置，供 NewConfig 加载。
func writeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return path
}

// TestNewConfig_EnvOverridesSigningKey 验证 AMPROBE_AUTH_SIGNINGKEY 能覆盖 toml 中的弱默认值。
// viper 对 AutomaticEnv + Unmarshal 有已知坑，必须靠显式 BindEnv 才能让 Unmarshal 拾取。
func TestNewConfig_EnvOverrideSigningKey(t *testing.T) {
	const toml = `
[Auth]
Enable = true
SigningMethod = "HS512"
SigningKey = "amprobe"
Expired = 7200
RefreshExpired = 86400
Prefix = "auth_"
`
	t.Setenv("AMPROBE_AUTH_SIGNINGKEY", "overridden-strong-secret-from-env")

	cfg, err := NewConfig(writeConfig(t, toml))
	if err != nil {
		t.Fatalf("NewConfig: %v", err)
	}
	if cfg.Auth.SigningKey != "overridden-strong-secret-from-env" {
		t.Fatalf("expected SigningKey overridden by env, got %q", cfg.Auth.SigningKey)
	}
}

// TestNewConfig_EnvOverridesJoinToken 验证控制通道 JoinToken 同样可被环境变量覆盖。
func TestNewConfig_EnvOverridesJoinToken(t *testing.T) {
	const toml = `
[Control]
Enable = true
Address = "0.0.0.0:8081"
JoinToken = "from-file"
`
	t.Setenv("AMPROBE_CONTROL_JOINTOKEN", "from-env")

	cfg, err := NewConfig(writeConfig(t, toml))
	if err != nil {
		t.Fatalf("NewConfig: %v", err)
	}
	if cfg.Control.JoinToken != "from-env" {
		t.Fatalf("expected JoinToken overridden by env, got %q", cfg.Control.JoinToken)
	}
}

// TestResolveSigningKey_EmptyGenerates 验证空密钥生成非空临时密钥。
func TestResolveSigningKey_EmptyGenerates(t *testing.T) {
	k := resolveSigningKey("")
	if k == "" || k == "amprobe" {
		t.Fatalf("expected generated non-weak key, got %q", k)
	}
}

// TestResolveSigningKey_WeakDefaultPreserved 验证弱默认值保留（仅告警，不破坏现有部署）。
func TestResolveSigningKey_WeakDefaultPreserved(t *testing.T) {
	if got := resolveSigningKey("amprobe"); got != "amprobe" {
		t.Fatalf("expected weak default preserved, got %q", got)
	}
}

// TestResolveSigningKey_CustomPreserved 验证用户自定义强密钥原样保留。
func TestResolveSigningKey_CustomPreserved(t *testing.T) {
	const custom = "a-very-long-random-production-secret-12345"
	if got := resolveSigningKey(custom); got != custom {
		t.Fatalf("expected custom key preserved, got %q", got)
	}
}
