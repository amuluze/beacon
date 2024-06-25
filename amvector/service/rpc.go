// Package service
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package service

import (
	"github.com/amuluze/amutool/database"
	"github.com/amuluze/amutool/docker"
	"github.com/amuluze/amvector/service/rpc"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/share"
	"net"
)

type Server struct {
	address string
	server  *server.Server
}

func NewRPCServer(config *Config, db *database.DB) (*Server, error) {
	srv := server.NewServer()
	p := server.NewFileTransfer("0.0.0.0:9527", uploadFileHandler, downloadFileHandler, 1000)
	srv.EnableFileTransfer(share.SendFileServiceName, p)
	manager, err := docker.NewManager()
	if err != nil {
		return nil, err
	}
	s := rpc.NewService(db, manager)

	err = srv.Register(s, "")
	if err != nil {
		return nil, err
	}

	return &Server{
		address: config.Rpc.Address,
		server:  srv,
	}, nil
}

func (s *Server) Start() error {
	return s.server.Serve("unix", s.address)
}

func (s *Server) Stop() error {
	return s.server.Close()
}

func uploadFileHandler(conn net.Conn, args *share.FileTransferArgs) {
	return
}

func downloadFileHandler(conn net.Conn, args *share.DownloadFileArgs) {
	return
}
