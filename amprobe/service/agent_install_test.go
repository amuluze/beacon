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
	if _, err := safeJoin("/tmp/amprobe", "linux", "amd64", "collia.install"); err != nil {
		t.Fatalf("expected safe path, got %v", err)
	}
	if _, err := safeJoin("/tmp/amprobe", "..", "collia.install"); err == nil {
		t.Fatal("expected path escape to be rejected")
	}
}

func TestBuildAgentInstallScriptUsesInstallTokenHeader(t *testing.T) {
	script := buildAgentInstallScript("http://127.0.0.1:1443", "1", "linux", 18080, true)
	if !strings.Contains(script, `X-Install-Token: $TOKEN`) {
		t.Fatal("expected script downloads to use X-Install-Token header")
	}
	if !strings.Contains(script, `/api/v1/host/install/package?node=$NODE&os_type=$OS_TYPE&arch=$ARCH`) {
		t.Fatal("expected script to download collia package")
	}
	if !strings.Contains(script, `/api/v1/host/install/certs?node=$NODE`) {
		t.Fatal("expected script to download certs when TLS is enabled")
	}
}
