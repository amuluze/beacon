// Package docker
// Date: 2024/07/09 14:13:56
// Author: Amu
// Description:
package docker

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/docker/docker/api/types/network"

	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/docker/libcompose/yaml"
	"github.com/tidwall/gjson"
	goyaml "gopkg.in/yaml.v3"
)

type ContainerSummary struct {
	ID           string            `json:"id"`      // ID
	Name         string            `json:"name"`    // Name
	Image        string            `json:"image"`   // Image
	Network      string            `json:"network"` // Network
	State        string            `json:"state"`   // State: created running paused restarting removing exited dead
	Created      string            `json:"created"` // create time
	Uptime       string            `json:"uptime"`  // uptime in seconds
	IP           string            `json:"ip"`      // ip
	Ports        []string          `json:"ports"`
	Volumes      []string          `json:"volumes"`
	Environments []string          `json:"environments"`
	Labels       map[string]string `json:"labels"`
}

type PortMapping struct {
	Proto         string
	IP            string
	HostPort      string
	ContainerPort string
}

// getUptime 获取指定容器的启动时间
func (m *Manager) getUptime(ctx context.Context, containerID string) string {
	inspect, _ := m.client.ContainerInspect(ctx, containerID)
	started, _ := time.Parse(time.RFC3339Nano, inspect.State.StartedAt)
	return started.Format("2006-01-02 15:04:05")
}

func (m *Manager) ListContainer(ctx context.Context) ([]ContainerSummary, error) {
	containers, err := m.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var containerSummaryList []ContainerSummary
	for _, c := range containers {
		var ports []string
		for _, port := range c.Ports {
			ports = append(ports, strconv.Itoa(int(port.PublicPort)))
		}
		var uptime string
		if c.State == "running" {
			uptime = m.getUptime(ctx, c.ID)
		}

		var ip string
		var networkName string
		for ntName, nt := range c.NetworkSettings.Networks {
			if nt.IPAddress != "" {
				ip = nt.IPAddress
				networkName = ntName
				break
			}
		}

		state := c.State
		inspect, err := m.client.ContainerInspect(ctx, c.ID)
		if err == nil {
			if inspect.ContainerJSONBase.State.Health != nil && inspect.ContainerJSONBase.State.Health.Status == "healthy" {
				state = "running"
			}
		}

		var volumes []string
		for _, mount := range c.Mounts {
			if mount.Driver != "" {
				volumes = append(volumes, fmt.Sprintf("%s:%s:%s", mount.Source, mount.Destination, mount.Driver))
			} else {
				volumes = append(volumes, fmt.Sprintf("%s:%s", mount.Source, mount.Destination))
			}
		}

		var environments []string
		inspect, _ = m.client.ContainerInspect(ctx, c.ID)
		environments = inspect.Config.Env

		cs := ContainerSummary{
			ID:           c.ID,
			Name:         strings.Trim(c.Names[0], "/"),
			Image:        c.Image,
			Network:      networkName,
			State:        state,
			Created:      time.Unix(c.Created, 0).Format("2006-01-02 15:04:05"),
			Uptime:       uptime,
			IP:           ip,
			Ports:        ports,
			Volumes:      volumes,
			Environments: environments,
			Labels:       c.Labels,
		}
		containerSummaryList = append(containerSummaryList, cs)
	}
	return containerSummaryList, nil
}

func (m *Manager) HasSameNameContainer(ctx context.Context, containerName string) (bool, error) {
	containers, err := m.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return false, err
	}
	for _, c := range containers {
		if strings.Trim(c.Names[0], "/") == containerName {
			return true, nil
		}
	}
	return false, nil
}

func (m *Manager) CreateContainer(ctx context.Context, containerName, imageName, networkName string, ports []string, vols []string, envs []string, commands []string, labels map[string]string) (string, error) {
	config := &container.Config{}
	config.Hostname = containerName
	config.Image = imageName
	config.Labels = labels
	config.Tty = true
	config.Env = envs
	if commands != nil {
		config.Cmd = commands
	}

	hostConfig := &container.HostConfig{}
	hostConfig.RestartPolicy = container.RestartPolicy{Name: "always"}
	hostConfig.PortBindings = make(nat.PortMap)

	nt, err := m.GetNetworkByName(ctx, networkName)
	if err != nil {
		return "", err
	}
	hostConfig.NetworkMode = container.NetworkMode(networkName)
	networkConfig := &network.NetworkingConfig{}
	networkConfig.EndpointsConfig = make(map[string]*network.EndpointSettings)
	networkConfig.EndpointsConfig[networkName] = &network.EndpointSettings{
		NetworkID: nt.ID,
	}

	for _, port := range ports {
		portsMapping, err := nat.ParsePortSpec(port)
		if err != nil {
			return "", err
		}
		for _, portMapping := range portsMapping {
			port, err := nat.NewPort(portMapping.Port.Proto(), portMapping.Port.Port())
			if err != nil {
				return "", err
			}
			hostIP := portMapping.Binding.HostIP
			if hostIP == "" {
				hostIP = "0.0.0.0"
			}
			hostConfig.PortBindings[port] = append(hostConfig.PortBindings[port], nat.PortBinding{
				HostIP:   hostIP,
				HostPort: portMapping.Binding.HostPort,
			})
		}
	}

	config.ExposedPorts = make(nat.PortSet)
	for port := range hostConfig.PortBindings {
		config.ExposedPorts[port] = struct{}{}
	}

	for _, vol := range vols {
		vol := "- " + vol
		volumes := &yaml.Volumes{}

		err := goyaml.Unmarshal([]byte(vol), volumes)
		if err != nil {
			return "", err
		}
		for _, volume := range volumes.Volumes {
			if volume.AccessMode != "ro" {
				volume.AccessMode = "rw"
			}
			volString := fmt.Sprintf("%s:%s:%s", volume.Source, volume.Destination, volume.AccessMode)
			hostConfig.Binds = append(hostConfig.Binds, volString)
		}
	}

	createResponse, err := m.client.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		return "", err
	}
	for _, w := range createResponse.Warnings {
		fmt.Printf("Container Create Warning: %s\n", w)
	}
	if err := m.client.NetworkConnect(ctx, nt.ID, createResponse.ID, nil); err != nil {
		return "", err
	}

	return createResponse.ID, nil
}

