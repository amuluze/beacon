package task

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"beacon/service/model"
	"common/database"
)

func newTaskTestDB(t *testing.T) *database.DB {
	t.Helper()
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	if err := db.AutoMigrate(
		new(model.Agent),
		new(model.AlarmThreshold),
		new(model.Audit),
		new(model.Mail),
		new(model.MonitorHost),
		new(model.MonitorCPU),
		new(model.MonitorMemory),
		new(model.MonitorDisk),
		new(model.MonitorContainer),
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	// Seed a minimal mail record so alarm tasks don't fail on sendMail.
	_ = db.Create(&model.Mail{Server: "localhost", Port: 0, Sender: "test@test.com", Password: "", Receiver: "test@test.com"})
	return db
}

func createTaskRecord(t *testing.T, db *database.DB, value interface{}) {
	t.Helper()
	if err := db.Create(value).Error; err != nil {
		t.Fatalf("create %T: %v", value, err)
	}
}

func taskAudits(t *testing.T, db *database.DB) []model.Audit {
	t.Helper()
	var audits []model.Audit
	if err := db.Order("id asc").Find(&audits).Error; err != nil {
		t.Fatalf("query audits: %v", err)
	}
	return audits
}

func seedTaskHost(t *testing.T, db *database.DB, agentID string, hostname string) {
	t.Helper()
	createTaskRecord(t, db, &model.MonitorHost{
		AgentID:   agentID,
		Timestamp: time.Now(),
		Hostname:  hostname,
	})
}

func TestCPUAlarmTaskScopesDataByAgent(t *testing.T) {
	db := newTaskTestDB(t)
	t.Cleanup(db.Close)

	now := time.Now()
	createTaskRecord(t, db, &model.AlarmThreshold{Type: "cpu", Duration: 10, Threshold: 80})
	seedTaskHost(t, db, "agent-a", "host-a")
	seedTaskHost(t, db, "agent-b", "host-b")
	createTaskRecord(t, db, &model.MonitorCPU{AgentID: "agent-a", Timestamp: now, CPUPercent: 0.95})
	createTaskRecord(t, db, &model.MonitorCPU{AgentID: "agent-b", Timestamp: now, CPUPercent: 0.10})

	if err := NewTask(db).CPUAlarmTask(context.Background()); err != nil {
		t.Fatalf("CPUAlarmTask: %v", err)
	}

	audits := taskAudits(t, db)
	if len(audits) != 1 {
		t.Fatalf("audit count = %d, want 1", len(audits))
	}
	msg := audits[0].Operate
	if !strings.Contains(msg, "agent-a") || !strings.Contains(msg, "host-a") {
		t.Fatalf("audit message %q does not identify agent-a/host-a", msg)
	}
	if strings.Contains(msg, "agent-b") || strings.Contains(msg, "host-b") {
		t.Fatalf("audit message %q unexpectedly includes agent-b", msg)
	}
}

func TestCPUAlarmTaskUsesCPUThresholdType(t *testing.T) {
	db := newTaskTestDB(t)
	t.Cleanup(db.Close)

	createTaskRecord(t, db, &model.AlarmThreshold{Type: "memory", Duration: 10, Threshold: 1})
	createTaskRecord(t, db, &model.AlarmThreshold{Type: "cpu", Duration: 10, Threshold: 80})
	seedTaskHost(t, db, "agent-a", "host-a")
	createTaskRecord(t, db, &model.MonitorCPU{AgentID: "agent-a", Timestamp: time.Now(), CPUPercent: 0.50})

	if err := NewTask(db).CPUAlarmTask(context.Background()); err != nil {
		t.Fatalf("CPUAlarmTask: %v", err)
	}

	if audits := taskAudits(t, db); len(audits) != 0 {
		t.Fatalf("audit count = %d, want 0; first audit = %q", len(audits), audits[0].Operate)
	}
}

func TestServiceTaskSeparatesContainerStateCacheByAgent(t *testing.T) {
	db := newTaskTestDB(t)
	t.Cleanup(db.Close)

	now := time.Now()
	seedTaskHost(t, db, "agent-a", "host-a")
	seedTaskHost(t, db, "agent-b", "host-b")
	createTaskRecord(t, db, &model.MonitorContainer{AgentID: "agent-a", Timestamp: now, Name: "app", State: "running"})
	createTaskRecord(t, db, &model.MonitorContainer{AgentID: "agent-b", Timestamp: now, Name: "app", State: "exited"})

	task := NewTask(db)
	if err := task.ServiceTask(context.Background()); err != nil {
		t.Fatalf("initial ServiceTask: %v", err)
	}
	if audits := taskAudits(t, db); len(audits) != 0 {
		t.Fatalf("initial audit count = %d, want 0", len(audits))
	}

	if err := db.Model(&model.MonitorContainer{}).Where("agent_id = ?", "agent-b").Update("state", "running").Error; err != nil {
		t.Fatalf("update agent-b state: %v", err)
	}
	if err := task.ServiceTask(context.Background()); err != nil {
		t.Fatalf("second ServiceTask: %v", err)
	}

	audits := taskAudits(t, db)
	if len(audits) != 1 {
		t.Fatalf("audit count = %d, want 1", len(audits))
	}
	msg := audits[0].Operate
	if !strings.Contains(msg, "agent-b") {
		t.Fatalf("audit message %q does not identify changed agent-b", msg)
	}
	if strings.Contains(msg, "agent-a") || strings.Contains(msg, "host-a") {
		t.Fatalf("audit message %q unexpectedly includes unchanged agent-a", msg)
	}
}

// TestCPUAlarmRecoversAfterMetricDrops 验证报警恢复：
// 当指标持续高于阈值时产生报警，指标回落后不再产生新报警。
// Domain I004约束：报警由可观测指标+阈值+恢复条件决定。
func TestCPUAlarmRecoversAfterMetricDrops(t *testing.T) {
	db := newTaskTestDB(t)
	t.Cleanup(db.Close)

	createTaskRecord(t, db, &model.AlarmThreshold{Type: "cpu", Duration: 5, Threshold: 80})
	seedTaskHost(t, db, "agent-a", "host-a")

	now := time.Now()
	// Simulate sustained high CPU for 5+ minutes
	for i := 0; i < 6; i++ {
		createTaskRecord(t, db, &model.MonitorCPU{
			AgentID:   "agent-a",
			Timestamp: now.Add(-time.Duration(5-i) * time.Minute),
			CPUPercent: 0.90, // 90% — above 80% threshold
		})
	}

	task := NewTask(db)
	if err := task.CPUAlarmTask(context.Background()); err != nil {
		t.Fatalf("CPUAlarmTask: %v", err)
	}

	// Expected: one alarm due to sustained high CPU
	audits := taskAudits(t, db)
	if len(audits) != 1 {
		t.Fatalf("alarm audit count = %d, want 1", len(audits))
	}
	if !strings.Contains(audits[0].Operate, "CPU") {
		t.Fatalf("expected CPU alarm message, got %q", audits[0].Operate)
	}

	// Now simulate recovery: CPU drops below threshold
	for i := 0; i < 3; i++ {
		createTaskRecord(t, db, &model.MonitorCPU{
			AgentID:   "agent-a",
			Timestamp: now.Add(time.Duration(i) * time.Minute),
			CPUPercent: 0.30, // 30% — well below 80% threshold
		})
	}

	// Run the task again — should NOT generate a new alarm
	if err := task.CPUAlarmTask(context.Background()); err != nil {
		t.Fatalf("CPUAlarmTask after recovery: %v", err)
	}

	// No additional audits beyond the original alarm
	audits = taskAudits(t, db)
	if len(audits) != 1 {
		t.Fatalf("audit count after recovery = %d, want 1 (no new alarm)", len(audits))
	}
}
