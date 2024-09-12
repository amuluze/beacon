// Package repository
// Date: 2024/06/11 19:38:25
// Author: Amu
// Description:
package repository

import (
	"context"

	"amprobe/pkg/rpc"
	rpcSchema "common/rpc/schema"

	"github.com/google/wire"
)

var ContainerServiceSet = wire.NewSet(NewContainerRepo, wire.Bind(new(IContainerRepo), new(*ContainerRepo)))

var _ IContainerRepo = (*ContainerRepo)(nil)

type IContainerRepo interface {
	Version(ctx context.Context, args rpcSchema.DockerArgs) (rpcSchema.DockerReply, error)

	ContainerList(ctx context.Context, args rpcSchema.ContainerQueryArgs) (rpcSchema.ContainerQueryReply, error)
	ContainersByImage(ctx context.Context, image string) (num int, err error)
	ContainerCount(ctx context.Context, args rpcSchema.ContainerCountArgs) (rpcSchema.ContainerCountReply, error)
	ContainerCreate(ctx context.Context, args rpcSchema.ContainerCreateArgs) (rpcSchema.ContainerCreateReply, error)
	ContainerUpdate(ctx context.Context, args rpcSchema.ContainerUpdateArgs) (rpcSchema.ContainerUpdateReply, error)
	ContainerDelete(ctx context.Context, args rpcSchema.ContainerDeleteArgs) error
	ContainerStart(ctx context.Context, args rpcSchema.ContainerStartArgs) error
	ContainerStop(ctx context.Context, args rpcSchema.ContainerStopArgs) error
	ContainerRestart(ctx context.Context, args rpcSchema.ContainerRestartArgs) error

	ImageList(ctx context.Context, args rpcSchema.ImageQueryArgs) (rpcSchema.ImageQueryReply, error)
	ImageCount(ctx context.Context, args rpcSchema.ImageCountArgs) (rpcSchema.ImageCountReply, error)
	ImagePull(ctx context.Context, args rpcSchema.ImagePullArgs) error
	ImageTag(ctx context.Context, args rpcSchema.ImageTagArgs) error
	ImageImport(ctx context.Context, args rpcSchema.ImageImportArgs) error
	ImageExport(ctx context.Context, args rpcSchema.ImageExportArgs) error
	ImageDelete(ctx context.Context, args rpcSchema.ImageDeleteArgs) error
	ImagesPrune(ctx context.Context) error

	NetworkList(ctx context.Context, args rpcSchema.NetworkQueryArgs) (rpcSchema.NetworkQueryReply, error)
	NetworkCount(ctx context.Context, args rpcSchema.NetworkCountArgs) (rpcSchema.NetworkCountReply, error)
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
	var reply rpcSchema.ContainerQueryReply
	err := c.RPCClient.Call(ctx, "ContainerList", args, &reply)
	return reply, err
}

func (c *ContainerRepo) ContainersByImage(ctx context.Context, image string) (num int, err error) {
	args := rpcSchema.ContainersByImageArgs{Image: image}
	var reply rpcSchema.ContainersByImageReply
	err = c.RPCClient.Call(ctx, "ContainersByImage", args, &reply)
	return reply.Num, err

}

func (c *ContainerRepo) ContainerCount(ctx context.Context, args rpcSchema.ContainerCountArgs) (rpcSchema.ContainerCountReply, error) {
	var reply rpcSchema.ContainerCountReply
	err := c.RPCClient.Call(ctx, "ContainerCount", args, &reply)
	return reply, err
}

func (c *ContainerRepo) ContainerCreate(ctx context.Context, args rpcSchema.ContainerCreateArgs) (rpcSchema.ContainerCreateReply, error) {
	var reply rpcSchema.ContainerCreateReply
	err := c.RPCClient.Call(ctx, "ContainerCreate", args, &reply)
	if err != nil {
		return rpcSchema.ContainerCreateReply{}, err
	}
	return reply, nil
}

func (c *ContainerRepo) ContainerUpdate(ctx context.Context, args rpcSchema.ContainerUpdateArgs) (rpcSchema.ContainerUpdateReply, error) {
	return rpcSchema.ContainerUpdateReply{}, nil
}

