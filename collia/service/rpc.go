// Package service
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package service

import (
	"collia/pkg/resources"
	"common/database"
	transporttls "common/transport/tlsconfig"
	"path/filepath"

	"collia/service/rpc"

	"github.com/amuluze/docker"
	"github.com/smallnest/rpcx/server"
)

type Server struct {
	network string
	address string
	server  *server.Server
}

func NewRPCServer(config *Config, db *database.DB) (*Server, error) {
	srv := server.NewServer()
	if config.Rpc.TLS.Enable {
		tlsConfig, err := transporttls.ServerConfig(config.Rpc.TLS.CertDir, config.Rpc.TLS.ClientNames)
		if err != nil {
			return nil, err
		}
		srv = server.NewServer(server.WithTLSConfig(tlsConfig))
	}
	manager, err := docker.NewManager()
	if err != nil {
		return nil, err
	}
	s := rpc.NewService(db, manager)

	err = srv.Register(s, "")
	if err != nil {
		return nil, err
	}

	network := config.Rpc.Network
	address := config.Rpc.Address
	if network == "" {
		network = "unix"
	}
	if address == "" {
		address = filepath.Join(string(config.prefix), resources.RootPath, resources.ColliaSockFile)
	}

	return &Server{
		network: network,
		address: address,
		server:  srv,
	}, nil
}

func (s *Server) Start() error {
	return s.server.Serve(s.network, s.address)
}

func (s *Server) Stop() error {
	return s.server.Close()
}
