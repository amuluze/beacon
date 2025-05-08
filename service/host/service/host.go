// Package service
// Date: 2024/3/6 12:53
// Author: Amu
// Description:
package service

import (
	"amprobe/pkg/psutil"
	"amprobe/service/host/repository"
	"amprobe/service/schema"
	"context"
	"os"
	"path/filepath"
	"time"

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
	NetUsage(context.Context, schema.NetUsageArgs) (schema.NetUsageReply, error)

	FilesSearch(context.Context, schema.FilesSearchArgs) (schema.FilesSearchReply, error)
	FileCreate(context.Context, schema.FileCreateArgs) error
	FileDelete(context.Context, schema.FileDeleteArgs) error
	FileDownload(context.Context, schema.FileDownloadArgs) (schema.FileDownloadReply, error)
	FolderCreate(context.Context, schema.FolderCreateArgs) error

	GetDNSSettings(context.Context) (schema.GetDNSSettingsReply, error)
	SetDNSSettings(context.Context, schema.SetDNSSettingsArgs) error
	GetSystemTime(context.Context) (schema.GetSystemTimeReply, error)
	SetSystemTime(context.Context, schema.SetSystemTimeArgs) error
	GetSystemTimeZone(context.Context) (schema.GetSystemTimeZoneReply, error)
	SetSystemTimeZone(context.Context, schema.SetSystemTimeZoneArgs) error
	GetSystemTimeZoneList(context.Context) (schema.GetSystemTimeZoneListReply, error)

	Reboot(context.Context) error
	Shutdown(context.Context) error
}

type HostService struct {
	HostRepo repository.IHostRepo
}

func NewHostService(hostRepo repository.IHostRepo) *HostService {
	return &HostService{HostRepo: hostRepo}
}

func (h *HostService) HostInfo(ctx context.Context) (schema.HostInfoReply, error) {
	var reply schema.HostInfoReply
	info, err := psutil.GetSystemInfo()
	if err != nil {
		return reply, err
	}

	reply.Timestamp = time.Now().Unix()
	reply.Hostname = info.Hostname
	reply.Uptime = info.Uptime
	reply.KernelArch = info.KernelArch
	reply.KernelVersion = info.KernelVersion
	reply.PlatformVersion = info.PlatformVersion
	reply.Platform = info.Platform
	reply.OS = info.Os
	return reply, nil
}

func (h *HostService) CPUInfo(ctx context.Context) (schema.CPUInfoReply, error) {
	var reply schema.CPUInfoReply
	cpuInfo, err := h.HostRepo.CPUInfo(ctx)
	if err != nil {
		return reply, err
	}
	reply.Percent = cpuInfo.CPUPercent
	return reply, nil
}

func (h *HostService) CPUUsage(ctx context.Context, args schema.CPUUsageArgs) (schema.CPUUsageReply, error) {
	var reply schema.CPUUsageReply
	results, err := h.HostRepo.CPUUsage(ctx, args)
	if err != nil {
		return reply, err
	}
	var list []schema.Usage
	for _, item := range results {
		list = append(list, schema.Usage{
			Timestamp: item.Timestamp.Unix(),
			Value:     item.CPUPercent,
		})
	}
	reply.Data = list
	return reply, nil
}

func (h *HostService) MemInfo(ctx context.Context) (schema.MemoryInfoReply, error) {
	var reply schema.MemoryInfoReply
	memInfo, err := h.HostRepo.MemInfo(ctx)
	if err != nil {
		return reply, err
	}
	reply.Used = memInfo.MemUsed
	reply.Total = memInfo.MemTotal
	reply.Percent = memInfo.MemPercent
	return reply, nil
}

func (h *HostService) MemUsage(ctx context.Context, args schema.MemoryUsageArgs) (schema.MemoryUsageReply, error) {
	var reply schema.MemoryUsageReply
	results, err := h.HostRepo.MemUsage(ctx, args)
	if err != nil {
		return reply, err
	}
	var list []schema.Usage
	for _, item := range results {
		list = append(list, schema.Usage{
			Timestamp: item.Timestamp.Unix(),
			Value:     item.MemPercent,
		})
	}
	reply.Data = list
	return reply, nil
}

func (h *HostService) DiskInfo(ctx context.Context) (schema.DiskInfoReply, error) {
	var reply schema.DiskInfoReply
	results, err := h.HostRepo.DiskInfo(ctx)
	if err != nil {
		return reply, err
	}
	var list []schema.DiskInfo
	for _, info := range results {
		list = append(list, schema.DiskInfo{
			Device:  info.Device,
			Percent: info.DiskPercent,
			Total:   info.DiskTotal,
			Used:    info.DiskUsed,
		})
	}
	reply.Info = list
	return reply, err
}

func (h *HostService) DiskUsage(ctx context.Context, args schema.DiskUsageArgs) (schema.DiskUsageReply, error) {
	var reply schema.DiskUsageReply
	results, err := h.HostRepo.DiskUsage(ctx, args)
	if err != nil {
		return reply, err
	}
	var data []schema.DiskUsage
	list := make(map[string][]schema.DiskIO)
	for _, item := range results {
		list[item.Device] = append(list[item.Device], schema.DiskIO{
			Timestamp: item.Timestamp.Unix(),
			IORead:    item.DiskRead,
			IOWrite:   item.DiskWrite,
		})
	}
	for device, l := range list {
		data = append(data, schema.DiskUsage{
			Device: device,
			Data:   l,
		})
	}
	reply.Usage = data

	return reply, nil
}

