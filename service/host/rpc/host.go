// Package rpc
// Date: 2024/06/11 19:26:34
// Author: Amu
// Description:
package rpc

import (
	"context"
	"github.com/amuluze/amprobe/service/model"

	"github.com/amuluze/amprobe/service"
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
	RPCClient *service.RPCClient
}

func NewHostService(client *service.RPCClient) *HostService {
	return &HostService{
		RPCClient: client,
	}
}

func (h HostService) HostInfo(ctx context.Context) (schema.HostInfoReply, error) {
	args := schema.HostArgs{}
	var reply model.Host
	err := h.RPCClient.Call(ctx, "HostInfo", args, &reply)
	if err != nil {
		return schema.HostInfoReply{}, err
	}
	return schema.HostInfoReply{
		Timestamp:       reply.Timestamp.Unix(),
		Uptime:          reply.Uptime,
		Hostname:        reply.Hostname,
		OS:              reply.Os,
		Platform:        reply.Platform,
		PlatformVersion: reply.PlatformVersion,
		KernelVersion:   reply.KernelVersion,
		KernelArch:      reply.KernelArch,
	}, nil
}

func (h HostService) CPUInfo(ctx context.Context) (schema.CPUInfoReply, error) {
	//TODO implement me
	panic("implement me")
}

func (h HostService) CPUUsage(ctx context.Context, args schema.CPUUsageArgs) (schema.CPUUsageReply, error) {
	//TODO implement me
	panic("implement me")
}

func (h HostService) MemInfo(ctx context.Context) (schema.MemoryInfoReply, error) {
	//TODO implement me
	panic("implement me")
}

func (h HostService) MemUsage(ctx context.Context, args schema.MemoryUsageArgs) (schema.MemoryUsageReply, error) {
	//TODO implement me
	panic("implement me")
}

func (h HostService) DiskInfo(ctx context.Context) (schema.DiskInfoReply, error) {
	//TODO implement me
	panic("implement me")
}

func (h HostService) DiskUsage(ctx context.Context, args schema.DiskUsageArgs) (schema.DiskUsageReply, error) {
	//TODO implement me
	panic("implement me")
}

func (h HostService) NetUsage(ctx context.Context, args schema.NetworkUsageArgs) (schema.NetworkUsageReply, error) {
	//TODO implement me
	panic("implement me")
}
