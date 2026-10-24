// Package repository
// Date:   2024/10/14 17:31
// Author: Amu
// Description:
package repository

import (
	"beacon/service/model"
	"beacon/service/schema"
	"common/database"
	"context"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var AlarmRepositorySet = wire.NewSet(NewAlarmRepository, wire.Bind(new(IAlarmRepository), new(*AlarmRepository)))

var _ IAlarmRepository = (*AlarmRepository)(nil)

type IAlarmRepository interface {
	AlarmQuery(ctx context.Context) ([]model.AlarmThreshold, error)
	AlarmUpdate(ctx context.Context, args schema.AlarmThresholdUpdateArgs) error
	AlarmMessageSettingQuery(ctx context.Context) (model.AlarmMessageSetting, error)
	AlarmMessageSettingUpdate(ctx context.Context, args schema.AlarmMessageSettingUpdateArgs) error
}

type AlarmRepository struct {
	DB *database.DB
}

func NewAlarmRepository(db *database.DB) *AlarmRepository {
	return &AlarmRepository{DB: db}
}

func (a AlarmRepository) AlarmQuery(ctx context.Context) ([]model.AlarmThreshold, error) {
	var reply []model.AlarmThreshold
	err := a.DB.WithContext(ctx).Model(&model.AlarmThreshold{}).Order("id ASC").Find(&reply).Error
	return reply, err
}

func (a AlarmRepository) AlarmUpdate(ctx context.Context, args schema.AlarmThresholdUpdateArgs) error {
	result := a.DB.WithContext(ctx).
		Model(&model.AlarmThreshold{}).
		Where("id = ? AND type = ?", args.ID, args.Type).
		Updates(map[string]interface{}{
			"duration":  args.Duration,
			"threshold": args.Threshold,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		return nil
	}

	var count int64
	if err := a.DB.WithContext(ctx).
		Model(&model.AlarmThreshold{}).
		Where("id = ? AND type = ?", args.ID, args.Type).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (a AlarmRepository) AlarmMessageSettingQuery(ctx context.Context) (model.AlarmMessageSetting, error) {
	var setting model.AlarmMessageSetting
	err := a.DB.WithContext(ctx).
		Where("key = ?", model.DefaultAlarmMessageSettingKey).
		First(&setting).Error
	return setting, err
}

func (a AlarmRepository) AlarmMessageSettingUpdate(ctx context.Context, args schema.AlarmMessageSettingUpdateArgs) error {
	result := a.DB.WithContext(ctx).
		Model(&model.AlarmMessageSetting{}).
		Where("key = ?", model.DefaultAlarmMessageSettingKey).
		Update("message_title", args.MessageTitle)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		return nil
	}

	var count int64
	if err := a.DB.WithContext(ctx).
		Model(&model.AlarmMessageSetting{}).
		Where("key = ?", model.DefaultAlarmMessageSettingKey).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
