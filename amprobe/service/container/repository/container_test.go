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

type fakeCaller struct {
	method string
}

func (f *fakeCaller) Call(ctx context.Context, method string, args interface{}, reply interface{}) error {
	f.method = method
	if r, ok := reply.(*rpcSchema.ContainerUpdateReply); ok {
		r.ContainerID = "updated"
	}
	return nil
}

func (f *fakeCaller) StreamCall(ctx context.Context, method string, args interface{}) (<-chan []byte, error) {
	return nil, nil
}

func (f *fakeCaller) Close() error {
	return nil
}

func TestContainerUpdateCallsAgentRPC(t *testing.T) {
	caller := &fakeCaller{}
	repo := &ContainerRepo{RPCClient: caller}

	reply, err := repo.ContainerUpdate(context.Background(), rpcSchema.ContainerUpdateArgs{ContainerID: "abc123"})
	if err != nil {
		t.Fatalf("ContainerUpdate returned error: %v", err)
	}
	if caller.method != "ContainerUpdate" {
		t.Fatalf("called method %q, want ContainerUpdate", caller.method)
	}
	if reply.ContainerID != "updated" {
		t.Fatalf("reply container id = %q, want updated", reply.ContainerID)
	}
}

// newContainerTestDB 构造一个临时 SQLite 并迁移容器监控相关模型。
func newContainerTestDB(t *testing.T) *database.DB {
	t.Helper()
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	if err := db.AutoMigrate(
		new(model.MonitorDocker),
		new(model.MonitorContainer),
		new(model.MonitorImage),
		new(model.MonitorNetwork),
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

// TestContainerQueriesRejectMissingAgentID 验证容器监控查询在 context 缺失 agentID 时
// 返回 ErrMissingAgentID，而非静默回退为全表查询（Domain I001/R001）。
func TestContainerQueriesRejectMissingAgentID(t *testing.T) {
	db := newContainerTestDB(t)
	t.Cleanup(db.Close)
	repo := &ContainerRepo{DB: db}
	ctx := context.Background()

	cases := []struct {
		name string
		call func() error
	}{
		{"Version", func() error { _, err := repo.Version(ctx, rpcSchema.DockerArgs{}); return err }},
		{"ContainerList", func() error { _, err := repo.ContainerList(ctx, rpcSchema.ContainerQueryArgs{Page: 1, Size: 10}); return err }},
		{"Usage", func() error { _, err := repo.Usage(ctx, rpcSchema.ContainerUsageArgs{}); return err }},
		{"ContainersByImage", func() error { _, err := repo.ContainersByImage(ctx, "nginx"); return err }},
		{"ContainerCount", func() error { _, err := repo.ContainerCount(ctx, rpcSchema.ContainerCountArgs{}); return err }},
		{"ImageList", func() error { _, err := repo.ImageList(ctx, rpcSchema.ImageQueryArgs{Page: 1, Size: 10}); return err }},
		{"ImageCount", func() error { _, err := repo.ImageCount(ctx, rpcSchema.ImageCountArgs{}); return err }},
		{"NetworkList", func() error { _, err := repo.NetworkList(ctx, rpcSchema.NetworkQueryArgs{Page: 1, Size: 10}); return err }},
		{"NetworkCount", func() error { _, err := repo.NetworkCount(ctx, rpcSchema.NetworkCountArgs{}); return err }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); err == nil {
				t.Fatalf("%s with missing agentID: expected error, got nil", tc.name)
			}
		})
	}
}

// TestContainerQueriesRejectInvalidAgentID 验证格式非法的 agentID 在读路径被拒绝。
func TestContainerQueriesRejectInvalidAgentID(t *testing.T) {
	db := newContainerTestDB(t)
	t.Cleanup(db.Close)
	repo := &ContainerRepo{DB: db}
	ctx := contextx.NewAgentID(context.Background(), "agent/evil")

	if _, err := repo.Version(ctx, rpcSchema.DockerArgs{}); err == nil {
		t.Fatal("Version with invalid agentID: expected error, got nil")
	}
}

// TestContainerQueriesScopedByAgentID 验证带 agentID 的容器查询只返回该 Agent 的数据，
// 不会跨 Agent 聚合（Domain I001/I003 的数据隔离约束）。
func TestContainerQueriesScopedByAgentID(t *testing.T) {
	db := newContainerTestDB(t)
	t.Cleanup(db.Close)

	// 为两个 agent 各写入一条 Docker 版本记录
	now := time.Now()
	if err := db.DB.Create(&model.MonitorDocker{AgentID: "agent-a", Timestamp: now, DockerVersion: "27.0"}).Error; err != nil {
		t.Fatalf("seed agent-a: %v", err)
	}
	if err := db.DB.Create(&model.MonitorDocker{AgentID: "agent-b", Timestamp: now, DockerVersion: "27.1"}).Error; err != nil {
		t.Fatalf("seed agent-b: %v", err)
	}

	repo := &ContainerRepo{DB: db}
	ctxA := contextx.NewAgentID(context.Background(), "agent-a")
	ctxB := contextx.NewAgentID(context.Background(), "agent-b")

	verA, err := repo.Version(ctxA, rpcSchema.DockerArgs{})
	if err != nil {
		t.Fatalf("Version agent-a: %v", err)
	}
	if verA.Data.DockerVersion != "27.0" {
		t.Fatalf("Version agent-a = %q, want 27.0", verA.Data.DockerVersion)
	}

	verB, err := repo.Version(ctxB, rpcSchema.DockerArgs{})
	if err != nil {
		t.Fatalf("Version agent-b: %v", err)
	}
	if verB.Data.DockerVersion != "27.1" {
		t.Fatalf("Version agent-b = %q, want 27.1", verB.Data.DockerVersion)
	}
}

