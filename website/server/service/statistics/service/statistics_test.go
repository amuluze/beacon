// Package service
// Date: 2026/07/16
// Author: Amu
// Description: tests for StatisticsService with mocked repository
package service

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"
	"server/service/model"
	"server/service/schema"
)

type mockRepo struct {
	stats     model.Statistics
	queryErr  error
	updateErr error
	reportErr error
}

func (m *mockRepo) StatisticsQuery(context.Context) (model.Statistics, error) {
	return m.stats, m.queryErr
}
func (m *mockRepo) StatisticsUpdate(context.Context, schema.StatisticsUpdateArgs) error {
	return m.updateErr
}
func (m *mockRepo) InstallationReport(context.Context, schema.InstallationReportArgs) error {
	return m.reportErr
}

func TestStatisticsService_Query(t *testing.T) {
	repo := &mockRepo{stats: model.Statistics{Model: gorm.Model{ID: 5}, Times: 7}}
	srv := NewStatisticsService(repo)

	rep, err := srv.StatisticsQuery(context.Background())
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}
	if rep.Data.ID != 5 || rep.Data.Times != 7 {
		t.Errorf("字段映射错误: %+v", rep.Data)
	}
}

func TestStatisticsService_QueryError(t *testing.T) {
	repo := &mockRepo{queryErr: errors.New("db down")}
	srv := NewStatisticsService(repo)

	if _, err := srv.StatisticsQuery(context.Background()); err == nil {
		t.Fatal("应向上透传 repository 错误")
	}
}

func TestStatisticsService_Update(t *testing.T) {
	repo := &mockRepo{}
	srv := NewStatisticsService(repo)

	if _, err := srv.StatisticsUpdate(context.Background(), schema.StatisticsUpdateArgs{ID: 1}); err != nil {
		t.Fatalf("正常更新不应报错: %v", err)
	}
}

func TestStatisticsService_UpdateError(t *testing.T) {
	repo := &mockRepo{updateErr: errors.New("not found")}
	srv := NewStatisticsService(repo)

	if _, err := srv.StatisticsUpdate(context.Background(), schema.StatisticsUpdateArgs{ID: 99}); err == nil {
		t.Fatal("应向上透传 repository 错误")
	}
}
