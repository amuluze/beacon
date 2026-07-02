package task

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"amprobe/service/model"
	"common/database"
)

func newTaskTestDB(t *testing.T) *database.DB {
	t.Helper()
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	if err := db.AutoMigrate(
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
	if !strings.Contains(msg, "agent-b") || !strings.Contains(msg, "host-b") {
		t.Fatalf("audit message %q does not identify changed agent-b/host-b", msg)
	}
	if strings.Contains(msg, "agent-a") || strings.Contains(msg, "host-a") {
		t.Fatalf("audit message %q unexpectedly includes unchanged agent-a", msg)
	}
}
