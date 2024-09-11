// Package repository
// Date: 2024/06/11 19:26:34
// Author: Amu
// Description:
package repository

import (
	"amprobe/pkg/rpc"
	rpcSchema "common/rpc/schema"
	"context"

	"github.com/google/wire"
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
	RPCClient *rpc.Client
}

func NewHostRepo(client *rpc.Client) *HostRepo {
	return &HostRepo{
		RPCClient: client,
	}
}

func (h *HostRepo) HostInfo(ctx context.Context, args rpcSchema.HostInfoArgs) (rpcSchema.HostInfoReply, error) {
	var reply rpcSchema.HostInfoReply
	err := h.RPCClient.Call(ctx, "HostInfo", args, &reply)
	return reply, err
}

func (h *HostRepo) CPUInfo(ctx context.Context, args rpcSchema.CPUInfoArgs) (rpcSchema.CPUInfoReply, error) {
	var reply rpcSchema.CPUInfoReply
	err := h.RPCClient.Call(ctx, "CPUInfo", args, &reply)
	return reply, err
}

func (h *HostRepo) CPUUsage(ctx context.Context, args rpcSchema.CPUUsageArgs) (rpcSchema.CPUUsageReply, error) {
	var reply rpcSchema.CPUUsageReply
	err := h.RPCClient.Call(ctx, "CPUUsage", args, &reply)
	return reply, err
}

func (h *HostRepo) MemInfo(ctx context.Context, args rpcSchema.MemoryInfoArgs) (rpcSchema.MemoryInfoReply, error) {
	var reply rpcSchema.MemoryInfoReply
	err := h.RPCClient.Call(ctx, "MemoryInfo", args, &reply)
	return reply, err
}

func (h *HostRepo) MemUsage(ctx context.Context, args rpcSchema.MemoryUsageArgs) (rpcSchema.MemoryUsageReply, error) {
	var reply rpcSchema.MemoryUsageReply
	err := h.RPCClient.Call(ctx, "MemoryUsage", args, &reply)
	return reply, err
}

func (h *HostRepo) DiskInfo(ctx context.Context, args rpcSchema.DiskInfoArgs) (rpcSchema.DiskInfoReply, error) {
	var reply rpcSchema.DiskInfoReply
	err := h.RPCClient.Call(ctx, "DiskInfo", args, &reply)
	return reply, err
}

func (h *HostRepo) DiskUsage(ctx context.Context, args rpcSchema.DiskUsageArgs) (rpcSchema.DiskUsageReply, error) {
	var reply rpcSchema.DiskUsageReply
	err := h.RPCClient.Call(ctx, "DiskUsage", args, &reply)
	return reply, err
}

func (h *HostRepo) NetUsage(ctx context.Context, args rpcSchema.NetUsageArgs) (rpcSchema.NetUsageReply, error) {
	var reply rpcSchema.NetUsageReply
	err := h.RPCClient.Call(ctx, "NetUsage", args, &reply)
	return reply, err
}

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
