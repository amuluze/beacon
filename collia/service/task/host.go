// Package task
package task

import (
	"errors"
	"fmt"

	"collia/pkg/psutil"
)

func (a *Task) HostTask() (*HostReport, error) {
	info, err := psutil.GetSystemInfo()
	if err != nil {
		return nil, err
	}
	return &HostReport{
		Uptime:          info.Uptime,
		Hostname:        info.Hostname,
		Os:              info.Os,
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		KernelArch:      info.KernelArch,
	}, nil
}

func (a *Task) CPUTask() (*CPUReport, error) {
	cpuPercent, err := psutil.GetCPUPercent()
	if err != nil {
		return nil, err
	}
	return &CPUReport{CPUPercent: cpuPercent}, nil
}

func (a *Task) MemoryTask() (*MemoryReport, error) {
	memPercent, memTotal, memUsed, err := psutil.GetMemInfo()
	if err != nil {
		return nil, err
	}
	return &MemoryReport{
		MemPercent: memPercent,
		MemTotal:   float64(memTotal),
		MemUsed:    float64(memUsed),
	}, nil
}

func (a *Task) DiskTask() ([]*DiskReport, error) {
	diskInfo, infoErr := psutil.GetDiskInfo(a.devices)
	diskIOMap, ioErr := psutil.GetDiskIO(a.devices)
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
	return reports, errors.Join(infoErr, ioErr)
}

func (a *Task) NetTask() ([]*NetReport, error) {
	netMap, err := psutil.GetNetworkIO(a.ethernet)
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
	if err != nil {
		return reports, fmt.Errorf("network io: %w", err)
	}
	return reports, nil
}
