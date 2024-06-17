// Package task
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package task

import (
	"github.com/amuluze/amutool/database"
	"github.com/amuluze/amutool/docker"
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	LatestDiskReadKey   = "latest_disk_io_read_"
	LatestDisKWriteKey  = "latest_disk_io_write_"
	LatestNetReceiveKey = "latest_net_io_receive_"
	LatestNetSendKey    = "latest_net_io_send_"
)

type Task struct {
	interval int
	db       *database.DB
	manager  *docker.Manager
	devices  map[string]struct{}
	ethernet map[string]struct{}
	cache    *cache.Cache
}

func NewTask(interval int, db *database.DB, manager *docker.Manager, dev map[string]struct{}, eth map[string]struct{}) *Task {
	return &Task{
		interval: interval,
		db:       db,
		manager:  manager,
		devices:  dev,
		ethernet: eth,
		cache:    cache.New(5*time.Minute, 60*time.Second),
	}
}
