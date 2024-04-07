// Package api
// Date: 2024/3/6 12:49
// Author: Amu
// Description:
package api

import (
	"log/slog"

	"github.com/amuluze/amprobe/pkg/fiberx"
	"github.com/amuluze/amprobe/pkg/validatex"
	"github.com/amuluze/amprobe/service/container/service"
	"github.com/amuluze/amprobe/service/schema"
	"github.com/gofiber/fiber/v2"
)

type ContainerAPI struct {
	ContainerService service.IContainerService
}

func NewContainerAPI(containerService service.IContainerService) *ContainerAPI {
	return &ContainerAPI{
		ContainerService: containerService,
	}
}

func (a *ContainerAPI) ContainerList(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	var args schema.ContainerQueryArgs
	if err := fiberx.ParseQuery(ctx, &args); err != nil {
		return fiberx.Failure(ctx, err)
	}
	slog.Info("args", "args", args)
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, err)
	}

	container, err := a.ContainerService.ContainerList(c, &args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, container)
}

func (a *ContainerAPI) ImageList(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	var args schema.ImageQueryArgs
	if err := fiberx.ParseQuery(ctx, &args); err != nil {
		return fiberx.Failure(ctx, err)
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, err)
	}
	images, err := a.ContainerService.ImageList(c, &args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, images)
}

func (a *ContainerAPI) Version(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	version, err := a.ContainerService.Version(c)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, version)
}

func (a *ContainerAPI) ContainerStart(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.ContainerStartArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, err)
	}
	slog.Info("args", "args", args)
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, err)
	}
	err := a.ContainerService.ContainerStart(c, &args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.NoContent(ctx)
}

func (a *ContainerAPI) ContainerStop(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	var args schema.ContainerStopArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, err)
	}

	slog.Info("args", "args", args)
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, err)
	}
	err := a.ContainerService.ContainerStop(c, &args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.NoContent(ctx)
}

func (a *ContainerAPI) ContainerRemove(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.ContainerRemoveArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, err)
	}
	slog.Info("args", "args", args)
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, err)
	}
	err := a.ContainerService.ContainerRemove(c, &args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.NoContent(ctx)
}

func (a *ContainerAPI) ContainerRestart(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.ContainerRestartArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, err)
	}
	slog.Info("args", "args", args)
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, err)
	}
	err := a.ContainerService.ContainerRestart(c, &args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.NoContent(ctx)
}

func (a *ContainerAPI) ImageRemove(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.ImageRemoveArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, err)
	}
	slog.Info("args", "args", args)
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, err)
	}
	err := a.ContainerService.ImageRemove(c, &args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.NoContent(ctx)
}

func (a *ContainerAPI) ImagesPrune(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	err := a.ContainerService.ImagesPrune(c)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.NoContent(ctx)
}
