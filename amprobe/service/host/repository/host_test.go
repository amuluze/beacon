package repository

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"amprobe/pkg/contextx"
	"amprobe/service/model"
	"common/database"
	rpcSchema "common/rpc/schema"
)

// newHostTestDB 构造一个临时 SQLite 并迁移所有 Host 监控模型。
func newHostTestDB(t *testing.T) *database.DB {
	t.Helper()
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	if err := db.AutoMigrate(
		new(model.MonitorHost),
		new(model.MonitorCPU),
		new(model.MonitorMemory),
		new(model.MonitorDisk),
		new(model.MonitorNet),
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func TestNetUsageReturnsDBError(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	if err := db.AutoMigrate(new(model.MonitorNet)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	db.Close()

	repo := &HostRepo{DB: db}
	ctx := contextx.NewAgentID(context.Background(), "agent-a")
	if _, err := repo.NetUsage(ctx, rpcSchema.NetUsageArgs{}); err == nil {
		t.Fatal("expected db error")
	}
}

// TestHostQueriesRejectMissingAgentID 验证监控查询在 context 缺失 agentID 时
// 返回 ErrMissingAgentID，而非静默回退为全表查询（Domain I001/R001）。
// 这是 P1-A 的核心约束：读路径与控制调用写路径行为一致。
func TestHostQueriesRejectMissingAgentID(t *testing.T) {
	db := newHostTestDB(t)
	t.Cleanup(db.Close)
	repo := &HostRepo{DB: db}
	// 不带 agentID 的 context
	ctx := context.Background()

	cases := []struct {
		name string
		call func() error
	}{
		{"HostInfo", func() error { _, err := repo.HostInfo(ctx, rpcSchema.HostInfoArgs{}); return err }},
		{"CPUInfo", func() error { _, err := repo.CPUInfo(ctx, rpcSchema.CPUInfoArgs{}); return err }},
		{"CPUUsage", func() error { _, err := repo.CPUUsage(ctx, rpcSchema.CPUUsageArgs{}); return err }},
		{"MemInfo", func() error { _, err := repo.MemInfo(ctx, rpcSchema.MemoryInfoArgs{}); return err }},
		{"MemUsage", func() error { _, err := repo.MemUsage(ctx, rpcSchema.MemoryUsageArgs{}); return err }},
		{"DiskInfo", func() error { _, err := repo.DiskInfo(ctx, rpcSchema.DiskInfoArgs{}); return err }},
		{"DiskUsage", func() error { _, err := repo.DiskUsage(ctx, rpcSchema.DiskUsageArgs{}); return err }},
		{"NetUsage", func() error { _, err := repo.NetUsage(ctx, rpcSchema.NetUsageArgs{}); return err }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); err == nil {
				t.Fatalf("%s with missing agentID: expected error, got nil", tc.name)
			}
		})
	}
}

// TestHostQueriesRejectInvalidAgentID 验证格式非法的 agentID 在读路径被拒绝。
func TestHostQueriesRejectInvalidAgentID(t *testing.T) {
	db := newHostTestDB(t)
	t.Cleanup(db.Close)
	repo := &HostRepo{DB: db}
	ctx := contextx.NewAgentID(context.Background(), "agent/evil")

	if _, err := repo.HostInfo(ctx, rpcSchema.HostInfoArgs{}); err == nil {
		t.Fatal("HostInfo with invalid agentID: expected error, got nil")
	}
}

// TestHostQueriesScopedByAgentID 验证带 agentID 的查询只返回该 Agent 的数据，
// 不会跨 Agent 聚合（Domain I001/I002 的数据隔离约束）。
func TestHostQueriesScopedByAgentID(t *testing.T) {
	db := newHostTestDB(t)
	t.Cleanup(db.Close)

	// 为两个 agent 各写入一条 CPU 记录
	now := time.Now()
	if err := db.DB.Create(&model.MonitorCPU{AgentID: "agent-a", Timestamp: now, CPUPercent: 11.1}).Error; err != nil {
		t.Fatalf("seed agent-a: %v", err)
	}
	if err := db.DB.Create(&model.MonitorCPU{AgentID: "agent-b", Timestamp: now, CPUPercent: 22.2}).Error; err != nil {
		t.Fatalf("seed agent-b: %v", err)
	}

	repo := &HostRepo{DB: db}
	ctxA := contextx.NewAgentID(context.Background(), "agent-a")
	ctxB := contextx.NewAgentID(context.Background(), "agent-b")

	// CPUInfo 取最新一条，应分别返回各自的百分比
	infoA, err := repo.CPUInfo(ctxA, rpcSchema.CPUInfoArgs{})
	if err != nil {
		t.Fatalf("CPUInfo agent-a: %v", err)
	}
	if infoA.Percent != 11.1 {
		t.Fatalf("CPUInfo agent-a = %v, want 11.1", infoA.Percent)
	}

	infoB, err := repo.CPUInfo(ctxB, rpcSchema.CPUInfoArgs{})
	if err != nil {
		t.Fatalf("CPUInfo agent-b: %v", err)
	}
	if infoB.Percent != 22.2 {
		t.Fatalf("CPUInfo agent-b = %v, want 22.2", infoB.Percent)
	}
}
