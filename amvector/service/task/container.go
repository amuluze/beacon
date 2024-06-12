// Package task
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package task

import (
	"context"
	"github.com/amuluze/amvector/service/model"
	"log/slog"
	"time"
)

func (a *Task) Container(timestamp time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cs, err := a.manager.ListContainer(ctx)
	if err != nil {
		slog.Error("failed to list containers", "error", err)
		return
	}
	var containers []model.Container
	for _, info := range cs {
		var d model.Container
		d.Timestamp = timestamp
		d.ContainerID = info.ID[:6]
		d.Name = info.Name
		d.State = info.State
		d.Image = info.Image
		d.Uptime = info.Uptime
		d.IP = info.IP

		cpuPercent, err := a.manager.GetContainerCPU(ctx, info.ID[:6])
		if err != nil {
			slog.Error("failed to get container cpu", "error", err)
		}
		d.CPUPercent = cpuPercent

		memPercent, used, limit, err := a.manager.GetContainerMem(ctx, info.ID[:6])
		if err != nil {
			slog.Error("failed to get container mem", "error", err)
		}
		d.MemPercent = memPercent

		d.MemUsage = used
		d.MemLimit = limit
		if _, ok := a.cache.Get(info.Image); !ok {
			a.cache.Set(info.Image, 1, 2*time.Minute)
		} else {
			count, err := a.cache.IncrementInt(info.Image, 1)
			slog.Info("container image cache", "image", info.Image, "count", count, "error", err)
		}
		containers = append(containers, d)
	}
	if err := a.db.Unscoped().Where("1 = 1").Delete(&model.Container{}).Error; err != nil {
		slog.Error("failed to delete container", "error", err)
	}
	a.db.Model(&model.Container{}).Create(&containers)
}

func (a *Task) Docker(timestamp time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	dockerVersion, err := a.manager.Version(ctx)
	if err != nil {
		slog.Error("failed to get docker version", "error", err)
		return
	}
	if err := a.db.Unscoped().Where("1 = 1").Delete(&model.Docker{}).Error; err != nil {
		slog.Error("failed to delete docker container", "error", err)
	}
	a.db.Model(&model.Docker{}).Create(&model.Docker{
		Timestamp:     timestamp,
		DockerVersion: dockerVersion.DockerVersion,
		APIVersion:    dockerVersion.APIVersion,
		MinAPIVersion: dockerVersion.MinAPIVersion,
		GitCommit:     dockerVersion.GitCommit,
		GoVersion:     dockerVersion.GoVersion,
		Os:            dockerVersion.OS,
		Arch:          dockerVersion.Arch,
	})
}

func (a *Task) Image(timestamp time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	images, err := a.manager.ListImage(ctx)
	if err != nil {
		slog.Error("failed to get version", "error", err)
		return
	}
	var list model.Images
	duplicateImage := make(map[string]struct{})
	for _, im := range images {
		val, ok := a.cache.Get(im.Name + ":" + im.Tag)
		if !ok {
			slog.Error("failed to get image cache", "error", err)
			val = 0
		}
		if _, ok := duplicateImage[im.ID]; !ok {
			duplicateImage[im.ID] = struct{}{}
		} else {
			if im.Tag != "latest" {
				continue
			}
		}
		list = append(list, model.Image{
			Timestamp: timestamp,
			ImageID:   im.ID[7:19],
			Name:      im.Name,
			Number:    val.(int),
			Tag:       im.Tag,
			Created:   im.Created,
			Size:      im.Size,
		})
		a.cache.Delete(im.Name + ":" + im.Tag)
	}
	if err := a.db.Unscoped().Where("1 = 1").Delete(&model.Image{}).Error; err != nil {
		slog.Error("failed to delete image", "error", err)
	}
	a.db.Model(&model.Image{}).Create(&list)
}

func (a *Task) ClearOldRecord() {
	a.db.Where("timestamp < ?", time.Now().Add(-time.Minute*5)).Delete(&model.Host{})
	a.db.Where("timestamp < ?", time.Now().Add(-time.Minute*5)).Delete(&model.Container{})
	a.db.Where("timestamp < ?", time.Now().Add(-time.Minute*5)).Delete(&model.Image{})
	a.db.Where("timestamp < ?", time.Now().Add(-time.Minute*5)).Delete(&model.Docker{})
	a.db.Where("timestamp < ?", time.Now().Add(-time.Hour*24*2)).Delete(&model.CPU{})
	a.db.Where("timestamp < ?", time.Now().Add(-time.Hour*24*2)).Delete(&model.Memory{})
	a.db.Where("timestamp < ?", time.Now().Add(-time.Hour*24*2)).Delete(&model.Disk{})
	a.db.Where("timestamp < ?", time.Now().Add(-time.Hour*24*2)).Delete(&model.Net{})
}
