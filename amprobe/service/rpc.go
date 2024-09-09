// Package service
// Date: 2024/06/11 19:29:11
// Author: Amu
// Description:
package service

import (
	"amprobe/pkg/rpc"
)

func NewRPCClient(config *Config) (*rpc.Client, error) {
	addr := config.Rpc.Address
	return rpc.NewClient(addr)
}
