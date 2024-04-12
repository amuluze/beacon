// Package service
// Date: 2024/3/6 12:49
// Author: Amu
// Description:
package service

import (
	"context"
	"fmt"
	"github.com/amuluze/amprobe/pkg/utils"
	"github.com/amuluze/amprobe/service/container/repository"
	"github.com/amuluze/amprobe/service/schema"
	"github.com/amuluze/amutool/errors"
	"github.com/google/wire"
	"log/slog"
)

var ContainerServiceSet = wire.NewSet(NewContainerService, wire.Bind(new(IContainerService), new(*ContainerService)))

type IContainerService interface {
	ContainerList(ctx context.Context, args *schema.ContainerQueryArgs) (*schema.ContainerQueryRely, error)
	ContainerStart(ctx context.Context, args *schema.ContainerStartArgs) error
	ContainerStop(ctx context.Context, args *schema.ContainerStopArgs) error
	ContainerRemove(ctx context.Context, args *schema.ContainerRemoveArgs) error
	ContainerRestart(ctx context.Context, args *schema.ContainerRestartArgs) error
	ImageList(ctx context.Context, args *schema.ImageQueryArgs) (*schema.ImageQueryReply, error)
	ImageRemove(ctx context.Context, args *schema.ImageRemoveArgs) error
	ImagesPrune(ctx context.Context) error
	Version(ctx context.Context) (*schema.Docker, error)
}

type ContainerService struct {
	ContainerRepo repository.IContainerRepo
}

func NewContainerService(containerRepo repository.IContainerRepo) *ContainerService {
	return &ContainerService{ContainerRepo: containerRepo}
}

func (a *ContainerService) ContainerList(ctx context.Context, args *schema.ContainerQueryArgs) (*schema.ContainerQueryRely, error) {
	mContainers, err := a.ContainerRepo.ContainerList(ctx, args)
	if err != nil {
		return nil, errors.New400Error(err.Error())
	}
	slog.Info("container list", "list", mContainers)
	var list []schema.Container
	for _, item := range mContainers {
		list = append(list, schema.Container{
			ID:            item.ContainerID[:6],
			Name:          item.Name,
			Image:         item.Image,
			State:         item.State,
			Uptime:        item.Uptime,
			IP:            item.IP,
			CPUPercent:    fmt.Sprintf("%.2f", item.CPUPercent) + " %",
			MemoryPercent: fmt.Sprintf("%.2f", item.MemPercent) + " %",
			MemoryLimit:   utils.ConvertBytesToReadable(item.MemLimit),
			MemoryUsage:   utils.ConvertBytesToReadable(item.MemUsage),
		})
	}
	total, err := a.ContainerRepo.ContainerCount(ctx)
	if err != nil {
		slog.Error("query container count error", "err", err)
		return nil, errors.New400Error(err.Error())
	}
	return &schema.ContainerQueryRely{Data: list, Total: total, Page: args.Page, Size: args.Size}, nil
}

func (a *ContainerService) ImageList(ctx context.Context, args *schema.ImageQueryArgs) (*schema.ImageQueryReply, error) {
	images, err := a.ContainerRepo.ImageList(ctx, args)
	if err != nil {
		return &schema.ImageQueryReply{}, errors.New400Error(err.Error())
	}
	var list []schema.Image
	for _, item := range images {
		list = append(list, schema.Image{
			ID:      item.ImageID,
			Name:    item.Name,
			Tag:     item.Tag,
			Created: item.Created,
			Size:    item.Size,
			Number:  item.Number,
		})
	}
	total, _ := a.ContainerRepo.ImageCount(ctx)
	return &schema.ImageQueryReply{Data: list, Total: total, Page: args.Page, Size: args.Size}, nil
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

func (a *ContainerService) ContainerStart(ctx context.Context, args *schema.ContainerStartArgs) error {
	return a.ContainerRepo.ContainerStart(ctx, args)
}

func (a *ContainerService) ContainerStop(ctx context.Context, args *schema.ContainerStopArgs) error {
	return a.ContainerRepo.ContainerStop(ctx, args)
}

func (a *ContainerService) ContainerRemove(ctx context.Context, args *schema.ContainerRemoveArgs) error {
	return a.ContainerRepo.ContainerRemove(ctx, args)
}

func (a *ContainerService) ContainerRestart(ctx context.Context, args *schema.ContainerRestartArgs) error {
	return a.ContainerRepo.ContainerRestart(ctx, args)
}

func (a *ContainerService) ImageRemove(ctx context.Context, args *schema.ImageRemoveArgs) error {
	return a.ContainerRepo.ImageRemove(ctx, args)
}

func (a *ContainerService) ImagesPrune(ctx context.Context) error {
	return a.ContainerRepo.ImagesPrune(ctx)
}
