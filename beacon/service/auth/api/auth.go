// Package api
// Date: 2024/3/27 16:39
// Author: Amu
// Description:
package api

import (
	stderrors "errors"

	"beacon/pkg/auth"
	"beacon/pkg/authcookie"
	"beacon/pkg/errors"
	"beacon/pkg/fiberx"
	"beacon/pkg/validatex"
	"beacon/service/auth/service"
	"beacon/service/schema"

	"github.com/gofiber/fiber/v2"
)

type AuthAPI struct {
	AuthService   *service.AuthService
	CookieManager *authcookie.Manager
}

func NewLoginAPI(userService *service.AuthService, cookieManager *authcookie.Manager) *AuthAPI {
	return &AuthAPI{AuthService: userService, CookieManager: cookieManager}
}

func (a *AuthAPI) CSRF(ctx *fiber.Ctx) error {
	return fiberx.NoContent(ctx)
}

func (a *AuthAPI) Login(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	oldAccessToken := a.CookieManager.AccessToken(ctx)
	oldRefreshToken := a.CookieManager.RefreshToken(ctx)

	var args schema.LoginArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}

	session, err := a.AuthService.Login(c, args)
	if err != nil {
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	if err := a.AuthService.Logout(c, oldAccessToken, oldRefreshToken); err != nil {
		_ = a.AuthService.Logout(c, session.AccessToken, session.RefreshToken)
		a.CookieManager.ClearSession(ctx)
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	if err := a.CookieManager.SetSession(ctx, session.AccessToken, session.RefreshToken); err != nil {
		_ = a.AuthService.Logout(c, session.AccessToken, session.RefreshToken)
		a.CookieManager.ClearSession(ctx)
		return fiberx.Failure(ctx, errors.New500Error("create secure session failed"))
	}
	return fiberx.Success(ctx, schema.LoginResult{MustChangePassword: session.MustChangePassword})
}

func (a *AuthAPI) Logout(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	accessToken := a.CookieManager.AccessToken(ctx)
	refreshToken := a.CookieManager.RefreshToken(ctx)
	a.CookieManager.ClearSession(ctx)

	if err := a.AuthService.Logout(c, accessToken, refreshToken); err != nil {
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}

	return fiberx.NoContent(ctx)
}

func (a *AuthAPI) PassUpdate(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	var args schema.PasswordUpdateArgs
	if err := fiberx.ParseBody(ctx, &args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		return fiberx.Failure(ctx, errors.New400Error(err.Error()))
	}

	err := a.AuthService.PassUpdate(c, args)
	if err != nil {
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	accessToken := a.CookieManager.AccessToken(ctx)
	refreshToken := a.CookieManager.RefreshToken(ctx)
	a.CookieManager.ClearSession(ctx)
	if err := a.AuthService.Logout(c, accessToken, refreshToken); err != nil {
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	return fiberx.NoContent(ctx)
}

func (a *AuthAPI) TokenUpdate(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	accessToken := a.CookieManager.AccessToken(ctx)
	refreshToken := a.CookieManager.RefreshToken(ctx)
	session, err := a.AuthService.TokenUpdate(c, accessToken, refreshToken)
	if err != nil {
		a.CookieManager.ClearSession(ctx)
		if stderrors.Is(err, auth.ErrInvalidToken) {
			return fiberx.Unauthorized(ctx)
		}
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	if err := a.CookieManager.SetSession(ctx, session.AccessToken, session.RefreshToken); err != nil {
		_ = a.AuthService.Logout(c, session.AccessToken, session.RefreshToken)
		a.CookieManager.ClearSession(ctx)
		return fiberx.Failure(ctx, errors.New500Error("rotate secure session failed"))
	}
	return fiberx.NoContent(ctx)
}

func (a *AuthAPI) UserInfo(ctx *fiber.Ctx) error {
	c := ctx.UserContext()
	res, err := a.AuthService.UserInfo(c)
	if err != nil {
		return fiberx.Failure(ctx, fiberx.ServiceError(err))
	}
	return fiberx.Success(ctx, res)
}
