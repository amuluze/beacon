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
	if _, err := safeJoin("/tmp/beacon", "linux", "amd64", "collia"); err != nil {
		t.Fatalf("expected safe path, got %v", err)
	}
	if _, err := safeJoin("/tmp/beacon", "..", "collia"); err == nil {
		t.Fatal("expected path escape to be rejected")
	}
}

func TestBuildAgentInstallScriptUsesInstallTokenHeader(t *testing.T) {
	script := buildAgentInstallScript("http://127.0.0.1:1443", "1", false)
	if !strings.Contains(script, `X-Install-Token: $TOKEN`) {
		t.Fatal("expected script downloads to use X-Install-Token header")
	}
	if !strings.Contains(script, `/api/v1/host/install/package?arch=$ARCH`) {
		t.Fatal("expected script to download collia binary selected by arch")
	}
}

func TestBuildAgentInstallScriptDownloadsCertsWhenTLSEnabled(t *testing.T) {
	script := buildAgentInstallScript("http://127.0.0.1:1443", "agent-a", true)
	if !strings.Contains(script, `/api/v1/host/install/certs?node=$NODE`) {
		t.Fatal("expected script to download collia cert package when TLS is enabled")
	}
}

func TestBuildColliaConfigUsesControlJoinToken(t *testing.T) {
	router := &Router{config: &Config{
		Control:      Control{JoinToken: "control-secret"},
		AgentInstall: AgentInstall{Token: "install-secret"},
	}}

	config := router.buildColliaConfig("agent-a")
	if !strings.Contains(config, `join_token: "control-secret"`) {
		t.Fatalf("expected control join token in config, got:\n%s", config)
	}
	if !strings.Contains(config, `token: "install-secret"`) {
		t.Fatalf("expected report/install token in config, got:\n%s", config)
	}
}

func TestBuildColliaConfigUsesControlTLSFlag(t *testing.T) {
	router := &Router{config: &Config{
		Control:      Control{TLS: ControlTLS{Enable: true}},
		AgentInstall: AgentInstall{CertDir: "/etc/collia/certs"},
	}}

	config := router.buildColliaConfig("agent-a")
	if !strings.Contains(config, `enable: true`) {
		t.Fatalf("expected TLS enable flag in config, got:\n%s", config)
	}
	if !strings.Contains(config, `cert_dir: /etc/collia/certs`) {
		t.Fatalf("expected TLS cert dir in config, got:\n%s", config)
	}
}
