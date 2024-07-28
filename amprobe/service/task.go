// Package service
// Date: 2024/3/6 13:34
// Author: Amu
// Description:
package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/amuluze/amprobe/pkg/psutil"
	"github.com/amuluze/amprobe/service/model"
	"github.com/amuluze/amutool/database"
	"github.com/amuluze/amutool/timex"
	"github.com/amuluze/docker"
)

const (
	LatestDiskReadKey   = "latest_disk_io_read_"
	LatestDisKWriteKey  = "latest_disk_io_write_"
	LatestNetReceiveKey = "latest_net_io_receive_"
	LatestNetSendKey    = "latest_net_io_send_"
)

type TimedTask struct {
	interval int
	db       *database.DB
	manager  *docker.Manager
	devices  map[string]struct{}
	ethernet map[string]struct{}
	ticker   timex.Ticker
	stopCh   chan struct{}
	cache    *cache.Cache
}

func NewTimedTask(conf *Config, db *database.DB) *TimedTask {
	interval := conf.Task.Interval
	tk := timex.NewTicker(time.Duration(interval) * time.Second)
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
		interval: interval,
		devices:  dev,
		ethernet: eth,
		ticker:   tk,
		stopCh:   make(chan struct{}),
		db:       db,
		manager:  manager,
		cache:    cache.New(5*time.Minute, 60*time.Second),
	}
}

