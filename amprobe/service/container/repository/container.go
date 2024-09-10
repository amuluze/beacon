// Package repository
// Date: 2024/06/11 19:38:25
// Author: Amu
// Description:
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"amprobe/pkg/rpc"
	"amprobe/service/model"
	rpcSchema "common/rpc/schema"

	"github.com/google/wire"
)

var ContainerServiceSet = wire.NewSet(NewContainerRepo, wire.Bind(new(IContainerRepo), new(*ContainerRepo)))

var _ IContainerRepo = (*ContainerRepo)(nil)

type IContainerRepo interface {
	Version(ctx context.Context, args rpcSchema.DockerArgs) (rpcSchema.DockerReply, error)

	ContainerList(ctx context.Context, args rpcSchema.ContainerQueryArgs) (rpcSchema.ContainerQueryReply, error)
	ContainerCreate(ctx context.Context, args rpcSchema.ContainerCreateArgs) (rpcSchema.ContainerCreateReply, error)
	ContainerStart(ctx context.Context, args rpcSchema.ContainerStartArgs) error
	ContainerStop(ctx context.Context, args rpcSchema.ContainerStopArgs) error
	ContainerRemove(ctx context.Context, args rpcSchema.ContainerDeleteArgs) error
	ContainerRestart(ctx context.Context, args rpcSchema.ContainerRestartArgs) error

	ImageList(ctx context.Context, args rpcSchema.ImageQueryArgs) (rpcSchema.ImageQueryReply, error)
	ImagePull(ctx context.Context, args rpcSchema.ImagePullArgs) error
	ImageImport(ctx context.Context, args rpcSchema.ImageImportArgs) error
	ImageExport(ctx context.Context, args rpcSchema.ImageExportArgs) error
	ImageRemove(ctx context.Context, args rpcSchema.ImageDeleteArgs) error
	ImagesPrune(ctx context.Context) error

	NetworkList(ctx context.Context, args rpcSchema.NetworkQueryArgs) (rpcSchema.NetworkQueryReply, error)
	NetworkCreate(ctx context.Context, args rpcSchema.NetworkCreateArgs) (rpcSchema.NetworkCreateReply, error)
	NetworkDelete(ctx context.Context, args rpcSchema.NetworkDeleteArgs) error

	GetDockerRegistryMirrors(ctx context.Context, args rpcSchema.GetDockerRegistryMirrorsArgs) (rpcSchema.GetDockerRegistryMirrorsReply, error)
	SetDockerRegistryMirrors(ctx context.Context, args rpcSchema.SetDockerRegistryMirrorsArgs) error
}

type ContainerRepo struct {
	RPCClient *rpc.Client
}

func NewContainerRepo(client *rpc.Client) *ContainerRepo {
	return &ContainerRepo{RPCClient: client}
}

func (c *ContainerRepo) Version(ctx context.Context, args rpcSchema.DockerArgs) (rpcSchema.DockerReply, error) {
	var reply rpcSchema.DockerReply
	err := c.RPCClient.Call(ctx, "Version", args, &reply)
	return reply, err
}

func (c *ContainerRepo) ContainerList(ctx context.Context, args rpcSchema.ContainerQueryArgs) (rpcSchema.ContainerQueryReply, error) {
	return rpcSchema.ContainerQueryReply{}, nil
}

func (c *ContainerRepo) ContainerCreate(ctx context.Context, args rpcSchema.ContainerCreateArgs) (rpcSchema.ContainerCreateReply, error) {
	var reply rpcSchema.ContainerCreateReply
	err := c.RPCClient.Call(ctx, "ContainerCreate", args, &reply)
	if err != nil {
		return rpcSchema.ContainerCreateReply{}, err
	}
	return reply, nil
}

