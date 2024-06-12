// Package rpc
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package rpc

import (
	"context"
	"time"

	"github.com/amuluze/amvector/service/model"
	"github.com/amuluze/amvector/service/schema"
)

func (s *Service) HostInfo(ctx context.Context, args schema.HostArgs, reply *model.Host) error {
	if err := s.DB.Order("timestamp desc").Take(reply).Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) CPUInfo(ctx context.Context, args schema.CPUArgs, reply *model.CPU) error {
	if err := s.DB.Order("timestamp desc").Take(reply).Error; err != nil {
		return err
	}
	return nil
}
func (s *Service) CPUUsage(ctx context.Context, args schema.CPUUsageArgs, reply *[]model.CPU) error {
	if err := s.DB.Model(&model.CPU{}).Where("timestamp > ? and timestamp < ?", time.Unix(args.StartTime, 0), time.Unix(args.EndTime, 0)).Order("timestamp asc").Find(reply).Error; err != nil {
		return nil
	}
	return nil
}
func (s *Service) MemInfo(ctx context.Context, args schema.MemoryArgs, reply *model.Memory) error {
	if err := s.DB.Order("timestamp desc").Take(reply).Error; err != nil {
		return err
	}
	return nil
}
func (s *Service) MemUsage(ctx context.Context, args schema.MemoryUsageArgs, reply *[]model.Memory) error {
	if err := s.DB.Model(&model.Memory{}).Where("timestamp > ?", time.Unix(args.StartTime, 0)).Order("timestamp asc").Find(reply).Error; err != nil {
		return nil
	}
	return nil
}
func (s *Service) DiskInfo(ctx context.Context, args schema.DiskArgs, reply *[]model.Disk) error {
	if err := s.DB.Model(&model.Disk{}).Group("device").Order("timestamp desc").Find(reply).Error; err != nil {
		return err
	}
	return nil
}
func (s *Service) DiskUsage(ctx context.Context, args schema.DiskUsageArgs, reply *[]model.Disk) error {
	if err := s.DB.Model(&model.Disk{}).Where("timestamp > ?", time.Unix(args.StartTime, 0)).Order("timestamp asc").Find(reply).Error; err != nil {
		return nil
	}
	return nil
}
func (s *Service) NetUsage(ctx context.Context, args schema.NetworkUsageArgs, reply *[]model.Net) error {
	if err := s.DB.Model(&model.Net{}).Where("timestamp > ?", time.Unix(args.StartTime, 0)).Order("timestamp asc").Find(reply).Error; err != nil {
		return nil
	}
	return nil
}
