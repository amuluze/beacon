// Package rpc
// Date: 2024/06/11 19:38:25
// Author: Amu
// Description:
package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/amuluze/amprobe/pkg/rpc"
	"github.com/amuluze/amprobe/pkg/utils"
	"github.com/amuluze/amprobe/service/model"
	"github.com/amuluze/amprobe/service/schema"
	"github.com/amuluze/docker"
	"github.com/google/wire"
)

var ContainerServiceSet = wire.NewSet(NewContainerService, wire.Bind(new(IContainerService), new(*ContainerService)))

type IContainerService interface {
	Version(ctx context.Context) (schema.Docker, error)
	ContainerList(ctx context.Context, args schema.ContainerQueryArgs) (schema.ContainerQueryRely, error)
	ContainerCreate(ctx context.Context, args schema.ContainerCreateArgs) (schema.ContainerCreateReply, error)
	ContainerStart(ctx context.Context, args schema.ContainerStartArgs) error
	ContainerStop(ctx context.Context, args schema.ContainerStopArgs) error
	ContainerRemove(ctx context.Context, args schema.ContainerRemoveArgs) error
	ContainerRestart(ctx context.Context, args schema.ContainerRestartArgs) error
	ImageList(ctx context.Context, args schema.ImageQueryArgs) (schema.ImageQueryReply, error)
	ImagePull(ctx context.Context, args schema.ImagePullArgs) error
	ImageImport(ctx context.Context, args schema.ImageImportArgs) error
	ImageExport(ctx context.Context, args schema.ImageExportArgs) error
	ImageRemove(ctx context.Context, args schema.ImageRemoveArgs) error
	ImagesPrune(ctx context.Context) error
	NetworkList(ctx context.Context, args schema.NetworkListArgs) (schema.NetworkListReply, error)
	NetworkCreate(ctx context.Context, args schema.NetworkCreateArgs) (schema.NetworkCreateReply, error)
	NetworkDelete(ctx context.Context, args schema.NetworkDeleteArgs) error
	GetDockerRegistryMirrors(ctx context.Context, args schema.GetDockerRegistryMirrorsArgs) (schema.GetDockerRegistryMirrorsReply, error)
	SetDockerRegistryMirrors(ctx context.Context, args schema.SetDockerRegistryMirrorsArgs) error
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
		var labels map[string]string
		if err := json.Unmarshal([]byte(v.Labels), &labels); err != nil {
			return schema.ContainerQueryRely{}, err
		}
		serverType := ""
		if _, ok := labels[docker.ServerTypeLabel]; ok {
			serverType = labels[docker.ServerTypeLabel]
		}
		data = append(data, schema.Container{
			ID:            v.ContainerID[:6],
			Name:          v.Name,
			Image:         v.Image,
			IP:            v.IP,
			Ports:         v.Ports,
			ServerType:    serverType,
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

func (c ContainerService) ContainerCreate(ctx context.Context, args schema.ContainerCreateArgs) (schema.ContainerCreateReply, error) {
	var reply schema.ContainerCreateReply
	err := c.RPCClient.Call(ctx, "ContainerCreate", args, &reply)
	if err != nil {
		return schema.ContainerCreateReply{}, err
	}
	return reply, nil
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
	err := c.RPCClient.Call(ctx, "ContainerDelete", args, &reply)
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

func (c ContainerService) ImagePull(ctx context.Context, args schema.ImagePullArgs) error {
	var reply schema.ImagePullReply
	err := c.RPCClient.Call(ctx, "ImagePull", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c ContainerService) ImageImport(ctx context.Context, args schema.ImageImportArgs) error {
	var reply schema.ImageImportReply
	err := c.RPCClient.Call(ctx, "ImageImport", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c ContainerService) ImageExport(ctx context.Context, args schema.ImageExportArgs) error {
	var reply schema.ImageExportReply
	rpcArgs := schema.ImageExportRPCArgs{
		ImageIDs:   []string{args.ImageID},
		TargetFile: fmt.Sprintf("/tmp/%s.tar", args.ImageName),
	}
	err := c.RPCClient.Call(ctx, "ImageExport", rpcArgs, &reply)
	if err != nil {
		return err
	}
	return nil
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

func (c ContainerService) NetworkList(ctx context.Context, args schema.NetworkListArgs) (schema.NetworkListReply, error) {
	var reply model.Networks
	err := c.RPCClient.Call(ctx, "NetworkList", args, &reply)
	if err != nil {
		slog.Error("NetworkList", "err", err)
		return schema.NetworkListReply{}, err
	}
	var result schema.NetworkListReply
	var data []schema.Network
	for _, v := range reply {
		labels := make(map[string]string)
		err := json.Unmarshal([]byte(v.Labels), &labels)
		if err != nil {
			slog.Error("json.Unmarshal", "err", err)
			return schema.NetworkListReply{}, err
		}
		data = append(data, schema.Network{
			ID:      v.NetworkID[:6],
			Name:    v.Name,
			Driver:  v.Driver,
			Created: v.Created,
			Subnet:  v.Subnet,
			Gateway: v.Gateway,
			Labels:  labels,
		})
	}
	countArgs := schema.NetworkCountArgs{}
	countReply := schema.NetworkCountReply{}
	err = c.RPCClient.Call(ctx, "NetworkCount", countArgs, &countReply)
	if err != nil {
		return schema.NetworkListReply{}, err
	}
	result.Data = data
	result.Page = args.Page
	result.Size = args.Size
	result.Total = countReply.Count
	return result, nil
}

func (c ContainerService) NetworkCreate(ctx context.Context, args schema.NetworkCreateArgs) (schema.NetworkCreateReply, error) {
	var reply schema.NetworkCreateReply
	err := c.RPCClient.Call(ctx, "NetworkCreate", args, &reply)
	if err != nil {
		return schema.NetworkCreateReply{}, err
	}
	return reply, nil
}

func (c ContainerService) NetworkDelete(ctx context.Context, args schema.NetworkDeleteArgs) error {
	var reply schema.NetworkDeleteReply
	err := c.RPCClient.Call(ctx, "NetworkDelete", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c ContainerService) GetDockerRegistryMirrors(ctx context.Context, args schema.GetDockerRegistryMirrorsArgs) (schema.GetDockerRegistryMirrorsReply, error) {
	var reply schema.GetDockerRegistryMirrorsReply
	// 调用rpc获取数据
	err := c.RPCClient.Call(ctx, "GetDockerRegistryMirrors", args, &reply)
	if err != nil {
		return schema.GetDockerRegistryMirrorsReply{}, err
	}
	return reply, nil
}

func (c ContainerService) SetDockerRegistryMirrors(ctx context.Context, args schema.SetDockerRegistryMirrorsArgs) error {
	var reply schema.SetDockerRegistryMirrorsReply
	// 调用rpc设置数据
	err := c.RPCClient.Call(ctx, "SetDockerRegistryMirrors", args, &reply)
	if err != nil {
		return err
	}
	return nil
}
