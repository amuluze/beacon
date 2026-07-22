// Package docker
// Date: 2024/07/09 14:13:06
// Author: Amu
// Description:
package docker

import (
	"context"
	"github.com/docker/docker/api/types/registry"
	"io"

	"github.com/docker/docker/client"
)

type Manager struct {
	client *client.Client
}

func NewManager() (*Manager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return &Manager{client: cli}, err
}

var _ IManager = (*Manager)(nil)

type IManager interface {
	Version(context.Context) (*Version, error)

	ListContainer(ctx context.Context) ([]ContainerSummary, error)
	HasSameNameContainer(ctx context.Context, containerName string) (bool, error)
	CreateContainer(ctx context.Context, containerName, imageName, networkName string, ports []string, vols []string, env []string, commands []string, labels map[string]string) (string, error)
	StartContainer(ctx context.Context, containerID string) error
	StopContainer(ctx context.Context, containerID string) error
	RestartContainer(ctx context.Context, containerID string) error
	DeleteContainer(ctx context.Context, containerID string) error
	CopyFileToContainer(ctx context.Context, containerID string, srcFile, dstFile string) error
	GetContainerMem(ctx context.Context, containerID string) (float64, float64, float64, error)
	GetContainerCpu(ctx context.Context, containerID string) (float64, error)
	GetContainerIDByContainerName(ctx context.Context, containerName string) (string, error)
	ContainerLogs(ctx context.Context, containerID string) (io.ReadCloser, error)
	RenameContainer(ctx context.Context, containerID, newName string) error

	ListImage(ctx context.Context) ([]ImageSummary, error)
	DeleteImage(ctx context.Context, imageID string) error
	PruneImages(ctx context.Context) error
	SearchImage(ctx context.Context, imageName string) ([]registry.SearchResult, error)
	PullImage(ctx context.Context, imageName string) error
	TagImage(ctx context.Context, oldTag, newTag string) error
	ImportImage(ctx context.Context, sourceFile string) error
	ExportImage(ctx context.Context, imageIDs []string, targetFile string) error
	GetImageByName(ctx context.Context, imageName string) (*ImageSummary, error)
	GetImageByID(ctx context.Context, imageID string) (*ImageSummary, error)

	ListNetwork(ctx context.Context) ([]NetworkSummary, error)
	HasSameNameNetwork(ctx context.Context, networkName string) (bool, error)
	CreateNetwork(ctx context.Context, name, driver, subnet, gateway string, labels map[string]string) (string, error)
	GetNetworkByID(ctx context.Context, networkID string) (*NetworkSummary, error)
	DeleteNetwork(ctx context.Context, networkID string) error
	PruneNetwork(ctx context.Context) error
	JoinNetwork(ctx context.Context, containerID, networkID string) error
	LeaveNetwork(ctx context.Context, containerID, networkID string) error
}
