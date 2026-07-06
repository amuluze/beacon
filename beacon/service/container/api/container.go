// Package api
// Date: 2024/3/6 12:49
// Author: Amu
// Description:
package api

import (
	"beacon/pkg/errors"
	"beacon/pkg/fiberx"
	"beacon/pkg/validatex"
	"beacon/service/container/service"
	"beacon/service/schema"
	"io"

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

func (a *ContainerAPI) Version(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	version, err := a.ContainerService.Version(c)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.Success(ctx, version)
}

func (a *ContainerAPI) ContainerCreate(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.ContainerCreateArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	container, err := a.ContainerService.ContainerCreate(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.Success(ctx, container)
}

func (a *ContainerAPI) Usage(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	var args schema.ContainerUsageArgs
	if err := fiberx.ParseQuery(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	usage, err := a.ContainerService.Usage(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.Success(ctx, usage)
}

func (a *ContainerAPI) ContainerList(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	var args schema.ContainerQueryArgs
	if err := fiberx.ParseQuery(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}

	container, err := a.ContainerService.ContainerList(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.Success(ctx, container)
}

func (a *ContainerAPI) ContainerStart(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.ContainerStartArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	err := a.ContainerService.ContainerStart(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.NoContent(ctx)
}

func (a *ContainerAPI) ContainerStop(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	var args schema.ContainerStopArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	err := a.ContainerService.ContainerStop(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.NoContent(ctx)
}

func (a *ContainerAPI) ContainerRemove(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.ContainerDeleteArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	err := a.ContainerService.ContainerDelete(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.NoContent(ctx)
}

func (a *ContainerAPI) ContainerRestart(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.ContainerRestartArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	err := a.ContainerService.ContainerRestart(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.NoContent(ctx)
}

func (a *ContainerAPI) ImageList(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	var args schema.ImageQueryArgs
	if err := fiberx.ParseQuery(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	images, err := a.ContainerService.ImageList(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.Success(ctx, images)
}

func (a *ContainerAPI) ImagePull(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.ImagePullArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	err := a.ContainerService.ImagePull(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.NoContent(ctx)
}

// ImageImport 包含两部分，文件上传和加载
func (a *ContainerAPI) ImageImport(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	file, err := ctx.FormFile("file")
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	src, err := file.Open()
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	defer src.Close()
	data, err := io.ReadAll(src)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}

	args := schema.ImageImportArgs{
		FileName: file.Filename,
		Data:     data,
	}
	if err := a.ContainerService.ImageImport(c, args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.NoContent(ctx)
}

// ImageExport 包含两步，将镜像导出为压缩文件，并提供下载
func (a *ContainerAPI) ImageExport(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.ImageExportArgs
	if err := fiberx.ParseQuery(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	reply, err := a.ContainerService.ImageExport(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	ctx.Attachment(reply.FileName)
	return ctx.Send(reply.Data)
}

func (a *ContainerAPI) ImageRemove(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.ImageDeleteArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	err := a.ContainerService.ImageDelete(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.NoContent(ctx)
}

func (a *ContainerAPI) ImagesPrune(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	err := a.ContainerService.ImagesPrune(c)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.NoContent(ctx)
}

func (a *ContainerAPI) NetworkList(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.NetworkQueryArgs
	if err := fiberx.ParseQuery(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	networks, err := a.ContainerService.NetworkList(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.Success(ctx, networks)
}

func (a *ContainerAPI) NetworkDelete(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.NetworkDeleteArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	err := a.ContainerService.NetworkDelete(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.NoContent(ctx)
}

func (a *ContainerAPI) NetworkCreate(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.NetworkCreateArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	reply, err := a.ContainerService.NetworkCreate(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.Success(ctx, reply)
}

func (a *ContainerAPI) GetDockerRegistryMirrors(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	args := schema.GetDockerRegistryMirrorsArgs{}
	res, err := a.ContainerService.GetDockerRegistryMirrors(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.Success(ctx, res)
}

func (a *ContainerAPI) SetDockerRegistryMirrors(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	args := schema.SetDockerRegistryMirrorsArgs{}
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	err := a.ContainerService.SetDockerRegistryMirrors(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.FromError(err))
	}
	return fiberx.NoContent(ctx)
}
