// Package service
// Date: 2024/06/11 19:29:11
// Author: Amu
// Description:
package service

import (
	"fmt"
	"log/slog"

	"beacon/pkg/rpc"
	"beacon/service/agent"
	tunnelpkg "common/rpc/tunnel"
	transporttls "common/transport/tlsconfig"
)

// globalTunnel guards the singleton tunnel instance.
// Set once during initialization; read safely by health probes and lifecycle hooks.
var globalTunnel *tunnelpkg.ServerTunnel

// NewRPCClient creates the tunnel-based RPC client.
// Server listens for reverse connections from Agents.
func NewRPCClient(config *Config, agentService *agent.Service, certManager *CertManager) (rpc.Caller, error) {
	addr := config.Control.Address
	if addr == "" {
		addr = ":17000"
	}
	slog.Info("starting reverse tunnel server", "addr", addr)

	if config.Control.TLS.Enable {
		if err := certManager.EnsureControlCerts(); err != nil {
			return nil, fmt.Errorf("ensure control certs: %w", err)
		}
	}

	tun := tunnelpkg.NewServerTunnel()
	tun.SetJoinToken(controlJoinToken(config))
	agentService.SetTunnel(tun)
	if config.Control.TLS.Enable {
		tlsCfg, err := transporttls.ServerConfig(config.Control.TLS.CertDir, config.Control.TLS.ClientNames)
		if err != nil {
			return nil, err
		}
		tun.SetTLSConfig(tlsCfg)
		slog.Info("reverse tunnel tls enabled", "cert_dir", config.Control.TLS.CertDir)
	}
	globalTunnel = tun

	go func() {
		if err := tun.Start(addr); err != nil {
			slog.Error("reverse tunnel server stopped", "err", err)
		}
	}()

	return rpc.NewTunnelClient(tun), nil
}

func controlJoinToken(config *Config) string {
	if config.Control.JoinToken != "" {
		return config.Control.JoinToken
	}
	return config.AgentInstall.Token
}

// ServerTunnelFromHolder returns the singleton tunnel instance held by globalTunnel.
// It exposes the live tunnel so other components (e.g. health.Probe) can wire
// into it without taking on the lifetime ownership.
func ServerTunnelFromHolder() *tunnelpkg.ServerTunnel {
	return globalTunnel
}
