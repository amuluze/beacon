// Package rpc
// Date: 2024/6/24 09:30
// Author: Amu
// Description:
package rpc

import (
	"context"
	"github.com/amuluze/amvector/service/schema"
	"github.com/docker/docker/libnetwork/resolvconf"
)

func (s *Service) GetDNS(ctx context.Context, args schema.GetDNSArgs, reply *schema.GetDNSReply) error {
	resolvConf, err := resolvconf.Get()
	if err != nil {
		return err
	}
	reply.DNS = resolvconf.GetNameservers(resolvConf.Content, resolvconf.IP)
	return nil
}

func (s *Service) SetDNS(ctx context.Context, args schema.SetDNSArgs, reply *schema.SetDNSReply) error {
	_, err := resolvconf.Build(resolvconf.Path(), args.DNS, []string{}, []string{})
	if err != nil {
		return err
	}
	return nil
}
