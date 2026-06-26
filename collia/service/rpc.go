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
	"google.golang.org/grpc/credentials"
)

// Server manages the reverse tunnel connection to the Server.
type Server struct {
	tunnel *rpctunnel.AgentTunnel
}

// NewRPCServer creates the tunnel connection to the Server.
func NewRPCServer(config *Config, db *database.DB) (*Server, error) {
	manager, err := docker.NewManager()
	if err != nil {
		return nil, err
	}
	s := rpc.NewService(db, manager)

	agentID := config.Control.AgentID
	if agentID == "" {
		agentID = "default"
	}

	tunnel := rpctunnel.NewAgentTunnel(config.Control.Server, agentID)
	if config.Control.TLS.Enable {
		tlsCfg, err := transporttls.ClientConfig(config.Control.TLS.CertDir, "amprobe/collia")
		if err != nil {
			return nil, err
		}
		tunnel = rpctunnel.NewAgentTunnel(config.Control.Server, agentID, rpctunnel.WithAgentTLS(credentials.NewTLS(tlsCfg)))
	}
	tunnel.SetHandler(buildRPCDispatcher(s))
	slog.Info("reverse tunnel configured", "server", config.Control.Server, "agent_id", agentID, "tls", config.Control.TLS.Enable)

	return &Server{
		tunnel: tunnel,
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
