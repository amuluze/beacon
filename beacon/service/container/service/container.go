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

	"beacon/pkg/utils"
	"beacon/service/container/repository"
	"beacon/service/schema"
	rpcSchema "common/rpc/schema"

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
	ImageDelete(ctx context.Context, args schema.ImageDeleteArgs) error
	ImagesPrune(ctx context.Context) error

	NetworkList(ctx context.Context, args schema.NetworkQueryArgs) (schema.NetworkQueryReply, error)
	NetworkCreate(ctx context.Context, args schema.NetworkCreateArgs) (schema.NetworkCreateReply, error)
	NetworkDelete(ctx context.Context, args schema.NetworkDeleteArgs) error

	GetDockerRegistryMirrors(ctx context.Context, args schema.GetDockerRegistryMirrorsArgs) (schema.GetDockerRegistryMirrorsReply, error)
	SetDockerRegistryMirrors(ctx context.Context, args schema.SetDockerRegistryMirrorsArgs) error
}

type ContainerService struct {
	ContainerRepo repository.IContainerRepo
}

func NewContainerService(containerRepo repository.IContainerRepo) *ContainerService {
	return &ContainerService{ContainerRepo: containerRepo}
}

func mapFreshness(f rpcSchema.Freshness) schema.Freshness {
	return schema.Freshness{
		CollectedAt: f.CollectedAt,
		AgeSeconds:  f.AgeSeconds,
		Stale:       f.Stale,
		Degraded:    f.Degraded,
	}
}

func (a *ContainerService) Version(ctx context.Context) (schema.Docker, error) {
	var args rpcSchema.DockerArgs
	version, err := a.ContainerRepo.Version(ctx, args)
	if err != nil {
		return schema.Docker{}, err
	}
	return schema.Docker{
		Timestamp:     version.Data.Timestamp,
		DockerVersion: version.Data.DockerVersion,
		APIVersion:    version.Data.APIVersion,
		MinAPIVersion: version.Data.MinAPIVersion,
		GitCommit:     version.Data.GitCommit,
		GoVersion:     version.Data.GoVersion,
		Os:            version.Data.Os,
		Arch:          version.Data.Arch,
		Freshness:     mapFreshness(version.Freshness),
	}, nil
}

func (a *ContainerService) ContainerList(ctx context.Context, args schema.ContainerQueryArgs) (schema.ContainerQueryRely, error) {
	var reply schema.ContainerQueryRely
	rpcArgs := rpcSchema.ContainerQueryArgs{
		Page: args.Page,
		Size: args.Size,
	}
	rpcReply, err := a.ContainerRepo.ContainerList(ctx, rpcArgs)
	if err != nil {
		return reply, err
	}

	var list []schema.Container
	for _, item := range rpcReply.Data {
		labels := make(map[string]string)
		if item.Labels != "" {
			if err := json.Unmarshal([]byte(item.Labels), &labels); err != nil {
				return reply, fmt.Errorf("decode container labels: %w", err)
			}
		}
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
			Labels:        labels,
		})
	}
	countReply, err := a.ContainerRepo.ContainerCount(ctx, rpcSchema.ContainerCountArgs{})
	if err != nil {
		return reply, err
	}
	reply.Data = list
	reply.Total = countReply.Count
	reply.Page = args.Page
	reply.Size = args.Size
	reply.Freshness = mapFreshness(rpcReply.Freshness)
	return reply, nil
}

func (a *ContainerService) Usage(ctx context.Context, args schema.ContainerUsageArgs) (schema.ContainerUsageReply, error) {
	var reply schema.ContainerUsageReply
	reply.Names = make([]string, 0)
	reply.CPUUsage = make(map[string][]schema.Usage)
	reply.MemUsage = make(map[string][]schema.Usage)
	rpcReply, err := a.ContainerRepo.Usage(ctx, rpcSchema.ContainerUsageArgs{
		StartTime: args.StartTime,
		EndTime:   args.EndTime,
	})
	if err != nil {
		return reply, err
	}
	reply.Names = rpcReply.Names
	for key, value := range rpcReply.CPUUsage {
		if _, ok := reply.CPUUsage[key]; !ok {
			reply.CPUUsage[key] = make([]schema.Usage, 0)
		}
		for _, item := range value {
			reply.CPUUsage[key] = append(reply.CPUUsage[key], schema.Usage{
				Timestamp: item.Timestamp,
				Value:     item.Value,
			})
		}
	}
	for key, value := range rpcReply.MemUsage {
		if _, ok := reply.MemUsage[key]; !ok {
			reply.MemUsage[key] = make([]schema.Usage, 0)
		}
		for _, item := range value {
			reply.MemUsage[key] = append(reply.MemUsage[key], schema.Usage{
				Timestamp: item.Timestamp,
				Value:     item.Value,
			})
		}
	}
	reply.Freshness = mapFreshness(rpcReply.Freshness)
	return reply, nil
}

