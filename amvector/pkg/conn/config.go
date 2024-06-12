// Package conn
// Date: 2024/5/16 15:43
// Author: Amu
// Description:
package conn

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
)

const (
	caCert  = "ca.pem"
	tlsCert = "tls.crt"
	tlsKey  = "tls.key"
)

// TLSConfig includes tls config and server info
type TLSConfig struct {
	*tls.Config
	ServerAddresses []string
}

// ClientConfig returns tls config for client
func ClientConfig(absDir string) (*TLSConfig, error) {
	return config(absDir, true)
}

// ServerConfig returns tls config for server
func ServerConfig(absDir string) (*TLSConfig, error) {
	return config(absDir, false)
}

// Config returns TLS config using local files
func config(absDir string, isClient bool) (*TLSConfig, error) {
	tlsCert, err := os.ReadFile(filepath.Join(absDir, tlsCert))
	if err != nil {
		return nil, err
	}
	tlsCertBlock, _ := pem.Decode(tlsCert)

	// parse private key
	tlsKey, err := os.ReadFile(filepath.Join(absDir, tlsKey))
	if err != nil {
		return nil, err
	}
	tlsKeyBlock, _ := pem.Decode(tlsKey)
	key, err := x509.ParsePKCS1PrivateKey(tlsKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	// generate tls certificate
	cert := tls.Certificate{
		Certificate: [][]byte{tlsCertBlock.Bytes},
		PrivateKey:  key,
	}

	// build ca pool
	caPool := x509.NewCertPool()
	caPem, err := os.ReadFile(filepath.Join(absDir, caCert))
	if err != nil {
		return nil, err
	}
	caPemBlock, _ := pem.Decode(caPem)
	ca, err := x509.ParseCertificate(caPemBlock.Bytes)
	if err != nil {
		return nil, err
	}
	caPool.AddCert(ca)

	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		Rand:         rand.Reader,
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		},
	}

	if isClient {
		cfg.InsecureSkipVerify = false
		cfg.RootCAs = caPool
		cfg.ServerName = "amvector"
	} else {
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
		cfg.ClientCAs = caPool
	}

	tlsConfig := &TLSConfig{
		Config:          cfg,
		ServerAddresses: []string{},
	}
	return tlsConfig, nil
}
