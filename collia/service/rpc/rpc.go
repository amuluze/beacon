// Package rpc
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package rpc

import (
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
