// Package rpc
// Date: 2024/06/11 19:26:34
// Author: Amu
// Description:
package rpc

import (
	"context"

	"github.com/amuluze/amprobe/pkg/rpc"
	"github.com/amuluze/amprobe/service/model"

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
	FilesSearch(ctx context.Context, args schema.FilesSearchArgs) (schema.FilesSearchReply, error)
	FileUpload(ctx context.Context, args schema.FileUploadArgs) error
	FileDownload(ctx context.Context, args schema.FileDownloadArgs) (schema.FileDownloadReply, error)
	FileDelete(ctx context.Context, args schema.FileDeleteArgs) error
	GetDNSSettings(ctx context.Context, args schema.GetDNSSettingsArgs) (schema.GetDNSSettingsReply, error)
	SetDNSSettings(ctx context.Context, args schema.SetDNSSettingsArgs) error
	GetSystemTime(ctx context.Context, args schema.GetSystemTimeArgs) (schema.GetSystemTimeReply, error)
	SetSystemTime(ctx context.Context, args schema.SetSystemTimeArgs) error
	GetSystemTimezoneList(ctx context.Context, args schema.GetSystemTimezoneListArgs) (schema.GetSystemTimezoneListReply, error)
	GetSystemTimezone(ctx context.Context, args schema.GetSystemTimezoneArgs) (schema.GetSystemTimezoneReply, error)
	SetSystemTimezone(ctx context.Context, args schema.SetSystemTimezoneArgs) error
	Reboot(ctx context.Context, args schema.RebootArgs) error
	Shutdown(ctx context.Context, args schema.ShutdownArgs) error
}

type HostService struct {
	RPCClient *rpc.Client
}

func NewHostService(client *rpc.Client) *HostService {
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
	args := schema.CPUArgs{}
	var reply model.CPU
	err := h.RPCClient.Call(ctx, "CPUInfo", args, &reply)
	if err != nil {
		return schema.CPUInfoReply{}, err
	}
	return schema.CPUInfoReply{
		Percent: reply.CPUPercent,
	}, nil
}

func (h HostService) CPUUsage(ctx context.Context, args schema.CPUUsageArgs) (schema.CPUUsageReply, error) {
	var reply []model.CPU
	err := h.RPCClient.Call(ctx, "CPUUsage", args, &reply)
	if err != nil {
		return schema.CPUUsageReply{}, err
	}
	var result schema.CPUUsageReply
	var data []schema.Usage
	for _, r := range reply {
		data = append(data, schema.Usage{
			Timestamp: r.Timestamp.Unix(),
			Value:     r.CPUPercent,
		})
	}
	result.Data = data
	return result, nil
}

func (h HostService) MemInfo(ctx context.Context) (schema.MemoryInfoReply, error) {
	args := schema.MemoryArgs{}
	var reply model.Memory
	err := h.RPCClient.Call(ctx, "MemInfo", args, &reply)
	if err != nil {
		return schema.MemoryInfoReply{}, err
	}
	return schema.MemoryInfoReply{
		Total:   reply.MemTotal,
		Used:    reply.MemUsed,
		Percent: reply.MemPercent,
	}, nil
}

func (h HostService) MemUsage(ctx context.Context, args schema.MemoryUsageArgs) (schema.MemoryUsageReply, error) {
	var reply []model.Memory
	err := h.RPCClient.Call(ctx, "MemUsage", args, &reply)
	if err != nil {
		return schema.MemoryUsageReply{}, err
	}
	var result schema.MemoryUsageReply
	var data []schema.Usage
	for _, r := range reply {
		data = append(data, schema.Usage{
			Timestamp: r.Timestamp.Unix(),
			Value:     r.MemPercent,
		})
	}
	result.Data = data
	return result, nil
}

func (h HostService) DiskInfo(ctx context.Context) (schema.DiskInfoReply, error) {
	args := schema.DiskArgs{}
	var reply []model.Disk
	err := h.RPCClient.Call(ctx, "DiskInfo", args, &reply)
	if err != nil {
		return schema.DiskInfoReply{}, err
	}
	var result schema.DiskInfoReply
	var data []schema.DiskInfo
	for _, r := range reply {
		data = append(data, schema.DiskInfo{
			Device:  r.Device,
			Total:   r.DiskTotal,
			Used:    r.DiskUsed,
			Percent: r.DiskPercent,
		})
	}
	result.Info = data
	return result, nil
}

