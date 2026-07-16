// Package repository
// Date: 2026/07/15
// Author: Amu
// Description: in-memory sqlite tests for StatisticsRepository
package repository

import (
	"context"
	"path/filepath"
	"testing"

	"server/pkg/database"
	"server/service/model"
	"server/service/schema"
)

func newTestRepo(t *testing.T) *StatisticsRepository {
	t.Helper()
	// 每个测试使用独立临时文件 DB，避免 :memory: 在连接池下的跨连接/跨用例污染
	dbFile := filepath.Join(t.TempDir(), "test.db")
	db, err := database.NewDB(
		database.WithType("sqlite"),
		database.WithDBName(dbFile),
	)
	if err != nil {
		t.Fatalf("打开测试 DB 失败: %v", err)
	}
	if err := db.AutoMigrate(&model.Statistics{}, &model.InstallationReport{}); err != nil {
		t.Fatalf("AutoMigrate 失败: %v", err)
	}
	return NewStatisticsRepository(db)
}

func TestStatisticsQuery_AutoSeed(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	// 首次查询：表为空时应自动写入一条初值记录并返回
	got, err := repo.StatisticsQuery(ctx)
	if err != nil {
		t.Fatalf("StatisticsQuery 返回错误: %v", err)
	}
	if got.Times != 0 {
		t.Errorf("初值 Times = %d, want 0", got.Times)
	}
	if got.ID == 0 {
		t.Errorf("自动写入的记录应带有效 ID")
	}
}

func TestStatisticsUpdate_Increment(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	seed, err := repo.StatisticsQuery(ctx)
	if err != nil {
		t.Fatalf("预热查询失败: %v", err)
	}

	// 连续自增 3 次
	for i := 1; i <= 3; i++ {
		if err := repo.StatisticsUpdate(ctx, schema.StatisticsUpdateArgs{ID: seed.ID}); err != nil {
			t.Fatalf("第 %d 次 StatisticsUpdate 失败: %v", i, err)
		}
	}

	got, err := repo.StatisticsQuery(ctx)
	if err != nil {
		t.Fatalf("最终查询失败: %v", err)
	}
	if got.Times != 3 {
		t.Errorf("自增后 Times = %d, want 3", got.Times)
	}
}

func TestInstallationReport_Persist(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	args := schema.InstallationReportArgs{
		InstallID:     "install-001",
		Image:         "registry.example.com/beacon:latest",
		Version:       "v3.0.4",
		PublicBaseURL: "https://example.com",
		InstallDir:    "/data/beacon",
		HTTPPort:      "1443",
		ControlPort:   "17000",
		ContainerName: "beacon",
		Hostname:      "prod-node-01",
	}
	if err := repo.InstallationReport(ctx, args); err != nil {
		t.Fatalf("InstallationReport 失败: %v", err)
	}

	var count int64
	if err := repo.DB.Model(&model.InstallationReport{}).Count(&count).Error; err != nil {
		t.Fatalf("Count 失败: %v", err)
	}
	if count != 1 {
		t.Errorf("入库条数 = %d, want 1", count)
	}
}