func (c *ContainerRepo) ContainerDelete(ctx context.Context, args rpcSchema.ContainerDeleteArgs) error {
	var reply rpcSchema.ContainerDeleteReply
	return c.RPCClient.Call(ctx, "ContainerDelete", args, &reply)
}

func (c *ContainerRepo) ContainerStart(ctx context.Context, args rpcSchema.ContainerStartArgs) error {
	var reply rpcSchema.ContainerStartReply
	return c.RPCClient.Call(ctx, "ContainerStart", args, &reply)
}

func (c *ContainerRepo) ContainerStop(ctx context.Context, args rpcSchema.ContainerStopArgs) error {
	var reply rpcSchema.ContainerStopReply
	return c.RPCClient.Call(ctx, "ContainerStop", args, &reply)
}

func (c *ContainerRepo) ContainerRestart(ctx context.Context, args rpcSchema.ContainerRestartArgs) error {
	var reply rpcSchema.ContainerRestartReply
	return c.RPCClient.Call(ctx, "ContainerRestart", args, &reply)
}

func (c *ContainerRepo) ImageList(ctx context.Context, args rpcSchema.ImageQueryArgs) (rpcSchema.ImageQueryReply, error) {
	var reply rpcSchema.ImageQueryReply
	err := c.RPCClient.Call(ctx, "ImageList", args, &reply)
	return reply, err
}

func (c *ContainerRepo) ImageCount(ctx context.Context, args rpcSchema.ImageCountArgs) (rpcSchema.ImageCountReply, error) {
	var reply rpcSchema.ImageCountReply
	err := c.RPCClient.Call(ctx, "ImageCount", args, &reply)
	return reply, err
}

func (c *ContainerRepo) ImagePull(ctx context.Context, args rpcSchema.ImagePullArgs) error {
	var reply rpcSchema.ImagePullReply
	return c.RPCClient.Call(ctx, "ImagePull", args, &reply)
}

func (c *ContainerRepo) ImageTag(ctx context.Context, args rpcSchema.ImageTagArgs) error {
	var reply rpcSchema.ImageTagReply
	return c.RPCClient.Call(ctx, "ImageTag", args, &reply)
}

func (c *ContainerRepo) ImageImport(ctx context.Context, args rpcSchema.ImageImportArgs) error {
	var reply rpcSchema.ImageImportReply
	return c.RPCClient.Call(ctx, "ImageImport", args, &reply)
}

func (c *ContainerRepo) ImageExport(ctx context.Context, args rpcSchema.ImageExportArgs) error {
	var reply rpcSchema.ImageExportReply
	return c.RPCClient.Call(ctx, "ImageExport", args, &reply)
}

func (c *ContainerRepo) ImageDelete(ctx context.Context, args rpcSchema.ImageDeleteArgs) error {
	var reply rpcSchema.ImageDeleteReply
	return c.RPCClient.Call(ctx, "ImageDelete", args, &reply)
}

func (c *ContainerRepo) ImagesPrune(ctx context.Context) error {
	return c.RPCClient.Call(ctx, "ImagesPrune", nil, nil)
}

func (c *ContainerRepo) NetworkList(ctx context.Context, args rpcSchema.NetworkQueryArgs) (rpcSchema.NetworkQueryReply, error) {
	var reply rpcSchema.NetworkQueryReply
	err := c.RPCClient.Call(ctx, "NetworkList", args, &reply)
	return reply, err
}

func (c *ContainerRepo) NetworkCount(ctx context.Context, args rpcSchema.NetworkCountArgs) (rpcSchema.NetworkCountReply, error) {
	var reply rpcSchema.NetworkCountReply
	err := c.RPCClient.Call(ctx, "NetworkCount", args, &reply)
	return reply, err
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
	return c.RPCClient.Call(ctx, "NetworkDelete", args, &reply)
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
	return c.RPCClient.Call(ctx, "SetDockerRegistryMirrors", args, &reply)
}
