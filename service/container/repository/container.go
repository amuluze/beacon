// Package repository
// Date: 2024/06/11 19:38:25
// Author: Amu
// Description:
package repository

import (
	"amprobe/pkg/command"
	"amprobe/service/constants"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"amprobe/pkg/database"
	"amprobe/service/model"
	"amprobe/service/schema"

	"github.com/amuluze/docker"
	"github.com/google/wire"
	"github.com/patrickmn/go-cache"
)

var ContainerServiceSet = wire.NewSet(NewContainerRepo, wire.Bind(new(IContainerRepo), new(*ContainerRepo)))

var _ IContainerRepo = (*ContainerRepo)(nil)

type IContainerRepo interface {
	Version(ctx context.Context) (model.Docker, error)

	ContainerList(ctx context.Context) ([]model.Container, error)
	Usage(ctx context.Context, args schema.ContainerUsageArgs) (schema.ContainerUsageReply, error)
	ContainersByImage(ctx context.Context, image string) (int64, error)
	ContainerCount(ctx context.Context) (int64, error)
	ContainerCreate(ctx context.Context, args schema.ContainerCreateArgs) (string, error)
	ContainerUpdate(ctx context.Context, args schema.ContainerUpdateArgs) (string, error)
	ContainerDelete(ctx context.Context, args schema.ContainerDeleteArgs) error
	ContainerStart(ctx context.Context, args schema.ContainerStartArgs) error
	ContainerStop(ctx context.Context, args schema.ContainerStopArgs) error
	ContainerRestart(ctx context.Context, args schema.ContainerRestartArgs) error

	ImageList(ctx context.Context) ([]model.Image, error)
	ImageCount(ctx context.Context) (int64, error)
	ImagePull(ctx context.Context, args schema.ImagePullArgs) error
	ImageTag(ctx context.Context, args schema.ImageTagArgs) error
	ImageImport(ctx context.Context, args schema.ImageImportArgs) error
	ImageExport(ctx context.Context, args schema.ImageExportArgs) error
	ImageDelete(ctx context.Context, args schema.ImageDeleteArgs) error
	ImagesPrune(ctx context.Context) error

	NetworkList(ctx context.Context) ([]model.Network, error)
	NetworkCount(ctx context.Context) (int64, error)
	NetworkCreate(ctx context.Context, args schema.NetworkCreateArgs) (string, error)
	NetworkDelete(ctx context.Context, args schema.NetworkDeleteArgs) error

	GetDockerRegistryMirrors(ctx context.Context) ([]string, error)
	SetDockerRegistryMirrors(ctx context.Context, args schema.SetDockerRegistryMirrorsArgs) error
}

type ContainerRepo struct {
	manager *docker.Manager
	cache   *cache.Cache
}

func NewContainerRepo(db *database.DB, localCache cache.Cache) *ContainerRepo {
	manager, err := docker.NewManager()
	if err != nil {
		return nil
	}
	return &ContainerRepo{manager: manager}
}

func (c *ContainerRepo) Version(ctx context.Context) (model.Docker, error) {
	var reply model.Docker
	version, err := c.manager.Version(ctx)
	if err != nil {
		return reply, err
	}
	reply.Timestamp = time.Now()
	reply.DockerVersion = version.DockerVersion
	reply.APIVersion = version.APIVersion
	reply.MinAPIVersion = version.MinAPIVersion
	reply.GitCommit = version.GitCommit
	reply.Os = version.OS
	reply.Arch = version.Arch
	reply.GoVersion = version.GoVersion
	return reply, nil
}

