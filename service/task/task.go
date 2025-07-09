// Package task
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package task

import (
	"amprobe/pkg/database"
	"time"

	"github.com/amuluze/docker"
	"github.com/patrickmn/go-cache"
)

var _ ITask = (*Task)(nil)

type ITask interface {
	CPUTask(timestamp time.Time) error
	MemoryTask(timestamp time.Time) error
	DiskTask(timestamp time.Time) error
	NetTask(timestamp time.Time) error
}

type Task struct {
	interval int
	maxAge   int
	db       *database.DB
	manager  *docker.Manager
	devices  map[string]struct{}
	ethernet map[string]struct{}
	cache    *cache.Cache
}

func NewTask(interval int, maxAge int, db *database.DB, manager *docker.Manager, localCache *cache.Cache, dev map[string]struct{}, eth map[string]struct{}) *Task {
	return &Task{
		interval: interval,
		maxAge:   maxAge,
		db:       db,
		manager:  manager,
		devices:  dev,
		ethernet: eth,
		cache:    localCache,
	}
}
