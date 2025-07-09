// Package api
// Date       : 2024/9/4 14:56
// Author     : Amu
// Description:
package api

import (
	"amprobe/pkg/errors"
	"amprobe/pkg/fiberx"
	"amprobe/pkg/validatex"
	"amprobe/service/account/service"
	"amprobe/service/schema"

	"github.com/gofiber/fiber/v2"
)

type AccountAPI struct {
	AccountService service.IAccountService
}

func NewAccountAPI(accountService service.IAccountService) *AccountAPI {
	return &AccountAPI{AccountService: accountService}
}

func (a *AccountAPI) UserQuery(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.UserQueryArgs
	if err := fiberx.ParseQuery(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}

	users, err := a.AccountService.UserQuery(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	return fiberx.Success(ctx, users)
}

func (a *AccountAPI) UserCreate(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.UserCreateArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	err := a.AccountService.UserCreate(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	return fiberx.NoContent(ctx)
}

func (a *AccountAPI) UserUpdate(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.UserUpdateArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	err := a.AccountService.UserUpdate(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	return fiberx.NoContent(ctx)
}

func (a *AccountAPI) UserDelete(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.UserDeleteArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	err := a.AccountService.UserDelete(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	return fiberx.NoContent(ctx)
}

func (a *AccountAPI) RoleQuery(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.RoleQueryArgs
	if err := fiberx.ParseQuery(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	res, err := a.AccountService.RoleQuery(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	return fiberx.Success(ctx, res)
}

func (a *AccountAPI) RoleCreate(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.RoleCreateArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	err := a.AccountService.RoleCreate(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	return fiberx.NoContent(ctx)
}

func (a *AccountAPI) RoleUpdate(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.RoleUpdateArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	err := a.AccountService.RoleUpdate(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	return fiberx.NoContent(ctx)
}

func (a *AccountAPI) RoleDelete(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.RoleDeleteArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	err := a.AccountService.RoleDelete(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	return fiberx.NoContent(ctx)
}

func (a *AccountAPI) ResourceQuery(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.ResourceQueryArgs
	if err := fiberx.ParseQuery(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	res, err := a.AccountService.ResourceQuery(c, args)
	if err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}
	return fiberx.Success(ctx, res)
}
