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

	"beacon/pkg/auth"
	"beacon/pkg/auth/jwtauth"

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

// minSecretLen 是生产模式下 Agent 准入凭据（JoinToken/installToken）的最小字节长度。
// 这些凭据用于 Agent 注册与监控上报鉴权，不直接签发长期 token，
// 因此阈值低于 JWT 签名密钥（32），但仍要求足够的熵防止穷举。
const minSecretLen = 16

// weakDefaultSecrets 是配置文件与示例中常见的弱默认值黑名单。
// 这些值出现在公开仓库中，任何人都能据此注册 Agent 或上报伪造数据，
// 因此在生产模式下必须拒绝。
var weakDefaultSecrets = map[string]struct{}{
	"amprobe":          {},
	"change-me":        {},
	"change-me-please": {},
	"secret":           {},
	"token":            {},
	"join-token":       {},
}

// isWeakDefaultSecret 报告 s 是否落在弱默认值黑名单中。
func isWeakDefaultSecret(s string) bool {
	_, ok := weakDefaultSecrets[s]
	return ok
}

// ErrInsecureSigningKey 表示生产模式下签名密钥不安全（缺失、弱默认或过短）。
var ErrInsecureSigningKey = errors.New("insecure JWT signing key for production")

// ErrInsecureControlToken 表示生产模式下控制通道（反向 tunnel 注册）JoinToken 不安全。
var ErrInsecureControlToken = errors.New("insecure control JoinToken for production")

// ErrInsecureInstallToken 表示生产模式下 Agent 安装/上报 installToken 不安全。
var ErrInsecureInstallToken = errors.New("insecure agent install token for production")

// resolveSigningKey 处理 JWT 签名密钥的安全降级：
//   - 生产模式（env == production）：空、弱默认或短于 minSigningKeyLen 的密钥一律拒绝，返回 ErrInsecureSigningKey。
//   - 非 production：空则随机生成临时密钥（重启即失效，所有 token 作废）并告警；弱默认保留但告警。
//
// 返回最终应使用的密钥；生产模式下不安全时返回 error，由调用方中止启动。
func resolveSigningKey(configured, env string) (string, error) {
	if strings.EqualFold(env, productionEnv) {
		if configured == "" {
			return "", fmt.Errorf("%w: SigningKey is empty, set BEACON_AUTH_SIGNING_KEY", ErrInsecureSigningKey)
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
		slog.Warn("auth: SigningKey not configured; generated an ephemeral key (all tokens invalid on restart). Set BEACON_AUTH_SIGNING_KEY in production.")
		return key, nil
	}
	if configured == weakDefaultSigningKey {
		slog.Warn("auth: SigningKey is the built-in weak default. Override it via BEACON_AUTH_SIGNING_KEY in production.")
	}
	return configured, nil
}

// resolveControlToken 处理控制通道（反向 tunnel 注册）JoinToken 的安全降级：
//   - 生产模式（env == production）：空、弱默认或短于 minSecretLen 的 token 一律拒绝，
//     返回 ErrInsecureControlToken。控制通道承载远程 shell 等高危调用，Agent 注册
//     必须强鉴权，未配置或弱凭据会允许任意节点注册。
//   - 非 production：空或弱默认仅告警，不破坏本地开发与现有部署。
//
// 返回最终应使用的 token；生产模式下不安全时返回 error，由调用方中止启动。
func resolveControlToken(configured, env string) (string, error) {
	if strings.EqualFold(env, productionEnv) {
		if configured == "" {
			return "", fmt.Errorf("%w: JoinToken is empty, set Control.JoinToken or BEACON_CONTROL_JOIN_TOKEN", ErrInsecureControlToken)
		}
		if isWeakDefaultSecret(configured) {
			return "", fmt.Errorf("%w: JoinToken is a known weak default", ErrInsecureControlToken)
		}
		if len(configured) < minSecretLen {
			return "", fmt.Errorf("%w: JoinToken shorter than %d bytes", ErrInsecureControlToken, minSecretLen)
		}
		return configured, nil
	}
	if configured == "" {
		slog.Warn("auth: control JoinToken not configured; agents can register without authentication. Set Control.JoinToken or BEACON_CONTROL_JOIN_TOKEN in production.")
	} else if isWeakDefaultSecret(configured) {
		slog.Warn("auth: control JoinToken is a known weak default. Override it before deploying to production.")
	}
	return configured, nil
}

// resolveInstallToken 处理 Agent 安装/监控上报 installToken 的安全降级。
// 语义与 resolveControlToken 一致：生产模式拒绝空/弱默认/过短，非生产仅告警。
// 该 token 同时用于安装包下载鉴权（/api/v1/host/install/*）和监控上报鉴权
// （/api/v1/host/report），任一泄露都能让攻击者投递伪造监控数据或下载 Agent 资产。
func resolveInstallToken(configured, env string) (string, error) {
	if strings.EqualFold(env, productionEnv) {
		if configured == "" {
			return "", fmt.Errorf("%w: Token is empty, set AgentInstall.Token or BEACON_AGENT_INSTALL_TOKEN", ErrInsecureInstallToken)
		}
		if isWeakDefaultSecret(configured) {
			return "", fmt.Errorf("%w: Token is a known weak default", ErrInsecureInstallToken)
		}
		if len(configured) < minSecretLen {
			return "", fmt.Errorf("%w: Token shorter than %d bytes", ErrInsecureInstallToken, minSecretLen)
		}
		return configured, nil
	}
	if configured == "" {
		slog.Warn("auth: agent install token not configured; report endpoints reject all uploads. Set AgentInstall.Token or BEACON_AGENT_INSTALL_TOKEN in production.")
	} else if isWeakDefaultSecret(configured) {
		slog.Warn("auth: agent install token is a known weak default. Override it before deploying to production.")
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
