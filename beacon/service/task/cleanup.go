// Package task
// Date: 2026/6/26
// Author: Amu
// Description: Cleanup task for expired monitoring time-series data.
package task

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"beacon/service/model"
	"common/database"
)

// CleanupTask deletes monitoring time-series data older than the configured retention period.
// Only time-series (append-strategy) models are cleaned; replace-strategy models
// (Host, Docker, Image, Network) never accumulate beyond the latest batch.
type CleanupTask struct {
	db     *database.DB
	days   int
	dbType string
}

// NewCleanupTask creates a cleanup task with the given retention days.
// dbType 决定是否在清理后执行 SQLite 专属的空间回收 PRAGMA。
func NewCleanupTask(db *database.DB, days int, dbType string) *CleanupTask {
	if days <= 0 {
		days = 7
	}
	return &CleanupTask{db: db, days: days, dbType: strings.ToLower(strings.TrimSpace(dbType))}
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

	// Session 表：只删除已关闭（ended_at 非空）且超过保留期的记录，保留活跃/异常会话供排查。
	sessionResult := t.db.Unscoped().Where("ended_at IS NOT NULL AND ended_at < ?", cutoff).Delete(&model.Session{})
	if sessionResult.Error != nil {
		slog.Error("cleanup: delete sessions failed", "err", sessionResult.Error)
	} else if sessionResult.RowsAffected > 0 {
		slog.Info("cleanup: deleted closed sessions", "rows", sessionResult.RowsAffected)
	}

	// SQLite 专属：删除后页面不会自动回收，db 文件只增不减，mmap 会让 RSS 缓慢上涨。
	// wal_checkpoint(PASSIVE) 把 WAL 合并回主库；incremental_vacuum 增量回收空闲页面
	//（需配合 db.go 启动期 auto_vacuum=INCREMENTAL，老库无副作用）。
	if t.isSQLite() {
		if err := t.db.Exec("PRAGMA wal_checkpoint(PASSIVE);").Error; err != nil {
			slog.Error("cleanup: sqlite wal_checkpoint failed", "err", err)
		}
		if err := t.db.Exec("PRAGMA incremental_vacuum;").Error; err != nil {
			slog.Error("cleanup: sqlite incremental_vacuum failed", "err", err)
		}
	}
	return nil
}

// isSQLite 判断当前是否使用 SQLite 后端。
func (t *CleanupTask) isSQLite() bool {
	return t.dbType == "" || t.dbType == "sqlite"
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
	case *model.Session:
		return v.TableName()
	default:
		return "unknown"
	}
}
