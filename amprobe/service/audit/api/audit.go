// Package api
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package api

import (
	"github.com/amuluze/amprobe/pkg/fiberx"
	"github.com/amuluze/amprobe/pkg/validatex"
	"github.com/amuluze/amprobe/service/audit/service"
	"github.com/amuluze/amprobe/service/schema"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

type AuditAPI struct {
	AuditService service.IAuditService
}

func NewAuditAPI(service service.IAuditService) *AuditAPI {
	return &AuditAPI{AuditService: service}
}

func (a *AuditAPI) AuditQuery(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	var args schema.AuditQueryArgs
	if err := fiberx.ParseQuery(ctx, &args); err != nil {
		fiberx.Failure(ctx, err)
	}

	slog.Info("audit query", "args", args)
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, err)
	}
	slog.Info("audit query validate", "args", args)
	audits, err := a.AuditService.AuditQuery(c, &args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, audits)
}
