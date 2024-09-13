// Package service
// Date: 2024/3/6 12:53
// Author: Amu
// Description:
package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"amprobe/service/host/repository"
	"amprobe/service/schema"
	rpcSchema "common/rpc/schema"

	"github.com/google/wire"
)

var HostServiceSet = wire.NewSet(NewHostService, wire.Bind(new(IHostService), new(*HostService)))

var _ IHostService = (*HostService)(nil)

type IHostService interface {
	HostInfo(context.Context) (schema.HostInfoReply, error)
	CPUInfo(context.Context) (schema.CPUInfoReply, error)
	CPUUsage(context.Context, schema.CPUUsageArgs) (schema.CPUUsageReply, error)
	MemInfo(context.Context) (schema.MemoryInfoReply, error)
	MemUsage(context.Context, schema.MemoryUsageArgs) (schema.MemoryUsageReply, error)
	DiskInfo(context.Context) (schema.DiskInfoReply, error)
	DiskUsage(context.Context, schema.DiskUsageArgs) (schema.DiskUsageReply, error)
	NetUsage(context.Context, schema.NetworkUsageArgs) (schema.NetworkUsageReply, error)

	FilesSearch(context.Context, schema.FilesSearchArgs) (schema.FilesSearchReply, error)
	FileCreate(context.Context, schema.FileCreateArgs) error
	FileDelete(context.Context, schema.FileDeleteArgs) error
	FileUpload(context.Context, schema.FileUploadArgs) error
	FileDownload(context.Context, schema.FileDownloadArgs) (schema.FileDownloadReply, error)
	FolderCreate(context.Context, schema.FolderCreateArgs) error
	Reboot(context.Context, schema.RebootArgs) error
	Shutdown(context.Context, schema.ShutdownArgs) error
	GetDNSSettings(context.Context, schema.GetDNSSettingsArgs) (schema.GetDNSSettingsReply, error)
	SetDNSSettings(context.Context, schema.SetDNSSettingsArgs) error
	GetSystemTime(context.Context, schema.GetSystemTimeArgs) (schema.GetSystemTimeReply, error)
	SetSystemTime(context.Context, schema.SetSystemTimeArgs) error
	GetSystemTimeZone(context.Context, schema.GetSystemTimeZoneArgs) (schema.GetSystemTimeZoneReply, error)
	SetSystemTimeZone(context.Context, schema.SetSystemTimeZoneArgs) error
	GetSystemTimeZoneList(context.Context, schema.GetSystemTimeZoneListArgs) (schema.GetSystemTimeZoneListReply, error)
}

type HostService struct {
	HostRepo repository.IHostRepo
}

func NewHostService(hostRepo repository.IHostRepo) *HostService {
	return &HostService{HostRepo: hostRepo}
}

func (h *HostService) HostInfo(ctx context.Context) (schema.HostInfoReply, error) {
	var reply schema.HostInfoReply
	info, err := h.HostRepo.HostInfo(ctx, rpcSchema.HostInfoArgs{})
	if err != nil {
		return reply, err
	}
	reply.Timestamp = info.Timestamp
	reply.Hostname = info.Hostname
	reply.Uptime = info.Uptime
	reply.KernelArch = info.KernelArch
	reply.KernelVersion = info.KernelVersion
	reply.PlatformVersion = info.PlatformVersion
	reply.Platform = info.Platform
	reply.OS = info.OS
	return reply, err
}

func (h *HostService) CPUInfo(ctx context.Context) (schema.CPUInfoReply, error) {
	var reply schema.CPUInfoReply
	cpuInfo, err := h.HostRepo.CPUInfo(ctx, rpcSchema.CPUInfoArgs{})
	if err != nil {
		return reply, err
	}
	reply.Percent = cpuInfo.Percent
	return reply, nil
}

func (h *HostService) CPUUsage(ctx context.Context, args schema.CPUUsageArgs) (schema.CPUUsageReply, error) {
	var reply schema.CPUUsageReply
	rpcArgs := rpcSchema.CPUUsageArgs{
		StartTime: args.StartTime,
		EndTime:   args.EndTime,
	}

	rpcReply, err := h.HostRepo.CPUUsage(ctx, rpcArgs)
	if err != nil {
		return reply, err
	}
	slog.Info("cpu usage rpc reply", "reply", rpcReply)
	var list []schema.Usage
	for _, item := range rpcReply.Data {
		list = append(list, schema.Usage{
			Timestamp: item.Timestamp,
			Value:     item.Value,
		})
	}
	reply.Data = list
	return reply, nil
}

