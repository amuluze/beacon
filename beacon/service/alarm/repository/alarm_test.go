package repository

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"beacon/service/model"
	"beacon/service/schema"
	"common/database"

	"gorm.io/gorm"
)

func newAlarmRepositoryTestDB(t *testing.T) *database.DB {
	t.Helper()
	db, err := database.NewDB(
		database.WithType("sqlite"),
		database.WithDBName(filepath.Join(t.TempDir(), "beacon")),
	)
	if err != nil {
		t.Fatalf("create test database: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(&model.AlarmThreshold{}, &model.AlarmMessageSetting{}); err != nil {
		t.Fatalf("migrate alarm models: %v", err)
	}
	return db
}

func TestAlarmMessageSettingQueryAndUpdate(t *testing.T) {
	db := newAlarmRepositoryTestDB(t)
	repository := NewAlarmRepository(db)

	_, err := repository.AlarmMessageSettingQuery(context.Background())
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("query missing setting error = %v, want record not found", err)
	}
	err = repository.AlarmMessageSettingUpdate(context.Background(), schema.AlarmMessageSettingUpdateArgs{MessageTitle: "新标题"})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("update missing setting error = %v, want record not found", err)
	}

	setting := model.AlarmMessageSetting{
		Key:          model.DefaultAlarmMessageSettingKey,
		MessageTitle: model.DefaultAlarmMessageTitle,
	}
	if err := db.Create(&setting).Error; err != nil {
		t.Fatalf("create setting: %v", err)
	}
	if err := repository.AlarmMessageSettingUpdate(context.Background(), schema.AlarmMessageSettingUpdateArgs{MessageTitle: "生产环境主机告警"}); err != nil {
		t.Fatalf("update setting: %v", err)
	}
	got, err := repository.AlarmMessageSettingQuery(context.Background())
	if err != nil {
		t.Fatalf("query setting: %v", err)
	}
	if got.ID != setting.ID || got.MessageTitle != "生产环境主机告警" {
		t.Fatalf("setting = %+v, want updated singleton", got)
	}

	if err := repository.AlarmMessageSettingUpdate(context.Background(), schema.AlarmMessageSettingUpdateArgs{MessageTitle: "生产环境主机告警"}); err != nil {
		t.Fatalf("unchanged update should succeed: %v", err)
	}
}

func TestAlarmUpdateUsesIDAndCanonicalType(t *testing.T) {
	db := newAlarmRepositoryTestDB(t)
	record := model.AlarmThreshold{Type: "memory", Duration: 4, Threshold: 73}
	if err := db.Create(&record).Error; err != nil {
		t.Fatalf("create threshold: %v", err)
	}
	repository := NewAlarmRepository(db)

	if err := repository.AlarmUpdate(context.Background(), schema.AlarmThresholdUpdateArgs{
		ID: record.ID, Type: "memory", Duration: 0, Threshold: 0,
	}); err != nil {
		t.Fatalf("update zero values: %v", err)
	}

	var updated model.AlarmThreshold
	if err := db.First(&updated, record.ID).Error; err != nil {
		t.Fatalf("query updated threshold: %v", err)
	}
	if updated.Duration != 0 || updated.Threshold != 0 || updated.Type != "memory" {
		t.Fatalf("updated threshold = %#v, want preserved type with zero values", updated)
	}

	if err := repository.AlarmUpdate(context.Background(), schema.AlarmThresholdUpdateArgs{
		ID: record.ID, Type: "memory", Duration: 0, Threshold: 0,
	}); err != nil {
		t.Fatalf("unchanged update should succeed: %v", err)
	}

	err := repository.AlarmUpdate(context.Background(), schema.AlarmThresholdUpdateArgs{
		ID: record.ID, Type: "disk", Duration: 2, Threshold: 85,
	})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("type mismatch error = %v, want record not found", err)
	}

	err = repository.AlarmUpdate(context.Background(), schema.AlarmThresholdUpdateArgs{
		ID: record.ID + 1, Type: "memory", Duration: 2, Threshold: 80,
	})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("missing target error = %v, want record not found", err)
	}
}
