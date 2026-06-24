// Package task
package task

import (
	"collia/pkg/psutil"
)

func (a *Task) HostTask() *HostReport {
	info, _ := psutil.GetSystemInfo()
	return &HostReport{
		Uptime:          info.Uptime,
		Hostname:        info.Hostname,
		Os:              info.Os,
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		KernelArch:      info.KernelArch,
	}
}

func (a *Task) CPUTask() *CPUReport {
	cpuPercent, _ := psutil.GetCPUPercent()
	return &CPUReport{CPUPercent: cpuPercent}
}

func (a *Task) MemoryTask() *MemoryReport {
	memPercent, memTotal, memUsed, _ := psutil.GetMemInfo()
	return &MemoryReport{
		MemPercent: memPercent,
		MemTotal:   float64(memTotal),
		MemUsed:    float64(memUsed),
	}
}

func (a *Task) DiskTask() []*DiskReport {
	diskInfo, _ := psutil.GetDiskInfo(a.devices)
	diskIOMap, _ := psutil.GetDiskIO(a.devices)
	var reports []*DiskReport

	for device, info := range diskInfo {
		r := &DiskReport{
			Device:      device,
			DiskPercent: info.Percent,
			DiskTotal:   float64(info.Total),
			DiskUsed:    float64(info.Used),
		}
		for dev, state := range diskIOMap {
			if dev == device {
				if latestReadBytes, ok := a.cache.Get(LatestDiskReadKey + device); ok {
					r.DiskRead = float64((state.Read - latestReadBytes.(uint64)) / uint64(a.interval))
					a.cache.Set(LatestDiskReadKey+device, state.Read, 0)
				} else {
					a.cache.Set(LatestDiskReadKey+device, state.Read, 0)
					r.DiskRead = 0
				}
				if latestWriteBytes, ok := a.cache.Get(LatestDisKWriteKey + device); ok {
					r.DiskWrite = float64((state.Write - latestWriteBytes.(uint64)) / uint64(a.interval))
					a.cache.Set(LatestDisKWriteKey+device, state.Write, 0)
				} else {
					a.cache.Set(LatestDisKWriteKey+device, state.Write, 0)
					r.DiskWrite = 0
				}
			}
		}
		reports = append(reports, r)
	}
	return reports
}

func (a *Task) NetTask() []*NetReport {
	netMap, _ := psutil.GetNetworkIO(a.ethernet)
	var reports []*NetReport
	for eth, info := range netMap {
		r := &NetReport{Ethernet: eth}
		if LatestNetReceiveBytes, ok := a.cache.Get(LatestNetReceiveKey + eth); ok {
			r.NetRecv = float64((info.Recv - LatestNetReceiveBytes.(uint64)) / uint64(a.interval))
			a.cache.Set(LatestNetReceiveKey+eth, info.Recv, 0)
		} else {
			a.cache.Set(LatestNetReceiveKey+eth, info.Recv, 0)
			r.NetRecv = 0
		}
		if LatestNetSendBytes, ok := a.cache.Get(LatestNetSendKey + eth); ok {
			r.NetSend = float64((info.Send - LatestNetSendBytes.(uint64)) / uint64(a.interval))
			a.cache.Set(LatestNetSendKey+eth, info.Send, 0)
		} else {
			a.cache.Set(LatestNetSendKey+eth, info.Send, 0)
			r.NetSend = 0
		}
		reports = append(reports, r)
	}
	return reports
}
