// Package api
// Date:   2024/10/14 17:31
// Author: Amu
// Description:
package api

import (
	"beacon/pkg/errors"
	"beacon/pkg/fiberx"
	"beacon/pkg/validatex"
	"beacon/service/alarm/service"
	"beacon/service/schema"

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
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
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
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	return fiberx.NoContent(ctx)
}
