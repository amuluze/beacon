package service

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	transporttls "common/transport/tlsconfig"
)

func TestCertManagerEnsureControlCertsGeneratesMissing(t *testing.T) {
	certDir := t.TempDir()
	cm := NewCertManager(&Config{
		Control: Control{TLS: ControlTLS{Enable: true, CertDir: certDir}},
	})

	if err := cm.EnsureControlCerts(); err != nil {
		t.Fatalf("EnsureControlCerts failed: %v", err)
	}

	for _, name := range []string{transporttls.CACertFile, transporttls.TLSCertFile, transporttls.TLSKeyFile} {
		if _, err := os.Stat(filepath.Join(certDir, name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
}

func TestCertManagerLeafCoversPublicBaseURLHostSAN(t *testing.T) {
	cases := []struct {
		name string
		base string
		ip   string // 期望出现的 IP SAN（空则断言不存在）
		dns  string // 期望出现的 DNS SAN（空则断言不存在）
	}{
		{name: "ip-host-uses-ip-san", base: "https://10.9.35.5:1443", ip: "10.9.35.5", dns: "beacon/collia"},
		{name: "domain-host-uses-dns-san", base: "https://beacon.example.com", dns: "beacon.example.com", ip: "127.0.0.1"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			certDir := t.TempDir()
			cm := NewCertManager(&Config{
				Control:      Control{TLS: ControlTLS{Enable: true, CertDir: certDir}},
				AgentInstall: AgentInstall{PublicBaseURL: tc.base},
			})

			if err := cm.EnsureControlCerts(); err != nil {
				t.Fatalf("EnsureControlCerts failed: %v", err)
			}

			raw, err := os.ReadFile(filepath.Join(certDir, transporttls.TLSCertFile))
			if err != nil {
				t.Fatalf("read server cert: %v", err)
			}
			block, _ := pem.Decode(raw)
			if block == nil {
				t.Fatal("server cert is not PEM")
			}
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				t.Fatalf("parse server cert: %v", err)
			}

			hasDNS := func(want string) bool {
				for _, dns := range cert.DNSNames {
					if dns == want {
						return true
					}
				}
				return false
			}
			hasIP := func(want string) bool {
				for _, ip := range cert.IPAddresses {
					if ip.String() == want {
						return true
					}
				}
				return false
			}

			if tc.ip != "" && !hasIP(tc.ip) {
				t.Fatalf("missing IP SAN %s (dns=%v ip=%v)", tc.ip, cert.DNSNames, cert.IPAddresses)
			}
			if tc.ip != "" && hasDNS(tc.ip) {
				t.Fatalf("IP host must not be emitted as DNS SAN (dns=%v)", cert.DNSNames)
			}
			if tc.dns != "" && !hasDNS(tc.dns) {
				t.Fatalf("missing DNS SAN %s (dns=%v ip=%v)", tc.dns, cert.DNSNames, cert.IPAddresses)
			}
			if tc.ip == "" && !hasDNS("beacon/collia") {
				t.Fatalf("missing control server name beacon/collia (dns=%v)", cert.DNSNames)
			}
			if !hasIP("127.0.0.1") {
				t.Fatalf("missing loopback IP SAN (ip=%v)", cert.IPAddresses)
			}
		})
	}
}

func TestCertManagerEnsureControlCertsPreservesExisting(t *testing.T) {
	certDir := t.TempDir()
	if err := transporttls.GenerateCA(certDir); err != nil {
		t.Fatal(err)
	}
	if err := transporttls.GenerateLeafCert(certDir, certDir, "existing", []string{"existing"}, nil); err != nil {
		t.Fatal(err)
	}

	// Capture original cert content to prove it is not overwritten.
	original, err := os.ReadFile(filepath.Join(certDir, transporttls.TLSCertFile))
	if err != nil {
		t.Fatal(err)
	}

	cm := NewCertManager(&Config{
		Control: Control{TLS: ControlTLS{Enable: true, CertDir: certDir}},
	})
	if err := cm.EnsureControlCerts(); err != nil {
		t.Fatalf("EnsureControlCerts failed: %v", err)
	}

	after, err := os.ReadFile(filepath.Join(certDir, transporttls.TLSCertFile))
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(original) {
		t.Fatal("existing server certificate was overwritten")
	}
}

func TestCertManagerEnsureControlCertsNoopWhenTLSDisabled(t *testing.T) {
	certDir := t.TempDir()
	cm := NewCertManager(&Config{
		Control: Control{TLS: ControlTLS{Enable: false, CertDir: certDir}},
	})

	if err := cm.EnsureControlCerts(); err != nil {
		t.Fatalf("EnsureControlCerts failed: %v", err)
	}

	entries, err := os.ReadDir(certDir)
	if err != nil {
		t.Fatal(err)
	}
	// CA is always generated when missing; server cert is skipped when TLS disabled.
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries (ca.pem, ca.key) when TLS disabled, got %d entries", len(entries))
	}
}

func TestCertManagerGenerateAgentCertPackage(t *testing.T) {
	certDir := t.TempDir()
	packageDir := t.TempDir()

	if err := transporttls.GenerateCA(certDir); err != nil {
		t.Fatal(err)
	}
	if err := transporttls.GenerateLeafCert(certDir, certDir, controlServerName(), []string{controlServerName()}, []string{"127.0.0.1"}); err != nil {
		t.Fatal(err)
	}

	cm := NewCertManager(&Config{
		Control: Control{TLS: ControlTLS{Enable: true, CertDir: certDir}},
		AgentInstall: AgentInstall{
			TLSEnable:  true,
			PackageDir: packageDir,
		},
	})

	path, err := cm.GenerateAgentCertPackage("node-01")
	if err != nil {
		t.Fatalf("GenerateAgentCertPackage failed: %v", err)
	}

	expected := filepath.Join(packageDir, "certs", "node-01.tar.gz")
	if path != expected {
		t.Fatalf("unexpected package path: %s, want %s", path, expected)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("package not created: %v", err)
	}
}

func TestCertManagerGenerateAgentCertPackageRejectsUnsafeNode(t *testing.T) {
	certDir := t.TempDir()
	packageDir := t.TempDir()

	if err := transporttls.GenerateCA(certDir); err != nil {
		t.Fatal(err)
	}

	cm := NewCertManager(&Config{
		Control: Control{TLS: ControlTLS{Enable: true, CertDir: certDir}},
		AgentInstall: AgentInstall{
			TLSEnable:  true,
			PackageDir: packageDir,
		},
	})

	if _, err := cm.GenerateAgentCertPackage("../etc"); err == nil {
		t.Fatal("expected path traversal to be rejected")
	}
}

func TestCertManagerGenerateAgentCertPackageRequiresTLS(t *testing.T) {
	cm := NewCertManager(&Config{
		Control:      Control{TLS: ControlTLS{Enable: true, CertDir: t.TempDir()}},
		AgentInstall: AgentInstall{TLSEnable: false},
	},
	)

	if _, err := cm.GenerateAgentCertPackage("node-01"); err == nil {
		t.Fatal("expected error when agent install TLS is disabled")
	}
}
