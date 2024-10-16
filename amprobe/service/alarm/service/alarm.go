// Package service
// Date:   2024/10/14 17:32
// Author: Amu
// Description:
package service

import (
	"amprobe/service/alarm/repository"
	"amprobe/service/schema"
	"context"
	"github.com/google/wire"
)

var AlarmServiceSet = wire.NewSet(NewAlarmService, wire.Bind(new(IAlarmService), new(*AlarmService)))

var _ IAlarmService = (*AlarmService)(nil)

type IAlarmService interface {
	AlarmQuery(ctx context.Context) (schema.AlarmThreshold, error)
	AlarmUpdate(ctx context.Context, update schema.AlarmThresholdUpdateArgs) error
}

type AlarmService struct {
	AlarmRepository repository.IAlarmRepository
}

func NewAlarmService(alarmRepository repository.IAlarmRepository) *AlarmService {
	return &AlarmService{AlarmRepository: alarmRepository}
}

func (a AlarmService) AlarmQuery(ctx context.Context) (schema.AlarmThreshold, error) {
	var result schema.AlarmThreshold
	reply, err := a.AlarmRepository.AlarmQuery(ctx)
	if err != nil {
		return result, err
	}
	result.ID = reply.ID
	result.Type = reply.Type
	result.Threshold = reply.Threshold
	result.Duration = reply.Duration
	return result, nil
}

func (a AlarmService) AlarmUpdate(ctx context.Context, args schema.AlarmThresholdUpdateArgs) error {
	return a.AlarmRepository.AlarmUpdate(ctx, args)
}
