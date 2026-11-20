package service

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	ErrHTTPSPublicBaseURLRequired = errors.New("production AgentInstall.PublicBaseURL must use HTTPS")
	ErrControlTLSRequired         = errors.New("production control channel requires mutual TLS")
	ErrAgentTLSRequired           = errors.New("production agent certificate provisioning must be enabled")
	ErrControlSocketMismatch      = errors.New("control unix socket mode requires both Control.UnixSocket (listen) and AgentInstall.ControlSocket (advertised)")
)

func validateTransportSecurity(config *Config) error {
	// 单机 UDS 模式的两个路径必须成对：容器监听路径 + 下发给 Agent 的宿主路径。
	// 只配一个会导致服务端监听与 Agent 目标不一致，必须快速失败。
	hasListenSocket := strings.TrimSpace(config.Control.UnixSocket) != ""
	hasAdvertisedSocket := strings.TrimSpace(config.AgentInstall.ControlSocket) != ""
	if hasListenSocket != hasAdvertisedSocket {
		return ErrControlSocketMismatch
	}
	if !strings.EqualFold(config.App.Env, productionEnv) {
		return nil
	}
	if config.AgentInstall.Enable {
		if err := validateHTTPSBaseURL(config.AgentInstall.PublicBaseURL); err != nil {
			return err
		}
	}
	// 单机 Unix Socket 模式：本机链路免 mTLS，也不下发 Agent 证书；
	// 多机 TCP 模式仍强制 mTLS 与证书下发（IU009）。
	socketMode := strings.TrimSpace(config.Control.UnixSocket) != ""
	if config.Control.Enable && !socketMode && !config.Control.TLS.Enable {
		return ErrControlTLSRequired
	}
	if config.Control.Enable && !socketMode && !config.AgentInstall.TLSEnable {
		return ErrAgentTLSRequired
	}
	return nil
}

func validateHTTPSBaseURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil || !strings.EqualFold(parsed.Scheme, "https") || parsed.Host == "" || parsed.User != nil {
		return fmt.Errorf("%w: set BEACON_PUBLIC_BASE_URL to the externally reachable HTTPS origin", ErrHTTPSPublicBaseURLRequired)
	}
	return nil
}
