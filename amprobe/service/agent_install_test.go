package service

import (
	"strings"
	"testing"
)

func TestSafeInstallName(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "1", want: true},
		{name: "host-a", want: true},
		{name: "linux", want: true},
		{name: ".", want: false},
		{name: "..", want: false},
		{name: "../host", want: false},
		{name: "", want: false},
	}

	for _, tt := range tests {
		if got := isSafeInstallName(tt.name); got != tt.want {
			t.Fatalf("isSafeInstallName(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSafeJoinRejectsEscapes(t *testing.T) {
	if _, err := safeJoin("/tmp/amprobe", "linux", "amd64", "collia"); err != nil {
		t.Fatalf("expected safe path, got %v", err)
	}
	if _, err := safeJoin("/tmp/amprobe", "..", "collia"); err == nil {
		t.Fatal("expected path escape to be rejected")
	}
}

func TestBuildAgentInstallScriptUsesInstallTokenHeader(t *testing.T) {
	script := buildAgentInstallScript("http://127.0.0.1:1443", "1")
	if !strings.Contains(script, `X-Install-Token: $TOKEN`) {
		t.Fatal("expected script downloads to use X-Install-Token header")
	}
	if !strings.Contains(script, `/api/v1/host/install/package?arch=$ARCH`) {
		t.Fatal("expected script to download collia binary selected by arch")
	}
}

// TestBuildColliaConfig_InjectsJoinToken 验证生成的 agent 配置携带控制通道 JoinToken，
// 保证 server 强制鉴权后新安装的 agent 能正常注册（server/agent token 自洽）。
func TestBuildColliaConfig_InjectsJoinToken(t *testing.T) {
	r := &Router{config: &Config{
		Control:      Control{JoinToken: "secret-join-token"},
		AgentInstall: AgentInstall{Token: "install-token", PublicBaseURL: "http://srv:1443"},
	}}
	cfg := r.buildColliaConfig("node-1")
	if !strings.Contains(cfg, `join_token: "secret-join-token"`) {
		t.Fatalf("expected join_token injected into agent config, got:\n%s", cfg)
	}
	if strings.Contains(cfg, `join_token: ""`) {
		t.Fatalf("expected non-empty join_token, got empty placeholder:\n%s", cfg)
	}
	// report token 仍来自 AgentInstall.Token，保持既有契约。
	if !strings.Contains(cfg, `token: "install-token"`) {
		t.Fatalf("expected report token preserved, got:\n%s", cfg)
	}
}

