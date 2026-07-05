// Package conn
// Date: 2024/5/16 15:43
// Author: Amu
// Description:
package conn

import (
	"crypto/tls"

	transporttls "common/transport/tlsconfig"
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
	cfg, err := transporttls.ClientConfig(absDir, "beacon/collia")
	if err != nil {
		return nil, err
	}
	return &TLSConfig{Config: cfg, ServerAddresses: []string{}}, nil
}

// ServerConfig returns tls config for server
func ServerConfig(absDir string) (*TLSConfig, error) {
	cfg, err := transporttls.ServerConfig(absDir, nil)
	if err != nil {
		return nil, err
	}
	return &TLSConfig{Config: cfg, ServerAddresses: []string{}}, nil
}
