// Package service
// Date: 2024/3/27 16:39
// Author: Amu
// Description:
package service

import (
	"beacon/pkg/auth"
	"beacon/pkg/contextx"
	"beacon/pkg/errors"
	"beacon/service/auth/repository"
	"beacon/service/schema"
	"context"

	"github.com/google/wire"
)

var AuthServiceSet = wire.NewSet(NewAuthService, wire.Bind(new(IAuthService), new(*AuthService)))

var _ IAuthService = (*AuthService)(nil)

type IAuthService interface {
	Login(ctx context.Context, args schema.LoginArgs) (Session, error)
	Logout(ctx context.Context, accessToken, refreshToken string) error
	PassUpdate(ctx context.Context, args schema.PasswordUpdateArgs) error
	MustChangePassword(ctx context.Context, userID string) (bool, error)
	TokenUpdate(ctx context.Context, accessToken, refreshToken string) (Session, error)
	UserInfo(ctx context.Context) (schema.UserInfo, error)
}

type Session struct {
	AccessToken        string
	RefreshToken       string
	MustChangePassword bool
}

type AuthService struct {
	Auth     auth.Auther
	AuthRepo repository.IAuthRepository
}

func NewAuthService(auth auth.Auther, authRepo repository.IAuthRepository) *AuthService {
	return &AuthService{Auth: auth, AuthRepo: authRepo}
}

func (a *AuthService) Login(ctx context.Context, args schema.LoginArgs) (Session, error) {
	var session Session
	u, err := a.AuthRepo.Login(ctx, args)
	if err != nil {
		return session, err
	}
	tokenInfo, err := a.Auth.GenerateToken(u.ID.String(), u.Username)
	if err != nil {
		return session, err
	}
	session.AccessToken = tokenInfo.GetAccessToken()
	session.RefreshToken = tokenInfo.GetRefreshToken()
	session.MustChangePassword = u.MustChangePassword
	return session, nil
}

func (a *AuthService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	var firstErr error
	for _, token := range []string{accessToken, refreshToken} {
		if token == "" {
			continue
		}
		if err := a.Auth.DestroyToken(token); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (a *AuthService) PassUpdate(ctx context.Context, args schema.PasswordUpdateArgs) error {
	if username := contextx.FromUsername(ctx); username != "" {
		args.Username = username
	}
	if args.OldPassword == args.NewPassword {
		return errors.New("same password")
	}
	if err := auth.ValidatePassword(args.NewPassword); err != nil {
		return errors.New400Error(err.Error())
	}
	return a.AuthRepo.PassUpdate(ctx, args)
}

func (a *AuthService) MustChangePassword(ctx context.Context, userID string) (bool, error) {
	return a.AuthRepo.MustChangePassword(ctx, userID)
}

func (a *AuthService) TokenUpdate(ctx context.Context, accessToken, refreshToken string) (Session, error) {
	var session Session
	userID, username, err := a.Auth.ParseToken(refreshToken, "refresh_token")
	if err != nil {
		return session, err
	}
	tokenInfo, err := a.Auth.GenerateToken(userID, username)
	if err != nil {
		return session, err
	}
	if err := a.Logout(ctx, accessToken, refreshToken); err != nil {
		_ = a.Auth.DestroyToken(tokenInfo.GetAccessToken())
		_ = a.Auth.DestroyToken(tokenInfo.GetRefreshToken())
		return session, err
	}
	session.AccessToken = tokenInfo.GetAccessToken()
	session.RefreshToken = tokenInfo.GetRefreshToken()
	return session, nil
}

func (a *AuthService) UserInfo(ctx context.Context) (schema.UserInfo, error) {
	var reply schema.UserInfo
	userID := contextx.FromUserID(ctx)
	if userID == "" {
		return reply, auth.ErrInvalidToken
	}
	user, err := a.AuthRepo.UserInfo(ctx, userID)
	if err != nil {
		return reply, err
	}
	reply.ID = user.ID.String()
	reply.Username = user.Username
	reply.Status = user.Status
	reply.IsAdmin = user.IsAdmin
	reply.MustChangePassword = user.MustChangePassword
	return reply, nil
}
