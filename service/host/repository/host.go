// Package repository
// Date: 2024/06/11 19:26:34
// Author: Amu
// Description:
package repository

import (
	"context"
	"time"

	"amprobe/pkg/command"
	"amprobe/pkg/database"
	"amprobe/pkg/timectl"
	"amprobe/service/model"
	"amprobe/service/schema"

	"github.com/google/wire"

	"github.com/docker/docker/libnetwork/resolvconf"
)

var HostRepoSet = wire.NewSet(NewHostRepo, wire.Bind(new(IHostRepo), new(*HostRepo)))

var _ IHostRepo = (*HostRepo)(nil)

type IHostRepo interface {
	CPUInfo(context.Context) (model.CPU, error)
	CPUUsage(context.Context, schema.CPUUsageArgs) ([]model.CPU, error)
	MemInfo(context.Context) (model.Memory, error)
	MemUsage(context.Context, schema.MemoryUsageArgs) ([]model.Memory, error)
	DiskInfo(context.Context) ([]model.Disk, error)
	DiskUsage(context.Context, schema.DiskUsageArgs) ([]model.Disk, error)
	NetUsage(context.Context, schema.NetUsageArgs) ([]model.Net, error)

	GetDNSSettings(context.Context) ([]string, error)
	SetDNSSettings(context.Context, schema.SetDNSSettingsArgs) error
	GetSystemTime(context.Context) (string, error)
	SetSystemTime(context.Context, schema.SetSystemTimeArgs) error
	GetSystemTimeZoneList(context.Context) ([]string, error)
	GetSystemTimeZone(context.Context) (string, error)
	SetSystemTimeZone(context.Context, schema.SetSystemTimeZoneArgs) error
	Reboot(context.Context) error
	Shutdown(context.Context) error
}

type HostRepo struct {
	db *database.DB
}

func NewHostRepo(db *database.DB) *HostRepo {
	return &HostRepo{
		db: db,
	}
}

func (h *HostRepo) CPUInfo(ctx context.Context) (model.CPU, error) {
	var reply model.CPU
	if err := h.db.Model(&model.CPU{}).Order("timestamp desc").First(&reply).Error; err != nil {
		return reply, err
	}
	return reply, nil
}

func (h *HostRepo) CPUUsage(ctx context.Context, args schema.CPUUsageArgs) ([]model.CPU, error) {
	var reply []model.CPU
	if err := h.db.Model(&model.CPU{}).Where("timestamp >= ?", time.Unix(args.StartTime, 0)).Order("timestamp asc").Find(&reply).Error; err != nil {
		return reply, err
	}
	return reply, nil
}

func (h *HostRepo) MemInfo(ctx context.Context) (model.Memory, error) {
	var reply model.Memory
	if err := h.db.Model(&model.Memory{}).Order("timestamp desc").First(&reply).Error; err != nil {
		return reply, err
	}
	return reply, nil
}

func (h *HostRepo) MemUsage(ctx context.Context, args schema.MemoryUsageArgs) ([]model.Memory, error) {
	var reply []model.Memory
	if err := h.db.Model(&model.Memory{}).Where("timestamp >= ?", time.Unix(args.StartTime, 0)).Order("timestamp asc").Find(&reply).Error; err != nil {
		return reply, err
	}
	return reply, nil
}

func (h *HostRepo) DiskInfo(ctx context.Context) ([]model.Disk, error) {
	var reply []model.Disk
	if err := h.db.Model(&model.Disk{}).Group("device").Order("timestamp desc").Find(&reply).Error; err != nil {
		return reply, err
	}
	return reply, nil
}

func (h *HostRepo) DiskUsage(ctx context.Context, args schema.DiskUsageArgs) ([]model.Disk, error) {
	var reply []model.Disk
	if err := h.db.Model(&model.Disk{}).Where("timestamp >= ?", time.Unix(args.StartTime, 0)).Order("timestamp asc").Find(&reply).Error; err != nil {
		return reply, err
	}
	return reply, nil
}

func (h *HostRepo) NetUsage(ctx context.Context, args schema.NetUsageArgs) ([]model.Net, error) {
	var reply []model.Net
	if err := h.db.Model(&model.Net{}).Where("timestamp >= ?", time.Unix(args.StartTime, 0)).Order("timestamp asc").Find(&reply).Error; err != nil {
		return reply, err
	}
	return reply, nil
}

func (h *HostRepo) GetDNSSettings(ctx context.Context) ([]string, error) {
	resolvConf, err := resolvconf.Get()
	if err != nil {
		return nil, err
	}
	return resolvconf.GetNameservers(resolvConf.Content, resolvconf.IP), nil
}

func (h *HostRepo) SetDNSSettings(ctx context.Context, args schema.SetDNSSettingsArgs) error {
	_, err := resolvconf.Build(resolvconf.Path(), args.DNS, []string{}, []string{})
	return err

}

func (h *HostRepo) GetSystemTime(ctx context.Context) (string, error) {
	systemTime, err := timectl.GetTime(ctx)
	if err != nil {
		return "", err
	}
	return time.Unix(systemTime, 0).Format("2006-01-02 15:04:05"), nil
}

func (h *HostRepo) SetSystemTime(ctx context.Context, args schema.SetSystemTimeArgs) error {
	tt, err := time.Parse(args.SystemTime, "2006-01-02 15:04:05")
	if err != nil {
		return err
	}
	return timectl.SetTime(ctx, tt.Unix())
}

func (h *HostRepo) GetSystemTimeZoneList(ctx context.Context) ([]string, error) {
	return timectl.GetTimeZoneList(ctx)
}

func (h *HostRepo) GetSystemTimeZone(ctx context.Context) (string, error) {
	return timectl.GetTimeZone(ctx)
}

func (h *HostRepo) SetSystemTimeZone(ctx context.Context, args schema.SetSystemTimeZoneArgs) error {
	return timectl.SetTimeZone(ctx, args.SystemTimeZone)
}

func (h *HostRepo) Reboot(ctx context.Context) error {
	if _, err := command.RunCommand(ctx, "reboot"); err != nil {
		return err
	}
	return nil
}

func (h *HostRepo) Shutdown(ctx context.Context) error {
	if _, err := command.RunCommand(ctx, "shutdown"); err != nil {
		return err
	}
	return nil
}
