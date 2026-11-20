package service

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrControlTLSRequired        = errors.New("control channel requires mutual TLS")
	ErrHTTPSReportTargetRequired = errors.New("authenticated report target must use HTTPS")
)

func validateTransportSecurity(config *Config) error {
	// 单机部署经 Unix Socket 接入时本机链路无需 mTLS；TCP 目标必须启用。
	if config.Control.Server != "" && !isUnixControlTarget(config.Control.Server) && !config.Control.TLS.Enable {
		return ErrControlTLSRequired
	}
	if config.Task.Report.URL != "" && config.Task.Report.Token != "" {
		parsed, err := url.Parse(config.Task.Report.URL)
		if err != nil || !strings.EqualFold(parsed.Scheme, "https") || parsed.Host == "" || parsed.User != nil {
			return ErrHTTPSReportTargetRequired
		}
	}
	return nil
}

// isUnixControlTarget 报告控制通道目标是否为 Unix Socket（unix:// 前缀）。
func isUnixControlTarget(server string) bool {
	return strings.HasPrefix(strings.TrimSpace(server), "unix://")
}
