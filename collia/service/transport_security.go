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
	if config.Control.Server != "" && !config.Control.TLS.Enable {
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
