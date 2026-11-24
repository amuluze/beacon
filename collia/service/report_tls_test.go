package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"collia/service/report"
	rpcSchema "common/rpc/schema"
)

// 上报端点是 Server 的公开 HTTPS 接口（X-Install-Token 鉴权），
// 证书必须按上报 URL 的主机名并对照系统信任库校验；控制通道不再
// 携带任何证书（Unix Socket），也不存在可复用的控制 CA。
func TestBuildReportTLSConfigUsesSystemTrustOnly(t *testing.T) {
	cfg := buildReportTLSConfig(&Config{})
	if cfg == nil {
		t.Fatal("expected tls config for https report client")
	}
	if cfg.ServerName != "" {
		t.Fatalf("report client must not pin a server name, got %q", cfg.ServerName)
	}
	if cfg.RootCAs == nil {
		t.Fatal("expected report client to trust the system store")
	}
	if cfg.InsecureSkipVerify {
		t.Fatal("report client must never skip certificate verification")
	}
}

// 端到端：系统信任库不认识的自签证书必须被拒绝——上报客户端
// 禁止为兼容自签部署而降级校验。
func TestReportClientRejectsUntrustedSelfSignedCert(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Install-Token") != "install-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := report.NewClient(srv.URL, "install-token", buildReportTLSConfig(&Config{}))
	if err := client.Push(context.Background(), rpcSchema.MonitorReportArgs{AgentID: "agent-1"}); err == nil {
		t.Fatal("report client must reject certificates not trusted by the system store")
	}
}
