// Package service
// Date: 2024/06/11 19:29:11
// Author: Amu
// Description:
package service

import (
	"context"
	"github.com/smallnest/rpcx/client"
)

type RPCClient struct {
	client client.XClient
}

func NewRPCClient(config *Config) (*RPCClient, error) {
	addr := config.Rpc.Address
	sf, err := client.NewPeer2PeerDiscovery("unix@"+addr, "")
	if err != nil {
		return nil, err
	}
	xclient := client.NewXClient("Service", client.Failtry, client.RandomSelect, sf, client.DefaultOption)
	return &RPCClient{
		client: xclient,
	}, nil
}

func (c *RPCClient) Call(ctx context.Context, method string, args interface{}, reply interface{}) error {
	return c.client.Call(ctx, method, args, reply)
}

func (c *RPCClient) Close() error {
	return c.client.Close()
}
