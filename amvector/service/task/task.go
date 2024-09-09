// Package task
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package task

import (
	"time"

	"github.com/amuluze/amutool/database"
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
	DockerSummary(timestamp time.Time) error
	ContainerSummary(timestamp time.Time) error
	ImageSummary(timestamp time.Time) error
	NetworkSummary(timestamp time.Time) error

	HostSummary(timestamp time.Time) error
	CPUSummary(timestamp time.Time) error
	MemorySummary(timestamp time.Time) error
	DiskSummary(timestamp time.Time) error
	NetSummary(timestamp time.Time) error
}

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
