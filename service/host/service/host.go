// Package service
// Date: 2024/3/6 12:53
// Author: Amu
// Description:
package service

import (
	"context"
	
	"github.com/amuluze/amprobe/service/host/repository"
	"github.com/amuluze/amprobe/service/schema"
	"github.com/google/wire"
)

var HostServiceSet = wire.NewSet(NewHostService, wire.Bind(new(IHostService), new(*HostService)))

type IHostService interface {
	HostInfo(ctx context.Context) (schema.HostInfoReply, error)
	CPUInfo(ctx context.Context) (schema.CPUInfoReply, error)
	CPUUsage(ctx context.Context, args schema.CPUUsageArgs) (schema.CPUUsageReply, error)
	MemInfo(ctx context.Context) (schema.MemoryInfoReply, error)
	MemUsage(ctx context.Context, args schema.MemoryUsageArgs) (schema.MemoryUsageReply, error)
	DiskInfo(ctx context.Context) (schema.DiskInfoReply, error)
	DiskUsage(ctx context.Context, args schema.DiskUsageArgs) (schema.DiskUsageReply, error)
	NetUsage(ctx context.Context, args schema.NetworkUsageArgs) (schema.NetworkUsageReply, error)
}

type HostService struct {
	HostRepo repository.IHostRepo
}

func NewHostService(hostRepo repository.IHostRepo) *HostService {
	return &HostService{HostRepo: hostRepo}
}

func (h HostService) HostInfo(ctx context.Context) (schema.HostInfoReply, error) {
	info, err := h.HostRepo.HostInfo(ctx)
	if err != nil {
		return schema.HostInfoReply{}, err
	}
	return schema.HostInfoReply{
		Timestamp:       info.Timestamp.Unix(),
		Uptime:          info.Uptime,
		Hostname:        info.Hostname,
		OS:              info.Os,
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		KernelArch:      info.KernelArch,
	}, err
}

func (h HostService) CPUInfo(ctx context.Context) (schema.CPUInfoReply, error) {
	cpuInfo, err := h.HostRepo.CPUInfo(ctx)
	if err != nil {
		return schema.CPUInfoReply{}, err
	}
	return schema.CPUInfoReply{Percent: cpuInfo.CPUPercent}, nil
}

func (h HostService) CPUUsage(ctx context.Context, args schema.CPUUsageArgs) (schema.CPUUsageReply, error) {
	mHosts, err := h.HostRepo.CPUUsage(ctx, args)
	if err != nil {
		return schema.CPUUsageReply{}, err
	}
	var list []schema.Usage
	for _, item := range mHosts {
		list = append(list, schema.Usage{
			Timestamp: item.Timestamp.Unix(),
			Value:     item.CPUPercent,
		})
	}
	return schema.CPUUsageReply{Data: list}, nil
}

func (h HostService) MemInfo(ctx context.Context) (schema.MemoryInfoReply, error) {
	memInfo, err := h.HostRepo.MemInfo(ctx)
	if err != nil {
		return schema.MemoryInfoReply{}, err
	}
	return schema.MemoryInfoReply{Percent: memInfo.MemPercent, Total: memInfo.MemTotal, Used: memInfo.MemUsed}, nil
}

func (h HostService) MemUsage(ctx context.Context, args schema.MemoryUsageArgs) (schema.MemoryUsageReply, error) {
	memInfos, err := h.HostRepo.MemUsage(ctx, args)
	if err != nil {
		return schema.MemoryUsageReply{}, err
	}
	var list []schema.Usage
	for _, item := range memInfos {
		list = append(list, schema.Usage{
			Timestamp: item.Timestamp.Unix(),
			Value:     item.MemPercent,
		})
	}
	return schema.MemoryUsageReply{Data: list}, nil
}

func (h HostService) DiskInfo(ctx context.Context) (schema.DiskInfoReply, error) {
	diskInfos, err := h.HostRepo.DiskInfo(ctx)
	if err != nil {
		return schema.DiskInfoReply{}, err
	}
	var list []schema.DiskInfo
	for _, di := range diskInfos {
		list = append(list, schema.DiskInfo{
			Device:  di.Device,
			Percent: di.DiskPercent,
			Total:   di.DiskTotal,
			Used:    di.DiskUsed,
		})
	}
	
	return schema.DiskInfoReply{Info: list}, err
}

func (h HostService) DiskUsage(ctx context.Context, args schema.DiskUsageArgs) (schema.DiskUsageReply, error) {
	diskInfos, err := h.HostRepo.DiskUsage(ctx, args)
	if err != nil {
		return schema.DiskUsageReply{}, err
	}
	
	mDisk := make([]schema.DiskIO, 0)
	device := ""
	for _, item := range diskInfos {
		device = item.Device
		mDisk = append(mDisk, schema.DiskIO{
			Timestamp: item.Timestamp.Unix(),
			IORead:    item.DiskRead,
			IOWrite:   item.DiskWrite,
		})
	}
	return schema.DiskUsageReply{Device: device, Data: mDisk}, nil
}

func (h HostService) NetUsage(ctx context.Context, args schema.NetworkUsageArgs) (schema.NetworkUsageReply, error) {
	netInfos, err := h.HostRepo.NetUsage(ctx, args)
	if err != nil {
		return schema.NetworkUsageReply{}, err
	}
	mNet := make([]schema.NetIO, 0)
	ethernet := ""
	for _, item := range netInfos {
		ethernet = item.Ethernet
		mNet = append(mNet, schema.NetIO{
			Timestamp: item.Timestamp.Unix(),
			BytesSent: item.NetSend,
			BytesRecv: item.NetRecv,
		})
	}
	return schema.NetworkUsageReply{Data: mNet, Ethernet: ethernet}, nil
}
