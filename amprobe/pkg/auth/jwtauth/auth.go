// Package jwtauth
// Date: 2024/3/27 15:02
// Author: Amu
// Description:
package jwtauth

import (
	"strings"
	"time"

	"amprobe/pkg/auth"
	"amprobe/service/model"
	"common/database"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuth struct {
	opts  *options
	store Storer
	db    *database.DB
}

func New(store Storer, db *database.DB, opts ...Option) *JWTAuth {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}
	return &JWTAuth{opts: &o, store: store, db: db}
}

func (a *JWTAuth) generateAccessToken(userID string, username string) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(a.opts.signingMethod, jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(a.opts.expired) * time.Second)),
		NotBefore: jwt.NewNumericDate(now),
		Subject:   userID + "." + username,
	})

	tokenString, err := token.SignedString(a.opts.signingKey)
	if err != nil {
		return "", err
	}

	err = a.callStore(func(storer Storer) error {
		return storer.Set(tokenString, time.Duration(a.opts.expired)*time.Second)
	})
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *JWTAuth) generateRefreshToken(userID string, username string) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(a.opts.signingMethod, jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(a.opts.refreshExpired) * time.Second)),
		NotBefore: jwt.NewNumericDate(now),
		Subject:   username + "." + userID,
	})

	tokenString, err := token.SignedString(a.opts.signingKey)
	if err != nil {
		return "", err
	}
	err = a.callStore(func(storer Storer) error {
		return storer.Set(tokenString, time.Duration(a.opts.refreshExpired)*time.Second)
	})
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GenerateToken 生成令牌
func (a *JWTAuth) GenerateToken(userID string, username string) (auth.TokenInfo, error) {
	accessToken, err := a.generateAccessToken(userID, username)
	if err != nil {
		return nil, err
	}
	refreshToken, err := a.generateRefreshToken(userID, username)
	if err != nil {
		return nil, err
	}
	tokenInfo := &tokenInfo{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return tokenInfo, nil
}

// parseToken parses the token string and returns registered claims.
func (a *JWTAuth) parseToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, a.opts.keyfunc)
	if err != nil || !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	return token.Claims.(*jwt.RegisteredClaims), nil
}

func (a *JWTAuth) callStore(fn func(Storer) error) error {
	if store := a.store; store != nil {
		return fn(store)
	}
	return nil
}

// DestroyToken 销毁令牌
func (a *JWTAuth) DestroyToken(tokenString string) error {
	return a.callStore(func(storer Storer) error {
		return storer.Set(tokenString, 1*time.Second)
	})
}

// ParseToken 解析用户 ID, username
func (a *JWTAuth) ParseToken(tokenString string, tokenType string) (string, string, error) {
	if tokenString == "" {
		return "", "", auth.ErrInvalidToken
	}

	err := a.callStore(func(storer Storer) error {
		if exists, err := storer.Check(tokenString); err != nil {
			return err
		} else if !exists {
			return auth.ErrInvalidToken
		}
		return nil
	})

	if err != nil {
		return "", "", err
	}

	claims, err := a.parseToken(tokenString)
	if err != nil {
		return "", "", err
	}

	switch tokenType {
	case "access_token":
		return strings.Split(claims.Subject, ".")[0], strings.Split(claims.Subject, ".")[1], nil
	case "refresh_token":
		return strings.Split(claims.Subject, ".")[1], strings.Split(claims.Subject, ".")[0], nil
	default:
		return "", "", auth.ErrInvalidToken
	}
}

// Release 释放资源
func (a *JWTAuth) Release() error {
	return a.callStore(func(storer Storer) error {
		return storer.Close()
	})
}

func (a *JWTAuth) RecordAudit(username string, operate string) {
	a.db.Model(&model.Audit{}).Create(&model.Audit{Username: username, Operate: operate})
}
