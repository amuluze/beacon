// Package service
// Date: 2024/3/6 13:34
// Author: Amu
// Description:
package service

import (
	"fmt"
	"github.com/amuluze/amprobe/pkg/database"
	"github.com/amuluze/amprobe/pkg/docker"
	"github.com/amuluze/amprobe/pkg/psutil"
	"github.com/amuluze/amprobe/pkg/ticker"
	"github.com/amuluze/amprobe/service/model"
	"log/slog"
	"time"
)

type TimedTask struct {
	db       *database.DB
	manager  *docker.Manager
	devices  map[string]struct{}
	ethernet map[string]struct{}
	ticker   ticker.Ticker
	stopCh   chan struct{}
}

func NewTimedTask(conf *Config, db *database.DB) *TimedTask {
	interval := conf.Task.Interval
	tk := ticker.NewTicker(time.Duration(interval) * time.Second)
	manager, err := docker.NewManager()
	if err != nil {
		return nil
	}

	dev := make(map[string]struct{})
	for _, d := range conf.Disk.Devices {
		dev[d] = struct{}{}
	}

	eth := make(map[string]struct{})
	for _, d := range conf.Ethernet.Names {
		eth[d] = struct{}{}
	}
	return &TimedTask{
		devices:  dev,
		ethernet: eth,
		ticker:   tk,
		stopCh:   make(chan struct{}),
		db:       db,
		manager:  manager,
	}
}

func (a *TimedTask) Execute() {
	timestamp := time.Now()
	slog.Debug("TimedTask execute debug")
	slog.Info("TimedTask execute info")
	slog.Error("TimedTask execute error")
	// 处理数组指标
	go a.host(timestamp)
	go a.cpu(timestamp)
	go a.memory(timestamp)
	go a.disk(timestamp)
	go a.network(timestamp)

	// 处理 Docker 容器指标
	go a.docker(timestamp)
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

func (a *TimedTask) host(timestamp time.Time) {
	uptime, _ := psutil.GetSystemUpTime()
	a.db.Model(&model.Host{}).Create(&model.Host{
		Timestamp: timestamp,
		Uptime:    uptime,
	})
}

func (a *TimedTask) cpu(timestamp time.Time) {
	cpuPercent, _ := psutil.GetCPUPercent()
	a.db.Model(&model.CPU{}).Create(&model.CPU{
		Timestamp:  timestamp,
		CPUPercent: cpuPercent,
	})
}

func (a *TimedTask) memory(timestamp time.Time) {
	memPercent, memTotal, memUsed, _ := psutil.GetMemInfo()
	a.db.Model(&model.Memory{}).Create(&model.Memory{
		Timestamp:  timestamp,
		MemPercent: memPercent,
		MemTotal:   float64(memTotal),
		MemUsed:    float64(memUsed),
	})
}

func (a *TimedTask) disk(timestamp time.Time) {
	diskInfo, _ := psutil.GetDiskInfo(a.devices)
	diskMap, _ := psutil.GetDiskIO(a.devices)
	for device, info := range diskInfo {
		disk := &model.Disk{Timestamp: timestamp}
		disk.Device = device
		disk.DiskPercent = info.Percent
		disk.DiskTotal = float64(info.Total)
		disk.DiskUsed = float64(info.Used)
		for dev, state := range diskMap {
			if dev == device {
				disk.DiskRead = float64(state.Read)
				disk.DiskWrite = float64(state.Write)
			}
		}
		a.db.Model(&model.Disk{}).Create(disk)
	}
}

func (a *TimedTask) network(timestamp time.Time) {
	netMap, _ := psutil.GetNetworkIO(a.ethernet)
	time.Sleep(1 * time.Second)
	netMapAfterSecond, _ := psutil.GetNetworkIO(a.ethernet)
	for eth, info := range netMap {
		for e, i := range netMapAfterSecond {
			if eth == e {
				net := &model.Net{Timestamp: timestamp}
				net.Ethernet = eth
				net.NetSend = float64(i.Send - info.Send)
				net.NetRecv = float64(i.Recv - info.Recv)
				a.db.Model(&model.Net{}).Create(net)
			}
		}
	}
}

func (a *TimedTask) docker(timestamp time.Time) {
	infos, err := a.manager.GetContainerInfos()
	if err != nil {
		return
	}
	var dockers []model.Container
	for _, info := range infos {
		var d model.Container
		d.Timestamp = timestamp
		d.ContainerID = info.ID
		d.Name = info.Name
		d.State = info.State
		d.Image = info.Image
		d.IP = info.IP
		d.Uptime = info.Uptime
		d.CPUPercent = info.CPUPercent
		d.MemPercent = info.MemPercent
		d.MemUsage = info.MemUsage
		d.MemLimit = info.MemLimit
		dockers = append(dockers, d)
	}
	a.db.Model(&model.Container{}).Create(&dockers)
}

func (a *TimedTask) clearOldRecord() {
	a.db.Where("timestamp < ?", time.Now().Add(-time.Minute*5)).Delete(model.Host{})
	a.db.Where("timestamp < ?", time.Now().Add(-time.Minute*5)).Delete(model.Container{})
	//a.db.Where("timestamp < ?", time.Now().Add(-time.Hour*24*2)).Delete(model.Host{})
	//a.db.Where("timestamp < ?", time.Now().Add(-time.Hour*24*2)).Delete(model.Docker{})
}
