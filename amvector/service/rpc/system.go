// Package rpc
// Date: 2024/6/24 09:30
// Author: Amu
// Description:
package rpc

import (
	rpcSchema "common/rpc/schema"
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"amvector/pkg/timectl"
	"amvector/pkg/utils"

	"github.com/docker/docker/libnetwork/resolvconf"
)

const (
	DaemonJsonPath = "/etc/docker/daemon.json"
)

// Reboot 重启
func (s *Service) Reboot(ctx context.Context, args rpcSchema.RebootArgs, reply *rpcSchema.RebootReply) error {
	if _, err := utils.RunCommand(ctx, "reboot"); err != nil {
		return errors.New("failed to reboot: " + err.Error())
	}
	return nil
}

// Shutdown 关机
func (s *Service) Shutdown(ctx context.Context, args rpcSchema.ShutdownArgs, reply *rpcSchema.ShutdownReply) error {
	if _, err := utils.RunCommand(ctx, "shutdown"); err != nil {
		return errors.New("failed to shutdown: " + err.Error())
	}
	return nil
}

// GetDNS 获取 DNS 配置
func (s *Service) GetDNS(ctx context.Context, args rpcSchema.GetDNSArgs, reply *rpcSchema.GetDNSReply) error {
	resolvConf, err := resolvconf.Get()
	if err != nil {
		return err
	}
	reply.DNS = resolvconf.GetNameservers(resolvConf.Content, resolvconf.IP)
	return nil
}

// SetDNS 更新 DNS 配置
func (s *Service) SetDNS(ctx context.Context, args rpcSchema.SetDNSArgs, reply *rpcSchema.SetDNSReply) error {
	_, err := resolvconf.Build(resolvconf.Path(), args.DNS, []string{}, []string{})
	if err != nil {
		return err
	}
	return nil
}

// GetSystemTime 获取系统时间
func (s *Service) GetSystemTime(ctx context.Context, args rpcSchema.GetSystemTimeArgs, reply *rpcSchema.GetSystemTimeReply) error {
	systemTime, err := timectl.GetTime(ctx)
	if err != nil {
		return err
	}
	reply.SystemTime = time.Unix(systemTime, 0).Format("2006-01-02 15:04:05")
	return nil
}

// SetSystemTime 设置系统时间
func (s *Service) SetSystemTime(ctx context.Context, args rpcSchema.SetSystemTimeArgs, reply *rpcSchema.SetSystemTimeReply) error {
	tt, err := time.Parse(args.SystemTime, "2006-01-02 15:04:05")
	if err != nil {
		return err
	}
	err = timectl.SetTime(ctx, tt.Unix())
	if err != nil {
		return err
	}
	return nil
}

// GetSystemTimeZoneList
func (s *Service) GetSystemTimeZoneList(ctx context.Context, args rpcSchema.GetSystemTimeZoneListArgs, reply *rpcSchema.GetSystemTimeZoneListReply) error {
	tzList, err := timectl.GetTimeZoneList(ctx)
	if err != nil {
		return err
	}
	reply.SystemTimeZoneList = tzList
	return nil
}

// GetSystemTimeZone 获取系统时区
func (s *Service) GetSystemTimeZone(ctx context.Context, args rpcSchema.GetSystemTimeZoneArgs, reply *rpcSchema.GetSystemTimeZoneReply) error {
	tz, err := timectl.GetTimeZone(ctx)
	if err != nil {
		return err
	}
	reply.SystemTimeZone = tz
	return nil
}

// SetSystemTimeZone 设置系统时区
func (s *Service) SetSystemTimeZone(ctx context.Context, args rpcSchema.SetSystemTimeZoneArgs, reply *rpcSchema.SetSystemTimeZoneReply) error {
	if err := timectl.SetTimeZone(ctx, args.SystemTimeZone); err != nil {
		return err
	}
	return nil
}

// GetDockerRegistryMirrors 读取 Docker daemon.json 配置文件
func (s *Service) GetDockerRegistryMirrors(ctx context.Context, args rpcSchema.GetDockerRegistryMirrorsArgs, reply *rpcSchema.GetDockerRegistryMirrorsReply) error {
	if _, err := os.Stat(DaemonJsonPath); err != nil {
		return err
	}
	file, err := os.ReadFile(DaemonJsonPath)
	if err != nil {
		return err
	}
	daemonMap := make(map[string]interface{})
	if err := json.Unmarshal(file, &daemonMap); err != nil {
		return err
	}
	if _, ok := daemonMap["registry-mirrors"]; ok {
		for _, val := range daemonMap["registry-mirrors"].([]interface{}) {
			reply.Mirrors = append(reply.Mirrors, val.(string))
		}
	}
	return nil
}

// SetDockerRegistryMirrors 更新 Docker daemon.json 配置文件
func (s *Service) SetDockerRegistryMirrors(ctx context.Context, args rpcSchema.SetDockerRegistryMirrorsArgs, reply *rpcSchema.SetDockerRegistryMirrorsReply) error {
	// 判断 constants.DaemonJsonPath 文件是否存在
	if _, err := os.Stat(DaemonJsonPath); err != nil {
		// 创建 constants.DaemonJsonPath
		daemonMap := map[string]interface{}{
			"registry-mirrors": args.Mirrors,
		}
		// 将 daemonMap 写入 constants.DaemonJsonPath 文件
		setting, err := json.MarshalIndent(daemonMap, "", "    ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(DaemonJsonPath, setting, 0644); err != nil {
			return err
		}
	} else {
		file, err := os.ReadFile(DaemonJsonPath)
		if err != nil {
			return err
		}
		daemonMap := make(map[string]interface{})
		if err := json.Unmarshal(file, &daemonMap); err != nil {
			return err
		}
		daemonMap["registry-mirrors"] = args.Mirrors
		setting, err := json.MarshalIndent(daemonMap, "", "    ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(DaemonJsonPath, setting, 0644); err != nil {
			return err
		}
	}
	// 重启 Docker 服务
	_, err := utils.RunCommand(ctx, "systemctl", "daemon-reload")
	if err != nil {
		return err
	}
	_, err = utils.RunCommand(ctx, "systemctl", "restart", "docker-daemon")
	if err != nil {
		return err
	}
	return nil
}
