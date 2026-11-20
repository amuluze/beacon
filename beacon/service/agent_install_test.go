package service

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

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
	if !strings.Contains(script, `--ca-file`) || !strings.Contains(script, `--cacert "$CA_FILE"`) {
		t.Fatal("install script must support an explicit CA file for pinned self-signed deployments")
	}
	caIndex := strings.Index(script, `--cacert "$CA_FILE"`)
	plainIndex := strings.Index(script, `curl -fsSL -H "X-Install-Token: $TOKEN"`)
	if caIndex < 0 || plainIndex < 0 || caIndex > plainIndex {
		t.Fatal("install script must try the pinned CA before system trust and must not leak the install token over an unverified first request")
	}
	if !strings.Contains(script, `-o "$2" 2>/dev/null`) {
		t.Fatal("install script must silence intermediate TLS probe failures")
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
	// 瞬时传输中断（curl 18）要自动重试，避免一次抖动导致整次安装失败。
	if !strings.Contains(script, `while [ "$attempt" -le 3 ]`) || !strings.Contains(script, "下载中断，2 秒后重试") {
		t.Fatal("expected script to retry transient download failures")
	}
}

func TestBuildAgentInstallScriptDownloadsCertsWhenTLSEnabled(t *testing.T) {
	script := buildAgentInstallScript("https://beacon.example.com", "agent-a", true)
	if !strings.Contains(script, `/api/v1/host/install/certs?node=$NODE`) {
		t.Fatal("expected script to download collia cert package when TLS is enabled")
	}
}

