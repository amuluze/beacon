// Package service
package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"collia/service/report"
	"collia/service/task"
	rpcSchema "common/rpc/schema"

	"github.com/amuluze/amutool/timex"
	"github.com/amuluze/docker"
)

type TimedTask struct {
	agentID      string
	task         task.ITask
	ticker       timex.Ticker
	stopCh       chan struct{}
	reportClient *report.Client
}

func NewTimedTask(conf *Config) *TimedTask {
	interval := conf.Task.Interval
	tk := timex.NewTicker(time.Duration(interval) * time.Second)
	manager, err := docker.NewManager()
	if err != nil {
		slog.Error("create docker manager failed", "error", err)
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

	newTask := task.NewTask(interval, manager, dev, eth)

	var rptClient *report.Client
	if conf.Task.Report.URL != "" {
		rptClient = report.NewClient(conf.Task.Report.URL, conf.Task.Report.Token)
	}

	return &TimedTask{
		agentID:      conf.Task.Report.AgentID,
		task:         newTask,
		ticker:       tk,
		stopCh:       make(chan struct{}),
		reportClient: rptClient,
	}
}

func (a *TimedTask) Execute() {
	timestamp := time.Now()
	r, err := a.task.Report(timestamp)
	if err != nil {
		slog.Error("report collection failed", "error", err)
		return
	}
	if a.reportClient != nil && r != nil {
		args := a.buildReportArgs(timestamp, r)
		if err := a.reportClient.Push(context.Background(), args); err != nil {
			slog.Error("push report failed", "error", err)
		}
	}
}

func (a *TimedTask) buildReportArgs(ts time.Time, r *task.MonitorReport) rpcSchema.MonitorReportArgs {
	args := rpcSchema.MonitorReportArgs{
		AgentID: a.agentID,
	}

	// Host
	if r.Host != nil {
		args.Host = rpcSchema.HostReport{
			Timestamp:       ts,
			Uptime:          r.Host.Uptime,
			Hostname:        r.Host.Hostname,
			Os:              r.Host.Os,
			Platform:        r.Host.Platform,
			PlatformVersion: r.Host.PlatformVersion,
			KernelArch:      r.Host.KernelArch,
			KernelVersion:   r.Host.KernelVersion,
		}
	}

	// CPU
	if r.CPU != nil {
		args.CPU = rpcSchema.CPUReport{Timestamp: ts, CPUPercent: r.CPU.CPUPercent}
	}

	// Memory
	if r.Memory != nil {
		args.Memory = rpcSchema.MemoryReport{
			Timestamp:  ts,
			MemPercent: r.Memory.MemPercent,
			MemTotal:   r.Memory.MemTotal,
			MemUsed:    r.Memory.MemUsed,
		}
	}

	// Disk
	for _, d := range r.Disks {
		args.Disks = append(args.Disks, rpcSchema.DiskReport{
			Timestamp:   ts,
			Device:      d.Device,
			DiskPercent: d.DiskPercent,
			DiskTotal:   d.DiskTotal,
			DiskUsed:    d.DiskUsed,
			DiskRead:    d.DiskRead,
			DiskWrite:   d.DiskWrite,
		})
	}

	// Net
	for _, n := range r.Nets {
		args.Nets = append(args.Nets, rpcSchema.NetReport{
			Timestamp: ts,
			Ethernet:  n.Ethernet,
			NetRecv:   n.NetRecv,
			NetSend:   n.NetSend,
		})
	}

	// Docker
	if r.Docker != nil {
		args.Docker = rpcSchema.DockerReport{
			Timestamp:     ts,
			DockerVersion: r.Docker.DockerVersion,
			APIVersion:    r.Docker.APIVersion,
			MinAPIVersion: r.Docker.MinAPIVersion,
			GitCommit:     r.Docker.GitCommit,
			GoVersion:     r.Docker.GoVersion,
			Os:            r.Docker.Os,
			Arch:          r.Docker.Arch,
		}
	}

	// Containers
	for _, c := range r.Containers {
		args.Containers = append(args.Containers, rpcSchema.ContainerReport{
			Timestamp:   ts,
			ContainerID: c.ContainerID,
			Name:        c.Name,
			Image:       c.Image,
			IP:          c.IP,
			Ports:       c.Ports,
			State:       c.State,
			Uptime:      c.Uptime,
			CPUPercent:  c.CPUPercent,
			MemPercent:  c.MemPercent,
			MemUsage:    c.MemUsage,
			MemLimit:    c.MemLimit,
			Labels:      c.Labels,
		})
	}

	// Images
	for _, im := range r.Images {
		args.Images = append(args.Images, rpcSchema.ImageReport{
			Timestamp: ts,
			ImageID:   im.ImageID,
			Name:      im.Name,
			Tag:       im.Tag,
			Created:   im.Created,
			Size:      im.Size,
			Number:    im.Number,
		})
	}

	// Networks
	for _, n := range r.Networks {
		args.Networks = append(args.Networks, rpcSchema.NetworkReport{
			Timestamp: ts,
			NetworkID: n.NetworkID,
			Name:      n.Name,
			Driver:    n.Driver,
			Scope:     n.Scope,
			Created:   n.Created,
			Internal:  n.Internal,
			Subnet:    n.Subnet,
			Gateway:   n.Gateway,
			Labels:    n.Labels,
		})
	}

	return args
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
	if a.reportClient != nil {
		a.reportClient.Close()
	}
}
