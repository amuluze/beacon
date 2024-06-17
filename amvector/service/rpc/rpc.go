// Package rpc
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package rpc

import (
	"github.com/amuluze/amutool/database"
	"github.com/amuluze/amutool/docker"
)

type Service struct {
	DB      *database.DB
	Manager *docker.Manager
}

func NewService(db *database.DB, manager *docker.Manager) *Service {
	return &Service{DB: db, Manager: manager}
}