func (c *ContainerRepo) ContainerList(ctx context.Context) ([]model.Container, error) {
	var reply []model.Container
	cs, err := c.manager.ListContainer(ctx)
	if err != nil {
		return reply, err
	}
	for _, container := range cs {
		cpuPercent, err := c.manager.GetContainerCpu(ctx, container.ID[:6])
		if err != nil {
			cpuPercent = 0
		}
		memPercent, memUsed, memLimit, err := c.manager.GetContainerMem(ctx, container.ID[:6])
		if err != nil {
			memPercent = 0
		}
		labels, _ := json.Marshal(container.Labels)
		reply = append(reply, model.Container{
			Timestamp:   time.Now(),
			ContainerID: container.ID[:6],
			Name:        container.Name,
			State:       container.State,
			Image:       container.Image,
			Uptime:      container.Uptime,
			IP:          container.IP,
			Ports:       strings.Join(container.Ports, ","),
			Labels:      string(labels),
			CPUPercent:  cpuPercent,
			MemPercent:  memPercent,
			MemUsage:    memUsed,
			MemLimit:    memLimit,
		})
	}
	return reply, nil
}

func (c *ContainerRepo) Usage(ctx context.Context, args schema.ContainerUsageArgs) (schema.ContainerUsageReply, error) {
	var reply schema.ContainerUsageReply
	return reply, nil
}

func (c *ContainerRepo) ContainersByImage(ctx context.Context, image string) (int64, error) {
	var count int64 = 0
	cs, err := c.manager.ListContainer(ctx)
	if err != nil {
		return count, err
	}
	for _, container := range cs {
		if container.Image == image {
			count++
		}
	}
	return count, nil
}

func (c *ContainerRepo) ContainerCount(ctx context.Context) (int64, error) {
	cs, err := c.manager.ListContainer(ctx)
	if err != nil {
		return 0, err
	}
	return int64(len(cs)), nil
}

func (c *ContainerRepo) ContainerCreate(ctx context.Context, args schema.ContainerCreateArgs) (string, error) {
	containerID, err := c.manager.CreateContainer(
		ctx,
		args.ContainerName,
		args.ImageName,
		args.NetworkName,
		args.Ports,
		args.Volumes,
		args.Environments,
		nil,
		args.Labels,
	)
	if err != nil {
		return containerID, err
	}
	return containerID, nil
}

func (c *ContainerRepo) ContainerUpdate(ctx context.Context, args schema.ContainerUpdateArgs) (string, error) {
	return "", nil
}

func (c *ContainerRepo) ContainerDelete(ctx context.Context, args schema.ContainerDeleteArgs) error {
	if err := c.manager.DeleteContainer(ctx, args.ContainerID); err != nil {
		return err
	}
	return nil
}

func (c *ContainerRepo) ContainerStart(ctx context.Context, args schema.ContainerStartArgs) error {
	if err := c.manager.StartContainer(ctx, args.ContainerID); err != nil {
		return err
	}
	return nil
}

func (c *ContainerRepo) ContainerStop(ctx context.Context, args schema.ContainerStopArgs) error {
	if err := c.manager.StopContainer(ctx, args.ContainerID); err != nil {
		return err
	}
	return nil
}

func (c *ContainerRepo) ContainerRestart(ctx context.Context, args schema.ContainerRestartArgs) error {
	if err := c.manager.RestartContainer(ctx, args.ContainerID); err != nil {
		return err
	}
	return nil
}

func (c *ContainerRepo) ImageList(ctx context.Context) ([]model.Image, error) {
	var reply []model.Image
	images, err := c.manager.ListImage(ctx)
	if err != nil {
		return reply, err
	}
	for _, im := range images {
		reply = append(reply, model.Image{
			Timestamp: time.Now(),
			ImageID:   im.ID[:6],
			Name:      im.Name,
			Tag:       im.Tag,
			Created:   im.Created,
			Size:      im.Size,
		})
	}
	return reply, nil
}

func (c *ContainerRepo) ImageCount(ctx context.Context) (int64, error) {
	images, err := c.manager.ListImage(ctx)
	if err != nil {
		return 0, err
	}
	return int64(len(images)), nil
}

func (c *ContainerRepo) ImagePull(ctx context.Context, args schema.ImagePullArgs) error {
	return c.manager.PullImage(ctx, args.ImageName)
}

