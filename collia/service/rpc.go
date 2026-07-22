// Package service
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package service

import (
	"context"
	"log/slog"

	"common/database"
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
func NewRPCServer(config *Config, db *database.DB, version Version) (*Server, error) {
	manager, err := docker.NewManager()
	if err != nil {
		return nil, err
	}
	rootDir := config.Variables.HostPrefix
	s := rpc.NewService(db, manager, rootDir)
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
	if config.Control.TLS.Enable {
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

func (s *Server) Stop() error {
	return s.tunnel.Close()
}
