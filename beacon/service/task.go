// Package service
// Date:   2024/10/14 17:14
// Author: Amu
// Description:
package service

import (
	"beacon/pkg/rpc"
	"beacon/service/task"
	"common/database"
	"context"
	"log/slog"
	"time"

	"github.com/amuluze/amutool/timex"
)

type TimedTask struct {
	task          task.ITask
	cleanup       *task.CleanupTask
	ticker        timex.Ticker
	cleanupTicker timex.Ticker
	stopCh        chan struct{}
}

func NewTimedTask(conf *Config, cli rpc.Caller, db *database.DB) *TimedTask {
	// 默认配置： 每 15s 执行一次
	interval := conf.Task.Interval
	tk := timex.NewTicker(time.Duration(interval) * time.Second)

	// Cleanup runs once per hour (3600 seconds)
	cleanupTk := timex.NewTicker(3600 * time.Second)

	newTask := task.NewTask(db)
	cleanupTask := task.NewCleanupTask(db, conf.Retention.Days)

	return &TimedTask{
		task:          newTask,
		cleanup:       cleanupTask,
		ticker:        tk,
		cleanupTicker: cleanupTk,
		stopCh:        make(chan struct{}),
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

func (a *TimedTask) ExecuteCleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := a.cleanup.Run(ctx); err != nil {
		slog.Error("cleanup task failed", "error", err)
	}
}

func (a *TimedTask) Run() {
	for {
		select {
		case <-a.ticker.Chan():
			go a.Execute()
		case <-a.cleanupTicker.Chan():
			go a.ExecuteCleanup()
		case <-a.stopCh:
			return
		}
	}
}

func (a *TimedTask) Stop() {
	close(a.stopCh)
}
