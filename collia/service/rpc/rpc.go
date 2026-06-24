// Package rpc
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package rpc

import (
	"common/database"
	"common/rpc"

	"github.com/amuluze/docker"
)

var _ rpc.IService = (*Service)(nil)

type Service struct {
	DB      *database.DB
	Manager *docker.Manager
}

func NewService(db *database.DB, manager *docker.Manager) *Service {
	return &Service{DB: db, Manager: manager}
}