func (a *ContainerService) ContainerCreate(ctx context.Context, args schema.ContainerCreateArgs) (schema.ContainerCreateReply, error) {
	rpcArgs := rpcSchema.ContainerCreateArgs{
		Name:          args.ContainerName,
		Image:         args.ImageName,
		Network:       args.NetworkName,
		RestartPolicy: args.RestartPolicy,
		Ports:         args.Ports,
		Volumes:       args.Volumes,
		Environments:  args.Environments,
		Commands:      nil,
		Labels:        args.Labels,
	}
	var reply schema.ContainerCreateReply
	rpcReply, err := a.ContainerRepo.ContainerCreate(ctx, rpcArgs)
	if err != nil {
		return reply, err
	}
	reply.ContainerID = rpcReply.ContainerID
	return reply, nil
}

func (a *ContainerService) ContainerUpdate(ctx context.Context, args schema.ContainerUpdateArgs) (schema.ContainerUpdateReply, error) {
	rpcArgs := rpcSchema.ContainerUpdateArgs{
		ContainerID:   args.ContainerID,
		Name:          args.ContainerName,
		Image:         args.ImageName,
		Network:       args.NetworkName,
		RestartPolicy: args.RestartPolicy,
		Ports:         args.Ports,
		Volumes:       args.Volumes,
		Environments:  args.Environments,
		Commands:      nil,
		Labels:        args.Labels,
	}
	var reply schema.ContainerUpdateReply
	rpcReply, err := a.ContainerRepo.ContainerUpdate(ctx, rpcArgs)
	if err != nil {
		return reply, err
	}
	reply.ContainerID = rpcReply.ContainerID
	return reply, nil
}

func (a *ContainerService) ContainerDelete(ctx context.Context, args schema.ContainerDeleteArgs) error {
	rpcArgs := rpcSchema.ContainerDeleteArgs{
		ContainerID: args.ContainerID,
	}
	return a.ContainerRepo.ContainerDelete(ctx, rpcArgs)
}

func (a *ContainerService) ContainerStart(ctx context.Context, args schema.ContainerStartArgs) error {
	rpcArgs := rpcSchema.ContainerStartArgs{
		ContainerID: args.ContainerID,
	}
	return a.ContainerRepo.ContainerStart(ctx, rpcArgs)
}

func (a *ContainerService) ContainerStop(ctx context.Context, args schema.ContainerStopArgs) error {
	rpcArgs := rpcSchema.ContainerStopArgs{
		ContainerID: args.ContainerID,
	}
	return a.ContainerRepo.ContainerStop(ctx, rpcArgs)
}

func (a *ContainerService) ContainerRestart(ctx context.Context, args schema.ContainerRestartArgs) error {
	rpcArgs := rpcSchema.ContainerRestartArgs{
		ContainerID: args.ContainerID,
	}
	return a.ContainerRepo.ContainerRestart(ctx, rpcArgs)
}

func (a *ContainerService) ImageList(ctx context.Context, args schema.ImageQueryArgs) (schema.ImageQueryReply, error) {
	var reply schema.ImageQueryReply
	rpcArgs := rpcSchema.ImageQueryArgs{
		Page: args.Page,
		Size: args.Size,
	}
	listReply, err := a.ContainerRepo.ImageList(ctx, rpcArgs)
	if err != nil {
		return reply, err
	}
	var list []schema.Image
	for _, item := range listReply.Data {
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
			Number:  num,
		})
	}
	countReply, err := a.ContainerRepo.ImageCount(ctx, rpcSchema.ImageCountArgs{})
	if err != nil {
		return reply, err
	}
	reply.Data = list
	reply.Total = countReply.Count
	reply.Page = args.Page
	reply.Size = args.Size
	reply.Freshness = mapFreshness(listReply.Freshness)
	return reply, nil
}

