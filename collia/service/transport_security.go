package service

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrUnixControlTargetRequired = errors.New("control channel target must be a unix:// socket (single-machine deployment)")
	ErrHTTPSReportTargetRequired = errors.New("authenticated report target must use HTTPS")
)

func validateTransportSecurity(config *Config) error {
	// 单机部署：控制通道只允许本机 Unix Socket 目标，不存在跨主机链路。
	if config.Control.Server != "" && !isUnixControlTarget(config.Control.Server) {
		return ErrUnixControlTargetRequired
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
