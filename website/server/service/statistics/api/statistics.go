// Package api
// Date:   2025/2/12 15:33
// Author: Amu
// Description:
package api

import (
	"github.com/gofiber/fiber/v2"
	"server/pkg/errors"
	"server/pkg/fiberx"
	"server/service/schema"
	"server/service/statistics/service"
)

type StatisticsAPI struct {
	StatisticsService service.IStatisticsService
}

func NewStatisticsAPI(srv service.IStatisticsService) *StatisticsAPI {
	return &StatisticsAPI{StatisticsService: srv}
}

func (api *StatisticsAPI) StatisticsQuery(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	result, err := api.StatisticsService.StatisticsQuery(c)
	if err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	return fiberx.Success(ctx, result)
}

func (api *StatisticsAPI) StatisticsUpdate(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.StatisticsUpdateArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	result, err := api.StatisticsService.StatisticsUpdate(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	return fiberx.Success(ctx, result)
}

func (api *StatisticsAPI) InstallationReport(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.InstallationReportArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	args.ClientIP = ctx.IP()
	args.UserAgent = ctx.Get("User-Agent")
	result, err := api.StatisticsService.InstallationReport(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	return fiberx.Success(ctx, result)
}
