// Package rpc
// Date: 2024/6/12 10:30
// Author: Amu
// Description:
package rpc

import (
	"context"
	"fmt"
	"os"

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

func (c *Client) SendFile(ctx context.Context, filename string, prefix string) error {
	return c.client.SendFile(ctx, filename, 0, map[string]string{"path": prefix})
}

func (c *Client) DownloadFile(ctx context.Context, downlaodFilepath string, filename string) error {
	file, _ := os.Create(fmt.Sprintf("/tmp/%s", filename))
	defer file.Close()
	return c.client.DownloadFile(ctx, downlaodFilepath, file, map[string]string{})
}

func (c *Client) Close() error {
	return c.client.Close()
}
