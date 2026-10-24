package service

import (
	"errors"
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

// TestNewConfig_EnvOverridesSigningKey 验证 BEACON_AUTH_SIGNINGKEY 能覆盖 toml 中的弱默认值。
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
	t.Setenv("BEACON_AUTH_SIGNINGKEY", "overridden-strong-secret-from-env")

	cfg, err := NewConfig(writeConfig(t, toml))
	if err != nil {
		t.Fatalf("NewConfig: %v", err)
	}
	if cfg.Auth.SigningKey != "overridden-strong-secret-from-env" {
		t.Fatalf("expected SigningKey overridden by env, got %q", cfg.Auth.SigningKey)
	}
}

func TestNewConfig_EnvOverrideAppEnv(t *testing.T) {
	const toml = `
[App]
Env = "development"
`
	t.Setenv("BEACON_APP_ENV", "production")

	cfg, err := NewConfig(writeConfig(t, toml))
	if err != nil {
		t.Fatalf("NewConfig: %v", err)
	}
	if cfg.App.Env != "production" {
		t.Fatalf("expected App.Env overridden by env, got %q", cfg.App.Env)
	}
}

func TestNewConfigRejectsInsecureProductionTransport(t *testing.T) {
	const toml = `
[App]
Env = "production"

[AgentInstall]
Enable = true
PublicBaseURL = "http://beacon.example.com"
TLSEnable = true

[Control]
Enable = true

[Control.TLS]
Enable = true
`
	_, err := NewConfig(writeConfig(t, toml))
	if !errors.Is(err, ErrHTTPSPublicBaseURLRequired) {
		t.Fatalf("NewConfig() error = %v, want ErrHTTPSPublicBaseURLRequired", err)
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
	t.Setenv("BEACON_CONTROL_JOINTOKEN", "from-env")

	cfg, err := NewConfig(writeConfig(t, toml))
	if err != nil {
		t.Fatalf("NewConfig: %v", err)
	}
	if cfg.Control.JoinToken != "from-env" {
		t.Fatalf("expected JoinToken overridden by env, got %q", cfg.Control.JoinToken)
	}
}

// TestNewConfig_BindEnvOverrides 验证 bindEnvs 绑定的 BEACON_* 字段可经环境变量覆盖，
// 这是替代容器启动 sed 改写 config.toml 的契约。
func TestNewConfig_BindEnvOverrides(t *testing.T) {
	const toml = `
[DB]
Host = ""
Port = ""
User = ""
Password = ""
DBName = "from-file"
SSLMode = ""

[Control]
Enable = true
Address = "0.0.0.0:8081"
JoinToken = ""

[AgentInstall]
Enable = true
Token = "change-me"
PublicBaseURL = ""
PackageDir = "/app/downloads/collia"
ControlPort = 18080
TLSEnable = true
CertDir = "/etc/collia/certs"

[AgentUpdate]
Enable = false
TargetVersion = "dev"
Timeout = 10
MaxConcurrent = 1
RetryInterval = 60
`
	t.Setenv("BEACON_DB_NAME", "/data/beacon-prod")
	t.Setenv("BEACON_CONTROL_ADDRESS", ":17000")
	t.Setenv("BEACON_PUBLIC_BASE_URL", "https://beacon.example.com")
	t.Setenv("BEACON_CONTROL_PORT", "17000")
	t.Setenv("BEACON_AGENT_UPDATE_ENABLE", "true")
	t.Setenv("BEACON_COLLIA_VERSION", "v3.0.6")
	t.Setenv("BEACON_AGENT_UPDATE_TIMEOUT", "300")
	t.Setenv("BEACON_AGENT_UPDATE_MAX_CONCURRENT", "2")
	t.Setenv("BEACON_AGENT_UPDATE_RETRY_INTERVAL", "1800")

	cfg, err := NewConfig(writeConfig(t, toml))
	if err != nil {
		t.Fatalf("NewConfig: %v", err)
	}
	if cfg.DB.DBName != "/data/beacon-prod" {
		t.Fatalf("expected DB.DBName overridden by env, got %q", cfg.DB.DBName)
	}
	if cfg.Control.Address != ":17000" {
		t.Fatalf("expected Control.Address overridden by env, got %q", cfg.Control.Address)
	}
	if cfg.AgentInstall.PublicBaseURL != "https://beacon.example.com" {
		t.Fatalf("expected AgentInstall.PublicBaseURL overridden by env, got %q", cfg.AgentInstall.PublicBaseURL)
	}
	if cfg.AgentInstall.ControlPort != 17000 {
		t.Fatalf("expected AgentInstall.ControlPort overridden by env, got %d", cfg.AgentInstall.ControlPort)
	}
	if !cfg.AgentUpdate.Enable || cfg.AgentUpdate.TargetVersion != "v3.0.6" {
		t.Fatalf("expected AgentUpdate target enabled at v3.0.6, got %#v", cfg.AgentUpdate)
	}
	if cfg.AgentUpdate.Timeout != 300 || cfg.AgentUpdate.MaxConcurrent != 2 || cfg.AgentUpdate.RetryInterval != 1800 {
		t.Fatalf("expected AgentUpdate limits overridden by env, got %#v", cfg.AgentUpdate)
	}
}

func TestNewConfig_AgentUpdateCanBeDisabledByEnv(t *testing.T) {
	const toml = `
[AgentUpdate]
Enable = true
TargetVersion = "v3.0.6"
Timeout = 300
MaxConcurrent = 2
RetryInterval = 1800
`
	t.Setenv("BEACON_AGENT_UPDATE_ENABLE", "false")

	cfg, err := NewConfig(writeConfig(t, toml))
	if err != nil {
		t.Fatalf("NewConfig: %v", err)
	}
	if cfg.AgentUpdate.Enable {
		t.Fatal("BEACON_AGENT_UPDATE_ENABLE=false must disable automatic updates")
	}
}

// TestResolveSigningKey_EmptyGenerates 验证非生产模式下空密钥生成非空临时密钥。
func TestResolveSigningKey_EmptyGenerates(t *testing.T) {
	k, err := resolveSigningKey("", "development")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if k == "" || k == "amprobe" {
		t.Fatalf("expected generated non-weak key, got %q", k)
	}
}

// TestResolveSigningKey_WeakDefaultPreserved 验证非生产模式下弱默认值保留（仅告警，不破坏现有部署）。
func TestResolveSigningKey_WeakDefaultPreserved(t *testing.T) {
	got, err := resolveSigningKey("amprobe", "development")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "amprobe" {
		t.Fatalf("expected weak default preserved, got %q", got)
	}
}

// TestResolveSigningKey_CustomPreserved 验证用户自定义强密钥原样保留。
func TestResolveSigningKey_CustomPreserved(t *testing.T) {
	const custom = "a-very-long-random-production-secret-12345"
	got, err := resolveSigningKey(custom, "development")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != custom {
		t.Fatalf("expected custom key preserved, got %q", got)
	}
}

// TestResolveSigningKey_ProductionRejectsInsecure 验证生产模式下空/弱/短密钥一律拒绝。
func TestResolveSigningKey_ProductionRejectsInsecure(t *testing.T) {
	cases := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{"empty", "", true},
		{"weak default", "amprobe", true},
		{"weak current default", "beacon", true},
		{"too short", "short", true},
		{"strong", "a-very-long-random-production-secret-12345", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := resolveSigningKey(tc.key, "production")
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for %q in production, got nil", tc.key)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error for strong key in production: %v", err)
			}
		})
	}
}

