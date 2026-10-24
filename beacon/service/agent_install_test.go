package service

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	servicemiddleware "beacon/service/middleware"

	"github.com/gofiber/fiber/v2"
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
	script := buildAgentInstallScript("https://beacon.example.com", "1", false)
	if !strings.Contains(script, `X-Install-Token: $TOKEN`) {
		t.Fatal("expected script downloads to use X-Install-Token header")
	}
	if !strings.Contains(script, `/api/v1/host/install/package?arch=$ARCH`) {
		t.Fatal("expected script to download collia binary selected by arch")
	}
	if strings.Contains(script, "curl -k") || strings.Contains(script, "--insecure") {
		t.Fatal("install script must not disable TLS certificate verification")
	}
	if !strings.Contains(script, "umask 077") {
		t.Fatal("install script must use an owner-only umask")
	}
	if !strings.Contains(script, `install -m 0600 "$WORKDIR/config.yml" /etc/collia/config.yml`) {
		t.Fatal("install script must install the collia config with owner-only permissions")
	}
}

func TestBuildAgentInstallScriptSurvivesStopAndStartFailures(t *testing.T) {
	script := buildAgentInstallScript("https://beacon.example.com", "1", false)
	// 首次安装时 collia stop 返回非 0，若不忽略会在 set -e 下中断脚本，探针永远不启动。
	if !strings.Contains(script, "collia stop >/dev/null 2>&1 || true") {
		t.Fatal("expected script to ignore collia stop failure so installation continues")
	}
	if strings.Contains(script, "if collia status >/dev/null 2>&1; then") {
		t.Fatal("expected script to stop guarding on collia status, which returns 0 even when stopped")
	}
	// start 失败必须给出可定位原因与修复方向（RR004）。
	if !strings.Contains(script, "探针启动失败") || !strings.Contains(script, "collia run") {
		t.Fatal("expected script to report actionable diagnostics when collia start fails")
	}
	// 无 systemd 的环境要提前提示，否则用户只看到 [FAILED]。
	if !strings.Contains(script, "未检测到 systemd") {
		t.Fatal("expected script to warn when systemd is unavailable")
	}
	// 下载失败要提示自签证书场景下的 CA 导入方式。
	if !strings.Contains(script, "下载失败") || !strings.Contains(script, "update-ca-certificates") {
		t.Fatal("expected script to explain download failures including self-signed CA guidance")
	}
}

func TestBuildAgentInstallScriptDownloadsCertsWhenTLSEnabled(t *testing.T) {
	script := buildAgentInstallScript("https://beacon.example.com", "agent-a", true)
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

func TestAgentInstallEndpointsResetSuccessStatus(t *testing.T) {
	packageDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(packageDir, "amd64"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(packageDir, "certs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(packageDir, "amd64", colliaBinaryName), []byte("collia"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(packageDir, "certs", "node-01.tar.gz"), []byte("certs"), 0o644); err != nil {
		t.Fatal(err)
	}

	router := &Router{config: &Config{
		Control: Control{},
		AgentInstall: AgentInstall{
			Enable:     true,
			Token:      "install-secret",
			PackageDir: packageDir,
		},
	}}
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Status(http.StatusNotFound)
		return c.Next()
	})
	app.Get("/api/v1/host/install", router.AgentInstallScript)
	app.Get("/api/v1/host/install/package", router.AgentInstallPackage)
	app.Get("/api/v1/host/install/config", router.AgentInstallConfig)
	app.Get("/api/v1/host/install/certs", router.AgentInstallCerts)

	tests := []struct {
		name  string
		path  string
		token bool
	}{
		{name: "script", path: "/api/v1/host/install?node=node-01"},
		{name: "package", path: "/api/v1/host/install/package?arch=amd64", token: true},
		{name: "config", path: "/api/v1/host/install/config?node=node-01", token: true},
		{name: "certs", path: "/api/v1/host/install/certs?node=node-01", token: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			if tt.token {
				req.Header.Set("X-Install-Token", "install-secret")
			}
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			if err := resp.Body.Close(); err != nil {
				t.Fatalf("close response body: %v", err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
			}
		})
	}
}

func TestCommonAuthSkipperAllowsAgentReport(t *testing.T) {
	var skipped bool
	app := fiber.New()
	app.Get("/api/v1/host/report", func(c *fiber.Ctx) error {
		skipped = servicemiddleware.SkipHandler(c, servicemiddleware.AllowPathPrefixSkipper(commonAuthSkipperPaths...))
		return c.SendStatus(http.StatusOK)
	})

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/v1/host/report", nil), -1)
	if err != nil {
		t.Fatal(err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Fatalf("close response body: %v", err)
	}
	if !skipped {
		t.Fatal("agent report endpoint must bypass user JWT middleware and use its install-token authentication")
	}
}

func TestAgentInstallRejectsQueryTokenFallback(t *testing.T) {
	router := &Router{config: &Config{AgentInstall: AgentInstall{Enable: true, Token: "install-secret"}}}
	app := fiber.New()
	app.Get("/api/v1/host/install/config", router.AgentInstallConfig)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/host/install/config?node=node-01&token=install-secret", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("query token status = %d, want 401", resp.StatusCode)
	}
}

func TestAgentInstallCertsGeneratesWhenTLSEnabled(t *testing.T) {
	certDir := t.TempDir()
	packageDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(packageDir, "amd64"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(packageDir, "amd64", colliaBinaryName), []byte("collia"), 0o755); err != nil {
		t.Fatal(err)
	}

	router := &Router{config: &Config{
		Control: Control{TLS: ControlTLS{Enable: true, CertDir: certDir}},
		AgentInstall: AgentInstall{
			Enable:     true,
			Token:      "install-secret",
			PackageDir: packageDir,
			TLSEnable:  true,
		},
	}}
	router.certManager = NewCertManager(router.config)

	app := fiber.New()
	app.Get("/api/v1/host/install/certs", router.AgentInstallCerts)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/host/install/certs?node=node-01", nil)
	req.Header.Set("X-Install-Token", "install-secret")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	pkgPath := filepath.Join(packageDir, "certs", "node-01.tar.gz")
	if _, err := os.Stat(pkgPath); err != nil {
		t.Fatalf("agent cert package not generated: %v", err)
	}
}
