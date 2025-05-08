// Package task
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package task

import (
	"time"

	"amprobe/pkg/psutil"
	"amprobe/service/model"
)

// cleanOldData 通用的数据清理函数，清理指定模型中超过保留期限的数据
func (a *Task) cleanOldData(model interface{}) error {
	ts := time.Now().Add(-time.Duration(a.maxAge) * 24 * time.Hour).Unix()
	if err := a.db.Model(model).Unscoped().Where("timestamp < ?", ts).Delete(model).Error; err != nil {
		return err
	}
	return nil
}

func (a *Task) CPUTask(timestamp time.Time) error {
	cpuPercent, _ := psutil.GetCPUPercent()
	if err := a.db.Model(&model.CPU{}).Create(&model.CPU{
		Timestamp:  timestamp,
		CPUPercent: cpuPercent,
	}).Error; err != nil {
		return err
	}

	// 清理旧数据
	return a.cleanOldData(&model.CPU{})
}

func (a *Task) MemoryTask(timestamp time.Time) error {
	memPercent, memTotal, memUsed, _ := psutil.GetMemInfo()
	if err := a.db.Model(&model.Memory{}).Create(&model.Memory{
		Timestamp:  timestamp,
		MemPercent: memPercent,
		MemTotal:   float64(memTotal),
		MemUsed:    float64(memUsed),
	}).Error; err != nil {
		return err
	}

	// 清理旧数据
	return a.cleanOldData(&model.Memory{})
}

/*
disk 函数用于获取磁盘指标，并将其存储到数据库中。
计算方法：两次采样间隔之间磁盘读写的平均速率
*/

func (a *Task) DiskTask(timestamp time.Time) error {
	diskInfo, _ := psutil.GetDiskInfo(a.devices)
	diskIOMap, _ := psutil.GetDiskIO(a.devices)
	var diskInfos []model.Disk

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
	if err := a.db.Model(&model.Disk{}).Create(diskInfos).Error; err != nil {
		return err
	}

	// 清理旧数据
	return a.cleanOldData(&model.Disk{})
}

/*
network 函数用于获取网络指标，并将其存储到数据库中。
计算方法：两次采样间隔之间发送、接收的平均速率
*/

func (a *Task) NetTask(timestamp time.Time) error {
	netMap, _ := psutil.GetNetworkIO(a.ethernet)
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
	if err := a.db.Model(&model.Net{}).Create(netInfos).Error; err != nil {
		return err
	}

	// 清理旧数据
	return a.cleanOldData(&model.Net{})
}
