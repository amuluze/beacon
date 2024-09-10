// Package repository
// Date: 2024/6/12 10:30
// Author: Amu
// Description:
package rpc

import (
	"context"
	"github.com/smallnest/rpcx/client"
)

type Client struct {
	client client.XClient
}

func NewClient(addr string) (*Client, error) {
	sf, err := client.NewPeer2PeerDiscovery("unix@"+addr, "")
	if err != nil {
		return nil, err
	}
	xclient := client.NewXClient("Service", client.Failtry, client.RandomSelect, sf, client.DefaultOption)
	return &Client{
		client: xclient,
	}, nil
}

func (c *Client) Call(ctx context.Context, method string, args interface{}, reply interface{}) error {
	return c.client.Call(ctx, method, args, reply)
}

func (c *Client) Close() error {
	return c.client.Close()
}
