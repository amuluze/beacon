// Package service
// Date: 2024/3/6 12:49
// Author: Amu
// Description:
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"amprobe/pkg/utils"
	"amprobe/service/container/repository"
	"amprobe/service/schema"

	"github.com/google/wire"
)

var ContainerServiceSet = wire.NewSet(NewContainerService, wire.Bind(new(IContainerService), new(*ContainerService)))

var _ IContainerService = (*ContainerService)(nil)

type IContainerService interface {
	Version(ctx context.Context) (schema.Docker, error)

	ContainerList(ctx context.Context, args schema.ContainerQueryArgs) (schema.ContainerQueryRely, error)
	Usage(ctx context.Context, args schema.ContainerUsageArgs) (schema.ContainerUsageReply, error)
	ContainerCreate(ctx context.Context, args schema.ContainerCreateArgs) (schema.ContainerCreateReply, error)
	ContainerUpdate(ctx context.Context, args schema.ContainerUpdateArgs) (schema.ContainerUpdateReply, error)
	ContainerStart(ctx context.Context, args schema.ContainerStartArgs) error
	ContainerStop(ctx context.Context, args schema.ContainerStopArgs) error
	ContainerDelete(ctx context.Context, args schema.ContainerDeleteArgs) error
	ContainerRestart(ctx context.Context, args schema.ContainerRestartArgs) error

	ImageList(ctx context.Context, args schema.ImageQueryArgs) (schema.ImageQueryReply, error)
	ImagePull(ctx context.Context, args schema.ImagePullArgs) error
	ImageImport(ctx context.Context, args schema.ImageImportArgs) error
	ImageExport(ctx context.Context, args schema.ImageExportArgs) error
	ImageDelete(ctx context.Context, args schema.ImageDeleteArgs) error
	ImagesPrune(ctx context.Context) error

	NetworkList(ctx context.Context, args schema.NetworkQueryArgs) (schema.NetworkQueryReply, error)
	NetworkCreate(ctx context.Context, args schema.NetworkCreateArgs) (schema.NetworkCreateReply, error)
	NetworkDelete(ctx context.Context, args schema.NetworkDeleteArgs) error

	GetDockerRegistryMirrors(ctx context.Context) (schema.GetDockerRegistryMirrorsReply, error)
	SetDockerRegistryMirrors(ctx context.Context, args schema.SetDockerRegistryMirrorsArgs) error
}

type ContainerService struct {
	ContainerRepo repository.IContainerRepo
}

func NewContainerService(containerRepo repository.IContainerRepo) *ContainerService {
	return &ContainerService{ContainerRepo: containerRepo}
}

