package service

import (
	"context"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"collia/service/report"
	rpcSchema "common/rpc/schema"
)

// 上报端点是 Server 的公开 HTTPS 接口（X-Install-Token 鉴权），
// 证书必须按上报 URL 的主机名校验；控制通道的 mTLS server_name
// 只标识隧道服务端证书身份（"beacon/collia"），不能用于上报。
func TestBuildReportTLSConfigDoesNotPinControlPlaneServerName(t *testing.T) {
	conf := &Config{}
	conf.Control.TLS.Enable = true
	conf.Control.TLS.CertDir = t.TempDir()
	conf.Control.TLS.ServerName = "beacon/collia"

	cfg := buildReportTLSConfig(conf)
	if cfg == nil {
		t.Fatal("expected tls config for https report client")
	}
	if cfg.ServerName != "" {
		t.Fatalf("report client must not pin control plane server name, got %q", cfg.ServerName)
	}
	if cfg.RootCAs == nil {
		t.Fatal("expected report client to trust the system store plus control plane ca")
	}
}

func TestBuildReportTLSConfigNilWhenControlTLSDisabled(t *testing.T) {
	conf := &Config{}
	if cfg := buildReportTLSConfig(conf); cfg != nil {
		t.Fatal("expected nil tls config when control plane tls is disabled")
	}
}

// 端到端：自建 CA 的 HTTPS 上报端点必须能收到带 X-Install-Token 的上报。
// 修复前上报客户端把控制通道 server_name 用作证书校验目标，会因
// "certificate is valid for <host>, not beacon/collia" 而失败。
func TestReportClientPushesToHTTPSEndpointWithControlPlaneCA(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Install-Token") != "install-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	certDir := t.TempDir()
	caPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: srv.Certificate().Raw})
	if err := os.WriteFile(filepath.Join(certDir, "ca.pem"), caPEM, 0o600); err != nil {
		t.Fatalf("write control plane ca: %v", err)
	}

	conf := &Config{}
	conf.Control.TLS.Enable = true
	conf.Control.TLS.CertDir = certDir
	conf.Control.TLS.ServerName = "beacon/collia"

	client := report.NewClient(srv.URL, "install-token", buildReportTLSConfig(conf))
	if err := client.Push(context.Background(), rpcSchema.MonitorReportArgs{AgentID: "agent-1"}); err != nil {
		t.Fatalf("push report over https: %v", err)
	}
}
