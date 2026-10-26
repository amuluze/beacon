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
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"collia/pkg/utils"
)

// DaemonJsonPath 是 Docker daemon 配置文件路径（可被测试替换）。
var DaemonJsonPath = "/etc/docker/daemon.json"

// dockerServiceUnits 列出常见的 Docker systemd 单元名，重启时逐一尝试。
var dockerServiceUnits = []string{"docker", "docker-daemon"}

// GetDockerRegistryMirrors 读取 Docker daemon.json 配置文件。
// 主机未配置 daemon.json 是常态（Docker 不会自动创建该文件），
// 此时返回空镜像列表而不是错误，避免控制台"镜像设置"整块报错。
func (s *Service) GetDockerRegistryMirrors(ctx context.Context, args rpcSchema.GetDockerRegistryMirrorsArgs, reply *rpcSchema.GetDockerRegistryMirrorsReply) error {
	file, err := os.ReadFile(DaemonJsonPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			reply.Mirrors = []string{}
			return nil
		}
		return err
	}
	daemonMap := make(map[string]interface{})
	if err := json.Unmarshal(file, &daemonMap); err != nil {
		return fmt.Errorf("parse %s: %w", DaemonJsonPath, err)
	}
	reply.Mirrors = []string{}
	raw, ok := daemonMap["registry-mirrors"]
	if !ok {
		return nil
	}
	list, ok := raw.([]interface{})
	if !ok {
		return fmt.Errorf("%s: registry-mirrors must be an array", DaemonJsonPath)
	}
	for _, val := range list {
		s, ok := val.(string)
		if !ok {
			return fmt.Errorf("%s: registry-mirrors must contain strings", DaemonJsonPath)
		}
		reply.Mirrors = append(reply.Mirrors, s)
	}
	return nil
}

// SetDockerRegistryMirrors 更新 Docker daemon.json 配置文件
func (s *Service) SetDockerRegistryMirrors(ctx context.Context, args rpcSchema.SetDockerRegistryMirrorsArgs, reply *rpcSchema.SetDockerRegistryMirrorsReply) error {
	daemonMap := make(map[string]interface{})
	file, err := os.ReadFile(DaemonJsonPath)
	switch {
	case err == nil:
		if err := json.Unmarshal(file, &daemonMap); err != nil {
			return fmt.Errorf("parse %s: %w", DaemonJsonPath, err)
		}
	case errors.Is(err, os.ErrNotExist):
		// 文件不存在：按空配置处理，并确保目录存在
		if err := os.MkdirAll(filepath.Dir(DaemonJsonPath), 0755); err != nil {
			return err
		}
	default:
		return err
	}
	mirrors := args.Mirrors
	if mirrors == nil {
		mirrors = []string{}
	}
	daemonMap["registry-mirrors"] = mirrors
	setting, err := json.MarshalIndent(daemonMap, "", "    ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(DaemonJsonPath, setting, 0600); err != nil {
		return err
	}
	// registry-mirrors 需重启 Docker 才生效。
	// 单元名在不同发行版不一致（docker / docker-daemon），且容器/微虚机环境没有 systemd：
	// 此时不判定为失败——配置已落盘，记录日志提示手动重启即可。
	for _, unit := range dockerServiceUnits {
		if _, err := utils.RunCommand(ctx, "systemctl", "restart", unit); err == nil {
			return nil
		}
	}
	slog.Warn("docker registry mirrors written but docker service not restarted",
		"path", DaemonJsonPath, "units", dockerServiceUnits)
	return nil
}

// ── Registration ──

func registerDockerHandlers(d *Dispatcher, svc *Service) {
	RegisterUnary[rpcSchema.GetDockerRegistryMirrorsArgs, rpcSchema.GetDockerRegistryMirrorsReply](d, "GetDockerRegistryMirrors", svc.GetDockerRegistryMirrors)
	RegisterUnary[rpcSchema.SetDockerRegistryMirrorsArgs, rpcSchema.SetDockerRegistryMirrorsReply](d, "SetDockerRegistryMirrors", svc.SetDockerRegistryMirrors)
}
