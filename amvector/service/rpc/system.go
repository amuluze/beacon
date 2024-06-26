// Package rpc
// Date: 2024/6/24 09:30
// Author: Amu
// Description:
package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/amuluze/amvector/pkg/utils"
	"os"

	"github.com/amuluze/amvector/pkg/timectl"

	"github.com/amuluze/amvector/service/constants"
	"github.com/amuluze/amvector/service/schema"
	"github.com/docker/docker/libnetwork/resolvconf"
)

// Reboot 重启
func (s *Service) Reboot(ctx context.Context, args schema.RebootArgs, reply *schema.RebootReply) error {
	if _, err := utils.RunCommand(ctx, "reboot"); err != nil {
		return errors.New("failed to reboot: " + err.Error())
	}
	return nil
}

// Shutdown 关机
func (s *Service) Shutdown(ctx context.Context, args schema.ShutdownArgs, reply *schema.ShutdownReply) error {
	if _, err := utils.RunCommand(ctx, "shutdown"); err != nil {
		return errors.New("failed to shutdown: " + err.Error())
	}
	return nil
}

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

// GetDockerRegistryMirrors 读取 Docker daemon.json 配置文件
func (s *Service) GetDockerRegistryMirrors(ctx context.Context, args schema.GetDockerRegistryMirrorsArgs, reply *schema.GetDockerRegistryMirrorsReply) error {
	if _, err := os.Stat(constants.DaemonJsonPath); err != nil {
		return err
	}
	file, err := os.ReadFile(constants.DaemonJsonPath)
	if err != nil {
		return err
	}
	daemonMap := make(map[string]interface{})
	if err := json.Unmarshal(file, &daemonMap); err != nil {
		return err
	}
	if daemonMap["registry-mirrors"] != nil {
		reply.Mirrors = daemonMap["registry-mirrors"].([]string)
	}
	return nil
}

// SetDockerRegistryMirrors 更新 Docker daemon.json 配置文件
func (s *Service) SetDockerRegistryMirrors(ctx context.Context, args schema.SetDockerRegistryMirrorsArgs, reply *schema.SetDockerRegistryMirrorsReply) error {
	// 判断 constants.DaemonJsonPath 文件是否存在
	if _, err := os.Stat(constants.DaemonJsonPath); err != nil {
		// 创建 constants.DaemonJsonPath
		daemonMap := map[string]interface{}{
			"registry-mirrors": args.Mirrors,
		}
		// 将 daemonMap 写入 constants.DaemonJsonPath 文件
		setting, err := json.Marshal(daemonMap)
		if err != nil {
			return err
		}
		if err := os.WriteFile(constants.DaemonJsonPath, setting, 0644); err != nil {
			return err
		}
	} else {
		file, err := os.ReadFile(constants.DaemonJsonPath)
		if err != nil {
			return err
		}
		daemonMap := make(map[string]interface{})
		if err := json.Unmarshal(file, &daemonMap); err != nil {
			return err
		}
		daemonMap["registry-mirrors"] = args.Mirrors
		setting, err := json.Marshal(daemonMap)
		if err != nil {
			return err
		}
		if err := os.WriteFile(constants.DaemonJsonPath, setting, 0644); err != nil {
			return err
		}
	}
	// TODO: 重启 Docker 服务
	return nil
}
