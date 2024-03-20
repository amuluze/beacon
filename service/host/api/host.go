// Package api
// Date: 2024/3/6 12:53
// Author: Amu
// Description:
package api

import (
	"github.com/amuluze/amprobe/pkg/errors"
	"github.com/amuluze/amprobe/pkg/fiberx"
	"github.com/amuluze/amprobe/pkg/validatex"
	"github.com/amuluze/amprobe/service/host/service"
	"github.com/amuluze/amprobe/service/schema"
	"github.com/gofiber/fiber/v2"
)

type HostAPI struct {
	HostService service.IHostService
}

func NewHostAPI(hostService service.IHostService) *HostAPI {
	return &HostAPI{HostService: hostService}
}

func (a *HostAPI) SystemUptime(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	uptime, err := a.HostService.SystemUptime(c)
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
		return fiberx.Failure(ctx, errors.ErrBadParameter)
	}
	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadParameter)
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
		return fiberx.Failure(ctx, errors.ErrBadParameter)
	}
	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadParameter)
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
		return fiberx.Failure(ctx, errors.ErrBadParameter)
	}
	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadParameter)
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
		return fiberx.Failure(ctx, errors.ErrBadParameter)
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