func (h *HostService) MemInfo(ctx context.Context) (schema.MemoryInfoReply, error) {
	var reply schema.MemoryInfoReply
	memInfo, err := h.HostRepo.MemInfo(ctx, rpcSchema.MemoryInfoArgs{})
	if err != nil {
		return reply, err
	}
	reply.Used = memInfo.Used
	reply.Total = memInfo.Total
	reply.Percent = memInfo.Percent
	return reply, nil
}

func (h *HostService) MemUsage(ctx context.Context, args schema.MemoryUsageArgs) (schema.MemoryUsageReply, error) {
	var reply schema.MemoryUsageReply
	rpcArgs := rpcSchema.MemoryUsageArgs{
		StartTime: args.StartTime,
		EndTime:   args.EndTime,
	}
	rpcReply, err := h.HostRepo.MemUsage(ctx, rpcArgs)
	if err != nil {
		return reply, err
	}
	slog.Info("memory usage rpc reply", "reply", rpcReply)
	var list []schema.Usage
	for _, item := range rpcReply.Data {
		list = append(list, schema.Usage{
			Timestamp: item.Timestamp,
			Value:     item.Value,
		})
	}
	reply.Data = list
	return reply, nil
}

func (h *HostService) DiskInfo(ctx context.Context) (schema.DiskInfoReply, error) {
	var reply schema.DiskInfoReply
	rpcReply, err := h.HostRepo.DiskInfo(ctx, rpcSchema.DiskInfoArgs{})
	if err != nil {
		return reply, err
	}
	var list []schema.DiskInfo
	for _, di := range rpcReply.Info {
		list = append(list, schema.DiskInfo{
			Device:  di.Device,
			Percent: di.DiskPercent,
			Total:   di.DiskTotal,
			Used:    di.DiskUsed,
		})
	}
	reply.Info = list
	return reply, err
}

func (h *HostService) DiskUsage(ctx context.Context, args schema.DiskUsageArgs) (schema.DiskUsageReply, error) {
	var reply schema.DiskUsageReply
	rpcArgs := rpcSchema.DiskUsageArgs{
		StartTime: args.StartTime,
		EndTime:   args.EndTime,
	}
	rpcReply, err := h.HostRepo.DiskUsage(ctx, rpcArgs)
	if err != nil {
		return reply, err
	}
	slog.Info("disk usage rpc reply", "reply", rpcReply)
	for _, item := range rpcReply.Usage {
		var data []schema.DiskIO
		for _, i := range item.Data {
			data = append(data, schema.DiskIO{
				Timestamp: i.Timestamp,
				IORead:    i.IORead,
				IOWrite:   i.IOWrite,
			})
		}
		reply.Usage = append(reply.Usage, schema.DiskUsage{
			Device: item.Device,
			Data:   data,
		})
	}

	return reply, nil
}

func (h *HostService) NetUsage(ctx context.Context, args schema.NetworkUsageArgs) (schema.NetworkUsageReply, error) {
	var reply schema.NetworkUsageReply
	rpcArgs := rpcSchema.NetUsageArgs{
		StartTime: args.StartTime,
		EndTime:   args.EndTime,
	}
	rpcReply, err := h.HostRepo.NetUsage(ctx, rpcArgs)
	if err != nil {
		return reply, err
	}
	slog.Info("net usage rpc reply", "reply", rpcReply)
	for _, item := range rpcReply.Usage {
		var data []schema.NetIO
		for _, i := range item.Data {
			data = append(data, schema.NetIO{
				Timestamp: i.Timestamp,
				BytesSent: i.BytesSent,
				BytesRecv: i.BytesRecv,
			})
			reply.Usage = append(reply.Usage, schema.NetUsage{
				Ethernet: item.Ethernet,
				Data:     data,
			})
		}
	}

	return reply, nil
}

func (h *HostService) FilesSearch(ctx context.Context, args schema.FilesSearchArgs) (schema.FilesSearchReply, error) {
	var reply schema.FilesSearchReply
	rpcArgs := rpcSchema.FilesSearchArgs{
		Path: args.Path,
	}
	rpcReply, err := h.HostRepo.FilesSearch(ctx, rpcArgs)
	if err != nil {
		return reply, err
	}
	for _, file := range rpcReply.Files {
		reply.Files = append(reply.Files, schema.FileInfo{
			Name:    file.Name,
			Size:    file.Size,
			Mode:    file.Mode,
			ModTime: file.ModTime,
			IsDir:   file.IsDir,
		})
	}
	return reply, nil
}

func (h *HostService) FileCreate(ctx context.Context, args schema.FileCreateArgs) error {
	rpcArgs := rpcSchema.FileCreateArgs{
		Path:     args.Path,
		FileName: args.FileName,
	}
	return h.HostRepo.FileCreate(ctx, rpcArgs)
}

