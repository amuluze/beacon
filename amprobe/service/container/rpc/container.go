// Package rpc
// Date: 2024/06/11 19:38:25
// Author: Amu
// Description:
package rpc

import (
	"context"
	"fmt"
	"github.com/amuluze/amprobe/pkg/rpc"
	"github.com/amuluze/amprobe/pkg/utils"
	"github.com/amuluze/amprobe/service/model"
	"github.com/amuluze/amprobe/service/schema"
	"github.com/google/wire"
)

var ContainerServiceSet = wire.NewSet(NewContainerService, wire.Bind(new(IContainerService), new(*ContainerService)))

type IContainerService interface {
	ContainerList(ctx context.Context, args schema.ContainerQueryArgs) (schema.ContainerQueryRely, error)
	ContainerStart(ctx context.Context, args schema.ContainerStartArgs) error
	ContainerStop(ctx context.Context, args schema.ContainerStopArgs) error
	ContainerRemove(ctx context.Context, args schema.ContainerRemoveArgs) error
	ContainerRestart(ctx context.Context, args schema.ContainerRestartArgs) error
	ImageList(ctx context.Context, args schema.ImageQueryArgs) (schema.ImageQueryReply, error)
	ImageRemove(ctx context.Context, args schema.ImageRemoveArgs) error
	ImagesPrune(ctx context.Context) error
	Version(ctx context.Context) (schema.Docker, error)
	GetDockerImageSettings(ctx context.Context, args schema.GetDockerImageSettingsArgs) (schema.GetDockerImageSettingsReply, error)
	SetDockerImageSettings(ctx context.Context, args schema.SetDockerImageSettingsArgs) error
}

type ContainerService struct {
	RPCClient *rpc.Client
}

func NewContainerService(client *rpc.Client) *ContainerService {
	return &ContainerService{RPCClient: client}
}

func (c ContainerService) Version(ctx context.Context) (schema.Docker, error) {
	args := schema.VersionArgs{}
	var reply model.Docker
	err := c.RPCClient.Call(ctx, "DockerVersion", args, &reply)
	if err != nil {
		return schema.Docker{}, err
	}
	return schema.Docker{
		Timestamp:     reply.Timestamp,
		DockerVersion: reply.DockerVersion,
		APIVersion:    reply.APIVersion,
		MinAPIVersion: reply.MinAPIVersion,
		GitCommit:     reply.GitCommit,
		Os:            reply.Os,
		Arch:          reply.Arch,
	}, nil
}

func (c ContainerService) ContainerList(ctx context.Context, args schema.ContainerQueryArgs) (schema.ContainerQueryRely, error) {
	var reply model.Containers
	err := c.RPCClient.Call(ctx, "ContainerList", args, &reply)
	if err != nil {
		return schema.ContainerQueryRely{}, err
	}
	var result schema.ContainerQueryRely
	var data []schema.Container
	for _, v := range reply {
		data = append(data, schema.Container{
			ID:            v.ContainerID[:6],
			Name:          v.Name,
			Image:         v.Image,
			IP:            v.IP,
			State:         v.State,
			Uptime:        v.Uptime,
			CPUPercent:    fmt.Sprintf("%.2f", v.CPUPercent) + " %",
			MemoryPercent: fmt.Sprintf("%.2f", v.MemPercent) + " %",
			MemoryLimit:   utils.ConvertBytesToReadable(v.MemLimit),
			MemoryUsage:   utils.ConvertBytesToReadable(v.MemUsage),
		})
	}
	countArgs := schema.QueryCountArgs{}
	var countReply schema.QueryCountReply
	err = c.RPCClient.Call(ctx, "ContainerCount", countArgs, &countReply)
	if err != nil {
		return schema.ContainerQueryRely{}, err
	}
	result.Data = data
	result.Page = args.Page
	result.Size = args.Size
	result.Total = countReply.Count
	return result, nil
}

func (c ContainerService) ContainerStart(ctx context.Context, args schema.ContainerStartArgs) error {
	var reply schema.ContainerStartReply
	err := c.RPCClient.Call(ctx, "ContainerStart", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c ContainerService) ContainerStop(ctx context.Context, args schema.ContainerStopArgs) error {
	var reply schema.ContainerStopReply
	err := c.RPCClient.Call(ctx, "ContainerStop", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c ContainerService) ContainerRemove(ctx context.Context, args schema.ContainerRemoveArgs) error {
	var reply schema.ContainerRemoveReply
	err := c.RPCClient.Call(ctx, "ContainerRemove", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c ContainerService) ContainerRestart(ctx context.Context, args schema.ContainerRestartArgs) error {
	var reply schema.ContainerRestartReply
	err := c.RPCClient.Call(ctx, "ContainerRestart", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c ContainerService) ImageList(ctx context.Context, args schema.ImageQueryArgs) (schema.ImageQueryReply, error) {
	var reply model.Images
	err := c.RPCClient.Call(ctx, "ImageList", args, &reply)
	if err != nil {
		return schema.ImageQueryReply{}, err
	}
	var result schema.ImageQueryReply
	var data []schema.Image
	for _, v := range reply {
		data = append(data, schema.Image{
			ID:      v.ImageID[:6],
			Name:    v.Name,
			Tag:     v.Tag,
			Created: v.Created,
			Size:    v.Size,
			Number:  v.Number,
		})
	}
	countArgs := schema.ImageCountArgs{}
	countReply := schema.ImageCountReply{}
	err = c.RPCClient.Call(ctx, "ImageCount", countArgs, &countReply)
	if err != nil {
		return schema.ImageQueryReply{}, err
	}
	result.Data = data
	result.Page = args.Page
	result.Size = args.Size
	result.Total = countReply.Count
	return result, nil
}

func (c ContainerService) ImageRemove(ctx context.Context, args schema.ImageRemoveArgs) error {
	var reply schema.ImageRemoveReply
	err := c.RPCClient.Call(ctx, "ImageDelete", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c ContainerService) ImagesPrune(ctx context.Context) error {
	err := c.RPCClient.Call(ctx, "ImagesPrune", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c ContainerService) GetDockerImageSettings(ctx context.Context, args schema.GetDockerImageSettingsArgs) (schema.GetDockerImageSettingsReply, error) {
	return schema.GetDockerImageSettingsReply{}, nil
}

func (c ContainerService) SetDockerImageSettings(ctx context.Context, args schema.SetDockerImageSettingsArgs) error {
	return nil
}
