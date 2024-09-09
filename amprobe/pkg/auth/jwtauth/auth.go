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

	"github.com/amuluze/amutool/database"
	"github.com/golang-jwt/jwt"
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

func (a *JWTAuth) generateAccessToken(userID string, username string, isAdmin string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(a.opts.expired) * time.Second).Unix()

	token := jwt.NewWithClaims(a.opts.signingMethod, &jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt,
		NotBefore: now.Unix(),
		Subject:   userID + "." + username + "." + isAdmin,
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

func (a *JWTAuth) generateRefreshToken(userID string, username string, isAdmin string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(a.opts.refreshExpired) * time.Second).Unix()

	token := jwt.NewWithClaims(a.opts.signingMethod, &jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt,
		NotBefore: now.Unix(),
		Subject:   username + "." + userID + "." + isAdmin,
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
func (a *JWTAuth) GenerateToken(userID string, username string, isAdmin string) (auth.TokenInfo, error) {
	accessToken, err := a.generateAccessToken(userID, username, isAdmin)
	if err != nil {
		return nil, err
	}
	refreshToken, err := a.generateRefreshToken(userID, username, isAdmin)
	if err != nil {
		return nil, err
	}
	tokenInfo := &tokenInfo{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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

// ParseToken 解析用户 ID, username
func (a *JWTAuth) ParseToken(tokenString string, tokenType string) (string, string, string, error) {
	if tokenString == "" {
		return "", "", "", auth.ErrInvalidToken
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
		return "", "", "", err
	}

	claims, err := a.parseToken(tokenString)
	if err != nil {
		return "", "", "", err
	}

	switch tokenType {
	case "access_token":
		return strings.Split(claims.Subject, ".")[0], strings.Split(claims.Subject, ".")[1], strings.Split(claims.Subject, ".")[2], nil
	case "refresh_token":
		return strings.Split(claims.Subject, ".")[1], strings.Split(claims.Subject, ".")[0], strings.Split(claims.Subject, ".")[2], nil
	default:
		return "", "", "", auth.ErrInvalidToken
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
