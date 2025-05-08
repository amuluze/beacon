// Package task
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package task

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"amprobe/service/model"
)

func (a *Task) DockerTask(timestamp time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dockerVersion, err := a.manager.Version(ctx)
	if err != nil {
		return err
	}
	if err := a.db.Unscoped().Where("1 = 1").Delete(&model.Docker{}).Error; err != nil {
		return err
	}
	if err := a.db.Model(&model.Docker{}).Create(&model.Docker{
		Timestamp:     timestamp,
		DockerVersion: dockerVersion.DockerVersion,
		APIVersion:    dockerVersion.APIVersion,
		MinAPIVersion: dockerVersion.MinAPIVersion,
		GitCommit:     dockerVersion.GitCommit,
		GoVersion:     dockerVersion.GoVersion,
		Os:            dockerVersion.OS,
		Arch:          dockerVersion.Arch,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (a *Task) ContainerTask(timestamp time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cs, err := a.manager.ListContainer(ctx)
	if err != nil {
		slog.Error("failed to list containers", "error", err)
		return err
	}
	var containers []model.Container
	for _, info := range cs {
		cpuPercent, err := a.manager.GetContainerCpu(ctx, info.ID[:6])
		if err != nil {
			return err
		}
		memPercent, used, limit, err := a.manager.GetContainerMem(ctx, info.ID[:6])
		if err != nil {
			return err
		}
		labels, _ := json.Marshal(info.Labels)
		containers = append(containers, model.Container{
			Timestamp:   timestamp,
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
	if err := a.db.Model(&model.Container{}).Create(&containers).Error; err != nil {
		return err
	}
	ts := time.Now().Add(-time.Duration(a.maxAge) * 24 * time.Hour).Unix()
	if err := a.db.Unscoped().Where("timestamp < ?", ts).Delete(&model.Container{}).Error; err != nil {
		return err
	}
	return nil
}

func (a *Task) ImageTask(timestamp time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	images, err := a.manager.ListImage(ctx)
	if err != nil {
		return err
	}
	var list []model.Image
	for _, im := range images {
		list = append(list, model.Image{
			Timestamp: timestamp,
			ImageID:   im.ID[7:19],
			Name:      im.Name,
			Tag:       im.Tag,
			Created:   im.Created,
			Size:      im.Size,
		})
	}
	if err := a.db.Unscoped().Where("1 = 1").Delete(&model.Image{}).Error; err != nil {
		slog.Error("failed to delete image", "error", err)
	}
	if err := a.db.Model(&model.Image{}).Create(&list).Error; err != nil {
		return err
	}
	return nil
}

func (a *Task) NetworkTask(timestamp time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	nets, err := a.manager.ListNetwork(ctx)
	if err != nil {
		return err
	}
	var list []model.Network
	for _, net := range nets {
		subnet := ""
		gateway := ""
		labels, _ := json.Marshal(net.Labels)
		if len(net.SubNet) > 0 {
			subnet = net.SubNet[0].Subnet
			gateway = net.SubNet[0].Gateway
		}
		list = append(list, model.Network{
			Timestamp: timestamp,
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
	if err := a.db.Unscoped().Where("1 = 1").Delete(&model.Network{}).Error; err != nil {
		return err
	}
	if err := a.db.Model(&model.Network{}).Create(&list).Error; err != nil {
		return err
	}
	return nil
}