func (h *HostService) FileDelete(ctx context.Context, args schema.FileDeleteArgs) error {
	rpcArgs := rpcSchema.FileDeleteArgs{
		Filepath: args.Filepath,
	}
	return h.HostRepo.FileDelete(ctx, rpcArgs)
}

func (h *HostService) FileUpload(ctx context.Context, args schema.FileUploadArgs) error {
	rpcArgs := rpcSchema.FileUploadArgs{
		SourceFilePath: args.SourceFilePath,
		TargetFilePath: args.TargetFilePath,
	}
	return h.HostRepo.FileUpload(ctx, rpcArgs)
}

func (h *HostService) FileDownload(ctx context.Context, args schema.FileDownloadArgs) (schema.FileDownloadReply, error) {
	var reply schema.FileDownloadReply
	filename := strings.Split(args.Filepath, "/")[len(strings.Split(args.Filepath, "/"))-1]
	rpcArgs := rpcSchema.FileDownloadArgs{
		SourceFilePath: args.Filepath,
		TargetFilePath: fmt.Sprintf("/tmp/%s", filename),
	}
	rpcReply, err := h.HostRepo.FileDownload(ctx, rpcArgs)
	if err != nil {
		return reply, err
	}
	reply.Filepath = rpcReply.Filepath
	return reply, nil
}

func (h *HostService) FolderCreate(ctx context.Context, args schema.FolderCreateArgs) error {
	rpcArgs := rpcSchema.FolderCreateArgs{
		Path:       args.Path,
		FolderName: args.FolderName,
	}
	return h.HostRepo.FolderCreate(ctx, rpcArgs)
}

func (h *HostService) Reboot(ctx context.Context, args schema.RebootArgs) error {
	return h.HostRepo.Reboot(ctx, rpcSchema.RebootArgs{})
}

func (h *HostService) Shutdown(ctx context.Context, args schema.ShutdownArgs) error {
	return h.HostRepo.Shutdown(ctx, rpcSchema.ShutdownArgs{})
}

func (h *HostService) GetDNSSettings(ctx context.Context, args schema.GetDNSSettingsArgs) (schema.GetDNSSettingsReply, error) {
	var reply schema.GetDNSSettingsReply
	rpcReply, err := h.HostRepo.GetDNSSettings(ctx, rpcSchema.GetDNSArgs{})
	if err != nil {
		return reply, err
	}
	reply.DNS = rpcReply.DNS
	return reply, nil
}

func (h *HostService) SetDNSSettings(ctx context.Context, args schema.SetDNSSettingsArgs) error {
	rpcArgs := rpcSchema.SetDNSArgs{
		DNS: args.DNS,
	}
	return h.HostRepo.SetDNSSettings(ctx, rpcArgs)
}

func (h *HostService) GetSystemTime(ctx context.Context, args schema.GetSystemTimeArgs) (schema.GetSystemTimeReply, error) {
	var reply schema.GetSystemTimeReply
	rpcReply, err := h.HostRepo.GetSystemTime(ctx, rpcSchema.GetSystemTimeArgs{})
	if err != nil {
		return reply, err
	}
	reply.SystemTime = rpcReply.SystemTime
	return reply, nil
}

func (h *HostService) SetSystemTime(ctx context.Context, args schema.SetSystemTimeArgs) error {
	rpcArgs := rpcSchema.SetSystemTimeArgs{
		SystemTime: args.SystemTime,
	}
	return h.HostRepo.SetSystemTime(ctx, rpcArgs)
}

func (h *HostService) GetSystemTimeZone(ctx context.Context, args schema.GetSystemTimeZoneArgs) (schema.GetSystemTimeZoneReply, error) {
	var reply schema.GetSystemTimeZoneReply
	rpcReply, err := h.HostRepo.GetSystemTimeZone(ctx, rpcSchema.GetSystemTimeZoneArgs{})
	if err != nil {
		return reply, err
	}
	reply.SystemTimeZone = rpcReply.SystemTimeZone
	return reply, nil
}

func (h *HostService) SetSystemTimeZone(ctx context.Context, args schema.SetSystemTimeZoneArgs) error {
	rpcArgs := rpcSchema.SetSystemTimeZoneArgs{
		SystemTimeZone: args.SystemTimeZone,
	}
	return h.HostRepo.SetSystemTimeZone(ctx, rpcArgs)
}

func (h *HostService) GetSystemTimeZoneList(ctx context.Context, args schema.GetSystemTimeZoneListArgs) (schema.GetSystemTimeZoneListReply, error) {
	var reply schema.GetSystemTimeZoneListReply
	rpcReply, err := h.HostRepo.GetSystemTimeZoneList(ctx, rpcSchema.GetSystemTimeZoneListArgs{})
	if err != nil {
		return reply, err
	}
	reply.SystemTimeZoneList = rpcReply.SystemTimeZoneList
	return reply, nil
}
