package service

import (
	"errors"
	"testing"
)

func TestValidateTransportSecurityProduction(t *testing.T) {
	secure := &Config{
		App:     App{Env: productionEnv},
		Control: Control{Enable: true, TLS: ControlTLS{Enable: true}},
		AgentInstall: AgentInstall{
			Enable:        true,
			PublicBaseURL: "https://beacon.example.com",
			TLSEnable:     true,
		},
	}
	if err := validateTransportSecurity(secure); err != nil {
		t.Fatalf("secure config rejected: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*Config)
		want   error
	}{
		{name: "missing public URL", mutate: func(c *Config) { c.AgentInstall.PublicBaseURL = "" }, want: ErrHTTPSPublicBaseURLRequired},
		{name: "HTTP public URL", mutate: func(c *Config) { c.AgentInstall.PublicBaseURL = "http://beacon.example.com" }, want: ErrHTTPSPublicBaseURLRequired},
		{name: "control TLS disabled", mutate: func(c *Config) { c.Control.TLS.Enable = false }, want: ErrControlTLSRequired},
		{name: "agent TLS disabled", mutate: func(c *Config) { c.AgentInstall.TLSEnable = false }, want: ErrAgentTLSRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configCopy := *secure
			controlCopy := secure.Control
			agentCopy := secure.AgentInstall
			configCopy.Control = controlCopy
			configCopy.AgentInstall = agentCopy
			tt.mutate(&configCopy)
			if err := validateTransportSecurity(&configCopy); !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
		})
	}
}
