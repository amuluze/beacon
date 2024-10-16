// Package repository
// Date:   2024/10/14 17:31
// Author: Amu
// Description:
package repository

import (
	"amprobe/service/model"
	"amprobe/service/schema"
	"common/database"
	"context"
	"github.com/google/wire"
)

var AlarmRepositorySet = wire.NewSet(NewAlarmRepository, wire.Bind(new(IAlarmRepository), new(*AlarmRepository)))

var _ IAlarmRepository = (*AlarmRepository)(nil)

type IAlarmRepository interface {
	AlarmQuery(ctx context.Context) (model.AlarmThreshold, error)
	AlarmUpdate(ctx context.Context, args schema.AlarmThresholdUpdateArgs) error
}

type AlarmRepository struct {
	DB *database.DB
}

func NewAlarmRepository(db *database.DB) *AlarmRepository {
	return &AlarmRepository{DB: db}
}

func (a AlarmRepository) AlarmQuery(ctx context.Context) (model.AlarmThreshold, error) {
	var reply model.AlarmThreshold
	err := a.DB.Model(&model.AlarmThreshold{}).First(&reply).Error
	return reply, err
}

func (a AlarmRepository) AlarmUpdate(ctx context.Context, args schema.AlarmThresholdUpdateArgs) error {
	params := make(map[string]interface{})
	if args.Threshold != 0 {
		params["threshold"] = args.Threshold
	}
	if args.Duration != 0 {
		params["duration"] = args.Duration
	}
	if args.Type != "" {
		params["type"] = args.Type
	}
	return a.DB.Model(&model.AlarmThreshold{}).Where("id = ?", args.ID).Updates(params).Error
}
