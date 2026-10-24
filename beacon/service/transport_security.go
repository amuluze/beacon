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
)

func validateTransportSecurity(config *Config) error {
	if !strings.EqualFold(config.App.Env, productionEnv) {
		return nil
	}
	if config.AgentInstall.Enable {
		if err := validateHTTPSBaseURL(config.AgentInstall.PublicBaseURL); err != nil {
			return err
		}
	}
	if config.Control.Enable && !config.Control.TLS.Enable {
		return ErrControlTLSRequired
	}
	if config.Control.Enable && !config.AgentInstall.TLSEnable {
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