// TestResolveControlToken_ProductionRejectsInsecure 验证控制通道 JoinToken 在生产模式下
// 空/弱默认/过短一律拒绝，强 token 通过。
func TestResolveControlToken_ProductionRejectsInsecure(t *testing.T) {
	cases := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{"empty", "", true},
		{"weak change-me", "change-me", true},
		{"weak change-me-please", "change-me-please", true},
		{"weak token", "token", true},
		{"too short", "short-token", true},
		{"strong", "a-strong-random-join-token-12345", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := resolveControlToken(tc.token, "production")
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for %q in production, got nil", tc.token)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error for strong token in production: %v", err)
			}
		})
	}
}

// TestResolveControlToken_NonProductionPreserves 验证非生产模式下空/弱默认保留（仅告警）。
func TestResolveControlToken_NonProductionPreserves(t *testing.T) {
	cases := []struct {
		name  string
		token string
	}{
		{"empty", ""},
		{"weak default", "change-me"},
		{"custom", "my-dev-token"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := resolveControlToken(tc.token, "development")
			if err != nil {
				t.Fatalf("unexpected error in development: %v", err)
			}
			if got != tc.token {
				t.Fatalf("expected %q preserved, got %q", tc.token, got)
			}
		})
	}
}

// TestResolveInstallToken_ProductionRejectsInsecure 验证 installToken 在生产模式下
// 空/弱默认/过短一律拒绝，强 token 通过。
func TestResolveInstallToken_ProductionRejectsInsecure(t *testing.T) {
	cases := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{"empty", "", true},
		{"weak change-me", "change-me", true},
		{"weak amprobe", "amprobe", true},
		{"too short", "short-token", true},
		{"strong", "a-strong-random-install-token-12345", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := resolveInstallToken(tc.token, "production")
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for %q in production, got nil", tc.token)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error for strong token in production: %v", err)
			}
		})
	}
}

// TestResolveInstallToken_NonProductionPreserves 验证非生产模式下空/弱默认保留（仅告警）。
func TestResolveInstallToken_NonProductionPreserves(t *testing.T) {
	got, err := resolveInstallToken("change-me", "development")
	if err != nil {
		t.Fatalf("unexpected error in development: %v", err)
	}
	if got != "change-me" {
		t.Fatalf("expected weak default preserved, got %q", got)
	}
}