func (h *HostService) NetUsage(ctx context.Context, args schema.NetUsageArgs) (schema.NetUsageReply, error) {
	var reply schema.NetUsageReply
	results, err := h.HostRepo.NetUsage(ctx, args)
	if err != nil {
		return reply, err
	}
	var data []schema.NetUsage
	list := make(map[string][]schema.NetIO)
	for _, item := range results {
		list[item.Ethernet] = append(list[item.Ethernet], schema.NetIO{
			Timestamp: item.Timestamp.Unix(),
			BytesRecv: item.NetRecv,
			BytesSent: item.NetSend,
		})
	}
	for eth, l := range list {
		data = append(data, schema.NetUsage{
			Ethernet: eth,
			Data:     l,
		})
	}
	reply.Usage = data

	return reply, nil
}

func (h *HostService) FilesSearch(ctx context.Context, args schema.FilesSearchArgs) (schema.FilesSearchReply, error) {
	var reply schema.FilesSearchReply
	files, err := os.ReadDir(args.Path)
	if err != nil {
		return reply, err
	}
	list := make([]schema.FileInfo, 0)
	for _, f := range files {
		info, _ := f.Info()
		list = append(list, schema.FileInfo{
			Name:    f.Name(),
			IsDir:   f.IsDir(),
			Size:    info.Size(),
			Mode:    info.Mode().String(),
			ModTime: info.ModTime().Unix(),
		})
	}
	reply.Files = list
	return reply, nil
}

func (h *HostService) FileCreate(ctx context.Context, args schema.FileCreateArgs) error {
	filePath := filepath.Join(args.Path, args.FileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0755))
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *HostService) FileDelete(ctx context.Context, args schema.FileDeleteArgs) error {
	if info, err := os.Stat(args.Filepath); err != nil {
		return err
	} else if info.IsDir() {
		return os.RemoveAll(args.Filepath)
	} else {
		return os.Remove(args.Filepath)
	}
}

func (h *HostService) FileDownload(ctx context.Context, args schema.FileDownloadArgs) (schema.FileDownloadReply, error) {
	var reply schema.FileDownloadReply
	if _, err := os.Stat(args.Filepath); os.IsNotExist(err) {
		return reply, err
	}
	reply.Filepath = args.Filepath
	return reply, nil
}

func (h *HostService) FolderCreate(ctx context.Context, args schema.FolderCreateArgs) error {
	folderPath := filepath.Join(args.Path, args.FolderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return os.MkdirAll(folderPath, os.ModePerm)
	}
	return nil
}

func (h *HostService) GetDNSSettings(ctx context.Context) (schema.GetDNSSettingsReply, error) {
	var reply schema.GetDNSSettingsReply
	dns, err := h.HostRepo.GetDNSSettings(ctx)
	if err != nil {
		return reply, err
	}
	reply.DNS = dns
	return reply, nil
}

func (h *HostService) SetDNSSettings(ctx context.Context, args schema.SetDNSSettingsArgs) error {
	return h.HostRepo.SetDNSSettings(ctx, args)
}

func (h *HostService) GetSystemTime(ctx context.Context) (schema.GetSystemTimeReply, error) {
	var reply schema.GetSystemTimeReply
	systemTime, err := h.HostRepo.GetSystemTime(ctx)
	if err != nil {
		return reply, err
	}
	reply.SystemTime = systemTime
	return reply, nil
}

func (h *HostService) SetSystemTime(ctx context.Context, args schema.SetSystemTimeArgs) error {
	return h.HostRepo.SetSystemTime(ctx, args)
}

func (h *HostService) GetSystemTimeZone(ctx context.Context) (schema.GetSystemTimeZoneReply, error) {
	var reply schema.GetSystemTimeZoneReply
	timezone, err := h.HostRepo.GetSystemTimeZone(ctx)
	if err != nil {
		return reply, err
	}
	reply.SystemTimeZone = timezone
	return reply, nil
}

func (h *HostService) SetSystemTimeZone(ctx context.Context, args schema.SetSystemTimeZoneArgs) error {
	return h.HostRepo.SetSystemTimeZone(ctx, args)
}

func (h *HostService) GetSystemTimeZoneList(ctx context.Context) (schema.GetSystemTimeZoneListReply, error) {
	var reply schema.GetSystemTimeZoneListReply
	timezoneList, err := h.HostRepo.GetSystemTimeZoneList(ctx)
	if err != nil {
		return reply, err
	}
	reply.SystemTimeZoneList = timezoneList
	return reply, nil
}

func (h *HostService) Reboot(ctx context.Context) error {
	return h.HostRepo.Reboot(ctx)
}

func (h *HostService) Shutdown(ctx context.Context) error {
	return h.HostRepo.Shutdown(ctx)
}