func (c *ContainerRepo) ImageTag(ctx context.Context, args schema.ImageTagArgs) error {
	return c.manager.TagImage(ctx, args.OldTag, args.NewTag)
}

func (c *ContainerRepo) ImageImport(ctx context.Context, args schema.ImageImportArgs) error {
	return c.manager.ImportImage(ctx, args.SourceFile)
}

func (c *ContainerRepo) ImageExport(ctx context.Context, args schema.ImageExportArgs) error {
	return c.manager.ExportImage(ctx, []string{args.ImageID}, fmt.Sprintf("/tmp/%s", args.ImageName))
}

func (c *ContainerRepo) ImageDelete(ctx context.Context, args schema.ImageDeleteArgs) error {
	return c.manager.DeleteImage(ctx, args.ImageID)
}

func (c *ContainerRepo) ImagesPrune(ctx context.Context) error {
	return c.manager.PruneImages(ctx)
}

func (c *ContainerRepo) NetworkList(ctx context.Context) ([]model.Network, error) {
	var reply []model.Network
	nets, err := c.manager.ListNetwork(ctx)
	if err != nil {
		return reply, err
	}
	for _, net := range nets {
		subnet := ""
		gateway := ""
		labels, _ := json.Marshal(net.Labels)
		if len(net.SubNet) > 0 {
			subnet = net.SubNet[0].Subnet
			gateway = net.SubNet[0].Gateway
		}
		reply = append(reply, model.Network{
			Timestamp: time.Now(),
			NetworkID: net.ID[:6],
			Name:      net.Name,
			Driver:    net.Driver,
			Created:   net.Created,
			Scope:     net.Scope,
			Internal:  net.Internal,
			Subnet:    subnet,
			Gateway:   gateway,
			Labels:    string(labels),
		})
	}
	return reply, nil
}

func (c *ContainerRepo) NetworkCount(ctx context.Context) (int64, error) {
	var count int64 = 0
	nets, err := c.manager.ListNetwork(ctx)
	if err != nil {
		return count, err
	}
	return int64(len(nets)), nil
}

func (c *ContainerRepo) NetworkCreate(ctx context.Context, args schema.NetworkCreateArgs) (string, error) {
	return c.manager.CreateNetwork(ctx, args.Name, args.Driver, args.Subnet, args.Gateway, args.Labels)
}

func (c *ContainerRepo) NetworkDelete(ctx context.Context, args schema.NetworkDeleteArgs) error {
	return c.manager.DeleteNetwork(ctx, args.NetworkID)
}

func (c *ContainerRepo) GetDockerRegistryMirrors(ctx context.Context) ([]string, error) {
	var reply []string
	if _, err := os.Stat(constants.DaemonJsonPath); err != nil {
		return reply, err
	}
	file, err := os.ReadFile(constants.DaemonJsonPath)
	if err != nil {
		return reply, err
	}
	daemonMap := make(map[string]interface{})
	if err := json.Unmarshal(file, &daemonMap); err != nil {
		return reply, err
	}
	if _, ok := daemonMap["registry-mirrors"]; ok {
		for _, val := range daemonMap["registry-mirrors"].([]interface{}) {
			reply = append(reply, val.(string))
		}
	}
	return reply, nil
}

func (c *ContainerRepo) SetDockerRegistryMirrors(ctx context.Context, args schema.SetDockerRegistryMirrorsArgs) error {
	daemonMap := map[string]interface{}{
		"registry-mirrors": args.Mirrors,
	}
	setting, err := json.MarshalIndent(daemonMap, "", "    ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(constants.DaemonJsonPath, setting, 0644); err != nil {
		return err
	}
	if _, err := command.RunCommand(ctx, "systemctl", "daemon-reload"); err != nil {
		return err
	}
	if _, err := command.RunCommand(ctx, "systemctl", "restart", "docker-daemon"); err != nil {
		return err
	}
	return nil
}
