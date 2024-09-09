// Package service
// Date: 2024/3/6 12:53
// Author: Amu
// Description:
package service

import (
	"amprobe/service/host/repository"
	"context"

	"amprobe/service/schema"

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
	var reply schema.DiskUsageReply
	usageMap := make(map[string][]schema.DiskIO)

	for _, item := range diskInfos {
		if _, ok := usageMap[item.Device]; !ok {
			usageMap[item.Device] = []schema.DiskIO{{
				Timestamp: item.Timestamp.Unix(),
				IORead:    item.DiskRead,
				IOWrite:   item.DiskWrite,
			}}
		} else {
			usageMap[item.Device] = append(usageMap[item.Device], schema.DiskIO{
				Timestamp: item.Timestamp.Unix(),
				IORead:    item.DiskRead,
				IOWrite:   item.DiskWrite,
			})
		}
	}
	for key, item := range usageMap {
		reply.Usage = append(reply.Usage, schema.DiskUsage{
			Device: key,
			Data:   item,
		})
	}
	return reply, nil
}

func (h HostService) NetUsage(ctx context.Context, args schema.NetworkUsageArgs) (schema.NetworkUsageReply, error) {
	netInfos, err := h.HostRepo.NetUsage(ctx, args)
	if err != nil {
		return schema.NetworkUsageReply{}, err
	}

	var reply schema.NetworkUsageReply
	usageMap := make(map[string][]schema.NetIO)

	for _, item := range netInfos {
		if _, ok := usageMap[item.Ethernet]; !ok {
			usageMap[item.Ethernet] = []schema.NetIO{{
				Timestamp: item.Timestamp.Unix(),
				BytesSent: item.NetSend,
				BytesRecv: item.NetRecv,
			}}
		} else {
			usageMap[item.Ethernet] = append(usageMap[item.Ethernet], schema.NetIO{
				Timestamp: item.Timestamp.Unix(),
				BytesSent: item.NetSend,
				BytesRecv: item.NetRecv,
			})
		}
	}
	for key, item := range usageMap {
		reply.Usage = append(reply.Usage, schema.NetUsage{
			Ethernet: key,
			Data:     item,
		})
	}
	return reply, nil
}
