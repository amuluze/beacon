// Package rpc
// Date: 2024/6/24 09:30
// Author: Amu
// Description:
package rpc

import (
	rpcSchema "common/rpc/schema"
	"context"
	"encoding/json"
	"os"

	"collia/pkg/utils"
)

const (
	DaemonJsonPath = "/etc/docker/daemon.json"
)

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
		if err := os.WriteFile(DaemonJsonPath, setting, 0600); err != nil {
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
		if err := os.WriteFile(DaemonJsonPath, setting, 0600); err != nil {
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

// ── Registration ──

func registerDockerHandlers(d *Dispatcher, svc *Service) {
	RegisterUnary[rpcSchema.GetDockerRegistryMirrorsArgs, rpcSchema.GetDockerRegistryMirrorsReply](d, "GetDockerRegistryMirrors", svc.GetDockerRegistryMirrors)
	RegisterUnary[rpcSchema.SetDockerRegistryMirrorsArgs, rpcSchema.SetDockerRegistryMirrorsReply](d, "SetDockerRegistryMirrors", svc.SetDockerRegistryMirrors)
}
