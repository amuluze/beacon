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
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(
		new(model.Agent),
		new(model.AlarmThreshold),
		new(model.MonitorHost),
		new(model.MonitorCPU),
		new(model.Audit),
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func TestCPUAlarmTaskEvaluatesAgentsIndependently(t *testing.T) {
	db := newTaskTestDB(t)
	now := time.Now()
	if err := db.Create(&model.AlarmThreshold{Type: "cpu", Duration: 2, Threshold: 80}).Error; err != nil {
		t.Fatalf("create threshold: %v", err)
	}
	if err := db.Create(&[]model.Agent{
		{AgentID: "agent-a", Status: "online", LastSeen: now},
		{AgentID: "agent-b", Status: "online", LastSeen: now},
	}).Error; err != nil {
		t.Fatalf("create agents: %v", err)
	}
	if err := db.Create(&[]model.MonitorHost{
		{AgentID: "agent-a", Timestamp: now, Hostname: "host-a"},
		{AgentID: "agent-b", Timestamp: now, Hostname: "host-b"},
	}).Error; err != nil {
		t.Fatalf("create hosts: %v", err)
	}
	if err := db.Create(&[]model.MonitorCPU{
		{AgentID: "agent-a", Timestamp: now, CPUPercent: 1.0},
		{AgentID: "agent-b", Timestamp: now, CPUPercent: 0.1},
	}).Error; err != nil {
		t.Fatalf("create cpu rows: %v", err)
	}

	task := NewTask(db)
	task.cache.Set("cpu:agent-a", "true", time.Minute)
	if err := task.CPUAlarmTask(context.Background()); err != nil {
		t.Fatalf("CPUAlarmTask returned error: %v", err)
	}

	var audits []model.Audit
	if err := db.Model(&model.Audit{}).Find(&audits).Error; err != nil {
		t.Fatalf("query audits: %v", err)
	}
	if len(audits) != 1 {
		t.Fatalf("audit count = %d, want 1", len(audits))
	}
	if !strings.Contains(audits[0].Operate, "[agent-a]") {
		t.Fatalf("audit operate = %q, want agent-a context", audits[0].Operate)
	}
	if strings.Contains(audits[0].Operate, "agent-b") {
		t.Fatalf("audit operate = %q, should not include non-triggering agent", audits[0].Operate)
	}
}