func (h HostService) DiskUsage(ctx context.Context, args schema.DiskUsageArgs) (schema.DiskUsageReply, error) {
	var reply []model.Disk
	err := h.RPCClient.Call(ctx, "DiskUsage", args, &reply)
	if err != nil {
		return schema.DiskUsageReply{}, err
	}
	var result schema.DiskUsageReply
	var data []schema.DiskUsage
	list := make(map[string][]schema.DiskIO)
	for _, r := range reply {
		list[r.Device] = append(list[r.Device], schema.DiskIO{
			Timestamp: r.Timestamp.Unix(),
			IORead:    r.DiskRead,
			IOWrite:   r.DiskWrite,
		})
	}
	for device, l := range list {
		data = append(data, schema.DiskUsage{
			Device: device,
			Data:   l,
		})
	}
	result.Usage = data
	return result, nil
}

func (h HostService) NetUsage(ctx context.Context, args schema.NetworkUsageArgs) (schema.NetworkUsageReply, error) {
	var reply []model.Net
	err := h.RPCClient.Call(ctx, "NetUsage", args, &reply)
	if err != nil {
		return schema.NetworkUsageReply{}, err
	}
	var result schema.NetworkUsageReply
	var data []schema.NetUsage
	list := make(map[string][]schema.NetIO)
	for _, r := range reply {
		list[r.Ethernet] = append(list[r.Ethernet], schema.NetIO{
			Timestamp: r.Timestamp.Unix(),
			BytesSent: r.NetSend,
			BytesRecv: r.NetRecv,
		})
	}
	for eth, l := range list {
		data = append(data, schema.NetUsage{
			Ethernet: eth,
			Data:     l,
		})
	}
	result.Usage = data
	return result, nil
}

func (h HostService) FilesSearch(ctx context.Context, args schema.FilesSearchArgs) (schema.FilesSearchReply, error) {
	var reply schema.FilesSearchReply
	err := h.RPCClient.Call(ctx, "FilesSearch", args, &reply)
	if err != nil {
		return schema.FilesSearchReply{}, err
	}
	return reply, nil
}

func (h HostService) FileUpload(ctx context.Context, args schema.FileUploadArgs) error {
	return nil
}

func (h HostService) FileDownload(ctx context.Context, args schema.FileDownloadArgs) (schema.FileDownloadReply, error) {
	return schema.FileDownloadReply{}, nil
}

func (h HostService) FileDelete(ctx context.Context, args schema.FileDeleteArgs) error {
	return nil
}

func (h HostService) GetDNSSettings(ctx context.Context, args schema.GetDNSSettingsArgs) (schema.GetDNSSettingsReply, error) {
	var reply schema.GetDNSSettingsReply
	err := h.RPCClient.Call(ctx, "GetDNS", args, &reply)
	if err != nil {
		return schema.GetDNSSettingsReply{}, err
	}
	return reply, nil
}

func (h HostService) SetDNSSettings(ctx context.Context, args schema.SetDNSSettingsArgs) error {
	var reply schema.SetDNSSettingsReply
	err := h.RPCClient.Call(ctx, "SetDNS", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (h HostService) GetSystemTime(ctx context.Context, args schema.GetSystemTimeArgs) (schema.GetSystemTimeReply, error) {
	var reply schema.GetSystemTimeReply
	err := h.RPCClient.Call(ctx, "GetSystemTime", args, &reply)
	if err != nil {
		return schema.GetSystemTimeReply{}, err
	}
	return reply, nil
}

func (h HostService) SetSystemTime(ctx context.Context, args schema.SetSystemTimeArgs) error {
	var reply schema.SetSystemTimeReply
	err := h.RPCClient.Call(ctx, "SetSystemTime", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (h HostService) GetSystemTimezoneList(ctx context.Context, args schema.GetSystemTimezoneListArgs) (schema.GetSystemTimezoneListReply, error) {
	var reply schema.GetSystemTimezoneListReply
	err := h.RPCClient.Call(ctx, "GetSystemTimeZoneList", args, &reply)
	if err != nil {
		return schema.GetSystemTimezoneListReply{}, err
	}
	return reply, nil
}

func (h HostService) GetSystemTimezone(ctx context.Context, args schema.GetSystemTimezoneArgs) (schema.GetSystemTimezoneReply, error) {
	var reply schema.GetSystemTimezoneReply
	err := h.RPCClient.Call(ctx, "GetSystemTimeZone", args, &reply)
	if err != nil {
		return schema.GetSystemTimezoneReply{}, err
	}
	return reply, nil
}
func (h HostService) SetSystemTimezone(ctx context.Context, args schema.SetSystemTimezoneArgs) error {
	var reply schema.SetSystemTimezoneReply
	err := h.RPCClient.Call(ctx, "SetSystemTimeZone", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (h HostService) Reboot(ctx context.Context, args schema.RebootArgs) error {
	var reply schema.RebootReply
	err := h.RPCClient.Call(ctx, "Reboot", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (h HostService) Shutdown(ctx context.Context, args schema.ShutdownArgs) error {
	var reply schema.ShutdownReply
	err := h.RPCClient.Call(ctx, "Shutdown", args, &reply)
	if err != nil {
		return err
	}
	return nil
}
