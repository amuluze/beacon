// Package docker
// Date: 2024/07/09 14:14:31
// Author: Amu
// Description:
package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"strings"
)

type NetworkSummary struct {
	ID         string
	Name       string
	Driver     string
	Scope      string
	Created    string
	Internal   bool
	SubNet     []SubNetworkConfig
	Containers map[string]string // map[cid]ipaddr
	Labels     map[string]string
}

type SubNetworkConfig struct {
	Subnet  string
	Gateway string
}

func (m *Manager) ListNetwork(ctx context.Context) ([]NetworkSummary, error) {
	nets, err := m.client.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return nil, err
	}

	var networkList []NetworkSummary
	for _, net := range nets {
		containers := make(map[string]string)
		for id, container := range net.Containers {
			ipAddr := container.IPv4Address
			if slashIdx := strings.IndexByte(ipAddr, '/'); slashIdx != -1 {
				ipAddr = ipAddr[:slashIdx]
			}
			containers[id] = ipAddr
		}
		subNet := make([]SubNetworkConfig, 0)
		for _, ncf := range net.IPAM.Config {
			subNet = append(subNet, SubNetworkConfig{
				Subnet:  ncf.Subnet,
				Gateway: ncf.Gateway,
			})
		}
		n := NetworkSummary{
			ID:         net.ID,
			Name:       net.Name,
			Driver:     net.Driver,
			Scope:      net.Scope,
			Created:    net.Created.Format("2006-01-02 15:04:05"),
			SubNet:     subNet,
			Containers: containers,
			Labels:     net.Labels,
		}
		networkList = append(networkList, n)
	}
	return networkList, nil
}

func (m *Manager) HasSameNameNetwork(ctx context.Context, networkName string) (bool, error) {
	nets, err := m.client.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, net := range nets {
		if net.Name == networkName {
			return true, nil
		}
	}
	return false, err
}

func (m *Manager) CreateNetwork(ctx context.Context, name, driver, subnet, gateway string, labels map[string]string) (string, error) {
	nt, err := m.client.NetworkCreate(ctx, name, network.CreateOptions{
		Driver:     driver,
		IPAM:       &network.IPAM{Config: []network.IPAMConfig{{Subnet: subnet, Gateway: gateway}}},
		Labels:     labels,
		Internal:   false,
		Attachable: true,
	})
	return nt.ID, err
}

func (m *Manager) GetNetworkByName(ctx context.Context, name string) (*NetworkSummary, error) {
	networks, err := m.ListNetwork(ctx)
	if err != nil {
		return nil, err
	}
	for _, nt := range networks {
		if nt.Name == name {
			return &NetworkSummary{
				ID:         nt.ID,
				Name:       nt.Name,
				Driver:     nt.Driver,
				Scope:      nt.Scope,
				Created:    nt.Created,
				Containers: nt.Containers,
				Labels:     nt.Labels,
				SubNet:     nt.SubNet,
			}, nil
		}
	}
	return nil, fmt.Errorf("nt %s not found", name)
}

func (m *Manager) GetNetworkByID(ctx context.Context, networkID string) (*NetworkSummary, error) {
	nr, err := m.client.NetworkInspect(ctx, networkID, network.InspectOptions{})
	if err != nil {
		return nil, err
	}
	containers := make(map[string]string)
	for id, container := range nr.Containers {
		ipAddr := container.IPv4Address
		if slashIdx := strings.IndexByte(ipAddr, '/'); slashIdx != -1 {
			ipAddr = ipAddr[:slashIdx]
		}
		containers[id] = ipAddr
	}
	subNet := make([]SubNetworkConfig, 0)
	for _, ncf := range nr.IPAM.Config {
		subNet = append(subNet, SubNetworkConfig{
			Subnet:  ncf.Subnet,
			Gateway: ncf.Gateway,
		})
	}
	nw := &NetworkSummary{
		ID:         nr.ID,
		Name:       nr.Name,
		Driver:     nr.Driver,
		Scope:      nr.Scope,
		Created:    nr.Created.Format("2006-01-02 15:04:05"),
		SubNet:     subNet,
		Containers: containers,
		Labels:     nr.Labels,
	}
	return nw, nil
}

func (m *Manager) DeleteNetwork(ctx context.Context, networkID string) error {
	return m.client.NetworkRemove(ctx, networkID)
}

func (m *Manager) PruneNetwork(ctx context.Context) error {
	_, err := m.client.NetworksPrune(ctx, filters.NewArgs(filters.Arg("until", "0")))
	return err
}

func (m *Manager) JoinNetwork(ctx context.Context, containerID, networkID string) error {
	if _, err := m.client.NetworkInspect(ctx, networkID, network.InspectOptions{}); err != nil {
		return err
	}
	if _, err := m.client.ContainerInspect(ctx, containerID); err != nil {
		return err
	}
	return m.client.NetworkConnect(ctx, networkID, containerID, &network.EndpointSettings{})
}

func (m *Manager) LeaveNetwork(ctx context.Context, containerID, networkID string) error {
	if _, err := m.client.NetworkInspect(ctx, networkID, network.InspectOptions{}); err != nil {
		return err
	}
	if _, err := m.client.ContainerInspect(ctx, containerID); err != nil {
		return err
	}
	return m.client.NetworkDisconnect(ctx, networkID, containerID, true)
}
