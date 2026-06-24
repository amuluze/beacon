// Package service
package service

import (
	"amprobe/service/report"
	"common/database"
)

func NewReportService(config *Config, db *database.DB) *report.Service {
	return report.NewService(db, config.AgentInstall.Token)
}
