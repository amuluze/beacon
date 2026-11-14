package api

import (
	"beacon/pkg/errors"
	"beacon/pkg/fiberx"
	"beacon/pkg/validatex"
	"beacon/service/feishu/service"
	"beacon/service/schema"

	"github.com/gofiber/fiber/v2"
)

type FeishuAPI struct {
	service service.IFeishuService
}

func NewFeishuAPI(service service.IFeishuService) *FeishuAPI {
	return &FeishuAPI{service: service}
}

func (a *FeishuAPI) Query(ctx *fiber.Ctx) error {
	setting, err := a.service.Query(ctx.UserContext())
	if err != nil {
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	return fiberx.Success(ctx, setting)
}

func (a *FeishuAPI) Update(ctx *fiber.Ctx) error {
	var args schema.FeishuUpdateArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	if err := a.service.Update(ctx.UserContext(), args); err != nil {
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	return fiberx.NoContent(ctx)
}

func (a *FeishuAPI) Test(ctx *fiber.Ctx) error {
	if err := a.service.Test(ctx.UserContext()); err != nil {
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	return fiberx.NoContent(ctx)
}