func (a *ContainerService) ImagePull(ctx context.Context, args schema.ImagePullArgs) error {
	rpcArgs := rpcSchema.ImagePullArgs{
		ImageName: args.ImageName,
	}
	return a.ContainerRepo.ImagePull(ctx, rpcArgs)
}

func (a *ContainerService) ImageTag(ctx context.Context, args schema.ImageTagArgs) error {
	rpcArgs := rpcSchema.ImageTagArgs{
		NewTag: args.NewTag,
		OldTag: args.OldTag,
	}
	return a.ContainerRepo.ImageTag(ctx, rpcArgs)
}

func (a *ContainerService) ImageDelete(ctx context.Context, args schema.ImageDeleteArgs) error {
	rpcArgs := rpcSchema.ImageDeleteArgs{
		ImageID: args.ImageID,
	}
	return a.ContainerRepo.ImageDelete(ctx, rpcArgs)
}

func (a *ContainerService) ImageImport(ctx context.Context, args schema.ImageImportArgs) error {
	rpcArgs := rpcSchema.ImageImportArgs{
		SourceFile: args.SourceFile,
		FileName:   args.FileName,
		Data:       args.Data,
	}
	return a.ContainerRepo.ImageImport(ctx, rpcArgs)
}

func (a *ContainerService) ImagesPrune(ctx context.Context) error {
	return a.ContainerRepo.ImagesPrune(ctx)
}

func (a *ContainerService) NetworkList(ctx context.Context, args schema.NetworkQueryArgs) (schema.NetworkQueryReply, error) {
	reply := schema.NetworkQueryReply{}
	rpcArgs := rpcSchema.NetworkQueryArgs{
		Page: args.Page,
		Size: args.Size,
	}
	rpcReply, err := a.ContainerRepo.NetworkList(ctx, rpcArgs)
	if err != nil {
		return reply, err
	}
	var list []schema.Network
	for _, item := range rpcReply.Data {
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
	countReply, err := a.ContainerRepo.NetworkCount(ctx, rpcSchema.NetworkCountArgs{})
	if err != nil {
		return reply, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Created > list[j].Created })
	reply.Data = list
	reply.Total = countReply.Count
	reply.Page = args.Page
	reply.Size = args.Size
	reply.Freshness = mapFreshness(rpcReply.Freshness)
	return reply, nil
}

func (a *ContainerService) NetworkCreate(ctx context.Context, args schema.NetworkCreateArgs) (schema.NetworkCreateReply, error) {
	var reply schema.NetworkCreateReply
	rpcArgs := rpcSchema.NetworkCreateArgs{
		Name:    args.Name,
		Driver:  args.Driver,
		Labels:  args.Labels,
		Subnet:  args.Subnet,
		Gateway: args.Gateway,
	}
	rpcReply, err := a.ContainerRepo.NetworkCreate(ctx, rpcArgs)
	if err != nil {
		return reply, err
	}
	reply.NetworkID = rpcReply.NetworkID
	return reply, nil
}

func (a *ContainerService) NetworkDelete(ctx context.Context, args schema.NetworkDeleteArgs) error {
	rpcArgs := rpcSchema.NetworkDeleteArgs{
		NetworkID: args.NetworkID,
	}
	return a.ContainerRepo.NetworkDelete(ctx, rpcArgs)
}

func (a *ContainerService) GetDockerRegistryMirrors(ctx context.Context, args schema.GetDockerRegistryMirrorsArgs) (schema.GetDockerRegistryMirrorsReply, error) {
	result := schema.GetDockerRegistryMirrorsReply{}
	rpcArgs := rpcSchema.GetDockerRegistryMirrorsArgs{}
	mirrors, err := a.ContainerRepo.GetDockerRegistryMirrors(ctx, rpcArgs)
	if err != nil {
		return result, err
	}
	result.Mirrors = mirrors.Mirrors
	return result, nil
}

func (a *ContainerService) SetDockerRegistryMirrors(ctx context.Context, args schema.SetDockerRegistryMirrorsArgs) error {
	rpcArgs := rpcSchema.SetDockerRegistryMirrorsArgs{
		Mirrors: args.Mirrors,
	}
	return a.ContainerRepo.SetDockerRegistryMirrors(ctx, rpcArgs)
}
