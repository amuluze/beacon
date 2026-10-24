package task

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	dingtalkclient "beacon/service/dingtalk/client"
	"beacon/service/model"
	"common/database"

	"gorm.io/gorm"
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
		new(model.AlarmMessageSetting),
		new(model.Audit),
		new(model.Mail),
		new(model.DingTalk),
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
	titles   []string
	err      error
}

func (s *stubMailSender) Send(title string, msg string) error {
	s.titles = append(s.titles, title)
	s.messages = append(s.messages, msg)
	return s.err
}

type stubDingTalkSender struct {
	configs  []dingtalkclient.Config
	messages []string
	err      error
}

func (s *stubDingTalkSender) Send(_ context.Context, config dingtalkclient.Config, msg string) error {
	s.configs = append(s.configs, config)
	s.messages = append(s.messages, msg)
	return s.err
}

// Domain Spec I004 — alarm tasks must evaluate each Agent independently,
// and recorded audit messages must carry the agent_id of the source Agent.
func TestCPUAlarmTaskEvaluatesAgentsIndependently(t *testing.T) {
	db := newTaskTestDB(t)
	now := time.Now()
	if err := db.Create(&model.AlarmThreshold{Type: "cpu", Duration: 0, Threshold: 80}).Error; err != nil {
		t.Fatalf("create threshold: %v", err)
	}
	seedTwoAgents(t, db, now)
	if err := db.Create(&[]model.MonitorCPU{
		{AgentID: "agent-a", Timestamp: now, CPUPercent: 100.0},
		{AgentID: "agent-b", Timestamp: now, CPUPercent: 10.0},
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
	if len(audits) != 0 {
		t.Fatalf("audit count = %d, want 0 while the alarm is inside its deduplication window", len(audits))
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
	createTaskRecord(t, db, &model.AlarmThreshold{Type: "cpu", Duration: 0, Threshold: 80})
	seedTaskHost(t, db, "agent-a", "host-a")
	seedTaskHost(t, db, "agent-b", "host-b")
	createTaskRecord(t, db, &model.MonitorCPU{AgentID: "agent-a", Timestamp: now, CPUPercent: 95})
	createTaskRecord(t, db, &model.MonitorCPU{AgentID: "agent-b", Timestamp: now, CPUPercent: 10})

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

	createTaskRecord(t, db, &model.AlarmThreshold{Type: "memory", Duration: 0, Threshold: 1})
	createTaskRecord(t, db, &model.AlarmThreshold{Type: "cpu", Duration: 0, Threshold: 80})
	seedTaskHost(t, db, "agent-a", "host-a")
	createTaskRecord(t, db, &model.MonitorCPU{AgentID: "agent-a", Timestamp: time.Now(), CPUPercent: 50})

	if err := NewTask(db).CPUAlarmTask(context.Background()); err != nil {
		t.Fatalf("CPUAlarmTask: %v", err)
	}

	if audits := taskAudits(t, db); len(audits) != 0 {
		t.Fatalf("audit count = %d, want 0; first audit = %q", len(audits), audits[0].Operate)
	}
}

func TestCPUAlarmTaskUsesReportedPercentScale(t *testing.T) {
	db := newTaskTestDB(t)
	t.Cleanup(db.Close)

	createTaskRecord(t, db, &model.AlarmThreshold{Type: "cpu", Duration: 0, Threshold: 80})
	seedTaskHost(t, db, "agent-a", "host-a")
	createTaskRecord(t, db, &model.MonitorCPU{AgentID: "agent-a", Timestamp: time.Now(), CPUPercent: 12})

	if err := NewTask(db).CPUAlarmTask(context.Background()); err != nil {
		t.Fatalf("CPUAlarmTask: %v", err)
	}

	if audits := taskAudits(t, db); len(audits) != 0 {
		t.Fatalf("audit count = %d, want 0 for reported CPU value 12%% under threshold 80%%", len(audits))
	}
}

func TestCPUAlarmTaskDeduplicatesRepeatedHighSamples(t *testing.T) {
	db := newTaskTestDB(t)
	t.Cleanup(db.Close)

	createTaskRecord(t, db, &model.AlarmThreshold{Type: "cpu", Duration: 0, Threshold: 80})
	seedTaskHost(t, db, "agent-a", "host-a")
	createTaskRecord(t, db, &model.MonitorCPU{AgentID: "agent-a", Timestamp: time.Now(), CPUPercent: 90})

	task := NewTask(db)
	if err := task.CPUAlarmTask(context.Background()); err != nil {
		t.Fatalf("first CPUAlarmTask: %v", err)
	}
	if err := task.CPUAlarmTask(context.Background()); err != nil {
		t.Fatalf("second CPUAlarmTask: %v", err)
	}

	if audits := taskAudits(t, db); len(audits) != 1 {
		t.Fatalf("audit count = %d, want 1 inside the alarm deduplication window", len(audits))
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

	now := time.Now().Truncate(time.Second)
	// Simulate sustained high CPU for 5+ minutes
	for i := 0; i < 6; i++ {
		createTaskRecord(t, db, &model.MonitorCPU{
			Model:      gorm.Model{CreatedAt: now.Add(-time.Duration(5-i) * time.Minute)},
			AgentID:    "agent-a",
			Timestamp:  now.Add(-time.Duration(5-i) * time.Minute),
			CPUPercent: 90, // 90% — above 80% threshold
		})
	}

	task := NewTask(db)
	task.now = func() time.Time { return now }
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

	// Now simulate a fresh recovery sample: one normal value makes the complete
	// window normal and emits exactly one recovery event.
	now = now.Add(time.Minute)
	createTaskRecord(t, db, &model.MonitorCPU{
		Model:      gorm.Model{CreatedAt: now},
		AgentID:    "agent-a",
		Timestamp:  now,
		CPUPercent: 30,
	})

	// Run the task again — should generate one recovery audit.
	if err := task.CPUAlarmTask(context.Background()); err != nil {
		t.Fatalf("CPUAlarmTask after recovery: %v", err)
	}

	audits = taskAudits(t, db)
	if len(audits) != 2 {
		t.Fatalf("audit count after recovery = %d, want 2 (alarm and recovery)", len(audits))
	}
	if !strings.Contains(audits[1].Operate, "恢复") {
		t.Fatalf("expected recovery audit, got %q", audits[1].Operate)
	}
	if err := task.CPUAlarmTask(context.Background()); err != nil {
		t.Fatalf("repeated CPUAlarmTask after recovery: %v", err)
	}
	if audits := taskAudits(t, db); len(audits) != 2 {
		t.Fatalf("repeated recovery audit count = %d, want 2", len(audits))
	}
	// I004: audit row must carry the source AgentID.
	if audits[0].AgentID != "agent-a" {
		t.Fatalf("audit AgentID = %q, want agent-a", audits[0].AgentID)
	}
}

func TestCPUAlarmRequiresCompleteContinuousWindow(t *testing.T) {
	tests := []struct {
		name    string
		offsets []time.Duration
		values  []float64
		want    int
	}{
		{name: "single high sample is insufficient", offsets: []time.Duration{0}, values: []float64{95}},
		{name: "complete high window fires", offsets: []time.Duration{-2 * time.Minute, -time.Minute, 0}, values: []float64{95, 94, 93}, want: 1},
		{name: "one normal sample keeps window normal", offsets: []time.Duration{-2 * time.Minute, -time.Minute, 0}, values: []float64{95, 50, 95}},
		{name: "sample gap makes window unknown", offsets: []time.Duration{-2 * time.Minute, 0}, values: []float64{95, 95}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTaskTestDB(t)
			now := time.Now().Truncate(time.Second)
			createTaskRecord(t, db, &model.AlarmThreshold{Type: "cpu", Duration: 2, Threshold: 80})
			seedTaskHost(t, db, "agent-a", "host-a")
			for index, offset := range tt.offsets {
				observedAt := now.Add(offset)
				createTaskRecord(t, db, &model.MonitorCPU{
					Model:      gorm.Model{CreatedAt: observedAt},
					AgentID:    "agent-a",
					Timestamp:  observedAt,
					CPUPercent: tt.values[index],
				})
			}
			task := NewTask(db)
			task.SetMailSender(&stubMailSender{})
			task.now = func() time.Time { return now }
			if err := task.CPUAlarmTask(context.Background()); err != nil {
				t.Fatalf("CPUAlarmTask: %v", err)
			}
			if audits := taskAudits(t, db); len(audits) != tt.want {
				t.Fatalf("audit count = %d, want %d", len(audits), tt.want)
			}
		})
	}
}

func TestCPUAlarmIgnoresStaleLatestSample(t *testing.T) {
	db := newTaskTestDB(t)
	now := time.Now().Truncate(time.Second)
	createTaskRecord(t, db, &model.AlarmThreshold{Type: "cpu", Duration: 0, Threshold: 80})
	seedTaskHost(t, db, "agent-a", "host-a")
	staleAt := now.Add(-2 * time.Minute)
	createTaskRecord(t, db, &model.MonitorCPU{
		Model:      gorm.Model{CreatedAt: staleAt},
		AgentID:    "agent-a",
		Timestamp:  staleAt,
		CPUPercent: 95,
	})
	task := NewTask(db)
	task.now = func() time.Time { return now }
	if err := task.CPUAlarmTask(context.Background()); err != nil {
		t.Fatalf("CPUAlarmTask: %v", err)
	}
	if audits := taskAudits(t, db); len(audits) != 0 {
		t.Fatalf("audit count = %d, want 0 for stale data", len(audits))
	}
}

func TestMemoryAlarmUsesZeroToOneHundredPercentScale(t *testing.T) {
	db := newTaskTestDB(t)
	now := time.Now().Truncate(time.Second)
	createTaskRecord(t, db, &model.AlarmThreshold{Type: "memory", Duration: 0, Threshold: 80})
	seedTaskHost(t, db, "agent-a", "host-a")
	createTaskRecord(t, db, &model.MonitorMemory{
		Model:      gorm.Model{CreatedAt: now},
		AgentID:    "agent-a",
		Timestamp:  now,
		MemPercent: 65.2,
	})
	task := NewTask(db)
	task.now = func() time.Time { return now }
	if err := task.MemoryAlarmTask(context.Background()); err != nil {
		t.Fatalf("MemoryAlarmTask: %v", err)
	}
	if audits := taskAudits(t, db); len(audits) != 0 {
		t.Fatalf("audit count = %d, want 0 for memory 65.2%% below 80%%", len(audits))
	}
}

func TestDiskAlarmUsesPercentScaleAndDurationWindow(t *testing.T) {
	t.Run("55 percent is below 85 percent", func(t *testing.T) {
		db := newTaskTestDB(t)
		now := time.Now().Truncate(time.Second)
		createTaskRecord(t, db, &model.AlarmThreshold{Type: "disk", Duration: 0, Threshold: 85})
		seedTaskHost(t, db, "agent-a", "host-a")
		createTaskRecord(t, db, &model.MonitorDisk{
			Model:       gorm.Model{CreatedAt: now},
			AgentID:     "agent-a",
			Timestamp:   now,
			Device:      "/dev/disk1",
			DiskPercent: 55,
		})
		task := NewTask(db)
		task.now = func() time.Time { return now }
		if err := task.DiskAlarmTask(context.Background()); err != nil {
			t.Fatalf("DiskAlarmTask: %v", err)
		}
		if audits := taskAudits(t, db); len(audits) != 0 {
			t.Fatalf("audit count = %d, want 0", len(audits))
		}
	})

	t.Run("complete sustained window fires", func(t *testing.T) {
		db := newTaskTestDB(t)
		now := time.Now().Truncate(time.Second)
		createTaskRecord(t, db, &model.AlarmThreshold{Type: "disk", Duration: 2, Threshold: 85})
		seedTaskHost(t, db, "agent-a", "host-a")
		for _, offset := range []time.Duration{-2 * time.Minute, -time.Minute, 0} {
			observedAt := now.Add(offset)
			createTaskRecord(t, db, &model.MonitorDisk{
				Model:       gorm.Model{CreatedAt: observedAt},
				AgentID:     "agent-a",
				Timestamp:   observedAt,
				Device:      "/dev/disk1",
				DiskPercent: 90,
			})
		}
		task := NewTask(db)
		task.SetMailSender(&stubMailSender{})
		task.now = func() time.Time { return now }
		if err := task.DiskAlarmTask(context.Background()); err != nil {
			t.Fatalf("DiskAlarmTask: %v", err)
		}
		audits := taskAudits(t, db)
		if len(audits) != 1 || !strings.Contains(audits[0].Operate, "连续 2 分钟") {
			t.Fatalf("audits = %+v, want one sustained disk alarm", audits)
		}
	})
}

func TestAlarmRecoveryNotifiesOnceAndAllowsRetrigger(t *testing.T) {
	db := newTaskTestDB(t)
	now := time.Now().Truncate(time.Second)
	createTaskRecord(t, db, &model.AlarmThreshold{Type: "cpu", Duration: 0, Threshold: 80})
	seedTaskHost(t, db, "agent-a", "host-a")
	task := NewTask(db)
	mail := &stubMailSender{}
	task.SetMailSender(mail)
	task.now = func() time.Time { return now }

	insertCPU := func(value float64) {
		createTaskRecord(t, db, &model.MonitorCPU{
			Model:      gorm.Model{CreatedAt: now},
			AgentID:    "agent-a",
			Timestamp:  now,
			CPUPercent: value,
		})
	}
	insertCPU(95)
	if err := task.CPUAlarmTask(context.Background()); err != nil {
		t.Fatalf("initial CPUAlarmTask: %v", err)
	}
	now = now.Add(time.Second)
	insertCPU(30)
	if err := task.CPUAlarmTask(context.Background()); err != nil {
		t.Fatalf("recovery CPUAlarmTask: %v", err)
	}
	if err := task.CPUAlarmTask(context.Background()); err != nil {
		t.Fatalf("repeated recovery CPUAlarmTask: %v", err)
	}
	now = now.Add(time.Second)
	insertCPU(96)
	if err := task.CPUAlarmTask(context.Background()); err != nil {
		t.Fatalf("retrigger CPUAlarmTask: %v", err)
	}

	audits := taskAudits(t, db)
	if len(audits) != 3 {
		t.Fatalf("audit count = %d, want alarm, one recovery, and retrigger", len(audits))
	}
	if !strings.Contains(audits[1].Operate, "恢复") {
		t.Fatalf("second audit = %q, want recovery", audits[1].Operate)
	}
	if len(mail.messages) != 3 {
		t.Fatalf("notification count = %d, want 3", len(mail.messages))
	}
}

func TestServiceTaskIgnoresStaleContainerState(t *testing.T) {
	db := newTaskTestDB(t)
	now := time.Now().Truncate(time.Second)
	seedTaskHost(t, db, "agent-a", "host-a")
	staleAt := now.Add(-3 * time.Minute)
	createTaskRecord(t, db, &model.MonitorContainer{
		Model:     gorm.Model{CreatedAt: staleAt},
		AgentID:   "agent-a",
		Timestamp: staleAt,
		Name:      "web",
		State:     "running",
	})
	task := NewTask(db)
	task.now = func() time.Time { return now }
	task.cache.Set("container:agent-a:web", "exited", time.Minute)
	if err := task.ServiceTask(context.Background()); err != nil {
		t.Fatalf("ServiceTask: %v", err)
	}
	if audits := taskAudits(t, db); len(audits) != 0 {
		t.Fatalf("audit count = %d, want 0 for stale container state", len(audits))
	}
}

func TestTriggerAlarmDeduplicatesConcurrentCalls(t *testing.T) {
	db := newTaskTestDB(t)
	task := NewTask(db)
	task.SetMailSender(&stubMailSender{})

	const callers = 20
	errs := make(chan error, callers)
	var wait sync.WaitGroup
	for range callers {
		wait.Add(1)
		go func() {
			defer wait.Done()
			errs <- task.triggerAlarm(context.Background(), "cpu:agent-a", "high CPU", "agent-a")
		}()
	}
	wait.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("triggerAlarm: %v", err)
		}
	}
	if audits := taskAudits(t, db); len(audits) != 1 {
		t.Fatalf("audit count = %d, want 1", len(audits))
	}
}

// I004 — Memory alarm should likewise be Agent-scoped: only the Agent whose
// memory usage is above threshold gets an audit row tagged with its id.
func TestMemoryAlarmTaskEvaluatesAgentsIndependently(t *testing.T) {
	db := newTaskTestDB(t)
	now := time.Now()
	if err := db.Create(&model.AlarmThreshold{Type: "memory", Duration: 0, Threshold: 50}).Error; err != nil {
		t.Fatalf("create threshold: %v", err)
	}
	seedTwoAgents(t, db, now)
	if err := db.Create(&[]model.MonitorMemory{
		{AgentID: "agent-a", Timestamp: now, MemPercent: 90},
		{AgentID: "agent-b", Timestamp: now, MemPercent: 20},
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
	if err := db.Create(&model.AlarmThreshold{Type: "disk", Duration: 0, Threshold: 50}).Error; err != nil {
		t.Fatalf("create threshold: %v", err)
	}
	seedTwoAgents(t, db, now)
	if err := db.Create(&[]model.MonitorDisk{
		{AgentID: "agent-a", Timestamp: now, Device: "/dev/sda1", DiskPercent: 90},
		{AgentID: "agent-b", Timestamp: now, Device: "/dev/sda1", DiskPercent: 10},
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

func TestResourceAlarmTasksSendDingTalk(t *testing.T) {
	tests := []struct {
		name         string
		alarmType    string
		threshold    int
		seed         func(*testing.T, *database.DB)
		run          func(*Task) error
		messageParts []string
	}{
		{
			name:      "cpu",
			alarmType: "cpu",
			threshold: 80,
			seed: func(t *testing.T, db *database.DB) {
				createTaskRecord(t, db, &model.MonitorCPU{AgentID: "agent-a", Timestamp: time.Now(), CPUPercent: 95})
			},
			run:          func(task *Task) error { return task.CPUAlarmTask(context.Background()) },
			messageParts: []string{"agent-a", "host-a", "CPU", "80%"},
		},
		{
			name:      "memory",
			alarmType: "memory",
			threshold: 80,
			seed: func(t *testing.T, db *database.DB) {
				createTaskRecord(t, db, &model.MonitorMemory{AgentID: "agent-a", Timestamp: time.Now(), MemPercent: 95})
			},
			run:          func(task *Task) error { return task.MemoryAlarmTask(context.Background()) },
			messageParts: []string{"agent-a", "host-a", "内存", "80%"},
		},
		{
			name:      "disk",
			alarmType: "disk",
			threshold: 80,
			seed: func(t *testing.T, db *database.DB) {
				createTaskRecord(t, db, &model.MonitorDisk{AgentID: "agent-a", Timestamp: time.Now(), Device: "/dev/disk1", DiskPercent: 95})
			},
			run:          func(task *Task) error { return task.DiskAlarmTask(context.Background()) },
			messageParts: []string{"agent-a", "host-a", "磁盘", "/dev/disk1", "80%"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTaskTestDB(t)
			createTaskRecord(t, db, &model.AlarmThreshold{Type: tt.alarmType, Duration: 0, Threshold: tt.threshold})
			seedTaskHost(t, db, "agent-a", "host-a")
			createTaskRecord(t, db, &model.DingTalk{
				Key:     model.DefaultDingTalkConfigKey,
				Enabled: true,
				Webhook: "https://oapi.dingtalk.com/robot/send?access_token=test-token",
				Secret:  "SEC-test",
				AtAll:   true,
			})
			createTaskRecord(t, db, &model.AlarmMessageSetting{
				Key:          model.DefaultAlarmMessageSettingKey,
				MessageTitle: "生产环境主机告警",
			})
			tt.seed(t, db)

			dingTalk := &stubDingTalkSender{}
			mail := &stubMailSender{}
			task := NewTask(db)
			task.SetMailSender(mail)
			task.SetDingTalkSender(dingTalk)
			if err := tt.run(task); err != nil {
				t.Fatalf("alarm task error = %v", err)
			}
			if len(dingTalk.messages) != 1 {
				t.Fatalf("DingTalk message count = %d, want 1", len(dingTalk.messages))
			}
			for _, part := range tt.messageParts {
				if !strings.Contains(dingTalk.messages[0], part) {
					t.Fatalf("DingTalk message %q missing %q", dingTalk.messages[0], part)
				}
			}
			if len(dingTalk.configs) != 1 || !dingTalk.configs[0].AtAll || dingTalk.configs[0].Secret != "SEC-test" {
				t.Fatalf("DingTalk config = %+v, want stored signed config", dingTalk.configs)
			}
			if dingTalk.configs[0].MessageTitle != "生产环境主机告警" {
				t.Fatalf("DingTalk title = %q, want configured title", dingTalk.configs[0].MessageTitle)
			}
			if len(mail.titles) != 1 || mail.titles[0] != "生产环境主机告警" {
				t.Fatalf("mail titles = %v, want configured title", mail.titles)
			}
			audits := taskAudits(t, db)
			if len(audits) != 1 || strings.Contains(audits[0].Operate, "生产环境主机告警") {
				t.Fatalf("audit messages must not be prefixed with the notification title: %+v", audits)
			}
		})
	}
}

func TestNotificationFallsBackToDefaultAlarmMessageTitle(t *testing.T) {
	db := newTaskTestDB(t)
	createTaskRecord(t, db, &model.AlarmMessageSetting{
		Key:          model.DefaultAlarmMessageSettingKey,
		MessageTitle: "invalid\nheader",
	})
	createTaskRecord(t, db, &model.DingTalk{
		Key:     model.DefaultDingTalkConfigKey,
		Enabled: true,
		Webhook: "https://oapi.dingtalk.com/robot/send?access_token=test-token",
	})
	mail := &stubMailSender{}
	dingTalk := &stubDingTalkSender{}
	task := NewTask(db)
	task.SetMailSender(mail)
	task.SetDingTalkSender(dingTalk)

	task.notify(context.Background(), "CPU 告警")

	if len(mail.titles) != 1 || mail.titles[0] != model.DefaultAlarmMessageTitle {
		t.Fatalf("mail titles = %v, want default title", mail.titles)
	}
	if len(dingTalk.configs) != 1 || dingTalk.configs[0].MessageTitle != model.DefaultAlarmMessageTitle {
		t.Fatalf("DingTalk configs = %+v, want default title", dingTalk.configs)
	}
}

func TestNotificationChannelsDoNotBlockEachOther(t *testing.T) {
	newConfiguredTask := func(t *testing.T) (*Task, *database.DB) {
		db := newTaskTestDB(t)
		createTaskRecord(t, db, &model.AlarmThreshold{Type: "cpu", Duration: 0, Threshold: 80})
		seedTaskHost(t, db, "agent-a", "host-a")
		createTaskRecord(t, db, &model.MonitorCPU{AgentID: "agent-a", Timestamp: time.Now(), CPUPercent: 95})
		createTaskRecord(t, db, &model.DingTalk{
			Key:     model.DefaultDingTalkConfigKey,
			Enabled: true,
			Webhook: "https://oapi.dingtalk.com/robot/send?access_token=test-token",
		})
		return NewTask(db), db
	}

	t.Run("mail failure does not block DingTalk", func(t *testing.T) {
		task, db := newConfiguredTask(t)
		mail := &stubMailSender{err: errors.New("smtp unavailable")}
		dingTalk := &stubDingTalkSender{}
		task.SetMailSender(mail)
		task.SetDingTalkSender(dingTalk)
		if err := task.CPUAlarmTask(context.Background()); err != nil {
			t.Fatalf("CPUAlarmTask() error = %v", err)
		}
		if len(mail.messages) != 1 || len(dingTalk.messages) != 1 {
			t.Fatalf("mail sends = %d, DingTalk sends = %d, want both attempted", len(mail.messages), len(dingTalk.messages))
		}
		if audits := taskAudits(t, db); len(audits) != 1 {
			t.Fatalf("audit count = %d, want 1", len(audits))
		}
	})

	t.Run("DingTalk failure does not roll back mail or audit", func(t *testing.T) {
		task, db := newConfiguredTask(t)
		mail := &stubMailSender{}
		dingTalk := &stubDingTalkSender{err: errors.New("webhook unavailable")}
		task.SetMailSender(mail)
		task.SetDingTalkSender(dingTalk)
		if err := task.CPUAlarmTask(context.Background()); err != nil {
			t.Fatalf("CPUAlarmTask() error = %v", err)
		}
		if len(mail.messages) != 1 || len(dingTalk.messages) != 1 {
			t.Fatalf("mail sends = %d, DingTalk sends = %d, want both attempted", len(mail.messages), len(dingTalk.messages))
		}
		if audits := taskAudits(t, db); len(audits) != 1 {
			t.Fatalf("audit count = %d, want 1", len(audits))
		}
	})
}
