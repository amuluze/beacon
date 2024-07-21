// Package container
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package container

import (
	"context"
	"github.com/amuluze/amprobe/amvector/pkg/profile"
)

func (m *Manager) CreateNetwork(profile *profile.NetworkProfile) error {
	// 如果存在同名 network，则直接返回
	network, err := m.manager.HasSameNameNetwork(context.TODO(), profile.Name)
	if err != nil {
		return err
	} else if network {
		return nil
	}

	// 创建 network
	ipam := profile.Ipam.Config[0]
	_, err = m.manager.CreateNetwork(context.TODO(), profile.Name, profile.Driver, ipam.Subnet, ipam.Gateway, profile.Labels)
	if err != nil {
		return err
	}
	return nil
}
