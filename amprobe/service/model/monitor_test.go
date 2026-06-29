package model

import (
	"path/filepath"
	"testing"

	"common/database"
)

// TestMonitorModels_AgentTimestampIndexes 验证监控表都有 (agent_id, timestamp)
// 复合索引，用于支撑按 Agent + 时间范围的监控查询与清理任务。
func TestMonitorModels_AgentTimestampIndexes(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)

	models := []interface{}{
		new(MonitorHost),
		new(MonitorCPU),
		new(MonitorMemory),
		new(MonitorDisk),
		new(MonitorNet),
		new(MonitorDocker),
		new(MonitorContainer),
		new(MonitorImage),
		new(MonitorNetwork),
	}
	if err := db.AutoMigrate(models...); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	indexes := []string{
		"idx_m_host_agent_time",
		"idx_m_cpu_agent_time",
		"idx_m_memory_agent_time",
		"idx_m_disk_agent_time",
		"idx_m_net_agent_time",
		"idx_m_docker_agent_time",
		"idx_m_container_agent_time",
		"idx_m_image_agent_time",
		"idx_m_network_agent_time",
	}
	for _, idx := range indexes {
		t.Run(idx, func(t *testing.T) {
			var count int64
			if err := db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type = 'index' AND name = ?", idx).Scan(&count).Error; err != nil {
				t.Fatalf("query sqlite_master: %v", err)
			}
			if count != 1 {
				t.Fatalf("index %s count = %d, want 1", idx, count)
			}
		})
	}
}
