// Package api
// Date: 2024/3/6 12:53
// Author: Amu
// Description:
package api

import (
	"github.com/amuluze/amprobe/pkg/fiberx"
	"github.com/amuluze/amprobe/pkg/validatex"
	"github.com/amuluze/amprobe/service/host/rpc"
	"github.com/amuluze/amprobe/service/schema"
	"github.com/amuluze/amutool/errors"
	"github.com/gofiber/fiber/v2"
)

type HostAPI struct {
	HostService rpc.IHostService
}

func NewHostAPI(hostService rpc.IHostService) *HostAPI {
	return &HostAPI{HostService: hostService}
}

func (a *HostAPI) HostInfo(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	uptime, err := a.HostService.HostInfo(c)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, uptime)
}

func (a *HostAPI) CPUInfo(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	cpuInfo, err := a.HostService.CPUInfo(c)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, cpuInfo)
}

func (a *HostAPI) CPUUsage(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.CPUUsageArgs
	if err := fiberx.ParseQuery(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	usage, err := a.HostService.CPUUsage(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, usage)
}

func (a *HostAPI) MemInfo(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	memInfo, err := a.HostService.MemInfo(c)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, memInfo)
}

func (a *HostAPI) MemUsage(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.MemoryUsageArgs
	if err := fiberx.ParseQuery(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	usage, err := a.HostService.MemUsage(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, usage)
}

func (a *HostAPI) DiskInfo(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	diskInfo, err := a.HostService.DiskInfo(c)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, diskInfo)
}

func (a *HostAPI) DiskUsage(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.DiskUsageArgs
	if err := fiberx.ParseQuery(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	usage, err := a.HostService.DiskUsage(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, usage)
}

func (a *HostAPI) NetUsage(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.NetworkUsageArgs
	if err := fiberx.ParseQuery(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, err)
	}
	usage, err := a.HostService.NetUsage(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, usage)
}

func (a *HostAPI) FilesSearch(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.FilesSearchArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	files, err := a.HostService.FilesSearch(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, files)
}

func (a *HostAPI) FileUpload(ctx *fiber.Ctx) error {
	return nil
}

func (a *HostAPI) FileDownload(ctx *fiber.Ctx) error {
	return nil
}

func (a *HostAPI) FileDelete(ctx *fiber.Ctx) error {
	return nil
}

func (a *HostAPI) GetDNSSettings(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	args := schema.GetDNSSettingsArgs{}
	settings, err := a.HostService.GetDNSSettings(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, settings)
}

func (a *HostAPI) SetDNSSettings(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	args := schema.SetDNSSettingsArgs{}
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	err := a.HostService.SetDNSSettings(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.NoContent(ctx)
}

func (a *HostAPI) GetSystemTime(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	args := schema.GetSystemTimeArgs{}
	time, err := a.HostService.GetSystemTime(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, time)
}

func (a *HostAPI) SetSystemTime(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	args := schema.SetSystemTimeArgs{}
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	err := a.HostService.SetSystemTime(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.NoContent(ctx)
}

func (a *HostAPI) GetSystemTimezone(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	args := schema.GetSystemTimezoneArgs{}
	zone, err := a.HostService.GetSystemTimezone(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, zone)
}

func (a *HostAPI) SetSystemTimezone(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	args := schema.SetSystemTimezoneArgs{}
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	err := a.HostService.SetSystemTimezone(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.NoContent(ctx)
}
