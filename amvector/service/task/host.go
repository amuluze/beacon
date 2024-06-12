// Package task
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package task

import (
	"github.com/amuluze/amvector/pkg/psutil"
	"github.com/amuluze/amvector/service/model"
	"log/slog"
	"time"
)

func (a *Task) Host(timestamp time.Time) {
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

func (a *Task) Cpu(timestamp time.Time) {
	cpuPercent, _ := psutil.GetCPUPercent()
	a.db.Model(&model.CPU{}).Create(&model.CPU{
		Timestamp:  timestamp,
		CPUPercent: cpuPercent,
	})
}

func (a *Task) Memory(timestamp time.Time) {
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
func (a *Task) Disk(timestamp time.Time) {
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
func (a *Task) Network(timestamp time.Time) {
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
