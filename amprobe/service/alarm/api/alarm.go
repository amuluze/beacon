// Package api
// Date:   2024/10/14 17:31
// Author: Amu
// Description:
package api

import (
	"amprobe/pkg/errors"
	"amprobe/pkg/fiberx"
	"amprobe/pkg/validatex"
	"amprobe/service/alarm/service"
	"amprobe/service/schema"

	"github.com/gofiber/fiber/v2"
)

type AlarmAPI struct {
	AlarmService service.IAlarmService
}

func NewAlarmAPI(alarmService service.IAlarmService) *AlarmAPI {
	return &AlarmAPI{AlarmService: alarmService}
}

func (a *AlarmAPI) AlarmQuery(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	alarmThreshold, err := a.AlarmService.AlarmQuery(c)
	if err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	return fiberx.Success(ctx, alarmThreshold)
}

func (a *AlarmAPI) AlarmUpdate(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.AlarmThresholdUpdateArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	if err := a.AlarmService.AlarmUpdate(c, args); err != nil {
		return fiberx.Failure(ctx, errors.New500Error(err.Error()))
	}
	return fiberx.NoContent(ctx)
}
