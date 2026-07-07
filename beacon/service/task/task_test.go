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
		new(model.MonitorMemory),
		new(model.MonitorDisk),
		new(model.MonitorContainer),
		new(model.Audit),
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func seedTwoAgents(t *testing.T, db *database.DB, now time.Time) {
	t.Helper()
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
}

// findAuditsByAgent returns audit rows recorded for a given agent_id.
func findAuditsByAgent(t *testing.T, db *database.DB, agentID string) []model.Audit {
	t.Helper()
	var audits []model.Audit
	if err := db.Model(&model.Audit{}).Where("agent_id = ?", agentID).Find(&audits).Error; err != nil {
		t.Fatalf("query audits for %s: %v", agentID, err)
	}
	return audits
}

// stubMailSender captures sent messages without dialing a real SMTP server.
type stubMailSender struct {
	messages []string
}

func (s *stubMailSender) Send(msg string) error {
	s.messages = append(s.messages, msg)
	return nil
}

// Domain Spec I004 — alarm tasks must evaluate each Agent independently,
// and recorded audit messages must carry the agent_id of the source Agent.
func TestCPUAlarmTaskEvaluatesAgentsIndependently(t *testing.T) {
	db := newTaskTestDB(t)
	now := time.Now()
	if err := db.Create(&model.AlarmThreshold{Type: "cpu", Duration: 2, Threshold: 80}).Error; err != nil {
		t.Fatalf("create threshold: %v", err)
	}
	seedTwoAgents(t, db, now)
	if err := db.Create(&[]model.MonitorCPU{
		{AgentID: "agent-a", Timestamp: now, CPUPercent: 1.0},
		{AgentID: "agent-b", Timestamp: now, CPUPercent: 0.1},
	}).Error; err != nil {
		t.Fatalf("create cpu rows: %v", err)
	}

	task := NewTask(db)
	// Pre-arm the cache so we can assert that a non-triggering Agent
	// does NOT pollute the audit log.
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
	// I004: audit row must carry the source AgentID.
	if audits[0].AgentID != "agent-a" {
		t.Fatalf("audit AgentID = %q, want agent-a", audits[0].AgentID)
	}
}

// I004 — Memory alarm should likewise be Agent-scoped: only the Agent whose
// memory usage is above threshold gets an audit row tagged with its id.
func TestMemoryAlarmTaskEvaluatesAgentsIndependently(t *testing.T) {
	db := newTaskTestDB(t)
	now := time.Now()
	if err := db.Create(&model.AlarmThreshold{Type: "memory", Duration: 2, Threshold: 50}).Error; err != nil {
		t.Fatalf("create threshold: %v", err)
	}
	seedTwoAgents(t, db, now)
	if err := db.Create(&[]model.MonitorMemory{
		{AgentID: "agent-a", Timestamp: now, MemPercent: 0.9}, // 90%
		{AgentID: "agent-b", Timestamp: now, MemPercent: 0.2}, // 20%
	}).Error; err != nil {
		t.Fatalf("create mem rows: %v", err)
	}

	task := NewTask(db)
	task.SetMailSender(&stubMailSender{})

	if err := task.MemoryAlarmTask(context.Background()); err != nil {
		t.Fatalf("MemoryAlarmTask: %v", err)
	}

	audits := findAuditsByAgent(t, db, "agent-a")
	if len(audits) != 1 {
		t.Fatalf("agent-a audit count = %d, want 1", len(audits))
	}
	if audits[0].AgentID != "agent-a" {
		t.Fatalf("audit AgentID = %q, want agent-a", audits[0].AgentID)
	}
	if audits := findAuditsByAgent(t, db, "agent-b"); len(audits) != 0 {
		t.Fatalf("agent-b audit count = %d, want 0", len(audits))
	}
}

// I004 — Disk alarm triggers once per (agent, device) pair; audit row must
// reference both the agent and the device via the cache key.
func TestDiskAlarmTaskScopesByAgent(t *testing.T) {
	db := newTaskTestDB(t)
	now := time.Now()
	if err := db.Create(&model.AlarmThreshold{Type: "disk", Duration: 2, Threshold: 50}).Error; err != nil {
		t.Fatalf("create threshold: %v", err)
	}
	seedTwoAgents(t, db, now)
	if err := db.Create(&[]model.MonitorDisk{
		{AgentID: "agent-a", Timestamp: now, Device: "/dev/sda1", DiskPercent: 0.9},
		{AgentID: "agent-b", Timestamp: now, Device: "/dev/sda1", DiskPercent: 0.1},
	}).Error; err != nil {
		t.Fatalf("create disk rows: %v", err)
	}

	task := NewTask(db)
	task.SetMailSender(&stubMailSender{})

	if err := task.DiskAlarmTask(context.Background()); err != nil {
		t.Fatalf("DiskAlarmTask: %v", err)
	}

	audits := findAuditsByAgent(t, db, "agent-a")
	if len(audits) != 1 {
		t.Fatalf("agent-a audit count = %d, want 1", len(audits))
	}
	if !strings.Contains(audits[0].Operate, "/dev/sda1") {
		t.Fatalf("disk audit should mention device, got %q", audits[0].Operate)
	}
	if audits := findAuditsByAgent(t, db, "agent-b"); len(audits) != 0 {
		t.Fatalf("agent-b audit count = %d, want 0", len(audits))
	}
}

// I004 — Service task state changes should be scoped to the source Agent so
// operators can grep audit messages back to a specific host.
func TestServiceTaskScopesByAgent(t *testing.T) {
	db := newTaskTestDB(t)
	now := time.Now()
	seedTwoAgents(t, db, now)
	if err := db.Create(&[]model.MonitorContainer{
		{AgentID: "agent-a", Timestamp: now, Name: "web", State: "running"},
		{AgentID: "agent-b", Timestamp: now, Name: "web", State: "running"},
	}).Error; err != nil {
		t.Fatalf("create containers: %v", err)
	}

	task := NewTask(db)
	task.SetMailSender(&stubMailSender{})

	// Pre-populate the cache. agent-a's stored state will differ from the DB
	// after the flip below (triggering an audit), while agent-b stays in sync
	// and must NOT produce an audit row.
	task.cache.Set("container:agent-a:web", "exited", time.Minute)
	task.cache.Set("container:agent-b:web", "running", time.Minute)

	// Flip agent-a's container state to "running"; agent-b stays unchanged.
	if err := db.Model(&model.MonitorContainer{}).
		Where("agent_id = ?", "agent-a").
		Update("state", "running").Error; err != nil {
		t.Fatalf("flip agent-a container: %v", err)
	}

	if err := task.ServiceTask(context.Background()); err != nil {
		t.Fatalf("ServiceTask: %v", err)
	}

	audits := findAuditsByAgent(t, db, "agent-a")
	if len(audits) != 1 {
		t.Fatalf("agent-a audit count = %d, want 1", len(audits))
	}
	if audits[0].AgentID != "agent-a" {
		t.Fatalf("audit AgentID = %q, want agent-a", audits[0].AgentID)
	}
	if audits := findAuditsByAgent(t, db, "agent-b"); len(audits) != 0 {
		t.Fatalf("agent-b audit count = %d, want 0 (state unchanged)", len(audits))
	}
}
