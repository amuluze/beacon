package report

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"beacon/service/model"
	"common/database"
	rpcSchema "common/rpc/schema"
)

func newTestDB(t *testing.T) *database.DB {
	t.Helper()

	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)

	if err := db.AutoMigrate(
		new(model.MonitorHost),
		new(model.MonitorCPU),
		new(model.MonitorMemory),
		new(model.MonitorDisk),
		new(model.MonitorNet),
		new(model.MonitorDocker),
		new(model.MonitorContainer),
		new(model.MonitorImage),
		new(model.MonitorNetwork),
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func TestStoreRejectsMissingAgentID(t *testing.T) {
	svc := NewService(newTestDB(t), "")

	if err := svc.Store(rpcSchema.MonitorReportArgs{}); !errors.Is(err, ErrMissingAgentID) {
		t.Fatalf("error = %v, want ErrMissingAgentID", err)
	}
}

// TestStoreRejectsInvalidAgentID 验证非法格式 agent_id 被拒绝。
// Domain R005 约束：缺失 Agent 标识的监控上报必须被拒绝。
func TestStoreRejectsInvalidAgentID(t *testing.T) {
	svc := NewService(newTestDB(t), "")

	// Invalid characters (spaces, slashes, etc.) should be rejected
	if err := svc.Store(rpcSchema.MonitorReportArgs{AgentID: "agent/evil"}); err == nil {
		t.Fatalf("expected error for invalid agent_id, got nil")
	}
	if err := svc.Store(rpcSchema.MonitorReportArgs{AgentID: "agent evil"}); err == nil {
		t.Fatalf("expected error for invalid agent_id with space, got nil")
	}
}

func TestStorePersistsReportBatch(t *testing.T) {
	db := newTestDB(t)
	svc := NewService(db, "")
	now := time.Now()

	err := svc.Store(rpcSchema.MonitorReportArgs{
		AgentID: "agent-a",
		Host: rpcSchema.HostReport{
			Timestamp: now,
			Hostname:  "host-a",
		},
		CPU: rpcSchema.CPUReport{
			Timestamp:  now,
			CPUPercent: 12.5,
		},
		Memory: rpcSchema.MemoryReport{
			Timestamp:  now,
			MemPercent: 33.3,
			MemTotal:   1024,
			MemUsed:    512,
		},
		Docker: rpcSchema.DockerReport{
			Timestamp:     now,
			DockerVersion: "27.0",
		},
		Disks: []rpcSchema.DiskReport{{
			Timestamp: now,
			Device:    "disk0",
		}},
		Nets: []rpcSchema.NetReport{{
			Timestamp: now,
			Ethernet:  "eth0",
		}},
		Containers: []rpcSchema.ContainerReport{{
			Timestamp:   now,
			ContainerID: "abc123",
			Name:        "app",
		}},
		Images: []rpcSchema.ImageReport{{
			Timestamp: now,
			ImageID:   "img123",
			Name:      "nginx",
		}},
		Networks: []rpcSchema.NetworkReport{{
			Timestamp: now,
			NetworkID: "net123",
			Name:      "bridge",
		}},
	})
	if err != nil {
		t.Fatalf("Store returned error: %v", err)
	}

	counts := map[string]struct {
		model interface{}
		want  int64
	}{
		"host":      {new(model.MonitorHost), 1},
		"cpu":       {new(model.MonitorCPU), 1},
		"memory":    {new(model.MonitorMemory), 1},
		"disk":      {new(model.MonitorDisk), 1},
		"net":       {new(model.MonitorNet), 1},
		"docker":    {new(model.MonitorDocker), 1},
		"container": {new(model.MonitorContainer), 1},
		"image":     {new(model.MonitorImage), 1},
		"network":   {new(model.MonitorNetwork), 1},
	}
	for name, tc := range counts {
		var got int64
		if err := db.Model(tc.model).Where("agent_id = ?", "agent-a").Count(&got).Error; err != nil {
			t.Fatalf("count %s: %v", name, err)
		}
		if got != tc.want {
			t.Fatalf("count %s = %d, want %d", name, got, tc.want)
		}
	}
}
