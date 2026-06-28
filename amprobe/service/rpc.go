// Package service
// Date: 2024/06/11 19:29:11
// Author: Amu
// Description:
package service

import (
	"log/slog"

	"amprobe/pkg/rpc"
	"amprobe/service/agent"
	tunnelpkg "common/rpc/tunnel"
)

// serverTunnelHolder holds the tunnel instance so other components can wire into it.
type serverTunnelHolder struct {
	tun *tunnelpkg.ServerTunnel
}

var globalTunnel serverTunnelHolder

// NewRPCClient creates the tunnel-based RPC client.
// Server listens for reverse connections from Agents.
func NewRPCClient(config *Config) (rpc.Caller, error) {
	addr := config.Control.Address
	if addr == "" {
		addr = ":17000"
	}
	slog.Info("starting reverse tunnel server", "addr", addr)

	tun := tunnelpkg.NewServerTunnel()
	tun.SetJoinToken(controlJoinToken(config))
	globalTunnel.tun = tun

	go func() {
		if err := tun.Start(addr); err != nil {
			slog.Error("reverse tunnel server stopped", "err", err)
		}
	}()

	defaultID := config.Control.DefaultAgentID
	if defaultID == "" {
		defaultID = rpc.DefaultAgentID
	}
	return rpc.NewTunnelClient(tun, defaultID), nil
}

func controlJoinToken(config *Config) string {
	if config.Control.JoinToken != "" {
		return config.Control.JoinToken
	}
	return config.AgentInstall.Token
}

// SetAgentLifecycle wires the agent service into the tunnel lifecycle hooks.
func SetAgentLifecycle(svc *agent.Service) {
	if globalTunnel.tun != nil {
		svc.SetTunnel(globalTunnel.tun)
	}
}
