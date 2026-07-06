// Package task
// Date: 2026/6/26
// Author: Amu
// Description: Cleanup task for expired monitoring time-series data.
package task

import (
	"context"
	"log/slog"
	"time"

	"amprobe/service/model"
	"common/database"
)

// CleanupTask deletes monitoring time-series data older than the configured retention period.
// Only time-series (append-strategy) models are cleaned; replace-strategy models
// (Host, Docker, Image, Network) never accumulate beyond the latest batch.
type CleanupTask struct {
	db   *database.DB
	days int
}

// NewCleanupTask creates a cleanup task with the given retention days.
func NewCleanupTask(db *database.DB, days int) *CleanupTask {
	if days <= 0 {
		days = 7
	}
	return &CleanupTask{db: db, days: days}
}

// Run deletes expired time-series records in a single pass.
func (t *CleanupTask) Run(ctx context.Context) error {
	cutoff := time.Now().AddDate(0, 0, -t.days)

	// Time-series models that accumulate with append strategy.
	timeSeriesModels := []interface{}{
		&model.MonitorCPU{},
		&model.MonitorMemory{},
		&model.MonitorDisk{},
		&model.MonitorNet{},
		&model.MonitorContainer{},
	}

	for _, m := range timeSeriesModels {
		result := t.db.Unscoped().Where("timestamp < ?", cutoff).Delete(m)
		if result.Error != nil {
			slog.Error("cleanup: delete failed", "model", modelName(m), "err", result.Error)
			continue
		}
		if result.RowsAffected > 0 {
			slog.Info("cleanup: deleted expired records", "model", modelName(m), "rows", result.RowsAffected)
		}
	}
	return nil
}

func modelName(m interface{}) string {
	switch v := m.(type) {
	case *model.MonitorCPU:
		return v.TableName()
	case *model.MonitorMemory:
		return v.TableName()
	case *model.MonitorDisk:
		return v.TableName()
	case *model.MonitorNet:
		return v.TableName()
	case *model.MonitorContainer:
		return v.TableName()
	default:
		return "unknown"
	}
}
