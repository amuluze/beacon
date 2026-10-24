// Package service
// Date:   2024/10/14 17:32
// Author: Amu
// Description:
package service

import (
	apperrors "beacon/pkg/errors"
	"beacon/service/alarm/repository"
	alarmtitle "beacon/service/alarm/title"
	"beacon/service/model"
	"beacon/service/schema"
	"context"
	"errors"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var AlarmServiceSet = wire.NewSet(NewAlarmService, wire.Bind(new(IAlarmService), new(*AlarmService)))

var _ IAlarmService = (*AlarmService)(nil)

type IAlarmService interface {
	AlarmQuery(ctx context.Context) (schema.AlarmThresholdQueryReply, error)
	AlarmUpdate(ctx context.Context, update schema.AlarmThresholdUpdateArgs) error
	AlarmMessageSettingQuery(ctx context.Context) (schema.AlarmMessageSetting, error)
	AlarmMessageSettingUpdate(ctx context.Context, update schema.AlarmMessageSettingUpdateArgs) error
}

type AlarmService struct {
	AlarmRepository repository.IAlarmRepository
}

func NewAlarmService(alarmRepository repository.IAlarmRepository) *AlarmService {
	return &AlarmService{AlarmRepository: alarmRepository}
}

func (a AlarmService) AlarmQuery(ctx context.Context) (schema.AlarmThresholdQueryReply, error) {
	var result schema.AlarmThresholdQueryReply
	reply, err := a.AlarmRepository.AlarmQuery(ctx)
	if err != nil {
		return result, err
	}
	for _, r := range reply {
		result.Data = append(result.Data, schema.AlarmThreshold{
			ID:        r.ID,
			Type:      r.Type,
			Threshold: r.Threshold,
			Duration:  r.Duration,
		})
	}
	return result, nil
}

func (a AlarmService) AlarmUpdate(ctx context.Context, args schema.AlarmThresholdUpdateArgs) error {
	return a.AlarmRepository.AlarmUpdate(ctx, args)
}

func (a AlarmService) AlarmMessageSettingQuery(ctx context.Context) (schema.AlarmMessageSetting, error) {
	setting, err := a.AlarmRepository.AlarmMessageSettingQuery(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return schema.AlarmMessageSetting{MessageTitle: model.DefaultAlarmMessageTitle}, nil
	}
	if err != nil {
		return schema.AlarmMessageSetting{}, err
	}
	return schema.AlarmMessageSetting{MessageTitle: alarmtitle.OrDefault(setting.MessageTitle)}, nil
}

func (a AlarmService) AlarmMessageSettingUpdate(ctx context.Context, args schema.AlarmMessageSettingUpdateArgs) error {
	title, err := alarmtitle.Normalize(args.MessageTitle)
	if err != nil {
		return apperrors.New400Error(err.Error())
	}
	args.MessageTitle = title
	return a.AlarmRepository.AlarmMessageSettingUpdate(ctx, args)
}
