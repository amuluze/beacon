package service

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"beacon/service/model"
	"beacon/service/report"
	"common/database"
	rpcSchema "common/rpc/schema"
)

// 隧道数据流写入监控批次：复用 HTTP 上报的 Store 路径（批次原子语义一致），
// 并以流注册的节点身份为准，拒绝载荷冒用其他 Agent 标识。

func newTunnelReportTestService(t *testing.T) *report.Service {
	t.Helper()
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	if err := db.AutoMigrate(
		new(model.MonitorHost), new(model.MonitorCPU), new(model.MonitorMemory),
		new(model.MonitorDisk), new(model.MonitorNet), new(model.MonitorDocker),
		new(model.MonitorContainer), new(model.MonitorImage), new(model.MonitorNetwork),
	); err != nil {
		db.Close()
		t.Fatalf("auto migrate: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return report.NewService(db, "secret-token")
}

func tunnelReportPayload(t *testing.T, agentID string) []byte {
	t.Helper()
	args := rpcSchema.MonitorReportArgs{
		AgentID: agentID,
		Host: rpcSchema.HostReport{
			Timestamp: time.Now(),
			Hostname:  "tunnel-host",
			Os:        "linux",
		},
		CPU: rpcSchema.CPUReport{
			Timestamp:  time.Now(),
			CPUPercent: 12.5,
			CPUCount:   2,
		},
		Memory: rpcSchema.MemoryReport{
			Timestamp:  time.Now(),
			MemPercent: 33.5,
			MemTotal:   1024,
			MemUsed:    343,
		},
	}
	payload, err := json.Marshal(args)
	if err != nil {
		t.Fatalf("marshal report args: %v", err)
	}
	return payload
}

func TestTunnelReportHandler_StoresBatch(t *testing.T) {
	svc := newTunnelReportTestService(t)
	handler := newTunnelReportHandler(svc)

	if err := handler(context.Background(), "node-1", rpcSchema.ReportMethod, tunnelReportPayload(t, "node-1")); err != nil {
		t.Fatalf("handle tunnel batch: %v", err)
	}

	var host model.MonitorHost
	if err := svc.DB.Where("agent_id = ?", "node-1").First(&host).Error; err != nil {
		t.Fatalf("host report not stored: %v", err)
	}
	if host.Hostname != "tunnel-host" {
		t.Fatalf("hostname = %q, want tunnel-host", host.Hostname)
	}
	var cpu model.MonitorCPU
	if err := svc.DB.Where("agent_id = ?", "node-1").First(&cpu).Error; err != nil {
		t.Fatalf("cpu report not stored: %v", err)
	}
	if cpu.CPUPercent != 12.5 {
		t.Fatalf("cpu percent = %v, want 12.5", cpu.CPUPercent)
	}
}

func TestTunnelReportHandler_RejectsAgentIDMismatch(t *testing.T) {
	svc := newTunnelReportTestService(t)
	handler := newTunnelReportHandler(svc)

	err := handler(context.Background(), "node-1", rpcSchema.ReportMethod, tunnelReportPayload(t, "node-2"))
	if err == nil {
		t.Fatal("expected mismatch rejection")
	}
	var count int64
	if err := svc.DB.Model(&model.MonitorHost{}).Count(&count).Error; err != nil {
		t.Fatalf("count host rows: %v", err)
	}
	if count != 0 {
		t.Fatalf("rejected batch stored %d rows", count)
	}
}

func TestTunnelReportHandler_RejectsUnknownMethod(t *testing.T) {
	svc := newTunnelReportTestService(t)
	handler := newTunnelReportHandler(svc)
	if err := handler(context.Background(), "node-1", "Whatever", tunnelReportPayload(t, "node-1")); err == nil {
		t.Fatal("expected unsupported method rejection")
	}
}
