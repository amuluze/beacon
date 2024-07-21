// Package profile
// Date: 2024/7/3 09:36
// Author: Amu
// Description:
package profile

import "github.com/mcuadros/go-defaults"

type NetworkProfile struct {
	Name   string            `yaml:"name" default:"app"`
	Driver string            `yaml:"driver" default:"bridge"`
	Ipam   *IPAM             `yaml:"ipam"`
	Labels map[string]string `yaml:"labels,omitempty"`
}

func (n *NetworkProfile) Entry() interface{} {
	return n
}

func (n *NetworkProfile) InitResources(prefix string) error {
	return nil
}

func (n *NetworkProfile) BuildProfile(prefix string) error {
	defaults.SetDefaults(n)
	n.Ipam.Driver = "default"
	n.Ipam.Config = append(n.Ipam.Config, &IPAMConfig{Gateway: BridgeGateway, Subnet: BridgeSubnet})
	return nil
}

type IPAMConfig struct {
	Gateway string `yaml:"gateway"`
	Subnet  string `yaml:"subnet"`
}

type IPAM struct {
	Driver string        `yaml:"driver"`
	Config []*IPAMConfig `yaml:"config"`
}

func NewNetworkBuilder() (ServiceBuilder, error) {
	return &NetworkProfile{Ipam: new(IPAM)}, nil
}
