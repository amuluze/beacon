package service

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	ErrHTTPSPublicBaseURLRequired = errors.New("production AgentInstall.PublicBaseURL must use HTTPS")
	ErrControlSocketRequired      = errors.New("control channel requires Control.UnixSocket (unix socket listen path)")
	ErrInstallSocketRequired      = errors.New("agent install requires AgentInstall.ControlSocket (host-visible unix socket path)")
	ErrControlSocketMismatch      = errors.New("control unix socket mode requires both Control.UnixSocket (listen) and AgentInstall.ControlSocket (advertised)")
)

func validateTransportSecurity(config *Config) error {
	// 控制通道只支持 Unix Socket（单机部署）：
	// 服务端监听路径与下发给 Agent 的宿主路径必须成对声明，
	// 缺失任一都会导致监听与 Agent 目标不一致，必须快速失败。
	if config.Control.Enable && strings.TrimSpace(config.Control.UnixSocket) == "" {
		return ErrControlSocketRequired
	}
	if config.AgentInstall.Enable && strings.TrimSpace(config.AgentInstall.ControlSocket) == "" {
		return ErrInstallSocketRequired
	}
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
	return nil
}

func validateHTTPSBaseURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil || !strings.EqualFold(parsed.Scheme, "https") || parsed.Host == "" || parsed.User != nil {
		return fmt.Errorf("%w: set BEACON_PUBLIC_BASE_URL to the externally reachable HTTPS origin", ErrHTTPSPublicBaseURLRequired)
	}
	return nil
}
