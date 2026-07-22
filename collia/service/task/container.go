// Package task
package task

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"
)

func (a *Task) DockerTask() *DockerReport {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	version, err := a.manager.Version(ctx)
	if err != nil {
		slog.Error("docker version failed", "error", err)
		return nil
	}
	return &DockerReport{
		DockerVersion: version.DockerVersion,
		APIVersion:    version.APIVersion,
		MinAPIVersion: version.MinAPIVersion,
		GitCommit:     version.GitCommit,
		GoVersion:     version.GoVersion,
		Os:            version.OS,
		Arch:          version.Arch,
	}
}

func (a *Task) ContainerTask() []*ContainerReport {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cs, err := a.manager.ListContainer(ctx)
	if err != nil {
		slog.Error("failed to list containers", "error", err)
		return nil
	}
	var containers []*ContainerReport
	for _, info := range cs {
		var cpuPercent, memPercent, used, limit float64
		if info.State == "running" {
			statsCtx, statsCancel := context.WithTimeout(context.Background(), 15*time.Second)
			cpuPercent, err = a.manager.GetContainerCpu(statsCtx, info.ID[:6])
			if err != nil {
				slog.Error("get container cpu failed", "id", info.ID[:6], "error", err)
			}
			memPercent, used, limit, err = a.manager.GetContainerMem(statsCtx, info.ID[:6])
			if err != nil {
				slog.Error("get container mem failed", "id", info.ID[:6], "error", err)
			}
			statsCancel()
		}
		labels, _ := json.Marshal(info.Labels)
		containers = append(containers, &ContainerReport{
			ContainerID: info.ID[:6],
			Name:        info.Name,
			State:       info.State,
			Image:       info.Image,
			Uptime:      info.Uptime,
			IP:          info.IP,
			Ports:       strings.Join(info.Ports, ","),
			Labels:      string(labels),
			CPUPercent:  cpuPercent,
			MemPercent:  memPercent,
			MemUsage:    used,
			MemLimit:    limit,
		})
	}
	return containers
}

func (a *Task) ImageTask() []*ImageReport {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	images, err := a.manager.ListImage(ctx)
	if err != nil {
		slog.Error("list images failed", "error", err)
		return nil
	}
	var reports []*ImageReport
	for _, im := range images {
		reports = append(reports, &ImageReport{
			ImageID: im.ID[7:19],
			Name:    im.Name,
			Tag:     im.Tag,
			Created: im.Created,
			Size:    im.Size,
		})
	}
	return reports
}

func (a *Task) NetworkTask() []*NetworkReport {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	nets, err := a.manager.ListNetwork(ctx)
	if err != nil {
		slog.Error("list networks failed", "error", err)
		return nil
	}
	var reports []*NetworkReport
	for _, net := range nets {
		subnet := ""
		gateway := ""
		labels, _ := json.Marshal(net.Labels)
		if len(net.SubNet) > 0 {
			subnet = net.SubNet[0].Subnet
			gateway = net.SubNet[0].Gateway
		}
		reports = append(reports, &NetworkReport{
			NetworkID: net.ID[:6],
			Name:      net.Name,
			Driver:    net.Driver,
			Created:   net.Created,
			Scope:     net.Scope,
			Internal:  net.Internal,
			Subnet:    subnet,
			Gateway:   gateway,
			Labels:    string(labels),
		})
	}
	return reports
}
