package api

import (
	"beacon/pkg/errors"
	"beacon/pkg/fiberx"
	"beacon/pkg/validatex"
	"beacon/service/schema"
	"beacon/service/wecom/service"

	"github.com/gofiber/fiber/v2"
)

type WeComAPI struct {
	service service.IWeComService
}

func NewWeComAPI(service service.IWeComService) *WeComAPI {
	return &WeComAPI{service: service}
}

func (a *WeComAPI) Query(ctx *fiber.Ctx) error {
	setting, err := a.service.Query(ctx.UserContext())
	if err != nil {
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	return fiberx.Success(ctx, setting)
}

func (a *WeComAPI) Update(ctx *fiber.Ctx) error {
	var args schema.WeComUpdateArgs
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

func (a *WeComAPI) Test(ctx *fiber.Ctx) error {
	if err := a.service.Test(ctx.UserContext()); err != nil {
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	return fiberx.NoContent(ctx)
}
