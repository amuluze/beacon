// Package repository
// Date: 2024/3/6 12:53
// Author: Amu
// Description:
package repository

import (
	"context"
	"time"

	"github.com/amuluze/amprobe/pkg/psutil"
	"github.com/amuluze/amprobe/service/model"
	"github.com/amuluze/amprobe/service/schema"
	"github.com/amuluze/amutool/database"
	"github.com/google/wire"
)

var HostRepoSet = wire.NewSet(NewHostRepo, wire.Bind(new(IHostRepo), new(*HostRepo)))

type IHostRepo interface {
	SystemUptime(ctx context.Context) (string, error)
	CPUInfo(ctx context.Context) (model.CPU, error)
	CPUUsage(ctx context.Context, args schema.CPUUsageArgs) ([]model.CPU, error)
	MemInfo(ctx context.Context) (model.Memory, error)
	MemUsage(ctx context.Context, args schema.MemoryUsageArgs) ([]model.Memory, error)
	DiskInfo(ctx context.Context) ([]model.Disk, error)
	DiskUsage(ctx context.Context, args schema.DiskUsageArgs) ([]model.Disk, error)
	NetUsage(ctx context.Context, args schema.NetworkUsageArgs) ([]model.Net, error)
}

type HostRepo struct {
	DB *database.DB
}

func NewHostRepo(db *database.DB) *HostRepo {
	return &HostRepo{DB: db}
}

func (h HostRepo) SystemUptime(ctx context.Context) (string, error) {
	uptime, err := psutil.GetSystemUpTime()
	if err != nil {
		return "", err
	}
	return uptime, nil
}

func (h HostRepo) CPUInfo(ctx context.Context) (model.CPU, error) {
	var cpuInfo model.CPU
	if err := h.DB.Order("timestamp desc").Take(&cpuInfo).Error; err != nil {
		return cpuInfo, err
	}
	return cpuInfo, nil
}

func (h HostRepo) CPUUsage(ctx context.Context, args schema.CPUUsageArgs) ([]model.CPU, error) {
	var cpuInfos []model.CPU
	if err := h.DB.Model(&model.CPU{}).Where("timestamp > ? and timestamp < ?", time.Unix(args.StartTime, 0), time.Unix(args.EndTime, 0)).Order("timestamp asc").Find(&cpuInfos).Error; err != nil {
		return cpuInfos, err
	}
	return cpuInfos, nil
}

func (h HostRepo) MemInfo(ctx context.Context) (model.Memory, error) {
	var memInfo model.Memory
	if err := h.DB.Order("timestamp desc").Take(&memInfo).Error; err != nil {
		return memInfo, err
	}
	return memInfo, nil
}

func (h HostRepo) MemUsage(ctx context.Context, args schema.MemoryUsageArgs) ([]model.Memory, error) {
	var memInfos []model.Memory
	if err := h.DB.Model(&model.Memory{}).Where("timestamp > ?", time.Unix(args.StartTime, 0)).Order("timestamp asc").Find(&memInfos).Error; err != nil {
		return memInfos, err
	}
	return memInfos, nil
}

func (h HostRepo) DiskInfo(ctx context.Context) ([]model.Disk, error) {
	var diskInfos []model.Disk
	if err := h.DB.Model(&model.Disk{}).Group("device").Order("timestamp desc").Find(&diskInfos).Error; err != nil {
		return diskInfos, err
	}
	return diskInfos, nil
}

func (h HostRepo) DiskUsage(ctx context.Context, args schema.DiskUsageArgs) ([]model.Disk, error) {
	var diskInfos []model.Disk
	if err := h.DB.Model(&model.Disk{}).Where("timestamp > ?", time.Unix(args.StartTime, 0)).Order("timestamp asc").Find(&diskInfos).Error; err != nil {
		return diskInfos, err
	}
	return diskInfos, nil
}

func (h HostRepo) NetUsage(ctx context.Context, args schema.NetworkUsageArgs) ([]model.Net, error) {
	var netInfos []model.Net
	if err := h.DB.Model(&model.Net{}).Where("timestamp > ?", time.Unix(args.StartTime, 0)).Order("timestamp asc").Find(&netInfos).Error; err != nil {
		return netInfos, err
	}
	return netInfos, nil
}
