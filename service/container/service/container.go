// Package service
// Date: 2024/3/6 12:49
// Author: Amu
// Description:
package service

import (
	"context"
	"time"
	
	"github.com/amuluze/amprobe/service/container/repository"
	"github.com/amuluze/amprobe/service/schema"
	"github.com/amuluze/amutool/errors"
	"github.com/google/wire"
)

var ContainerServiceSet = wire.NewSet(NewContainerService, wire.Bind(new(IContainerService), new(*ContainerService)))

type IContainerService interface {
	ContainerList(ctx context.Context) (*schema.ContainerQueryRely, error)
	ImageList(ctx context.Context) (*schema.ImageQueryReply, error)
	Version(ctx context.Context) (*schema.Docker, error)
}

type ContainerService struct {
	ContainerRepo repository.IContainerRepo
}

func NewContainerService(containerRepo repository.IContainerRepo) *ContainerService {
	return &ContainerService{ContainerRepo: containerRepo}
}

func (a *ContainerService) ContainerList(ctx context.Context) (*schema.ContainerQueryRely, error) {
	mContainers, err := a.ContainerRepo.ContainerList(ctx)
	if err != nil {
		return nil, errors.New400Error(err.Error())
	}
	var list []schema.Container
	var ts time.Time
	for _, item := range mContainers {
		if len(list) == 0 {
			ts = item.Timestamp
			list = append(list, schema.Container{
				ID:            item.ContainerID[:6],
				Name:          item.Name,
				Image:         item.Image,
				State:         item.State,
				Uptime:        item.Uptime,
				IP:            item.IP,
				CPUPercent:    item.CPUPercent,
				MemoryPercent: item.MemPercent,
				MemoryLimit:   item.MemLimit,
				MemoryUsage:   item.MemUsage,
			})
		} else {
			if ts != item.Timestamp {
				continue
			}
			list = append(list, schema.Container{
				ID:            item.ContainerID[:6],
				Name:          item.Name,
				Image:         item.Image,
				State:         item.State,
				Uptime:        item.Uptime,
				IP:            item.IP,
				CPUPercent:    item.CPUPercent,
				MemoryPercent: item.MemPercent,
				MemoryLimit:   item.MemLimit,
				MemoryUsage:   item.MemUsage,
			})
		}
	}
	return &schema.ContainerQueryRely{Containers: list}, nil
}

func (a *ContainerService) ImageList(ctx context.Context) (*schema.ImageQueryReply, error) {
	images, err := a.ContainerRepo.ImageList(ctx)
	if err != nil {
		return nil, errors.New400Error(err.Error())
	}
	var list []schema.Image
	for _, item := range images {
		list = append(list, schema.Image{
			ID:      item.ID,
			Name:    item.Name,
			Tag:     item.Tag,
			Created: item.Created,
			Size:    item.Size,
		})
	}
	return &schema.ImageQueryReply{Images: list}, nil
}

func (a *ContainerService) Version(ctx context.Context) (*schema.Docker, error) {
	version, err := a.ContainerRepo.Version(ctx)
	if err != nil {
		return nil, errors.New400Error(err.Error())
	}
	return &schema.Docker{
		Timestamp:     version.Timestamp,
		DockerVersion: version.DockerVersion,
		APIVersion:    version.APIVersion,
		MinAPIVersion: version.MinAPIVersion,
		GitCommit:     version.GitCommit,
		GoVersion:     version.GoVersion,
		Os:            version.Os,
		Arch:          version.Arch,
	}, nil
}
