// Package service
// Date: 2024/06/11 19:29:11
// Author: Amu
// Description:
package service

import (
	"fmt"
	"log/slog"

	"amprobe/pkg/rpc"
	tunnelpkg "common/rpc/tunnel"
	transporttls "common/transport/tlsconfig"
)

// TunnelResult encapsulates both the RPC caller and the underlying tunnel instance,
// enabling Wire to extract each as a separate dependency.
type TunnelResult struct {
	Caller rpc.Caller
	Tunnel *tunnelpkg.ServerTunnel
}

// NewRPCClient creates the tunnel-based RPC client.
// Server listens for reverse connections from Agents.
func NewRPCClient(config *Config) (*TunnelResult, error) {
	addr := config.Control.Address
	if addr == "" {
		addr = ":17000"
	}

	// 控制通道承载远程 shell 等高危调用，agent 注册必须强鉴权。
	// 不允许 JoinToken 为空的默认开放模式；未配置则拒绝启动。
	// 可通过环境变量 AMPROBE_CONTROL_JOINTOKEN 注入。
	if config.Control.Enable && config.Control.JoinToken == "" {
		return nil, fmt.Errorf("control JoinToken is not configured; set Control.JoinToken or AMPROBE_CONTROL_JOINTOKEN before enabling the reverse tunnel")
	}

	slog.Info("starting reverse tunnel server", "addr", addr, "auth", true)

	var opts []tunnelpkg.ServerOption
	opts = append(opts, tunnelpkg.WithJoinToken(config.Control.JoinToken))
	if config.Control.TLSEnable {
		cfg, err := transporttls.ServerConfig(config.Control.TLSCertDir, nil)
		if err != nil {
			return nil, err
		}
		opts = append(opts, tunnelpkg.WithServerTLS(cfg))
	}

	tun := tunnelpkg.NewServerTunnel(opts...)

	go func() {
		if err := tun.Start(addr); err != nil {
			slog.Error("reverse tunnel server stopped", "err", err)
		}
	}()

	defaultID := config.Control.DefaultAgentID
	if defaultID == "" {
		defaultID = rpc.DefaultAgentID
	}

	return &TunnelResult{
		Caller: rpc.NewTunnelClient(tun, defaultID),
		Tunnel: tun,
	}, nil
}

// NewRPCCaller extracts the RPC Caller from TunnelResult for Wire injection.
func NewRPCCaller(result *TunnelResult) rpc.Caller {
	return result.Caller
}

// NewServerTunnelFromResult extracts the ServerTunnel from TunnelResult for Wire injection.
func NewServerTunnelFromResult(result *TunnelResult) *tunnelpkg.ServerTunnel {
	return result.Tunnel
}
