// Package service
// Date: 2024/3/27 16:06
// Author: Amu
// Description:
package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"amprobe/pkg/auth"
	"amprobe/pkg/auth/jwtauth"

	"common/database"

	"github.com/golang-jwt/jwt/v5"
	"github.com/patrickmn/go-cache"
)

// weakDefaultSigningKey 是配置文件内置的弱签名密钥，仅用于本地开发。
const weakDefaultSigningKey = "amprobe"

// productionEnv 标识生产运行模式，触发敏感配置的严格校验。
const productionEnv = "production"

// minSigningKeyLen 是生产模式下签名密钥的最小字节长度，低于此值视为不安全。
const minSigningKeyLen = 32

// ErrInsecureSigningKey 表示生产模式下签名密钥不安全（缺失、弱默认或过短）。
var ErrInsecureSigningKey = errors.New("insecure JWT signing key for production")

// resolveSigningKey 处理 JWT 签名密钥的安全降级：
//   - 生产模式（env == production）：空、弱默认或短于 minSigningKeyLen 的密钥一律拒绝，返回 ErrInsecureSigningKey。
//   - 非 production：空则随机生成临时密钥（重启即失效，所有 token 作废）并告警；弱默认保留但告警。
//
// 返回最终应使用的密钥；生产模式下不安全时返回 error，由调用方中止启动。
func resolveSigningKey(configured, env string) (string, error) {
	if strings.EqualFold(env, productionEnv) {
		if configured == "" {
			return "", fmt.Errorf("%w: SigningKey is empty, set AMPROBE_AUTH_SIGNINGKEY", ErrInsecureSigningKey)
		}
		if configured == weakDefaultSigningKey {
			return "", fmt.Errorf("%w: SigningKey is the built-in weak default", ErrInsecureSigningKey)
		}
		if len(configured) < minSigningKeyLen {
			return "", fmt.Errorf("%w: SigningKey shorter than %d bytes", ErrInsecureSigningKey, minSigningKeyLen)
		}
		return configured, nil
	}

	if configured == "" {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			// rand.Read 失败极少见；用时间兜底避免空密钥，仍触发告警。
			slog.Error("auth: generate ephemeral signing key failed, falling back", "err", err)
			return time.Now().String(), nil
		}
		key := base64.StdEncoding.EncodeToString(b)
		slog.Warn("auth: SigningKey not configured; generated an ephemeral key (all tokens invalid on restart). Set AMPROBE_AUTH_SIGNINGKEY in production.")
		return key, nil
	}
	if configured == weakDefaultSigningKey {
		slog.Warn("auth: SigningKey is the built-in weak default. Override it via AMPROBE_AUTH_SIGNINGKEY in production.")
	}
	return configured, nil
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
	signingKey, err := resolveSigningKey(config.Auth.SigningKey, config.App.Env)
	if err != nil {
		return nil, nil, err
	}
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
	jAuth := jwtauth.New(authStore, db, opts...)
	cleanFunc := func() {
		err = jAuth.Release()
	}
	return jAuth, cleanFunc, err
}