func (m *Manager) StartContainer(ctx context.Context, containerID string) error {
	return m.client.ContainerStart(ctx, containerID, container.StartOptions{})
}

func (m *Manager) StopContainer(ctx context.Context, containerID string) error {
	return m.client.ContainerStop(ctx, containerID, container.StopOptions{})
}

func (m *Manager) RestartContainer(ctx context.Context, containerID string) error {
	return m.client.ContainerRestart(ctx, containerID, container.StopOptions{})
}

func (m *Manager) DeleteContainer(ctx context.Context, containerID string) error {
	return m.client.ContainerRemove(ctx, containerID, container.RemoveOptions{
		Force:         true,
		RemoveLinks:   false,
		RemoveVolumes: false,
	})
}

func (m *Manager) CopyFileToContainer(ctx context.Context, containerID string, srcFile, dstFile string) error {
	file, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	return m.client.CopyToContainer(ctx, containerID, dstFile, file, container.CopyToContainerOptions{})
}

func (m *Manager) GetContainerMem(ctx context.Context, containerID string) (float64, float64, float64, error) {
	stats, err := m.client.ContainerStats(ctx, containerID, false)
	if err != nil {
		return 0.0, 0.0, 0.0, err
	}
	defer stats.Body.Close()
	body, err := io.ReadAll(stats.Body)
	if err != nil {
		return 0.0, 0.0, 0.0, err
	}
	memUsage := gjson.Get(string(body), "memory_stats.usage").Float()
	memLimit := gjson.Get(string(body), "memory_stats.limit").Float()
	memPercent := percentage(memUsage, memLimit)
	return memPercent, memUsage, memLimit, nil
}

func (m *Manager) GetContainerCpu(ctx context.Context, containerID string) (float64, error) {
	stats, err := m.client.ContainerStats(ctx, containerID, false)
	if err != nil {
		return 0.0, err
	}
	defer stats.Body.Close()
	body, err := io.ReadAll(stats.Body)
	if err != nil {
		return 0.0, err
	}

	cpuDelta := gjson.Get(string(body), "cpu_stats.cpu_usage.total_usage").Float() - gjson.Get(string(body), "precpu_stats.cpu_usage.total_usage").Float()
	systemDelta := gjson.Get(string(body), "cpu_stats.system_cpu_usage").Float() - gjson.Get(string(body), "precpu_stats.system_cpu_usage").Float()
	cpuPercent := percentage(cpuDelta, systemDelta)
	return cpuPercent, nil
}

// percentage converts an Engine stats ratio to a JSON-safe percentage. Docker
// returns zero-valued stats for containers that have not started yet, so a
// direct division would produce NaN and make the Agent's whole report fail to
// marshal.
func percentage(value, total float64) float64 {
	if value <= 0 || total <= 0 {
		return 0
	}
	result := value / total * 100
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0
	}
	return result
}

func (m *Manager) GetContainerIDByContainerName(ctx context.Context, containerName string) (string, error) {
	containers, err := m.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return "", err
	}
	for _, ct := range containers {
		if ct.Names[0] == containerName {
			return ct.ID, nil
		}
	}
	return "", nil
}

func (m *Manager) ContainerLogs(ctx context.Context, containerID string) (io.ReadCloser, error) {
	reader, err := m.client.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
		Tail:       "any",
	})
	return reader, err
}

func (m *Manager) RenameContainer(ctx context.Context, containerID, newName string) error {
	return m.client.ContainerRename(ctx, containerID, newName)
}

func (m *Manager) ContainerExists(ctx context.Context, containerID string) (bool, error) {
	_, err := m.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return false, err
	}
	return true, nil
}
