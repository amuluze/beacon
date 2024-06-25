// Package rpc
// Date: 2024/6/24 09:30
// Author: Amu
// Description:
package rpc

import (
	"context"
	"github.com/amuluze/amvector/pkg/timectl"

	"github.com/amuluze/amvector/service/schema"
	"github.com/docker/docker/libnetwork/resolvconf"
)

// GetDNS 获取 DNS 配置
func (s *Service) GetDNS(ctx context.Context, args schema.GetDNSArgs, reply *schema.GetDNSReply) error {
	resolvConf, err := resolvconf.Get()
	if err != nil {
		return err
	}
	reply.DNS = resolvconf.GetNameservers(resolvConf.Content, resolvconf.IP)
	return nil
}

// SetDNS 更新 DNS 配置
func (s *Service) SetDNS(ctx context.Context, args schema.SetDNSArgs, reply *schema.SetDNSReply) error {
	_, err := resolvconf.Build(resolvconf.Path(), args.DNS, []string{}, []string{})
	if err != nil {
		return err
	}
	return nil
}

// GetSystemTime 获取系统时间
func (s *Service) GetSystemTime(ctx context.Context, args schema.GetSystemTimeArgs, reply *schema.GetSystemTimeReply) error {
	time, err := timectl.GetTime(ctx)
	if err != nil {
		return err
	}
	reply.SystemTime = time
	return nil
}

// SetSystemTime 设置系统时间
func (s *Service) SetSystemTime(ctx context.Context, args schema.SetSystemTimeArgs, reply *schema.SetSystemTimeReply) error {
	err := timectl.SetTime(ctx, args.SystemTime)
	if err != nil {
		return err
	}
	return nil
}

// GetSystemTimeZone 获取系统时区
func (s *Service) GetSystemTimeZone(ctx context.Context, args schema.GetSystemTimeZoneArgs, reply *schema.GetSystemTimeZoneReply) error {
	tz, err := timectl.GetTimeZone(ctx)
	if err != nil {
		return err
	}
	reply.SystemTimeZone = tz
	return nil
}

// SetSystemTimeZone 设置系统时区
func (s *Service) SetSystemTimeZone(ctx context.Context, args schema.SetSystemTimeZoneArgs, reply *schema.SetSystemTimeZoneReply) error {
	if err := timectl.SetTimeZone(ctx, args.SystemTimeZone); err != nil {
		return err
	}
	return nil
}
