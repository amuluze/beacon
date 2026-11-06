// Package service
// Description: 提供 Beacon 服务端自身版本与更新状态的对内读取入口。
package service

import (
	"beacon/service/schema"
)

// GlobalSystemVersion 返回当前构建版本与最近一次更新检查结果的对内表示。
// 未启用更新检查时只含当前版本，update_available 恒为 false。
func GlobalSystemVersion() schema.SystemVersion {
	status := GlobalUpdateStatus()
	return schema.SystemVersion{
		Version:         status.CurrentVersion,
		UpdateAvailable: status.UpdateAvailable,
		LatestVersion:   status.LatestVersion,
		ReleaseNotes:    status.ReleaseNotes,
		CheckedAt:       status.CheckedAt,
	}
}

// ProvideSystemVersion 供依赖注入使用的版本提供者。
func ProvideSystemVersion() func() schema.SystemVersion {
	return GlobalSystemVersion
}
