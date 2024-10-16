// Package rpc
// Date:   2024/10/15 10:57
// Author: Amu
// Description:
package rpc

import (
	"amvector/service/model"
	rpcSchema "common/rpc/schema"
	"context"
	"time"
)

func (s *Service) CPUAlarmQuery(ctx context.Context, args rpcSchema.CPUAlarmQueryArgs, reply *rpcSchema.CPUAlarmQueryReply) error {
	var result []model.CPU
	if err := s.DB.Model(&model.CPU{}).Where("timestamp > ?", time.Unix(args.StartTime, 0)).Where("timestamp <= ?", time.Unix(args.EndTime, 0)).Order("timestamp desc").Find(&result).Error; err != nil {
		return err
	}
	var list []rpcSchema.Usage
	for _, item := range result {
		list = append(list, rpcSchema.Usage{
			Timestamp: item.Timestamp.Unix(),
			Value:     item.CPUPercent,
		})
	}
	reply.Data = list
	return nil
}

func (s *Service) MemoryAlarmQuery(ctx context.Context, args rpcSchema.MemoryAlarmQueryArgs, reply *rpcSchema.MemoryAlarmQueryReply) error {
	var result []model.Memory
	if err := s.DB.Model(&model.Memory{}).Where("timestamp > ?", time.Unix(args.StartTime, 0)).Where("timestamp <= ?", time.Unix(args.EndTime, 0)).Or("timestamp desc").Find(&result).Error; err != nil {
		return err
	}
	var list []rpcSchema.Usage
	for _, item := range result {
		list = append(list, rpcSchema.Usage{
			Timestamp: item.Timestamp.Unix(),
			Value:     item.MemPercent,
		})
	}
	reply.Data = list
	return nil
}

func (s *Service) DiskAlarmQuery(ctx context.Context, args rpcSchema.DiskAlarmQueryArgs, reply *rpcSchema.DiskAlarmQueryReply) error {
	var results []model.Disk
	if err := s.DB.Model(&model.Disk{}).Where("timestamp > ?", time.Unix(args.StartTime, 0)).Where("timestamp <= ?", time.Unix(args.EndTime, 0)).Order("timestamp asc").Find(&results).Error; err != nil {
		return err
	}
	var list []rpcSchema.Disk
	for _, item := range results {
		list = append(list, rpcSchema.Disk{
			Device:      item.Device,
			DiskPercent: item.DiskPercent,
			DiskTotal:   item.DiskTotal,
			DiskUsed:    item.DiskUsed,
		})
	}
	reply.Data = list
	return nil
}

func (s *Service) ServiceAlarmQuery(ctx context.Context, args rpcSchema.ServiceAlarmQueryArgs, reply *rpcSchema.ServiceAlarmQueryReply) error {
	var containers []model.Container
	if err := s.DB.Model(&model.Container{}).Order("created_at desc").Find(&containers).Error; err != nil {
		return err
	}
	var results []rpcSchema.Container
	for _, container := range containers {
		results = append(results, rpcSchema.Container{
			Timestamp:   container.Timestamp,
			ContainerID: container.ContainerID,
			Name:        container.Name,
			Image:       container.Image,
			IP:          container.IP,
			Ports:       container.Ports,
			State:       container.State,
			Uptime:      container.Uptime,
			CPUPercent:  container.CPUPercent,
			MemPercent:  container.MemPercent,
			MemUsage:    container.MemUsage,
			MemLimit:    container.MemLimit,
			Labels:      container.Labels,
		})
	}
	reply.Data = results
	return nil
}
