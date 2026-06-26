package report

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"amprobe/service/model"
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

	if err := svc.Store(rpcSchema.MonitorReportArgs{}); err == nil {
		t.Fatal("expected missing agent_id error")
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

// TestIsValidAgentID 覆盖 agent_id 格式校验的边界。
func TestIsValidAgentID(t *testing.T) {
	tests := []struct {
		id   string
		want bool
	}{
		{"", false},
		{"agent-a", true},
		{"node_1.example", true},
		{"HOST", true},
		{"123", true},
		// 畸形值：含空格、斜杠、特殊符号
		{"agent a", false},
		{"agent/../../etc", false},
		{"agent;drop", false},
		{"agent\x00null", false},
		{"中文", false},
		// 超长
		{strings.Repeat("a", maxAgentIDLen), true},
		{strings.Repeat("a", maxAgentIDLen+1), false},
	}
	for _, tc := range tests {
		if got := isValidAgentID(tc.id); got != tc.want {
			t.Errorf("isValidAgentID(%q) = %v, want %v", tc.id, got, tc.want)
		}
	}
}

// TestStore_RejectsInvalidAgentID 验证畸形 agent_id 被拒绝入库。
func TestStore_RejectsInvalidAgentID(t *testing.T) {
	db := newTestDB(t)
	svc := &Service{DB: db, Token: "t"}
	err := svc.Store(rpcSchema.MonitorReportArgs{AgentID: "agent/evil"})
	if err == nil {
		t.Fatal("expected error for invalid agent_id")
	}
}
