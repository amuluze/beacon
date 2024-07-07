// Package api
// Date: 2024/3/6 12:53
// Author: Amu
// Description:
package api

import (
	"fmt"
	"log/slog"
	"strings"

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
	if err := fiberx.ParseQuery(ctx, &args); err != nil {
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
	c := ctx.UserContext()
	slog.Info("file upload")
	file, err := ctx.FormFile("file")
	if err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	if err := ctx.SaveFile(file, fmt.Sprintf("/tmp/%s", file.Filename)); err != nil {
		slog.Error("save file error", "err", err, "filepath", fmt.Sprintf("/tmp/%s", file.Filename))
		return fiberx.Failure(ctx, err)
	}
	// var args schema.FileUploadArgs
	prefix := ctx.FormValue("prefix", "")
	// slog.Info("get prefix", "prefix", prefix)
	// if err := fiberx.ParseBody(ctx, &args); err != nil {
	// 	return fiberx.Failure(ctx, errors.ErrBadRequest)
	// }
	// if err := validatex.ValidateStruct(args); err != nil {
	// 	return fiberx.Failure(ctx, errors.ErrBadRequest)
	// }
	slog.Info("file upload", "filepath", fmt.Sprintf("/tmp/%s", file.Filename), "prefix", prefix)
	if err := a.HostService.FileUpload(c, fmt.Sprintf("/tmp/%s", file.Filename), prefix); err != nil {
		slog.Error("file upload error", "err", err)
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	return fiberx.NoContent(ctx)
}

func (a *HostAPI) FileDownload(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.FileDownloadArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}

	filename := strings.Split(args.Filepath, "/")[len(strings.Split(args.Filepath, "/"))-1]
	res, err := a.HostService.FileDownload(c, args.Filepath, filename)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return ctx.Download(res.Filepath)
}
func (a *HostAPI) FileDelete(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.FileDeleteArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	err := a.HostService.FileDelete(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.NoContent(ctx)
}

func (a *HostAPI) FileCreate(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.FileCreateArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	err := a.HostService.FileCreate(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.NoContent(ctx)
}

func (a *HostAPI) FolderCreate(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.FolderCreateArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	err := a.HostService.FolderCreate(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.NoContent(ctx)
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

func (a *HostAPI) GetSystemTimezoneList(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	args := schema.GetSystemTimezoneListArgs{}
	list, err := a.HostService.GetSystemTimezoneList(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, list)
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

func (a *HostAPI) Reboot(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	args := schema.RebootArgs{}
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	err := a.HostService.Reboot(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.NoContent(ctx)
}

func (a *HostAPI) Shutdown(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	args := schema.ShutdownArgs{}
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.ErrBadRequest)
	}
	err := a.HostService.Shutdown(c, args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.NoContent(ctx)
}
