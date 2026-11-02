// Package auth
// Date: 2024/3/27 15:02
// Author: Amu
// Description:
package auth

import (
	"context"
	"errors"
)

var ErrInvalidToken = errors.New("invalid token")

type TokenInfo interface {
	GetAccessToken() string
	GetRefreshToken() string
}

type Auther interface {
	GenerateToken(userID string, username string) (TokenInfo, error)
	DestroyToken(token string) error
	ParseToken(token string, tokenType string) (string, string, error)
	Release() error
	RecordAudit(ctx context.Context, username string, operate string, object string)
}