func (c *ContainerRepo) ContainerStart(ctx context.Context, args rpcSchema.ContainerStartArgs) error {
	var reply rpcSchema.ContainerStartReply
	err := c.RPCClient.Call(ctx, "ContainerStart", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContainerRepo) ContainerStop(ctx context.Context, args rpcSchema.ContainerStopArgs) error {
	var reply rpcSchema.ContainerStopReply
	err := c.RPCClient.Call(ctx, "ContainerStop", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContainerRepo) ContainerRemove(ctx context.Context, args rpcSchema.ContainerDeleteArgs) error {
	var reply rpcSchema.ContainerRemoveReply
	err := c.RPCClient.Call(ctx, "ContainerDelete", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContainerRepo) ContainerRestart(ctx context.Context, args rpcSchema.ContainerRestartArgs) error {
	var reply rpcSchema.ContainerRestartReply
	err := c.RPCClient.Call(ctx, "ContainerRestart", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContainerRepo) ImageList(ctx context.Context, args rpcSchema.ImageQueryArgs) (rpcSchema.ImageQueryReply, error) {
	var reply model.Images
	err := c.RPCClient.Call(ctx, "ImageList", args, &reply)
	if err != nil {
		return rpcSchema.ImageQueryReply{}, err
	}
	var result rpcSchema.ImageQueryReply
	var data []rpcSchema.Image
	for _, v := range reply {
		data = append(data, rpcSchema.Image{
			ID:      v.ImageID[:6],
			Name:    v.Name,
			Tag:     v.Tag,
			Created: v.Created,
			Size:    v.Size,
			Number:  v.Number,
		})
	}
	countArgs := rpcSchema.ImageCountArgs{}
	countReply := rpcSchema.ImageCountReply{}
	err = c.RPCClient.Call(ctx, "ImageCount", countArgs, &countReply)
	if err != nil {
		return rpcSchema.ImageQueryReply{}, err
	}
	result.Data = data
	result.Page = args.Page
	result.Size = args.Size
	result.Total = countReply.Count
	return result, nil
}

func (c *ContainerRepo) ImagePull(ctx context.Context, args rpcSchema.ImagePullArgs) error {
	var reply rpcSchema.ImagePullReply
	err := c.RPCClient.Call(ctx, "ImagePull", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContainerRepo) ImageImport(ctx context.Context, args rpcSchema.ImageImportArgs) error {
	var reply rpcSchema.ImageImportReply
	err := c.RPCClient.Call(ctx, "ImageImport", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContainerRepo) ImageExport(ctx context.Context, args rpcSchema.ImageExportArgs) error {
	var reply rpcSchema.ImageExportReply
	rpcArgs := rpcSchema.ImageExportRPCArgs{
		ImageIDs:   []string{args.ImageID},
		TargetFile: fmt.Sprintf("/tmp/%s.tar", args.ImageName),
	}
	err := c.RPCClient.Call(ctx, "ImageExport", rpcArgs, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContainerRepo) ImageRemove(ctx context.Context, args rpcSchema.ImageDeleteArgs) error {
	var reply rpcSchema.ImageRemoveReply
	err := c.RPCClient.Call(ctx, "ImageDelete", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContainerRepo) ImagesPrune(ctx context.Context) error {
	err := c.RPCClient.Call(ctx, "ImagesPrune", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContainerRepo) NetworkList(ctx context.Context, args rpcSchema.NetworkQueryArgs) (rpcSchema.NetworkListReply, error) {
	var reply model.Networks
	err := c.RPCClient.Call(ctx, "NetworkList", args, &reply)
	if err != nil {
		slog.Error("NetworkList", "err", err)
		return rpcSchema.NetworkListReply{}, err
	}
	var result rpcSchema.NetworkListReply
	var data []rpcSchema.Network
	for _, v := range reply {
		labels := make(map[string]string)
		err := json.Unmarshal([]byte(v.Labels), &labels)
		if err != nil {
			slog.Error("json.Unmarshal", "err", err)
			return rpcSchema.NetworkListReply{}, err
		}
		data = append(data, rpcSchema.Network{
			ID:      v.NetworkID[:6],
			Name:    v.Name,
			Driver:  v.Driver,
			Created: v.Created,
			Subnet:  v.Subnet,
			Gateway: v.Gateway,
			Labels:  labels,
		})
	}
	countArgs := rpcSchema.NetworkCountArgs{}
	countReply := rpcSchema.NetworkCountReply{}
	err = c.RPCClient.Call(ctx, "NetworkCount", countArgs, &countReply)
	if err != nil {
		return rpcSchema.NetworkListReply{}, err
	}
	result.Data = data
	result.Page = args.Page
	result.Size = args.Size
	result.Total = countReply.Count
	return result, nil
}

func (c *ContainerRepo) NetworkCreate(ctx context.Context, args rpcSchema.NetworkCreateArgs) (rpcSchema.NetworkCreateReply, error) {
	var reply rpcSchema.NetworkCreateReply
	err := c.RPCClient.Call(ctx, "NetworkCreate", args, &reply)
	if err != nil {
		return rpcSchema.NetworkCreateReply{}, err
	}
	return reply, nil
}

func (c *ContainerRepo) NetworkDelete(ctx context.Context, args rpcSchema.NetworkDeleteArgs) error {
	var reply rpcSchema.NetworkDeleteReply
	err := c.RPCClient.Call(ctx, "NetworkDelete", args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContainerRepo) GetDockerRegistryMirrors(ctx context.Context, args rpcSchema.GetDockerRegistryMirrorsArgs) (rpcSchema.GetDockerRegistryMirrorsReply, error) {
	var reply rpcSchema.GetDockerRegistryMirrorsReply
	// 调用rpc获取数据
	err := c.RPCClient.Call(ctx, "GetDockerRegistryMirrors", args, &reply)
	if err != nil {
		return rpcSchema.GetDockerRegistryMirrorsReply{}, err
	}
	return reply, nil
}

func (c *ContainerRepo) SetDockerRegistryMirrors(ctx context.Context, args rpcSchema.SetDockerRegistryMirrorsArgs) error {
	var reply rpcSchema.SetDockerRegistryMirrorsReply
	// 调用rpc设置数据
	err := c.RPCClient.Call(ctx, "SetDockerRegistryMirrors", args, &reply)
	if err != nil {
		return err
	}
	return nil
}
