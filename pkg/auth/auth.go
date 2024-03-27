// Package auth
// Date: 2024/3/27 15:02
// Author: Amu
// Description:
package auth

import (
	"errors"
)

var ErrInvalidToken = errors.New("invalid token")

type TokenInfo interface {
	GetAccessToken() string
	GetTokenType() string
	GetExpiresAt() string
	EncodeToJSON() ([]byte, error)
}

type Auther interface {
	GenerateToken(userID string) (TokenInfo, error)
	DestroyToken(token string) error
	ParseUserID(token string) (string, error)
	Release() error
}
