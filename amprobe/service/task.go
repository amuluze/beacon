// Package service
// Date:   2024/10/14 17:14
// Author: Amu
// Description:
package service

import (
	"amprobe/pkg/rpc"
	"amprobe/service/task"
	"common/database"
	"context"
	"fmt"
	"github.com/amuluze/amutool/timex"
	"log/slog"
	"time"
)

type TimedTask struct {
	task   task.ITask
	ticker timex.Ticker
	stopCh chan struct{}
}

func NewTimedTask(conf *Config, cli *rpc.Client, db *database.DB) *TimedTask {
	// 默认配置： 每 15s 执行一次
	interval := conf.Task.Interval
	tk := timex.NewTicker(time.Duration(interval) * time.Second)

	newTask := task.NewTask(cli, db)

	return &TimedTask{
		task:   newTask,
		ticker: tk,
		stopCh: make(chan struct{}),
	}
}

func (a *TimedTask) Execute() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.task.CPUAlarmTask(ctx); err != nil {
		slog.Error("cpu alarm task failed", "error", err)
	}
	if err := a.task.MemoryAlarmTask(ctx); err != nil {
		slog.Error("memory alarm task failed", "error", err)
	}
	if err := a.task.DiskAlarmTask(ctx); err != nil {
		slog.Error("disk alarm task failed", "error", err)
	}
	if err := a.task.ServiceTask(ctx); err != nil {
		slog.Error("service task failed", "error", err)
	}
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
