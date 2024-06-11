// Package rpc
// Date: 2024/06/11 19:38:25
// Author: Amu
// Description:
package rpc

import (
	"context"

	"github.com/amuluze/amprobe/service"
	"github.com/amuluze/amprobe/service/schema"
	"github.com/google/wire"
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
	RPCClient *service.RPCClient
}

func NewContainerService(client *service.RPCClient) *ContainerService {
	return &ContainerService{RPCClient: client}
}