func (a *ContainerService) Version(ctx context.Context) (schema.Docker, error) {
	version, err := a.ContainerRepo.Version(ctx)
	if err != nil {
		return schema.Docker{}, err
	}
	return schema.Docker{
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

func (a *ContainerService) ContainerList(ctx context.Context, args schema.ContainerQueryArgs) (schema.ContainerQueryRely, error) {
	var reply schema.ContainerQueryRely
	results, err := a.ContainerRepo.ContainerList(ctx)
	if err != nil {
		return reply, err
	}

	var list []schema.Container
	for _, item := range results {
		list = append(list, schema.Container{
			ID:            item.ContainerID,
			Name:          item.Name,
			Image:         item.Image,
			State:         item.State,
			Uptime:        item.Uptime,
			IP:            item.IP,
			Ports:         item.Ports,
			CPUPercent:    fmt.Sprintf("%.2f", item.CPUPercent) + " %",
			MemoryPercent: fmt.Sprintf("%.2f", item.MemPercent) + " %",
			MemoryLimit:   utils.ConvertBytesToReadable(item.MemLimit),
			MemoryUsage:   utils.ConvertBytesToReadable(item.MemUsage),
		})
	}
	count, err := a.ContainerRepo.ContainerCount(ctx)
	if err != nil {
		return reply, err
	}
	reply.Data = list
	reply.Total = int(count)
	reply.Page = args.Page
	reply.Size = args.Size
	return reply, nil
}

func (a *ContainerService) Usage(ctx context.Context, args schema.ContainerUsageArgs) (schema.ContainerUsageReply, error) {
	var reply schema.ContainerUsageReply
	return reply, nil
}

func (a *ContainerService) ContainerCreate(ctx context.Context, args schema.ContainerCreateArgs) (schema.ContainerCreateReply, error) {
	var reply schema.ContainerCreateReply
	containerID, err := a.ContainerRepo.ContainerCreate(ctx, args)
	if err != nil {
		return reply, err
	}
	reply.ContainerID = containerID
	return reply, nil
}

func (a *ContainerService) ContainerUpdate(ctx context.Context, args schema.ContainerUpdateArgs) (schema.ContainerUpdateReply, error) {
	var reply schema.ContainerUpdateReply
	containerID, err := a.ContainerRepo.ContainerUpdate(ctx, args)
	if err != nil {
		return reply, err
	}
	reply.ContainerID = containerID
	return reply, nil
}

func (a *ContainerService) ContainerDelete(ctx context.Context, args schema.ContainerDeleteArgs) error {
	return a.ContainerRepo.ContainerDelete(ctx, args)
}

func (a *ContainerService) ContainerStart(ctx context.Context, args schema.ContainerStartArgs) error {
	return a.ContainerRepo.ContainerStart(ctx, args)
}

func (a *ContainerService) ContainerStop(ctx context.Context, args schema.ContainerStopArgs) error {
	return a.ContainerRepo.ContainerStop(ctx, args)
}

func (a *ContainerService) ContainerRestart(ctx context.Context, args schema.ContainerRestartArgs) error {
	return a.ContainerRepo.ContainerRestart(ctx, args)
}

func (a *ContainerService) ImageList(ctx context.Context, args schema.ImageQueryArgs) (schema.ImageQueryReply, error) {
	var reply schema.ImageQueryReply
	results, err := a.ContainerRepo.ImageList(ctx)
	if err != nil {
		return reply, err
	}
	var list []schema.Image
	for _, item := range results {
		num, err := a.ContainerRepo.ContainersByImage(ctx, fmt.Sprintf("%s:%s", item.Name, item.Tag))
		if err != nil {
			return reply, err
		}
		list = append(list, schema.Image{
			ID:      item.ImageID,
			Name:    item.Name,
			Tag:     item.Tag,
			Created: item.Created,
			Size:    item.Size,
			Number:  int(num),
		})
	}
	count, err := a.ContainerRepo.ImageCount(ctx)
	if err != nil {
		return reply, err
	}
	reply.Data = list
	reply.Total = int(count)
	reply.Page = args.Page
	reply.Size = args.Size
	return reply, nil
}

func (a *ContainerService) ImagePull(ctx context.Context, args schema.ImagePullArgs) error {
	return a.ContainerRepo.ImagePull(ctx, args)
}

func (a *ContainerService) ImageTag(ctx context.Context, args schema.ImageTagArgs) error {
	return a.ContainerRepo.ImageTag(ctx, args)
}

func (a *ContainerService) ImageDelete(ctx context.Context, args schema.ImageDeleteArgs) error {
	return a.ContainerRepo.ImageDelete(ctx, args)
}

func (a *ContainerService) ImageImport(ctx context.Context, args schema.ImageImportArgs) error {
	return a.ContainerRepo.ImageImport(ctx, args)
}

func (a *ContainerService) ImageExport(ctx context.Context, args schema.ImageExportArgs) error {
	return a.ContainerRepo.ImageExport(ctx, args)
}

func (a *ContainerService) ImagesPrune(ctx context.Context) error {
	return a.ContainerRepo.ImagesPrune(ctx)
}

func (a *ContainerService) NetworkList(ctx context.Context, args schema.NetworkQueryArgs) (schema.NetworkQueryReply, error) {
	reply := schema.NetworkQueryReply{}
	results, err := a.ContainerRepo.NetworkList(ctx)
	if err != nil {
		return reply, err
	}
	var list []schema.Network
	for _, item := range results {
		labels := make(map[string]string)
		err := json.Unmarshal([]byte(item.Labels), &labels)
		if err != nil {
			return reply, err
		}
		list = append(list, schema.Network{
			ID:      item.NetworkID,
			Name:    item.Name,
			Driver:  item.Driver,
			Created: item.Created,
			Subnet:  item.Subnet,
			Gateway: item.Gateway,
			Labels:  labels,
		})
	}
	count, err := a.ContainerRepo.NetworkCount(ctx)
	if err != nil {
		return reply, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Created > list[j].Created })
	reply.Data = list
	reply.Total = int(count)
	reply.Page = args.Page
	reply.Size = args.Size
	return reply, nil
}

func (a *ContainerService) NetworkCreate(ctx context.Context, args schema.NetworkCreateArgs) (schema.NetworkCreateReply, error) {
	var reply schema.NetworkCreateReply
	networkID, err := a.ContainerRepo.NetworkCreate(ctx, args)
	if err != nil {
		return reply, err
	}
	reply.NetworkID = networkID
	return reply, nil
}

func (a *ContainerService) NetworkDelete(ctx context.Context, args schema.NetworkDeleteArgs) error {
	return a.ContainerRepo.NetworkDelete(ctx, args)
}

func (a *ContainerService) GetDockerRegistryMirrors(ctx context.Context) (schema.GetDockerRegistryMirrorsReply, error) {
	var reply schema.GetDockerRegistryMirrorsReply
	mirrors, err := a.ContainerRepo.GetDockerRegistryMirrors(ctx)
	if err != nil {
		return reply, err
	}
	reply.Mirrors = mirrors
	return reply, nil
}

func (a *ContainerService) SetDockerRegistryMirrors(ctx context.Context, args schema.SetDockerRegistryMirrorsArgs) error {
	rpcArgs := schema.SetDockerRegistryMirrorsArgs{
		Mirrors: args.Mirrors,
	}
	return a.ContainerRepo.SetDockerRegistryMirrors(ctx, rpcArgs)
}
