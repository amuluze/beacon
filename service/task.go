// Package service
// Date:   2024/10/14 17:14
// Author: Amu
// Description:
package service

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/patrickmn/go-cache"

	"amprobe/pkg/database"
	"amprobe/pkg/timex"
	"amprobe/service/task"

	"github.com/amuluze/docker"
)

type TimedTask struct {
	task   task.ITask
	ticker timex.Ticker
	stopCh chan struct{}
}

func NewTimedTask(conf *Config, db *database.DB, localCache *cache.Cache) *TimedTask {
	interval := conf.Task.Interval
	tk := timex.NewTicker(time.Duration(interval) * time.Second)
	manager, err := docker.NewManager()
	if err != nil {
		return nil
	}

	dev := make(map[string]struct{})
	for _, d := range conf.Task.Disk.Devices {
		dev[d] = struct{}{}
	}

	eth := make(map[string]struct{})
	for _, d := range conf.Task.Ethernet.Names {
		eth[d] = struct{}{}
	}
	slog.Info("task init success", "devices", dev, "ethernet", eth, "interval", interval)
	newTask := task.NewTask(interval, conf.Task.MaxAge, db, manager, localCache, dev, eth)

	return &TimedTask{
		task:   newTask,
		ticker: tk,
		stopCh: make(chan struct{}),
	}
}

func (a *TimedTask) Execute() {
	timestamp := time.Now()
	slog.Info("task execute...", "timestamp", time.Now().Unix())
	// 处理宿主机指标
	if err := a.task.CPUTask(timestamp); err != nil {
		slog.Error("cpu summary failed: ", "error", err)
	}
	if err := a.task.MemoryTask(timestamp); err != nil {
		slog.Error("memory summary failed: ", "error", err)
	}
	if err := a.task.DiskTask(timestamp); err != nil {
		slog.Error("disk summary failed: ", "error", err)
	}
	if err := a.task.NetTask(timestamp); err != nil {
		slog.Error("net summary failed: ", "error", err)
	}
	slog.Info("task execute success", "timestamp", time.Now().Unix())
}

func (a *TimedTask) Run() {
	for {
		select {
		case <-a.ticker.Chan():
			go a.Execute()
		case <-a.stopCh:
			fmt.Println("task exit")
			return
		}
	}
}

func (a *TimedTask) Stop() {
	slog.Info("task stop...")
	close(a.stopCh)
}
