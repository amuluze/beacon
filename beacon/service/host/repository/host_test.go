package repository

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"beacon/pkg/contextx"
	"beacon/service/model"
	"common/database"
	rpcSchema "common/rpc/schema"
)

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

func TestHostInfoRequiresAgentID(t *testing.T) {
	repo := &HostRepo{}
	if _, err := repo.HostInfo(context.Background(), rpcSchema.HostInfoArgs{}); !errors.Is(err, contextx.ErrAgentIDRequired) {
		t.Fatalf("error = %v, want ErrAgentIDRequired", err)
	}
}

func TestHostInfoMarksStaleData(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.MonitorHost)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	old := time.Now().Add(-3 * time.Minute)
	if err := db.Create(&model.MonitorHost{AgentID: "agent-a", Timestamp: old, Hostname: "host-a"}).Error; err != nil {
		t.Fatalf("create host: %v", err)
	}

	repo := &HostRepo{DB: db}
	reply, err := repo.HostInfo(contextx.NewAgentID(context.Background(), "agent-a"), rpcSchema.HostInfoArgs{})
	if err != nil {
		t.Fatalf("HostInfo returned error: %v", err)
	}
	if !reply.Freshness.Stale || !reply.Freshness.Degraded {
		t.Fatalf("freshness = %+v, want stale degraded data", reply.Freshness)
	}
}

func TestDiskInfoReturnsLatestPerDeviceForSelectedAgent(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.MonitorDisk)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	now := time.Now()
	rows := []model.MonitorDisk{
		{AgentID: "agent-a", Device: "disk0", Timestamp: now.Add(-time.Minute), DiskPercent: 10},
		{AgentID: "agent-a", Device: "disk0", Timestamp: now, DiskPercent: 20},
		{AgentID: "agent-b", Device: "disk0", Timestamp: now.Add(time.Minute), DiskPercent: 90},
	}
	if err := db.Create(&rows).Error; err != nil {
		t.Fatalf("create disks: %v", err)
	}

	repo := &HostRepo{DB: db}
	reply, err := repo.DiskInfo(contextx.NewAgentID(context.Background(), "agent-a"), rpcSchema.DiskInfoArgs{})
	if err != nil {
		t.Fatalf("DiskInfo returned error: %v", err)
	}
	if len(reply.Info) != 1 {
		t.Fatalf("len(reply.Info) = %d, want 1", len(reply.Info))
	}
	if got := reply.Info[0].DiskPercent; got != 20 {
		t.Fatalf("disk percent = %v, want latest selected agent value 20", got)
	}
}
