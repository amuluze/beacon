// Package service
package service

import (
	"amprobe/service/report"
	"common/database"
)

// NewReportService 构造 Agent 监控上报服务。
// installToken 同时用于安装包下载鉴权与监控上报鉴权，生产模式下必须通过
// resolveInstallToken 强校验（拒绝空/弱默认/过短），未通过则拒绝启动。
func NewReportService(config *Config, db *database.DB) (*report.Service, error) {
	token, err := resolveInstallToken(config.AgentInstall.Token, config.App.Env)
	if err != nil {
		return nil, err
	}
	return report.NewService(db, token), nil
}
