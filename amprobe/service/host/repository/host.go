// Package repository
package repository

import (
	"context"
	"time"

	"amprobe/pkg/contextx"
	"amprobe/pkg/rpc"
	"amprobe/service/model"
	"common/database"
	rpcSchema "common/rpc/schema"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var HostRepoSet = wire.NewSet(NewHostRepo, wire.Bind(new(IHostRepo), new(*HostRepo)))

var _ IHostRepo = (*HostRepo)(nil)

type IHostRepo interface {
	HostInfo(context.Context, rpcSchema.HostInfoArgs) (rpcSchema.HostInfoReply, error)
	CPUInfo(context.Context, rpcSchema.CPUInfoArgs) (rpcSchema.CPUInfoReply, error)
	CPUUsage(context.Context, rpcSchema.CPUUsageArgs) (rpcSchema.CPUUsageReply, error)
	MemInfo(context.Context, rpcSchema.MemoryInfoArgs) (rpcSchema.MemoryInfoReply, error)
	MemUsage(context.Context, rpcSchema.MemoryUsageArgs) (rpcSchema.MemoryUsageReply, error)
	DiskInfo(context.Context, rpcSchema.DiskInfoArgs) (rpcSchema.DiskInfoReply, error)
	DiskUsage(context.Context, rpcSchema.DiskUsageArgs) (rpcSchema.DiskUsageReply, error)
	NetUsage(context.Context, rpcSchema.NetUsageArgs) (rpcSchema.NetUsageReply, error)

	FilesSearch(context.Context, rpcSchema.FilesSearchArgs) (rpcSchema.FilesSearchReply, error)
	FileUpload(context.Context, rpcSchema.FileUploadArgs) error
	FileDownload(context.Context, rpcSchema.FileDownloadArgs) (rpcSchema.FileDownloadReply, error)
	FileDelete(context.Context, rpcSchema.FileDeleteArgs) error
	FileCreate(context.Context, rpcSchema.FileCreateArgs) error
	FolderCreate(context.Context, rpcSchema.FolderCreateArgs) error
	GetDNSSettings(context.Context, rpcSchema.GetDNSArgs) (rpcSchema.GetDNSReply, error)
	SetDNSSettings(context.Context, rpcSchema.SetDNSArgs) error
	GetSystemTime(context.Context, rpcSchema.GetSystemTimeArgs) (rpcSchema.GetSystemTimeReply, error)
	SetSystemTime(context.Context, rpcSchema.SetSystemTimeArgs) error
	GetSystemTimeZoneList(context.Context, rpcSchema.GetSystemTimeZoneListArgs) (rpcSchema.GetSystemTimeZoneListReply, error)
	GetSystemTimeZone(context.Context, rpcSchema.GetSystemTimeZoneArgs) (rpcSchema.GetSystemTimeZoneReply, error)
	SetSystemTimeZone(context.Context, rpcSchema.SetSystemTimeZoneArgs) error
	Reboot(context.Context, rpcSchema.RebootArgs) error
	Shutdown(context.Context, rpcSchema.ShutdownArgs) error
}

type HostRepo struct {
	RPCClient rpc.Caller
	DB        *database.DB
}

func NewHostRepo(client rpc.Caller, db *database.DB) *HostRepo {
	return &HostRepo{
		RPCClient: client,
		DB:        db,
	}
}

// agentDB returns a scoped DB query filtered by the agent in context.
// 与控制调用（rpc.Call）保持一致：agentID 缺失或格式非法时返回错误，
// 不再静默回退为全表查询，避免跨 Agent 数据聚合泄漏。
func (h *HostRepo) agentDB(ctx context.Context) (*gorm.DB, error) {
	agentID, err := contextx.ResolveAgentID(ctx)
	if err != nil {
		return nil, err
	}
	return h.DB.DB.Where("agent_id = ?", agentID), nil
}

// ── Monitoring queries (local DB) ──

func (h *HostRepo) HostInfo(ctx context.Context, args rpcSchema.HostInfoArgs) (rpcSchema.HostInfoReply, error) {
	q, err := h.agentDB(ctx)
	if err != nil {
		return rpcSchema.HostInfoReply{}, err
	}
	var info model.MonitorHost
	if err := q.Model(&model.MonitorHost{}).Order("timestamp desc").First(&info).Error; err != nil {
		return rpcSchema.HostInfoReply{}, err
	}
	return rpcSchema.HostInfoReply{
		Timestamp:       info.Timestamp.Unix(),
		Uptime:          info.Uptime,
		Hostname:        info.Hostname,
		OS:              info.Os,
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		KernelArch:      info.KernelArch,
	}, nil
}

func (h *HostRepo) CPUInfo(ctx context.Context, args rpcSchema.CPUInfoArgs) (rpcSchema.CPUInfoReply, error) {
	q, err := h.agentDB(ctx)
	if err != nil {
		return rpcSchema.CPUInfoReply{}, err
	}
	var info model.MonitorCPU
	if err := q.Model(&model.MonitorCPU{}).Order("timestamp desc").First(&info).Error; err != nil {
		return rpcSchema.CPUInfoReply{}, err
	}
	return rpcSchema.CPUInfoReply{Percent: info.CPUPercent}, nil
}

func (h *HostRepo) CPUUsage(ctx context.Context, args rpcSchema.CPUUsageArgs) (rpcSchema.CPUUsageReply, error) {
	q, err := h.agentDB(ctx)
	if err != nil {
		return rpcSchema.CPUUsageReply{}, err
	}
	var results []model.MonitorCPU
	if err := q.Model(&model.MonitorCPU{}).
		Where("timestamp > ?", time.Unix(args.StartTime, 0)).
		Order("timestamp asc").Find(&results).Error; err != nil {
		return rpcSchema.CPUUsageReply{}, err
	}
	var list []rpcSchema.Usage
	for _, item := range results {
		list = append(list, rpcSchema.Usage{Timestamp: item.Timestamp.Unix(), Value: item.CPUPercent})
	}
	return rpcSchema.CPUUsageReply{Data: list}, nil
}

func (h *HostRepo) MemInfo(ctx context.Context, args rpcSchema.MemoryInfoArgs) (rpcSchema.MemoryInfoReply, error) {
	q, err := h.agentDB(ctx)
	if err != nil {
		return rpcSchema.MemoryInfoReply{}, err
	}
	var info model.MonitorMemory
	if err := q.Model(&model.MonitorMemory{}).Order("timestamp desc").First(&info).Error; err != nil {
		return rpcSchema.MemoryInfoReply{}, err
	}
	return rpcSchema.MemoryInfoReply{Percent: info.MemPercent, Total: info.MemTotal, Used: info.MemUsed}, nil
}

func (h *HostRepo) MemUsage(ctx context.Context, args rpcSchema.MemoryUsageArgs) (rpcSchema.MemoryUsageReply, error) {
	q, err := h.agentDB(ctx)
	if err != nil {
		return rpcSchema.MemoryUsageReply{}, err
	}
	var results []model.MonitorMemory
	if err := q.Model(&model.MonitorMemory{}).
		Where("timestamp > ?", time.Unix(args.StartTime, 0)).
		Order("timestamp asc").Find(&results).Error; err != nil {
		return rpcSchema.MemoryUsageReply{}, err
	}
	var list []rpcSchema.Usage
	for _, item := range results {
		list = append(list, rpcSchema.Usage{Timestamp: item.Timestamp.Unix(), Value: item.MemPercent})
	}
	return rpcSchema.MemoryUsageReply{Data: list}, nil
}

func (h *HostRepo) DiskInfo(ctx context.Context, args rpcSchema.DiskInfoArgs) (rpcSchema.DiskInfoReply, error) {
	q, err := h.agentDB(ctx)
	if err != nil {
		return rpcSchema.DiskInfoReply{}, err
	}
	var infos []model.MonitorDisk
	if err := q.Model(&model.MonitorDisk{}).Group("device").Order("timestamp desc").Find(&infos).Error; err != nil {
		return rpcSchema.DiskInfoReply{}, err
	}
	var list []rpcSchema.Disk
	for _, info := range infos {
		list = append(list, rpcSchema.Disk{
			Device:      info.Device,
			DiskPercent: info.DiskPercent,
			DiskTotal:   info.DiskTotal,
			DiskUsed:    info.DiskUsed,
		})
	}
	return rpcSchema.DiskInfoReply{Info: list}, nil
}

func (h *HostRepo) DiskUsage(ctx context.Context, args rpcSchema.DiskUsageArgs) (rpcSchema.DiskUsageReply, error) {
	q, err := h.agentDB(ctx)
	if err != nil {
		return rpcSchema.DiskUsageReply{}, err
	}
	var results []model.MonitorDisk
	if err := q.Model(&model.MonitorDisk{}).
		Where("timestamp > ?", time.Unix(args.StartTime, 0)).
		Order("timestamp asc").Find(&results).Error; err != nil {
		return rpcSchema.DiskUsageReply{}, err
	}
	list := make(map[string][]rpcSchema.DiskIO)
	for _, item := range results {
		list[item.Device] = append(list[item.Device], rpcSchema.DiskIO{
			Timestamp: item.Timestamp.Unix(), IORead: item.DiskRead, IOWrite: item.DiskWrite,
		})
	}
	var data []rpcSchema.DiskUsage
	for device, l := range list {
		data = append(data, rpcSchema.DiskUsage{Device: device, Data: l})
	}
	return rpcSchema.DiskUsageReply{Usage: data}, nil
}

func (h *HostRepo) NetUsage(ctx context.Context, args rpcSchema.NetUsageArgs) (rpcSchema.NetUsageReply, error) {
	q, err := h.agentDB(ctx)
	if err != nil {
		return rpcSchema.NetUsageReply{}, err
	}
	var results []model.MonitorNet
	if err := q.Model(&model.MonitorNet{}).
		Where("timestamp > ?", time.Unix(args.StartTime, 0)).
		Order("timestamp asc").Find(&results).Error; err != nil {
		return rpcSchema.NetUsageReply{}, err
	}
	list := make(map[string][]rpcSchema.NetIO)
	for _, item := range results {
		list[item.Ethernet] = append(list[item.Ethernet], rpcSchema.NetIO{
			Timestamp: item.Timestamp.Unix(), BytesSent: item.NetSend, BytesRecv: item.NetRecv,
		})
	}
	var data []rpcSchema.NetUsage
	for eth, l := range list {
		data = append(data, rpcSchema.NetUsage{Ethernet: eth, Data: l})
	}
	return rpcSchema.NetUsageReply{Usage: data}, nil
}

// ── Operational commands (RPC to Agent) ──

func (h *HostRepo) FilesSearch(ctx context.Context, args rpcSchema.FilesSearchArgs) (rpcSchema.FilesSearchReply, error) {
	var reply rpcSchema.FilesSearchReply
	err := h.RPCClient.Call(ctx, "FilesSearch", args, &reply)
	return reply, err
}

func (h *HostRepo) FileUpload(ctx context.Context, args rpcSchema.FileUploadArgs) error {
	var reply rpcSchema.FileUploadReply
	return h.RPCClient.Call(ctx, "FileUpload", args, &reply)
}

func (h *HostRepo) FileDownload(ctx context.Context, args rpcSchema.FileDownloadArgs) (rpcSchema.FileDownloadReply, error) {
	var reply rpcSchema.FileDownloadReply
	err := h.RPCClient.Call(ctx, "FileDownload", args, &reply)
	return reply, err
}

func (h *HostRepo) FileDelete(ctx context.Context, args rpcSchema.FileDeleteArgs) error {
	var reply rpcSchema.FileDeleteReply
	return h.RPCClient.Call(ctx, "FileDelete", args, &reply)
}

func (h *HostRepo) FileCreate(ctx context.Context, args rpcSchema.FileCreateArgs) error {
	var reply rpcSchema.FileCreateReply
	return h.RPCClient.Call(ctx, "FileCreate", args, &reply)
}

func (h *HostRepo) FolderCreate(ctx context.Context, args rpcSchema.FolderCreateArgs) error {
	var reply rpcSchema.FolderCreateReply
	return h.RPCClient.Call(ctx, "FolderCreate", args, &reply)
}

func (h *HostRepo) GetDNSSettings(ctx context.Context, args rpcSchema.GetDNSArgs) (rpcSchema.GetDNSReply, error) {
	var reply rpcSchema.GetDNSReply
	err := h.RPCClient.Call(ctx, "GetDNS", args, &reply)
	return reply, err
}

func (h *HostRepo) SetDNSSettings(ctx context.Context, args rpcSchema.SetDNSArgs) error {
	var reply rpcSchema.SetDNSReply
	return h.RPCClient.Call(ctx, "SetDNS", args, &reply)
}

func (h *HostRepo) GetSystemTime(ctx context.Context, args rpcSchema.GetSystemTimeArgs) (rpcSchema.GetSystemTimeReply, error) {
	var reply rpcSchema.GetSystemTimeReply
	err := h.RPCClient.Call(ctx, "GetSystemTime", args, &reply)
	return reply, err
}

func (h *HostRepo) SetSystemTime(ctx context.Context, args rpcSchema.SetSystemTimeArgs) error {
	var reply rpcSchema.SetSystemTimeReply
	return h.RPCClient.Call(ctx, "SetSystemTime", args, &reply)
}

func (h *HostRepo) GetSystemTimeZoneList(ctx context.Context, args rpcSchema.GetSystemTimeZoneListArgs) (rpcSchema.GetSystemTimeZoneListReply, error) {
	var reply rpcSchema.GetSystemTimeZoneListReply
	err := h.RPCClient.Call(ctx, "GetSystemTimeZoneList", args, &reply)
	return reply, err
}

func (h *HostRepo) GetSystemTimeZone(ctx context.Context, args rpcSchema.GetSystemTimeZoneArgs) (rpcSchema.GetSystemTimeZoneReply, error) {
	var reply rpcSchema.GetSystemTimeZoneReply
	err := h.RPCClient.Call(ctx, "GetSystemTimeZone", args, &reply)
	return reply, err
}

func (h *HostRepo) SetSystemTimeZone(ctx context.Context, args rpcSchema.SetSystemTimeZoneArgs) error {
	var reply rpcSchema.SetSystemTimeZoneReply
	return h.RPCClient.Call(ctx, "SetSystemTimeZone", args, &reply)
}

func (h *HostRepo) Reboot(ctx context.Context, args rpcSchema.RebootArgs) error {
	var reply rpcSchema.RebootReply
	return h.RPCClient.Call(ctx, "Reboot", args, &reply)
}

func (h *HostRepo) Shutdown(ctx context.Context, args rpcSchema.ShutdownArgs) error {
	var reply rpcSchema.ShutdownReply
	return h.RPCClient.Call(ctx, "Shutdown", args, &reply)
}
