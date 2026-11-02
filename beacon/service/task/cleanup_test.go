package task

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"beacon/service/model"
	"common/database"
)

// newCleanupTestDB 构造一个带 Session/MonitorNet 表的内存文件 SQLite，
// 供 CleanupTask 的 Session 清理与 PRAGMA 回收验证使用。
func newCleanupTestDB(t *testing.T) *database.DB {
	t.Helper()
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "cleanup")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(
		new(model.MonitorCPU),
		new(model.MonitorMemory),
		new(model.MonitorDisk),
		new(model.MonitorNet),
		new(model.MonitorContainer),
		new(model.Session),
		new(model.Audit),
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

// TestCleanupDeletesExpiredTimeSeries 验证时序表保留期内保留、过期删除。
func TestCleanupDeletesExpiredTimeSeries(t *testing.T) {
	db := newCleanupTestDB(t)
	now := time.Now()
	// 1 条过期 + 1 条保留
	cutoff := now.AddDate(0, 0, -7)
	if err := db.Create(&[]model.MonitorCPU{
		{AgentID: "agent-a", Timestamp: cutoff.Add(-time.Hour), CPUPercent: 90},
		{AgentID: "agent-a", Timestamp: now, CPUPercent: 10},
	}).Error; err != nil {
		t.Fatalf("seed cpu: %v", err)
	}

	// retention=7 天
	task := NewCleanupTask(db, 7, 90, "sqlite")
	if err := task.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}

	var remaining []model.MonitorCPU
	if err := db.Find(&remaining).Error; err != nil {
		t.Fatalf("query cpu: %v", err)
	}
	if len(remaining) != 1 {
		t.Fatalf("remaining cpu rows = %d, want 1 (only the in-retention sample)", len(remaining))
	}
}

// TestCleanupDeletesClosedSessions 验证仅清理已关闭且超过保留期的会话，
// 保留活跃/异常与保留期内会话。
func TestCleanupDeletesClosedSessions(t *testing.T) {
	db := newCleanupTestDB(t)
	now := time.Now()
	old := now.AddDate(0, 0, -10) // 超过 7 天保留期

	// 1. 已关闭 + 过期 → 应删除
	endedOld := old
	// 2. 已关闭 + 保留期内 → 应保留
	endedRecent := now.Add(-time.Hour)
	// 3. 仍活跃（ended_at 为 nil）→ 应保留
	var activeEnded *time.Time

	sessions := []model.Session{
		{SessionID: "closed-old", AgentID: "a", Status: "closed", StartedAt: old, EndedAt: &endedOld},
		{SessionID: "closed-recent", AgentID: "a", Status: "closed", StartedAt: endedRecent, EndedAt: &endedRecent},
		{SessionID: "active", AgentID: "a", Status: "active", StartedAt: now, EndedAt: activeEnded},
	}
	for i := range sessions {
		if err := db.Create(&sessions[i]).Error; err != nil {
			t.Fatalf("seed session: %v", err)
		}
	}

	task := NewCleanupTask(db, 7, 90, "sqlite")
	if err := task.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}

	var remaining []model.Session
	if err := db.Find(&remaining).Error; err != nil {
		t.Fatalf("query sessions: %v", err)
	}
	got := map[string]bool{}
	for _, s := range remaining {
		got[s.SessionID] = true
	}
	if got["closed-old"] {
		t.Fatalf("closed-old session should have been deleted (closed + past retention)")
	}
	if !got["closed-recent"] {
		t.Fatalf("closed-recent session should be retained (within retention)")
	}
	if !got["active"] {
		t.Fatalf("active session should be retained (ended_at IS NULL)")
	}
}

// TestCleanupSQLitePragmasDoNotError 验证 SQLite 路径下 wal_checkpoint 与
// incremental_vacuum 不会报错（老库 auto_vacuum 未启用，incremental_vacuum 是 no-op）。
func TestCleanupSQLitePragmasDoNotError(t *testing.T) {
	db := newCleanupTestDB(t)
	task := NewCleanupTask(db, 7, 90, "sqlite")
	if err := task.Run(context.Background()); err != nil {
		t.Fatalf("Run with sqlite pragmas: %v", err)
	}
}

// TestCleanupNonSQLite skips PRAGMA 验证：非 SQLite 类型不应执行 PRAGMA（也不报错）。
func TestCleanupNonSQLiteSkipsPragmas(t *testing.T) {
	db := newCleanupTestDB(t)
	task := NewCleanupTask(db, 7, 90, "postgres")
	if err := task.Run(context.Background()); err != nil {
		t.Fatalf("Run on postgres-type: %v", err)
	}
}

// TestCleanupDeletesExpiredAudit 验证审计日志按独立保留期清理：
// 超过 AuditDays 的记录被删除，保留期内记录保留。
func TestCleanupDeletesExpiredAudit(t *testing.T) {
	db := newCleanupTestDB(t)
	now := time.Now()

	old := model.Audit{Username: "admin", Operate: "启动容器"}
	old.CreatedAt = now.AddDate(0, 0, -120) // 超过 90 天
	fresh := model.Audit{Username: "admin", Operate: "停止容器"}
	fresh.CreatedAt = now.AddDate(0, 0, -10) // 保留期内

	if err := db.Create(&[]model.Audit{old, fresh}).Error; err != nil {
		t.Fatalf("seed audit: %v", err)
	}

	// 时序数据保留 7 天，审计保留 90 天
	task := NewCleanupTask(db, 7, 90, "sqlite")
	if err := task.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}

	var remaining []model.Audit
	if err := db.Unscoped().Find(&remaining).Error; err != nil {
		t.Fatalf("query audit: %v", err)
	}
	if len(remaining) != 1 {
		t.Fatalf("remaining audit rows = %d, want 1 (only the in-retention entry)", len(remaining))
	}
	if remaining[0].Operate != "停止容器" {
		t.Fatalf("kept audit = %q, want 停止容器", remaining[0].Operate)
	}
}

// TestCleanupUsesSeparateRetentionForAuditAndTimeSeries 验证两类数据的保留期互不影响：
// 时序数据按 Days 清理，审计按 AuditDays 清理。
func TestCleanupUsesSeparateRetentionForAuditAndTimeSeries(t *testing.T) {
	db := newCleanupTestDB(t)
	now := time.Now()

	// 30 天前的时序数据：超过 7 天保留期 → 删除
	if err := db.Create(&model.MonitorCPU{
		AgentID: "agent-a", Timestamp: now.AddDate(0, 0, -30), CPUPercent: 50,
	}).Error; err != nil {
		t.Fatalf("seed cpu: %v", err)
	}
	// 30 天前的审计：在 90 天保留期内 → 保留
	audit := model.Audit{Username: "admin", Operate: "登录"}
	audit.CreatedAt = now.AddDate(0, 0, -30)
	if err := db.Create(&audit).Error; err != nil {
		t.Fatalf("seed audit: %v", err)
	}

	task := NewCleanupTask(db, 7, 90, "sqlite")
	if err := task.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}

	var cpuRows []model.MonitorCPU
	if err := db.Find(&cpuRows).Error; err != nil {
		t.Fatalf("query cpu: %v", err)
	}
	if len(cpuRows) != 0 {
		t.Fatalf("remaining cpu rows = %d, want 0", len(cpuRows))
	}

	var auditRows []model.Audit
	if err := db.Unscoped().Find(&auditRows).Error; err != nil {
		t.Fatalf("query audit: %v", err)
	}
	if len(auditRows) != 1 {
		t.Fatalf("remaining audit rows = %d, want 1", len(auditRows))
	}
}
