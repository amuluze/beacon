// Package jwtauth
// Date: 2024/3/27 15:02
// Author: Amu
// Description:
package jwtauth

import (
	"time"

	"beacon/pkg/auth"
	"beacon/service/model"
	"common/database"

	"github.com/golang-jwt/jwt/v5"
)

// accessTokenType / refreshTokenType 标识 token 类型，存放在 userClaims.TokenType 中，
// 取代历史上用 Subject 字段顺序（access: uid.name / refresh: name.uid）区分的脆弱设计。
const (
	accessTokenType  = "access"
	refreshTokenType = "refresh"
)

// userClaims 是 JWT 的自定义 claims，内嵌标准 RegisteredClaims，
// 用结构化字段携带 UserID/Username/TokenType。
// 相比旧的 "userID + \".\" + username" 字符串拼接 Subject：
//   - 消除 username 含 "." 时的身份错乱（旧实现 strings.Split 会取错分段）；
//   - 消除越界 panic 风险（旧实现无长度校验）；
//   - TokenType 显式区分 access/refresh，不依赖 Subject 字段顺序。
type userClaims struct {
	jwt.RegisteredClaims
	UserID    string `json:"uid"`
	Username  string `json:"uname"`
	TokenType string `json:"ttype"`
}

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

	claims := userClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(a.opts.expired) * time.Second)),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID:    userID,
		Username:  username,
		TokenType: accessTokenType,
	}

	token := jwt.NewWithClaims(a.opts.signingMethod, claims)

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

	claims := userClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(a.opts.refreshExpired) * time.Second)),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID:    userID,
		Username:  username,
		TokenType: refreshTokenType,
	}

	token := jwt.NewWithClaims(a.opts.signingMethod, claims)

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

// parseToken 解析 token 字符串并返回自定义 claims。
// 用结构化 userClaims 直接读取 UserID/Username/TokenType，
// 不再依赖 Subject 字符串拼接与 strings.Split。
func (a *JWTAuth) parseToken(tokenString string) (*userClaims, error) {
	claims := &userClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, a.opts.keyfunc)
	if err != nil || !token.Valid {
		return nil, auth.ErrInvalidToken
	}
	return claims, nil
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

// ParseToken 解析用户 ID, username。
// 根据 claims 的 TokenType 字段校验与请求的 tokenType 一致，
// 防止 refresh token 被当作 access token 使用（或反之）。
// 空字段或类型不匹配返回 ErrInvalidToken，不再有越界 panic 风险。
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
	if claims.UserID == "" || claims.Username == "" || claims.TokenType == "" {
		return "", "", auth.ErrInvalidToken
	}

	var wantType string
	switch tokenType {
	case "access_token":
		wantType = accessTokenType
	case "refresh_token":
		wantType = refreshTokenType
	default:
		return "", "", auth.ErrInvalidToken
	}
	if claims.TokenType != wantType {
		return "", "", auth.ErrInvalidToken
	}

	return claims.UserID, claims.Username, nil
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