// 安装脚本必须能通过 POSIX shell 语法检查，重试循环等改动不得引入语法错误。
func TestBuildAgentInstallScriptIsValidShell(t *testing.T) {
	script := buildAgentInstallScript("https://beacon.example.com", "agent-a", true)
	cmd := exec.Command("sh", "-n")
	cmd.Stdin = strings.NewReader(script)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("install script has invalid shell syntax: %v\n%s", err, out)
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

func TestAgentInstallCAEndpointAndFingerprint(t *testing.T) {
	certDir := t.TempDir()
	config := &Config{
		Control: Control{TLS: ControlTLS{Enable: true, CertDir: certDir}},
		AgentInstall: AgentInstall{
			Enable:     true,
			Token:      "install-secret",
			TLSEnable:  true,
			PackageDir: t.TempDir(),
		},
	}
	router := &Router{config: config, certManager: NewCertManager(config)}
	if err := router.certManager.EnsureControlCerts(); err != nil {
		t.Fatalf("ensure control certs: %v", err)
	}

	app := fiber.New()
	// 模拟生产环境的静态文件中间件：未命中的路径会先写入 404，
	// 后续的 SendFile 处理器必须显式恢复 200。
	app.Use(func(c *fiber.Ctx) error {
		c.Status(http.StatusNotFound)
		return c.Next()
	})
	app.Get("/api/v1/host/install/ca", router.AgentInstallCA)
	app.Post("/api/v1/host/get_install_token", router.AgentInstallToken)

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/v1/host/install/ca", nil), -1)
	if err != nil {
		t.Fatal(err)
	}
	caBody, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK || !strings.Contains(string(caBody), "BEGIN CERTIFICATE") {
		t.Fatalf("CA endpoint status=%d body=%q", resp.StatusCode, string(caBody[:min(len(caBody), 80)]))
	}

	resp, err = app.Test(httptest.NewRequest(http.MethodPost, "/api/v1/host/get_install_token", nil), -1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("token endpoint status = %d, want 200", resp.StatusCode)
	}
	var payload struct {
		Token         string `json:"token"`
		CAFingerprint string `json:"ca_fingerprint"`
		ExpiresAt     int64  `json:"expires_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode token response: %v", err)
	}
	// 默认 TTL 下端点签发短时令牌，静态令牌不再直接暴露。
	if !strings.HasPrefix(payload.Token, "bpit.") {
		t.Fatalf("token payload must be ephemeral, got %q", payload.Token)
	}
	if payload.ExpiresAt <= 0 {
		t.Fatalf("ephemeral token must carry expires_at, got %d", payload.ExpiresAt)
	}
	if len(payload.CAFingerprint) != 64 {
		t.Fatalf("unexpected fingerprint payload: %+v", payload)
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

func TestIssueAndVerifyEphemeralInstallToken(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	token, expiresAt, err := issueEphemeralInstallToken("install-secret", 1800, now)
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}
	if expiresAt != now.Unix()+1800 {
		t.Fatalf("expiresAt = %d, want %d", expiresAt, now.Unix()+1800)
	}
	if !strings.HasPrefix(token, "bpit.") {
		t.Fatalf("token %q must carry the ephemeral prefix", token)
	}
	if !verifyEphemeralInstallToken("install-secret", token, now.Add(1799*time.Second)) {
		t.Fatal("valid token within TTL must verify")
	}
	if verifyEphemeralInstallToken("install-secret", token, now.Add(1801*time.Second)) {
		t.Fatal("expired token must be rejected")
	}
	if verifyEphemeralInstallToken("other-secret", token, now) {
		t.Fatal("token minted under a different static token must be rejected")
	}
	if verifyEphemeralInstallToken("install-secret", token+"x", now) {
		t.Fatal("tampered token must be rejected")
	}
	if verifyEphemeralInstallToken("install-secret", "bpit.oops", now) {
		t.Fatal("malformed token must be rejected")
	}
	if verifyEphemeralInstallToken("install-secret", "plain-static-token", now) {
		t.Fatal("non-ephemeral format must not verify as ephemeral")
	}
}

func TestAgentInstallTokenEndpointIssuesEphemeralToken(t *testing.T) {
	router := &Router{config: &Config{AgentInstall: AgentInstall{Enable: true, Token: "install-secret", TokenTTL: 600}}}
	app := fiber.New()
	app.Post("/api/v1/host/get_install_token", router.AgentInstallToken)
	app.Get("/api/v1/host/install/config", router.AgentInstallConfig)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/host/get_install_token", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("token endpoint status = %d, want 200", resp.StatusCode)
	}
	var payload struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode token response: %v", err)
	}
	if !strings.HasPrefix(payload.Token, "bpit.") {
		t.Fatalf("token = %q, want ephemeral format", payload.Token)
	}
	if payload.ExpiresAt <= time.Now().Unix() {
		t.Fatalf("expires_at = %d, must be in the future", payload.ExpiresAt)
	}

	// 短时令牌必须能通过安装期受信端点。
	installReq := httptest.NewRequest(http.MethodGet, "/api/v1/host/install/config?node=node-01", nil)
	installReq.Header.Set("X-Install-Token", payload.Token)
	installResp, err := app.Test(installReq, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = installResp.Body.Close() }()
	if installResp.StatusCode != http.StatusOK {
		t.Fatalf("install config with ephemeral token = %d, want 200", installResp.StatusCode)
	}
}

func TestAgentInstallTokenLegacyModeReturnsStaticToken(t *testing.T) {
	router := &Router{config: &Config{AgentInstall: AgentInstall{Enable: true, Token: "install-secret", TokenTTL: -1}}}
	app := fiber.New()
	app.Post("/api/v1/host/get_install_token", router.AgentInstallToken)

	resp, err := app.Test(httptest.NewRequest(http.MethodPost, "/api/v1/host/get_install_token", nil), -1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	var payload struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode token response: %v", err)
	}
	if payload.Token != "install-secret" {
		t.Fatalf("legacy mode token = %q, want static token", payload.Token)
	}
	if payload.ExpiresAt != 0 {
		t.Fatalf("legacy mode must omit expires_at, got %d", payload.ExpiresAt)
	}
}

func TestBuildColliaConfigUsesUnixSocket(t *testing.T) {
	const socket = "/data/beacon/data/control/collia.sock"
	router := &Router{config: &Config{
		AgentInstall: AgentInstall{
			Enable:        true,
			Token:         "install-secret",
			PublicBaseURL: "https://beacon.example.com",
			ControlSocket: socket,
		},
		Control: Control{TLS: ControlTLS{Enable: true}},
	}}
	cfg := router.buildColliaConfig("node-01")
	if !strings.Contains(cfg, "server: unix://"+socket) {
		t.Fatalf("agent config must use the unix socket target:\n%s", cfg)
	}
	if !strings.Contains(cfg, "enable: false") {
		t.Fatalf("unix target must disable agent mTLS:\n%s", cfg)
	}
	if !strings.Contains(cfg, "url: \"https://beacon.example.com/api/v1/host/report\"") {
		t.Fatalf("HTTP report fallback must stay available:\n%s", cfg)
	}
	if strings.Contains(cfg, "enable: true") {
		t.Fatalf("unix target must not request agent certificates:\n%s", cfg)
	}
}

func TestBuildColliaConfigEmbedsUpgradePublicKey(t *testing.T) {
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	encoded := base64.StdEncoding.EncodeToString(pub)
	router := &Router{config: &Config{
		AgentInstall: AgentInstall{Enable: true, Token: "install-secret", PublicBaseURL: "https://beacon.example.com"},
		AgentUpdate:  AgentUpdate{PublicKey: encoded},
	}}
	cfg := router.buildColliaConfig("node-01")
	if !strings.Contains(cfg, fmt.Sprintf("upgrade_public_key: %q", encoded)) {
		t.Fatalf("agent config must pin the release public key, got:\n%s", cfg)
	}

	router.config.AgentUpdate.PublicKey = ""
	cfg = router.buildColliaConfig("node-01")
	if strings.Contains(cfg, "upgrade_public_key") {
		t.Fatalf("unconfigured public key must not be written into agent config:\n%s", cfg)
	}

	// 轮换重叠窗口：逗号分隔的两把公钥写成 YAML 列表并逐项引用。
	second, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	secondEncoded := base64.StdEncoding.EncodeToString(second)
	router.config.AgentUpdate.PublicKey = encoded + "," + secondEncoded
	cfg = router.buildColliaConfig("node-01")
	if !strings.Contains(cfg, "upgrade_public_keys:") {
		t.Fatalf("multi-key config must use the list field, got:\n%s", cfg)
	}
	if !strings.Contains(cfg, fmt.Sprintf("- %q", encoded)) || !strings.Contains(cfg, fmt.Sprintf("- %q", secondEncoded)) {
		t.Fatalf("multi-key config must list both keys, got:\n%s", cfg)
	}
	if strings.Contains(cfg, "upgrade_public_key: \"") {
		t.Fatalf("multi-key config must not emit the single-key field:\n%s", cfg)
	}
}