func (a *TimedTask) Execute() {
	timestamp := time.Now()
	// 处理数组指标
	go a.host(timestamp)
	go a.cpu(timestamp)
	go a.memory(timestamp)
	go a.disk(timestamp)
	go a.network(timestamp)

	// // 处理 Docker 容器指标
	go a.container(timestamp)
	go func() {
		a.docker(timestamp)
		a.image(timestamp)
	}()

	go a.clearOldRecord()
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
	info, _ := psutil.GetSystemInfo()
	if err := a.db.Unscoped().Where("1 = 1").Delete(&model.Host{}).Error; err != nil {
		slog.Error("Error deleting host table", "error", err)
	}
	a.db.Model(&model.Host{}).Create(&model.Host{
		Timestamp:       timestamp,
		Uptime:          info.Uptime,
		Hostname:        info.Hostname,
		Os:              info.Os,
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		KernelArch:      info.KernelArch,
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

/*
disk 函数用于获取磁盘指标，并将其存储到数据库中。
计算方法：两次采样间隔之间磁盘读写的平均速率
*/
func (a *TimedTask) disk(timestamp time.Time) {
	diskInfo, _ := psutil.GetDiskInfo(a.devices)
	diskIOMap, _ := psutil.GetDiskIO(a.devices)
	var diskInfos []model.Disk
	slog.Info("disk info: ", "diskInfo", diskInfo)
	slog.Info("disk io map: ", "diskIOMap", diskIOMap)
	for device, info := range diskInfo {
		disk := model.Disk{Timestamp: timestamp}
		disk.Device = device
		disk.DiskPercent = info.Percent
		disk.DiskTotal = float64(info.Total)
		disk.DiskUsed = float64(info.Used)
		for dev, state := range diskIOMap {
			if dev == device {
				if latestReadBytes, ok := a.cache.Get(LatestDiskReadKey + device); ok {
					disk.DiskRead = float64((state.Read - latestReadBytes.(uint64)) / uint64(a.interval))
					a.cache.Set(LatestDiskReadKey+device, state.Read, 0)
				} else {
					a.cache.Set(LatestDiskReadKey+device, state.Read, 0)
					disk.DiskRead = 0
				}
				if latestWriteBytes, ok := a.cache.Get(LatestDisKWriteKey + device); ok {
					disk.DiskWrite = float64((state.Write - latestWriteBytes.(uint64)) / uint64(a.interval))
					a.cache.Set(LatestDisKWriteKey+device, state.Write, 0)
				} else {
					a.cache.Set(LatestDisKWriteKey+device, state.Write, 0)
					disk.DiskWrite = 0
				}
			}
		}
		diskInfos = append(diskInfos, disk)
	}
	a.db.Model(&model.Disk{}).Create(diskInfos)
}

/*
network 函数用于获取网络指标，并将其存储到数据库中。
计算方法：两次采样间隔之间发送、接收的平均速率
*/
func (a *TimedTask) network(timestamp time.Time) {
	netMap, _ := psutil.GetNetworkIO(a.ethernet)
	slog.Error("net map: ", "netMap", netMap)
	var netInfos []model.Net
	for eth, info := range netMap {
		net := model.Net{Timestamp: timestamp}
		net.Ethernet = eth
		if LatestNetReceiveBytes, ok := a.cache.Get(LatestNetReceiveKey + eth); ok {
			net.NetRecv = float64((info.Recv - LatestNetReceiveBytes.(uint64)) / uint64(a.interval))
			a.cache.Set(LatestNetReceiveKey+eth, info.Recv, 0)
		} else {
			a.cache.Set(LatestNetReceiveKey+eth, info.Recv, 0)
			net.NetRecv = 0
		}
		if LatestNetSendBytes, ok := a.cache.Get(LatestNetSendKey + eth); ok {
			net.NetSend = float64((info.Send - LatestNetSendBytes.(uint64)) / uint64(a.interval))
			a.cache.Set(LatestNetSendKey+eth, info.Send, 0)
		} else {
			a.cache.Set(LatestNetSendKey+eth, info.Send, 0)
			net.NetSend = 0
		}
		netInfos = append(netInfos, net)
	}
	a.db.Model(&model.Net{}).Create(netInfos)
}

func (a *TimedTask) container(timestamp time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cs, err := a.manager.ListContainer(ctx)
	if err != nil {
		slog.Error("failed to list containers", "error", err)
		return
	}
	var containers []model.Container
	for _, info := range cs {
		var d model.Container
		d.Timestamp = timestamp
		d.ContainerID = info.ID[:6]
		d.Name = info.Name
		d.State = info.State
		d.Image = info.Image
		d.Uptime = info.Uptime
		d.IP = info.IP

		cpuPercent, err := a.manager.GetContainerCpu(ctx, info.ID[:6])
		if err != nil {
			slog.Error("failed to get container cpu", "error", err)
		}
		d.CPUPercent = cpuPercent

		memPercent, used, limit, err := a.manager.GetContainerMem(ctx, info.ID[:6])
		if err != nil {
			slog.Error("failed to get container mem", "error", err)
		}
		d.MemPercent = memPercent

		d.MemUsage = used
		d.MemLimit = limit
		if _, ok := a.cache.Get(info.Image); !ok {
			a.cache.Set(info.Image, 1, 2*time.Minute)
		} else {
			count, err := a.cache.IncrementInt(info.Image, 1)
			slog.Info("container image cache", "image", info.Image, "count", count, "error", err)
		}
		containers = append(containers, d)
	}
	if err := a.db.Unscoped().Where("1 = 1").Delete(&model.Container{}).Error; err != nil {
		slog.Error("failed to delete container", "error", err)
	}
	a.db.Model(&model.Container{}).Create(&containers)
}

func (a *TimedTask) docker(timestamp time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	dockerVersion, err := a.manager.Version(ctx)
	if err != nil {
		slog.Error("failed to get docker version", "error", err)
		return
	}
	if err := a.db.Unscoped().Where("1 = 1").Delete(&model.Docker{}).Error; err != nil {
		slog.Error("failed to delete docker container", "error", err)
	}
	a.db.Model(&model.Docker{}).Create(&model.Docker{
		Timestamp:     timestamp,
		DockerVersion: dockerVersion.DockerVersion,
		APIVersion:    dockerVersion.APIVersion,
		MinAPIVersion: dockerVersion.MinAPIVersion,
		GitCommit:     dockerVersion.GitCommit,
		GoVersion:     dockerVersion.GoVersion,
		Os:            dockerVersion.OS,
		Arch:          dockerVersion.Arch,
	})
}

func (a *TimedTask) image(timestamp time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	images, err := a.manager.ListImage(ctx)
	if err != nil {
		slog.Error("failed to get version", "error", err)
		return
	}
	var list model.Images
	duplicateImage := make(map[string]struct{})
	for _, im := range images {
		val, ok := a.cache.Get(im.Name + ":" + im.Tag)
		if !ok {
			slog.Error("failed to get image cache", "error", err)
			val = 0
		}
		if _, ok := duplicateImage[im.ID]; !ok {
			duplicateImage[im.ID] = struct{}{}
		} else {
			if im.Tag != "latest" {
				continue
			}
		}
		list = append(list, model.Image{
			Timestamp: timestamp,
			ImageID:   im.ID[7:19],
			Name:      im.Name,
			Number:    val.(int),
			Tag:       im.Tag,
			Created:   im.Created,
			Size:      im.Size,
		})
		a.cache.Delete(im.Name + ":" + im.Tag)
	}
	if err := a.db.Unscoped().Where("1 = 1").Delete(&model.Image{}).Error; err != nil {
		slog.Error("failed to delete image", "error", err)
	}
	a.db.Model(&model.Image{}).Create(&list)
}

func (a *TimedTask) clearOldRecord() {
	a.db.Where("timestamp < ?", time.Now().Add(-time.Minute*5)).Delete(&model.Host{})
	a.db.Where("timestamp < ?", time.Now().Add(-time.Minute*5)).Delete(&model.Container{})
	a.db.Where("timestamp < ?", time.Now().Add(-time.Minute*5)).Delete(&model.Image{})
	a.db.Where("timestamp < ?", time.Now().Add(-time.Minute*5)).Delete(&model.Docker{})
	a.db.Where("timestamp < ?", time.Now().Add(-time.Hour*24*2)).Delete(&model.CPU{})
	a.db.Where("timestamp < ?", time.Now().Add(-time.Hour*24*2)).Delete(&model.Memory{})
	a.db.Where("timestamp < ?", time.Now().Add(-time.Hour*24*2)).Delete(&model.Disk{})
	a.db.Where("timestamp < ?", time.Now().Add(-time.Hour*24*2)).Delete(&model.Net{})
}
