package docker

import "github.com/docker/docker/errdefs"

// IsNameConflict 判定错误是否为 Docker 的名称冲突（409 Conflict）。
// 容器名唯一性以运行时为唯一事实源，调用方据此把冲突归一为可定位文案。
func IsNameConflict(err error) bool {
	return errdefs.IsConflict(err)
}
