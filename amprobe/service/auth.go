// Package service
// Date: 2024/3/27 16:06
// Author: Amu
// Description:
package service

import (
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"time"

	"amprobe/pkg/auth"
	"amprobe/pkg/auth/jwtauth"

	"common/database"

	"github.com/golang-jwt/jwt/v5"
	"github.com/patrickmn/go-cache"
)

// weakDefaultSigningKey 是配置文件内置的弱签名密钥，仅用于本地开发。
const weakDefaultSigningKey = "amprobe"

// resolveSigningKey 处理 JWT 签名密钥的安全降级：
//   - 空：随机生成临时密钥（重启即失效，所有 token 作废），并告警要求生产通过环境变量固化。
//   - 弱默认值：保留但告警，提示生产必须覆盖。
//
// 返回最终应使用的密钥。
func resolveSigningKey(configured string) string {
	if configured == "" {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			// rand.Read 失败极少见；用时间兜底避免空密钥，仍触发告警。
			slog.Error("auth: generate ephemeral signing key failed, falling back", "err", err)
			return time.Now().String()
		}
		key := base64.StdEncoding.EncodeToString(b)
		slog.Warn("auth: SigningKey not configured; generated an ephemeral key (all tokens invalid on restart). Set AMPROBE_AUTH_SIGNINGKEY in production.")
		return key
	}
	if configured == weakDefaultSigningKey {
		slog.Warn("auth: SigningKey is the built-in weak default. Override it via AMPROBE_AUTH_SIGNINGKEY in production.")
	}
	return configured
}

func InitAuthStore(config *Config) (*jwtauth.Store, func(), error) {
	var err error
	authStore := &jwtauth.Store{
		Storage: cache.New(5*time.Minute, 60*time.Second),
		Prefix:  config.Auth.Prefix,
	}
	cleanFunc := func() { err = authStore.Close() }
	return authStore, cleanFunc, err
}

func InitAuth(config *Config, authStore *jwtauth.Store, db *database.DB) (auth.Auther, func(), error) {
	signingKey := resolveSigningKey(config.Auth.SigningKey)
	var opts []jwtauth.Option
	opts = append(opts, jwtauth.SetExpired(config.Auth.Expired))
	opts = append(opts, jwtauth.SetRefreshExpired(config.Auth.RefreshExpired))
	opts = append(opts, jwtauth.SetSigningKey([]byte(signingKey)))
	opts = append(opts, jwtauth.SetKeyfunc(func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidToken
		}
		return []byte(signingKey), nil
	}))
	var method jwt.SigningMethod
	switch config.Auth.SigningMethod {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	default:
		method = jwt.SigningMethodHS512
	}
	opts = append(opts, jwtauth.SetSigningMethod(method))
	var err error
	jAuth := jwtauth.New(authStore, db, opts...)
	cleanFunc := func() {
		err = jAuth.Release()
	}
	return jAuth, cleanFunc, err
}
