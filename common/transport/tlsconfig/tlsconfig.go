// Package tlsconfig builds mutual TLS configs shared by beacon and collia.
package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"path/filepath"
)

const (
	CACertFile  = "ca.pem"
	TLSCertFile = "tls.crt"
	TLSKeyFile  = "tls.key"
)

// ClientConfig loads a client certificate and verifies the remote server name.
func ClientConfig(certDir string, serverName string) (*tls.Config, error) {
	if serverName == "" {
		return nil, fmt.Errorf("empty tls server name")
	}

	cert, caPool, err := loadMaterial(certDir)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		ServerName:   serverName,
	}, nil
}

// ServerConfig requires and verifies client certificates. If allowedClientNames
// is not empty, the client certificate must contain one of those identities.
func ServerConfig(certDir string, allowedClientNames []string) (*tls.Config, error) {
	cert, caPool, err := loadMaterial(certDir)
	if err != nil {
		return nil, err
	}

	cfg := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caPool,
	}

	if len(allowedClientNames) > 0 {
		cfg.VerifyConnection = func(state tls.ConnectionState) error {
			if len(state.PeerCertificates) == 0 {
				return fmt.Errorf("missing client certificate")
			}
			if !certificateMatchesAnyName(state.PeerCertificates[0], allowedClientNames) {
				return fmt.Errorf("client certificate identity is not allowed")
			}
			return nil
		}
	}

	return cfg, nil
}

func loadMaterial(certDir string) (tls.Certificate, *x509.CertPool, error) {
	if certDir == "" {
		return tls.Certificate{}, nil, fmt.Errorf("empty tls cert dir")
	}

	cert, err := tls.LoadX509KeyPair(
		filepath.Join(certDir, TLSCertFile),
		filepath.Join(certDir, TLSKeyFile),
	)
	if err != nil {
		return tls.Certificate{}, nil, err
	}

	caPEM, err := os.ReadFile(filepath.Join(certDir, CACertFile))
	if err != nil {
		return tls.Certificate{}, nil, err
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caPEM) {
		return tls.Certificate{}, nil, fmt.Errorf("load ca cert from %s failed", filepath.Join(certDir, CACertFile))
	}

	return cert, caPool, nil
}

func certificateMatchesAnyName(cert *x509.Certificate, names []string) bool {
	for _, name := range names {
		if name == "" {
			continue
		}
		if certificateMatchesName(cert, name) {
			return true
		}
	}
	return false
}

func certificateMatchesName(cert *x509.Certificate, name string) bool {
	if cert.Subject.CommonName == name {
		return true
	}
	for _, dnsName := range cert.DNSNames {
		if dnsName == name {
			return true
		}
	}
	for _, ip := range cert.IPAddresses {
		if ip.String() == name {
			return true
		}
	}
	for _, uri := range cert.URIs {
		if uri.String() == name {
			return true
		}
	}
	if net.ParseIP(name) != nil {
		for _, ip := range cert.IPAddresses {
			if ip.Equal(net.ParseIP(name)) {
				return true
			}
		}
	}
	return cert.VerifyHostname(name) == nil
}
