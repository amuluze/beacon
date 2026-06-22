// Package service
// Date: 2024/06/11 19:29:11
// Author: Amu
// Description:
package service

import (
	"amprobe/pkg/rpc"
	transporttls "common/transport/tlsconfig"
)

func NewRPCClient(config *Config) (*rpc.Client, error) {
	if len(config.Rpc.Agents) == 0 {
		network := config.Rpc.Network
		if network == "" {
			network = rpc.DefaultNetwork
		}
		agent, err := buildAgent(rpc.DefaultAgentID, network, config.Rpc.Address, config.Rpc.TLS, TLS{})
		if err != nil {
			return nil, err
		}
		return rpc.NewMultiClient(rpc.DefaultAgentID, []rpc.Agent{agent})
	}

	agents := make([]rpc.Agent, 0, len(config.Rpc.Agents))
	for _, agent := range config.Rpc.Agents {
		network := agent.Network
		if network == "" {
			network = config.Rpc.Network
		}
		if network == "" {
			network = rpc.DefaultNetwork
		}
		rpcAgent, err := buildAgent(agent.ID, network, agent.Address, config.Rpc.TLS, agent.TLS)
		if err != nil {
			return nil, err
		}
		agents = append(agents, rpcAgent)
	}
	return rpc.NewMultiClient(config.Rpc.DefaultAgentID, agents)
}

func buildAgent(id string, network string, address string, globalTLS TLS, agentTLS TLS) (rpc.Agent, error) {
	tlsSettings := mergeTLS(globalTLS, agentTLS)
	agent := rpc.Agent{
		ID:      id,
		Network: network,
		Address: address,
	}
	if !tlsSettings.Enable {
		return agent, nil
	}

	tlsConfig, err := transporttls.ClientConfig(tlsSettings.CertDir, tlsSettings.ServerName)
	if err != nil {
		return rpc.Agent{}, err
	}
	agent.TLSConfig = tlsConfig
	return agent, nil
}

func mergeTLS(globalTLS TLS, agentTLS TLS) TLS {
	out := globalTLS
	if agentTLS.Enable {
		out.Enable = true
	}
	if agentTLS.CertDir != "" {
		out.CertDir = agentTLS.CertDir
	}
	if agentTLS.ServerName != "" {
		out.ServerName = agentTLS.ServerName
	}
	return out
}
