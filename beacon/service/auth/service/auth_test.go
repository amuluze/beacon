// Package service
// Date: 2026/6/26
// Author: Amu
// Description: unit tests for AuthService
package service

import (
	"context"
	"testing"

	"amprobe/pkg/auth"
	"amprobe/service/model"
	"amprobe/service/schema"
	"amprobe/service/testutil"

	"github.com/google/uuid"
)

func newTestAuthService(auther *testutil.FakeAuther, repo *testutil.FakeAuthRepo) *AuthService {
	return &AuthService{Auth: auther, AuthRepo: repo}
}

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*testutil.FakeAuther, *testutil.FakeAuthRepo)
		args    schema.LoginArgs
		wantErr bool
	}{
		{
			name: "success",
			args: schema.LoginArgs{Username: "admin", Password: "pass123"},
			setup: func(a *testutil.FakeAuther, r *testutil.FakeAuthRepo) {
				r.LoginFn = func(ctx context.Context, args schema.LoginArgs) (model.User, error) {
					return model.User{ID: uuid.New(), Username: args.Username}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "repo error",
			args: schema.LoginArgs{Username: "admin", Password: "wrong"},
			setup: func(a *testutil.FakeAuther, r *testutil.FakeAuthRepo) {
				r.LoginFn = func(ctx context.Context, args schema.LoginArgs) (model.User, error) {
					return model.User{}, auth.ErrInvalidToken
				}
			},
			wantErr: true,
		},
		{
			name: "token generation error",
			args: schema.LoginArgs{Username: "admin", Password: "pass123"},
			setup: func(a *testutil.FakeAuther, r *testutil.FakeAuthRepo) {
				r.LoginFn = func(ctx context.Context, args schema.LoginArgs) (model.User, error) {
					return model.User{ID: uuid.New(), Username: args.Username}, nil
				}
				a.GenerateTokenFn = func(userID, username string) (auth.TokenInfo, error) {
					return nil, testutil.ErrTest
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auther := testutil.NewFakeAuther()
			repo := testutil.NewFakeAuthRepo()
			tt.setup(auther, repo)

			svc := newTestAuthService(auther, repo)
			result, err := svc.Login(context.Background(), tt.args)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.AccessToken == "" {
				t.Error("access_token is empty")
			}
			if result.RefreshToken == "" {
				t.Error("refresh_token is empty")
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*testutil.FakeAuther)
		wantErr bool
	}{
		{
			name:    "success",
			setup:   func(a *testutil.FakeAuther) {},
			wantErr: false,
		},
		{
			name: "destroy error",
			setup: func(a *testutil.FakeAuther) {
				a.DestroyTokenFn = func(token string) error { return testutil.ErrTest }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auther := testutil.NewFakeAuther()
			repo := testutil.NewFakeAuthRepo()
			tt.setup(auther)

			svc := newTestAuthService(auther, repo)
			err := svc.Logout(context.Background(), "some-token")

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestAuthService_PassUpdate(t *testing.T) {
	tests := []struct {
		name    string
		args    schema.PasswordUpdateArgs
		setup   func(*testutil.FakeAuthRepo)
		wantErr bool
	}{
		{
			name: "success",
			args: schema.PasswordUpdateArgs{Username: "admin", OldPassword: "old", NewPassword: "new"},
			setup: func(r *testutil.FakeAuthRepo) {},
			wantErr: false,
		},
		{
			name: "same password rejected",
			args: schema.PasswordUpdateArgs{Username: "admin", OldPassword: "same", NewPassword: "same"},
			setup: func(r *testutil.FakeAuthRepo) {},
			wantErr: true,
		},
		{
			name: "repo error",
			args: schema.PasswordUpdateArgs{Username: "admin", OldPassword: "old", NewPassword: "new"},
			setup: func(r *testutil.FakeAuthRepo) {
				r.PassUpdateFn = func(ctx context.Context, args schema.PasswordUpdateArgs) error { return testutil.ErrTest }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auther := testutil.NewFakeAuther()
			repo := testutil.NewFakeAuthRepo()
			tt.setup(repo)

			svc := newTestAuthService(auther, repo)
			err := svc.PassUpdate(context.Background(), tt.args)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestAuthService_TokenUpdate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*testutil.FakeAuther)
		wantErr bool
	}{
		{
			name:    "success",
			setup:   func(a *testutil.FakeAuther) {},
			wantErr: false,
		},
		{
			name: "parse token error",
			setup: func(a *testutil.FakeAuther) {
				a.ParseTokenFn = func(token, tokenType string) (string, string, error) {
					return "", "", auth.ErrInvalidToken
				}
			},
			wantErr: true,
		},
		{
			name: "generate token error",
			setup: func(a *testutil.FakeAuther) {
				a.GenerateTokenFn = func(userID, username string) (auth.TokenInfo, error) {
					return nil, testutil.ErrTest
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auther := testutil.NewFakeAuther()
			repo := testutil.NewFakeAuthRepo()
			tt.setup(auther)

			svc := newTestAuthService(auther, repo)
			result, err := svc.TokenUpdate(context.Background(), "refresh-token")

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.AccessToken == "" {
				t.Error("access_token is empty")
			}
		})
	}
}

func TestAuthService_UserInfo(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*testutil.FakeAuther, *testutil.FakeAuthRepo)
		wantErr bool
	}{
		{
			name:    "success",
			setup:   func(a *testutil.FakeAuther, r *testutil.FakeAuthRepo) {},
			wantErr: false,
		},
		{
			name: "parse token error",
			setup: func(a *testutil.FakeAuther, r *testutil.FakeAuthRepo) {
				a.ParseTokenFn = func(token, tokenType string) (string, string, error) {
					return "", "", auth.ErrInvalidToken
				}
			},
			wantErr: true,
		},
		{
			name: "repo error",
			setup: func(a *testutil.FakeAuther, r *testutil.FakeAuthRepo) {
				r.UserInfoFn = func(ctx context.Context, userID string) (model.User, error) {
					return model.User{}, testutil.ErrTest
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auther := testutil.NewFakeAuther()
			repo := testutil.NewFakeAuthRepo()
			tt.setup(auther, repo)

			svc := newTestAuthService(auther, repo)
			info, err := svc.UserInfo(context.Background(), "access-token")

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if info.Username == "" {
				t.Error("username is empty")
			}
		})
	}
}
