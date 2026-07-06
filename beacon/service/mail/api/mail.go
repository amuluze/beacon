// Package api
// Date:   2024/10/14 16:08
// Author: Amu
// Description:
package api

import (
	"beacon/pkg/errors"
	"beacon/pkg/fiberx"
	"beacon/pkg/validatex"
	"beacon/service/mail/service"
	"beacon/service/schema"

	"github.com/gofiber/fiber/v2"
)

type MailAPI struct {
	MailService service.IMailService
}

func NewMailAPI(mailService service.IMailService) *MailAPI {
	return &MailAPI{MailService: mailService}
}

func (m *MailAPI) MailQuery(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	mail, err := m.MailService.MailQuery(c)
	if err != nil {
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	return fiberx.Success(ctx, mail)
}

func (m *MailAPI) MailCreate(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.MailCreateArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	err := m.MailService.MailCreate(c, args)
	if err != nil {
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	return fiberx.NoContent(ctx)
}

func (m *MailAPI) MailUpdate(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.MailUpdateArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	err := m.MailService.MailUpdate(c, args)
	if err != nil {
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	return fiberx.NoContent(ctx)
}

func (m *MailAPI) MailDelete(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.MailDeleteArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	err := m.MailService.MailDelete(c, args)
	if err != nil {
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	return fiberx.NoContent(ctx)
}

func (m *MailAPI) MailTest(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.MailTestArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	err := m.MailService.MailTest(c, args)
	if err != nil {
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	return fiberx.NoContent(ctx)
}
