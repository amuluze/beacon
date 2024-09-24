// Package rpc
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package rpc

import (
	rpcSchema "common/rpc/schema"
	"context"
	"time"

	"amvector/service/model"
)

func (s *Service) HostInfo(ctx context.Context, args rpcSchema.HostInfoArgs, reply *rpcSchema.HostInfoReply) error {
	var info model.Host
	if err := s.DB.Model(&model.Host{}).Order("timestamp desc").First(&info).Error; err != nil {
		return err
	}
	reply.Timestamp = info.Timestamp.Unix()
	reply.Uptime = info.Uptime
	reply.Hostname = info.Hostname
	reply.OS = info.Os
	reply.Platform = info.Platform
	reply.PlatformVersion = info.PlatformVersion
	reply.KernelVersion = info.KernelVersion
	reply.KernelArch = info.KernelArch
	return nil
}

func (s *Service) CPUInfo(ctx context.Context, args rpcSchema.CPUInfoArgs, reply *rpcSchema.CPUInfoReply) error {
	var info model.CPU
	if err := s.DB.Model(&model.CPU{}).Order("timestamp desc").First(&info).Error; err != nil {
		return err
	}
	reply.Percent = info.CPUPercent
	return nil
}

func (s *Service) CPUUsage(ctx context.Context, args rpcSchema.CPUUsageArgs, reply *rpcSchema.CPUUsageReply) error {
	var results []model.CPU
	if err := s.DB.Model(&model.CPU{}).Where("timestamp > ?", time.Unix(args.StartTime, 0)).Order("timestamp asc").Find(&results).Error; err != nil {
		return nil
	}
	var list []rpcSchema.Usage
	for _, item := range results {
		list = append(list, rpcSchema.Usage{
			Timestamp: item.Timestamp.Unix(),
			Value:     item.CPUPercent,
		})
	}
	reply.Data = list
	return nil
}

func (s *Service) MemoryInfo(ctx context.Context, args rpcSchema.MemoryInfoArgs, reply *rpcSchema.MemoryInfoReply) error {
	var info model.Memory
	if err := s.DB.Model(&model.Memory{}).Order("timestamp desc").First(&info).Error; err != nil {
		return err
	}
	reply.Percent = info.MemPercent
	reply.Total = info.MemTotal
	reply.Used = info.MemUsed
	return nil
}

func (s *Service) MemoryUsage(ctx context.Context, args rpcSchema.MemoryUsageArgs, reply *rpcSchema.MemoryUsageReply) error {
	var results []model.Memory
	if err := s.DB.Model(&model.Memory{}).Where("timestamp > ?", time.Unix(args.StartTime, 0)).Order("timestamp asc").Find(&results).Error; err != nil {
		return nil
	}
	var list []rpcSchema.Usage
	for _, item := range results {
		list = append(list, rpcSchema.Usage{
			Timestamp: item.Timestamp.Unix(),
			Value:     item.MemPercent,
		})
	}
	reply.Data = list
	return nil
}

func (s *Service) DiskInfo(ctx context.Context, args rpcSchema.DiskInfoArgs, reply *rpcSchema.DiskInfoReply) error {
	var infos []model.Disk
	if err := s.DB.Model(&model.Disk{}).Group("device").Order("timestamp desc").Find(&infos).Error; err != nil {
		return err
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
	reply.Info = list
	return nil
}

func (s *Service) DiskUsage(ctx context.Context, args rpcSchema.DiskUsageArgs, reply *rpcSchema.DiskUsageReply) error {
	var results []model.Disk
	if err := s.DB.Model(&model.Disk{}).Where("timestamp > ?", time.Unix(args.StartTime, 0)).Order("timestamp asc").Find(&results).Error; err != nil {
		return nil
	}
	var data []rpcSchema.DiskUsage
	list := make(map[string][]rpcSchema.DiskIO)
	for _, item := range results {
		list[item.Device] = append(list[item.Device], rpcSchema.DiskIO{
			Timestamp: item.Timestamp.Unix(),
			IORead:    item.DiskRead,
			IOWrite:   item.DiskWrite,
		})
	}
	for device, l := range list {
		data = append(data, rpcSchema.DiskUsage{
			Device: device,
			Data:   l,
		})
	}
	reply.Usage = data
	return nil
}

func (s *Service) NetUsage(ctx context.Context, args rpcSchema.NetUsageArgs, reply *rpcSchema.NetUsageReply) error {
	var results []model.Net
	if err := s.DB.Model(&model.Net{}).Where("timestamp > ?", time.Unix(args.StartTime, 0)).Order("timestamp asc").Find(&results).Error; err != nil {
		return nil
	}
	var data []rpcSchema.NetUsage
	list := make(map[string][]rpcSchema.NetIO)
	for _, item := range results {
		list[item.Ethernet] = append(list[item.Ethernet], rpcSchema.NetIO{
			Timestamp: item.Timestamp.Unix(),
			BytesRecv: item.NetRecv,
			BytesSent: item.NetSend,
		})
	}
	for eth, l := range list {
		data = append(data, rpcSchema.NetUsage{
			Ethernet: eth,
			Data:     l,
		})
	}
	reply.Usage = data
	return nil
}
