// Package rpc
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package rpc

import (
	"common/database"
	"common/rpc"
	"github.com/amuluze/docker"
	"github.com/patrickmn/go-cache"
)

var _ rpc.IService = (*Service)(nil)

type Service struct {
	DB      *database.DB
	Manager *docker.Manager
	cache   *cache.Cache
}

func NewService(db *database.DB, manager *docker.Manager) *Service {
	return &Service{DB: db, Manager: manager, cache: cache.New(0, 0)}
}
