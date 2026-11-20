// Package rpc
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package rpc

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"strings"

	"common/database"
	"common/rpc"

	"github.com/amuluze/docker"
)

var _ rpc.IService = (*Service)(nil)

// defaultRootDir 是未配置根目录时的防御性回退，与 config.yml 默认 host_prefix 一致。
const defaultRootDir = "/data/amprobe"

type Service struct {
	DB                   *database.DB
	Manager              *docker.Manager
	restartPolicyUpdater restartPolicyUpdater
	rootDir              string
	binaryPath           string       // 自更新目标二进制绝对路径（默认 /usr/sbin/collia）
	restartFn            func() error // 自更新成功后触发的重启回调（由 cmd/main 注入）
	upgradePubKeys       []ed25519.PublicKey
}

// NewService 构造 Agent RPC 服务。rootDir 是文件操作的根目录沙箱，
// 为空时回退 defaultRootDir，所有文件操作必须在该前缀内。
func NewService(db *database.DB, manager *docker.Manager, rootDir string) *Service {
	if rootDir == "" {
		rootDir = defaultRootDir
	}
	return &Service{
		DB:                   db,
		Manager:              manager,
		restartPolicyUpdater: dockerRestartPolicyUpdater{},
		rootDir:              rootDir,
		binaryPath:           defaultColliaBinaryPath,
	}
}

// SetUpgradePublicKeys 固定发布签名公钥集合（base64 Ed25519）；空集合表示未固定。
// 支持轮换重叠窗口：升级包签名匹配集合中任一公钥即放行。非法编码或长度直接
// 报错，不静默降级；重复键按首次出现去重。
func (s *Service) SetUpgradePublicKeys(encodedKeys []string) error {
	keys := make([]ed25519.PublicKey, 0, len(encodedKeys))
	seen := make(map[string]struct{}, len(encodedKeys))
	for _, encoded := range encodedKeys {
		trimmed := strings.TrimSpace(encoded)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		raw, err := base64.StdEncoding.DecodeString(trimmed)
		if err != nil {
			return fmt.Errorf("decode upgrade public key: %w", err)
		}
		if len(raw) != ed25519.PublicKeySize {
			return fmt.Errorf("upgrade public key must be %d bytes, got %d", ed25519.PublicKeySize, len(raw))
		}
		seen[trimmed] = struct{}{}
		keys = append(keys, ed25519.PublicKey(raw))
	}
	if len(keys) == 0 {
		s.upgradePubKeys = nil
		return nil
	}
	s.upgradePubKeys = keys
	return nil
}

// SetSelfUpdateConfig 注入自更新所需依赖：目标二进制路径与重启回调。
// 由 cmd/main 在启动时注入（systemd 场景下 restartFn 触发 systemctl restart）。
func (s *Service) SetSelfUpdateConfig(binaryPath string, restartFn func() error) {
	if binaryPath != "" {
		s.binaryPath = binaryPath
	}
	s.restartFn = restartFn
}

// RootDir 返回文件操作沙箱根目录（已 Clean），供 file handler 复用。
func (s *Service) RootDir() string { return s.rootDir }
