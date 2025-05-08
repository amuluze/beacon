// Package api
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package api

import (
	"amprobe/pkg/errors"
	"amprobe/pkg/fiberx"

	"amprobe/pkg/validatex"

	"amprobe/service/audit/service"

	"amprobe/service/schema"

	"github.com/gofiber/fiber/v2"
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
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}

	audits, err := a.AuditService.AuditQuery(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	return fiberx.Success(ctx, audits)
}
