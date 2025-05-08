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

const (
	LatestDiskReadKey   = "latest_disk_io_read_"
	LatestDisKWriteKey  = "latest_disk_io_write_"
	LatestNetReceiveKey = "latest_net_io_receive_"
	LatestNetSendKey    = "latest_net_io_send_"
)

var _ ITask = (*Task)(nil)

type ITask interface {
	DockerTask(timestamp time.Time) error
	ContainerTask(timestamp time.Time) error
	ImageTask(timestamp time.Time) error
	NetworkTask(timestamp time.Time) error

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

func NewTask(interval int, maxAge int, db *database.DB, manager *docker.Manager, dev map[string]struct{}, eth map[string]struct{}) *Task {
	return &Task{
		interval: interval,
		maxAge:   maxAge,
		db:       db,
		manager:  manager,
		devices:  dev,
		ethernet: eth,
		cache:    cache.New(5*time.Minute, 60*time.Second),
	}
}
