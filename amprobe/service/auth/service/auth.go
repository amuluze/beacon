// Package service
// Date: 2024/3/27 16:39
// Author: Amu
// Description:
package service

import (
	"amprobe/pkg/auth"
	"amprobe/pkg/errors"
	"amprobe/service/auth/repository"
	"amprobe/service/schema"
	"context"

	"github.com/google/wire"
)

var AuthServiceSet = wire.NewSet(NewAuthService, wire.Bind(new(IAuthService), new(*AuthService)))

var _ IAuthService = (*AuthService)(nil)

type IAuthService interface {
	Login(ctx context.Context, args schema.LoginArgs) (schema.LoginResult, error)
	Logout(ctx context.Context, token string) error
	PassUpdate(ctx context.Context, args schema.PasswordUpdateArgs) error
	TokenUpdate(ctx context.Context, token string) (schema.LoginResult, error)
	UserInfo(ctx context.Context, token string) (schema.UserInfo, error)
}

type AuthService struct {
	Auth     auth.Auther
	AuthRepo repository.IAuthRepository
}

func NewAuthService(auth auth.Auther, authRepo repository.IAuthRepository) *AuthService {
	return &AuthService{Auth: auth, AuthRepo: authRepo}
}

func (a *AuthService) Login(ctx context.Context, args schema.LoginArgs) (schema.LoginResult, error) {
	var reply schema.LoginResult
	u, err := a.AuthRepo.Login(ctx, args)
	if err != nil {
		return reply, err
	}
	tokenInfo, err := a.Auth.GenerateToken(u.ID.String(), u.Username)
	if err != nil {
		return reply, err
	}
	reply.AccessToken = tokenInfo.GetAccessToken()
	reply.RefreshToken = tokenInfo.GetRefreshToken()
	return reply, nil
}

func (a *AuthService) Logout(ctx context.Context, token string) error {
	return a.Auth.DestroyToken(token)
}

func (a *AuthService) PassUpdate(ctx context.Context, args schema.PasswordUpdateArgs) error {
	if args.OldPassword == args.NewPassword {
		return errors.New("same password")
	}
	return a.AuthRepo.PassUpdate(ctx, args)
}

func (a *AuthService) TokenUpdate(ctx context.Context, token string) (schema.LoginResult, error) {
	var reply schema.LoginResult
	userID, username, err := a.Auth.ParseToken(token, "refresh_token")
	if err != nil {
		return reply, err
	}
	tokenInfo, err := a.Auth.GenerateToken(userID, username)
	if err != nil {
		return reply, err
	}
	reply.AccessToken = tokenInfo.GetAccessToken()
	reply.RefreshToken = tokenInfo.GetRefreshToken()
	return reply, nil
}

func (a *AuthService) UserInfo(ctx context.Context, token string) (schema.UserInfo, error) {
	var reply schema.UserInfo
	userID, _, err := a.Auth.ParseToken(token, "access_token")
	if err != nil {
		return reply, err
	}
	user, err := a.AuthRepo.UserInfo(ctx, userID)
	if err != nil {
		return reply, err
	}
	reply.ID = user.ID.String()
	reply.Username = user.Username
	reply.Status = user.Status
	reply.IsAdmin = user.IsAdmin
	return reply, nil
}
