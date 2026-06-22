// Package rpc
// Date: 2024/6/12 10:30
// Author: Amu
// Description:
package rpc

import (
	"context"
	"crypto/tls"
	"fmt"

	"amprobe/pkg/contextx"

	"github.com/smallnest/rpcx/client"
)

const (
	DefaultAgentID = "default"
	DefaultNetwork = "unix"
)

type Agent struct {
	ID        string
	Network   string
	Address   string
	TLSConfig *tls.Config
}

type Client struct {
	defaultAgentID string
	clients        map[string]client.XClient
}

func NewClient(addr string) (*Client, error) {
	return NewClientWithNetwork(DefaultNetwork, addr)
}

func NewClientWithNetwork(network string, addr string) (*Client, error) {
	return NewMultiClient(DefaultAgentID, []Agent{
		{
			ID:      DefaultAgentID,
			Network: network,
			Address: addr,
		},
	})
}

func NewMultiClient(defaultAgentID string, agents []Agent) (*Client, error) {
	if defaultAgentID == "" {
		defaultAgentID = DefaultAgentID
	}
	clients := make(map[string]client.XClient, len(agents))
	for _, agent := range agents {
		if agent.ID == "" {
			agent.ID = DefaultAgentID
		}
		if agent.Network == "" {
			agent.Network = DefaultNetwork
		}
		if agent.Address == "" {
			return nil, fmt.Errorf("empty rpc address for agent %q", agent.ID)
		}
		xclient, err := newXClient(agent.Network, agent.Address, agent.TLSConfig)
		if err != nil {
			closeClients(clients)
			return nil, err
		}
		clients[agent.ID] = xclient
	}
	if _, ok := clients[defaultAgentID]; !ok {
		closeClients(clients)
		return nil, fmt.Errorf("default agent %q not configured", defaultAgentID)
	}
	return &Client{
		defaultAgentID: defaultAgentID,
		clients:        clients,
	}, nil
}

func newXClient(network string, addr string, tlsConfig *tls.Config) (client.XClient, error) {
	sf, err := client.NewPeer2PeerDiscovery(network+"@"+addr, "")
	if err != nil {
		return nil, err
	}
	opt := client.DefaultOption
	opt.TLSConfig = tlsConfig
	return client.NewXClient("Service", client.Failtry, client.RandomSelect, sf, opt), nil
}

func closeClients(clients map[string]client.XClient) {
	for _, xclient := range clients {
		_ = xclient.Close()
	}
}

func (c *Client) Call(ctx context.Context, method string, args interface{}, reply interface{}) error {
	agentID := contextx.FromAgentID(ctx)
	if agentID == "" {
		agentID = c.defaultAgentID
	}
	xclient, ok := c.clients[agentID]
	if !ok {
		return fmt.Errorf("agent %q not configured", agentID)
	}
	return xclient.Call(ctx, method, args, reply)
}

func (c *Client) Close() error {
	var firstErr error
	for _, xclient := range c.clients {
		if err := xclient.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
