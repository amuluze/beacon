// Package api
// Date: 2024/3/27 16:39
// Author: Amu
// Description:
package api

import (
	"github.com/amuluze/amprobe/pkg/contextx"
	"github.com/amuluze/amprobe/pkg/fiberx"
	"github.com/amuluze/amprobe/pkg/validatex"
	"github.com/amuluze/amprobe/service/auth/service"
	"github.com/amuluze/amprobe/service/schema"
	"github.com/gofiber/fiber/v2"
)

type AuthAPI struct {
	AuthService *service.AuthService
}

func NewLoginAPI(userService *service.AuthService) *AuthAPI {
	return &AuthAPI{AuthService: userService}
}

func (a *AuthAPI) Login(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	var args schema.LoginArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, err)
	}

	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, err)
	}

	res, err := a.AuthService.Login(c, &args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, res)
}

func (a *AuthAPI) Logout(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	userID := contextx.FromUserID(c)
	tokenString := fiberx.GetToken(ctx)
	if err := a.AuthService.Logout(c, userID, tokenString); err != nil {
		return fiberx.Failure(ctx, err)
	}

	return fiberx.NoContent(ctx)
}

func (a *AuthAPI) PassUpdate(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.PasswordUpdateArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, err)
	}

	if err := validatex.ValidateStruct(args); err != nil {
		return fiberx.Failure(ctx, err)
	}

	err := a.AuthService.PassUpdate(c, &args)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, nil)
}

func (a *AuthAPI) TokenUpdate(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	tokenString := fiberx.GetToken(ctx)
	res, err := a.AuthService.TokenUpdate(c, tokenString)
	if err != nil {
		return fiberx.Failure(ctx, err)
	}
	return fiberx.Success(ctx, res)
}
