// Package jwtauth
// Date: 2024/3/27 15:02
// Author: Amu
// Description:
package jwtauth

import (
	"github.com/amuluze/amprobe/pkg/auth"
	jwt "github.com/dgrijalva/jwt-go"
	"log/slog"
	"time"
)

type JWTAuth struct {
	opts  *options
	store Storer
}

func New(store Storer, opts ...Option) *JWTAuth {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}
	return &JWTAuth{opts: &o, store: store}
}

// GenerateToken 生成令牌
func (a *JWTAuth) GenerateToken(userID string) (auth.TokenInfo, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(a.opts.expired) * time.Second).Unix()
	expiresTime := now.Add(time.Duration(a.opts.expired) * time.Second).Format("2006-01-02T15:04:05.616Z")

	token := jwt.NewWithClaims(a.opts.signingMethod, &jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt,
		NotBefore: now.Unix(),
		Subject:   userID,
	})

	tokenString, err := token.SignedString(a.opts.signingKey)
	if err != nil {
		return nil, err
	}

	err = a.callStore(func(storer Storer) error {
		return storer.Set(tokenString, time.Duration(a.opts.expired)*time.Second)
	})
	if err != nil {
		slog.Error("store token string error: %#v", "err", err)
		return nil, err
	}

	tokenInfo := &tokenInfo{
		ExpiresAt: expiresTime,
		TokenType: a.opts.tokenType,
		Token:     tokenString,
	}
	return tokenInfo, nil
}

// parseToke 解析令牌
func (a *JWTAuth) parseToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, a.opts.keyfunc)
	if err != nil || !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	return token.Claims.(*jwt.StandardClaims), nil
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

// ParseUserID 解析用户 ID
func (a *JWTAuth) ParseUserID(tokenString string) (string, error) {
	if tokenString == "" {
		return "", auth.ErrInvalidToken
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
		return "", err
	}

	claims, err := a.parseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Subject, nil
}

// Release 释放资源
func (a *JWTAuth) Release() error {
	return a.callStore(func(storer Storer) error {
		return storer.Close()
	})
}
