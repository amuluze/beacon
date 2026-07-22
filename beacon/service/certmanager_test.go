package service

import (
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
