// Package service
// Date: 2024/06/10 18:30:23
// Author: Amu
// Description:
package service

import (
	"common/database"
	"fmt"
	"log/slog"
	"time"

	"amvector/service/task"

	"github.com/amuluze/amutool/timex"
	"github.com/amuluze/docker"
)

type TimedTask struct {
	task   task.ITask
	ticker timex.Ticker
	stopCh chan struct{}
}

func NewTimedTask(conf *Config, db *database.DB) *TimedTask {
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

	newTask := task.NewTask(interval, conf.Task.MaxAge, db, manager, dev, eth)

	return &TimedTask{
		task:   newTask,
		ticker: tk,
		stopCh: make(chan struct{}),
	}
}

func (a *TimedTask) Execute() {
	timestamp := time.Now()
	// 处理宿主机指标
	if err := a.task.HostTask(timestamp); err != nil {
		slog.Error("host summary failed: ", "error", err)
	}
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

	// 处理 Docker 容器指标
	if err := a.task.DockerTask(timestamp); err != nil {
		slog.Error("docker summary failed", "error", err)
	}
	if err := a.task.ContainerTask(timestamp); err != nil {
		slog.Error("containers summary failed", "error", err)
	}
	if err := a.task.ImageTask(timestamp); err != nil {
		slog.Error("image summary failed", "error", err)
	}
	if err := a.task.NetworkTask(timestamp); err != nil {
		slog.Error("network summary failed", "error", err)
	}

	a.task.CleanTask()
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
	close(a.stopCh)
}
