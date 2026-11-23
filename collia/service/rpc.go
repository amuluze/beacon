// Package service
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package service

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"

	rpctunnel "common/rpc/tunnel"
	transporttls "common/transport/tlsconfig"

	"collia/service/rpc"

	"github.com/amuluze/docker"
)

// Server manages the reverse tunnel connection to the Server.
type Server struct {
	tunnel *rpctunnel.AgentTunnel
	svc    *rpc.Service
}

// NewRPCServer creates the tunnel connection to the Server.
func NewRPCServer(config *Config, version Version) (*Server, error) {
	manager, err := docker.NewManager()
	if err != nil {
		return nil, err
	}
	rootDir := config.Variables.HostPrefix
	s := rpc.NewService(manager, rootDir)
	// 固定发布签名公钥（安装期由 Server config 通道写入 config.yml）；
	// 单键与轮换集合取并集，非法配置直接终止启动，不静默降级为兼容模式。
	upgradeKeys := make([]string, 0, 1+len(config.Control.UpgradePublicKeys))
	if config.Control.UpgradePublicKey != "" {
		upgradeKeys = append(upgradeKeys, config.Control.UpgradePublicKey)
	}
	upgradeKeys = append(upgradeKeys, config.Control.UpgradePublicKeys...)
	if err := s.SetUpgradePublicKeys(upgradeKeys); err != nil {
		return nil, fmt.Errorf("configure upgrade public key: %w", err)
	}
	// 接线自更新：二进制路径 + systemctl restart 回调；并清理上次更新残留的 .bak
	s.SetSelfUpdateConfig("", rpc.DefaultRestartFn("collia"))
	rpc.CleanupBackup("")
	agentID := config.Control.AgentID
	if agentID == "" {
		agentID = "default"
	}

	slog.Info("collia agent", "version", version.String(), "agent_id", agentID)

	tunnel := rpctunnel.NewAgentTunnel(config.Control.Server, agentID)
	tunnel.SetJoinToken(config.Control.JoinToken)
	tunnel.SetMetadata(version.String(), runtime.GOOS, runtime.GOARCH)
	if config.Control.TLS.Enable && !isUnixControlTarget(config.Control.Server) {
		serverName := config.Control.TLS.ServerName
		if serverName == "" {
			serverName = "beacon/collia"
		}
		tlsCfg, err := transporttls.ClientConfig(config.Control.TLS.CertDir, serverName)
		if err != nil {
			return nil, err
		}
		tunnel.SetTLSConfig(tlsCfg)
		slog.Info("reverse tunnel tls enabled", "cert_dir", config.Control.TLS.CertDir, "server_name", serverName)
	} else if isUnixControlTarget(config.Control.Server) {
		slog.Info("reverse tunnel over unix socket", "socket", config.Control.Server)
	}
	tunnel.SetHandler(buildRPCDispatcher(s))
	slog.Info("reverse tunnel configured", "server", config.Control.Server, "agent_id", agentID)

	return &Server{
		tunnel: tunnel,
		svc:    s,
	}, nil
}

// buildRPCDispatcher creates a dispatch handler that routes RPC calls
// to the methods on the Service struct based on method name.
func buildRPCDispatcher(svc *rpc.Service) rpctunnel.Handler {
	dispatch := rpc.NewDispatcher(svc)
	return func(ctx context.Context, method string, payload []byte, streamSender func(*rpctunnel.Frame)) ([]byte, error) {
		return dispatch.Call(ctx, method, payload, streamSender)
	}
}

func (s *Server) Start() error {
	slog.Info("starting reverse tunnel connection")
	return s.tunnel.Start(context.Background())
}

// DataPusher exposes the tunnel data stream so the monitoring task can push
// batches over the same connection as the control channel.
func (s *Server) DataPusher() DataPusher {
	return s.tunnel
}

func (s *Server) Stop() error {
	return s.tunnel.Close()
}
