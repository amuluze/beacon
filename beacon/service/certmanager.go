package service

import (
	"fmt"
	"os"
	"path/filepath"

	transporttls "common/transport/tlsconfig"
)

// CertManager handles runtime generation and lifecycle of the self-signed CA
// used for the Beacon control plane and Collia agent mTLS.
type CertManager struct {
	config *Config
}

// NewCertManager creates a cert manager from config.
func NewCertManager(config *Config) *CertManager {
	return &CertManager{config: config}
}

// EnsureControlCerts checks whether the configured control-plane TLS directory
// contains a CA and server certificate. If TLS is enabled and any required file
// is missing, it generates a self-signed CA and a server certificate signed by
// that CA. Existing certificates are never overwritten.
func (m *CertManager) EnsureControlCerts() error {
	certDir := m.config.Control.TLS.CertDir
	if certDir == "" {
		return fmt.Errorf("control tls cert dir is empty")
	}

	caCertPath := filepath.Join(certDir, transporttls.CACertFile)

	if !fileExists(caCertPath) {
		if err := transporttls.GenerateCA(certDir); err != nil {
			return fmt.Errorf("generate control ca: %w", err)
		}
	}

	if !m.config.Control.TLS.Enable {
		return nil
	}

	serverCertPath := filepath.Join(certDir, transporttls.TLSCertFile)
	serverKeyPath := filepath.Join(certDir, transporttls.TLSKeyFile)

	if fileExists(serverCertPath) && fileExists(serverKeyPath) {
		return nil
	}

	dnsNames, ipAddresses := m.serverNames()
	if err := transporttls.GenerateLeafCert(certDir, certDir, controlServerName(), dnsNames, ipAddresses); err != nil {
		return fmt.Errorf("generate control server cert: %w", err)
	}

	return nil
}

// GenerateAgentCertPackage creates a client certificate for the given node,
// packages it with the runtime CA into a tar.gz, and returns the package path.
// The package is written under AgentInstall.PackageDir/certs/{node}.tar.gz.
func (m *CertManager) GenerateAgentCertPackage(node string) (string, error) {
	if !m.config.AgentInstall.TLSEnable {
		return "", fmt.Errorf("agent install tls is disabled")
	}

	if err := m.EnsureControlCerts(); err != nil {
		return "", fmt.Errorf("ensure control certs: %w", err)
	}

	caDir := m.config.Control.TLS.CertDir
	if caDir == "" {
		return "", fmt.Errorf("control tls cert dir is empty")
	}

	leafDir, err := safeJoin(m.agentCertsBaseDir(), node)
	if err != nil {
		return "", fmt.Errorf("invalid node name: %w", err)
	}

	if err := os.MkdirAll(leafDir, 0o700); err != nil {
		return "", fmt.Errorf("create agent cert dir: %w", err)
	}

	if err := transporttls.GenerateLeafCert(leafDir, caDir, node, []string{node}, nil); err != nil {
		return "", fmt.Errorf("generate agent cert: %w", err)
	}

	outPath, err := safeJoin(m.agentCertsBaseDir(), node+".tar.gz")
	if err != nil {
		return "", fmt.Errorf("invalid package path: %w", err)
	}

	if err := transporttls.CreateCertPackage(caDir, leafDir, outPath); err != nil {
		return "", fmt.Errorf("create cert package: %w", err)
	}

	return outPath, nil
}

func (m *CertManager) agentCertsBaseDir() string {
	base := m.config.AgentInstall.PackageDir
	if base == "" {
		base = defaultAgentInstallPackageDir
	}
	return filepath.Join(base, "certs")
}

func (m *CertManager) serverNames() (dnsNames, ipAddresses []string) {
	if base := m.config.AgentInstall.PublicBaseURL; base != "" {
		host := extractHost(base)
		if host != "" && host != "127.0.0.1" {
			dnsNames = append(dnsNames, host)
		}
	}
	return append(dnsNames, controlServerName()), []string{"127.0.0.1"}
}

func controlServerName() string {
	return "beacon/collia"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
