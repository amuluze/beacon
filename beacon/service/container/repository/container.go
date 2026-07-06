// Package repository
package repository

import (
	"context"
	"time"

	"beacon/pkg/contextx"
	"beacon/pkg/rpc"
	"beacon/service/model"
	"common/database"
	rpcSchema "common/rpc/schema"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var ContainerServiceSet = wire.NewSet(NewContainerRepo, wire.Bind(new(IContainerRepo), new(*ContainerRepo)))

var _ IContainerRepo = (*ContainerRepo)(nil)

type IContainerRepo interface {
	Version(ctx context.Context, args rpcSchema.DockerArgs) (rpcSchema.DockerReply, error)

	ContainerList(ctx context.Context, args rpcSchema.ContainerQueryArgs) (rpcSchema.ContainerQueryReply, error)
	Usage(ctx context.Context, args rpcSchema.ContainerUsageArgs) (rpcSchema.ContainerUsageReply, error)
	ContainersByImage(ctx context.Context, image string) (num int, err error)
	ContainerCount(ctx context.Context, args rpcSchema.ContainerCountArgs) (rpcSchema.ContainerCountReply, error)
	ContainerCreate(ctx context.Context, args rpcSchema.ContainerCreateArgs) (rpcSchema.ContainerCreateReply, error)
	ContainerUpdate(ctx context.Context, args rpcSchema.ContainerUpdateArgs) (rpcSchema.ContainerUpdateReply, error)
	ContainerDelete(ctx context.Context, args rpcSchema.ContainerDeleteArgs) error
	ContainerStart(ctx context.Context, args rpcSchema.ContainerStartArgs) error
	ContainerStop(ctx context.Context, args rpcSchema.ContainerStopArgs) error
	ContainerRestart(ctx context.Context, args rpcSchema.ContainerRestartArgs) error
	ContainerLogs(ctx context.Context, args rpcSchema.ContainerLogsArgs) (rpcSchema.ContainerLogsReply, error)

	ImageList(ctx context.Context, args rpcSchema.ImageQueryArgs) (rpcSchema.ImageQueryReply, error)
	ImageCount(ctx context.Context, args rpcSchema.ImageCountArgs) (rpcSchema.ImageCountReply, error)
	ImagePull(ctx context.Context, args rpcSchema.ImagePullArgs) error
	ImageTag(ctx context.Context, args rpcSchema.ImageTagArgs) error
	ImageImport(ctx context.Context, args rpcSchema.ImageImportArgs) error
	ImageExport(ctx context.Context, args rpcSchema.ImageExportArgs) (rpcSchema.ImageExportReply, error)
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
	RPCClient rpc.Caller
	DB        *database.DB
}

func NewContainerRepo(client rpc.Caller, db *database.DB) *ContainerRepo {
	return &ContainerRepo{RPCClient: client, DB: db}
}

func (c *ContainerRepo) agentDB(ctx context.Context) (*gorm.DB, error) {
	agentID, err := contextx.RequireAgentID(ctx)
	if err != nil {
		return nil, err
	}
	return c.DB.DB.Where("agent_id = ?", agentID), nil
}

// ── Monitoring queries (local DB) ──

func (c *ContainerRepo) Version(ctx context.Context, args rpcSchema.DockerArgs) (rpcSchema.DockerReply, error) {
	var result model.MonitorDocker
	db, err := c.agentDB(ctx)
	if err != nil {
		return rpcSchema.DockerReply{}, err
	}
	if err := db.Model(&model.MonitorDocker{}).First(&result).Error; err != nil {
		return rpcSchema.DockerReply{}, err
	}
	return rpcSchema.DockerReply{
		Data: rpcSchema.Docker{
			Timestamp:     result.Timestamp,
			DockerVersion: result.DockerVersion,
			APIVersion:    result.APIVersion,
			MinAPIVersion: result.MinAPIVersion,
			GitCommit:     result.GitCommit,
			GoVersion:     result.GoVersion,
			Os:            result.Os,
			Arch:          result.Arch,
		},
		Freshness: rpcSchema.ComputeFreshness(result.Timestamp),
	}, nil
}

func (c *ContainerRepo) ContainerList(ctx context.Context, args rpcSchema.ContainerQueryArgs) (rpcSchema.ContainerQueryReply, error) {
	var containers []model.MonitorContainer
	agentID, err := contextx.RequireAgentID(ctx)
	if err != nil {
		return rpcSchema.ContainerQueryReply{}, err
	}
	latest := c.DB.DB.Model(&model.MonitorContainer{}).
		Where("ports != ?", "").
		Where("agent_id = ?", agentID).
		Select("agent_id, name, MAX(timestamp) AS timestamp").
		Group("agent_id, name")
	query := c.DB.DB.Model(&model.MonitorContainer{}).
		Where("m_container.ports != ?", "").
		Where("m_container.agent_id = ?", agentID)
	if err := query.
		Joins("JOIN (?) latest ON latest.agent_id = m_container.agent_id AND latest.name = m_container.name AND latest.timestamp = m_container.timestamp", latest).
		Order("m_container.name asc").
		Offset((args.Page - 1) * args.Size).Limit(args.Size).Find(&containers).Error; err != nil {
		return rpcSchema.ContainerQueryReply{}, err
	}
	var results []rpcSchema.Container
	var latestTimestamp time.Time
	for _, container := range containers {
		if container.Timestamp.After(latestTimestamp) {
			latestTimestamp = container.Timestamp
		}
		results = append(results, rpcSchema.Container{
			Timestamp:   container.Timestamp,
			ContainerID: container.ContainerID,
			Name:        container.Name,
			Image:       container.Image,
			IP:          container.IP,
			Ports:       container.Ports,
			State:       container.State,
			Uptime:      container.Uptime,
			CPUPercent:  container.CPUPercent,
			MemPercent:  container.MemPercent,
			MemUsage:    container.MemUsage,
			MemLimit:    container.MemLimit,
			Labels:      container.Labels,
		})
	}
	return rpcSchema.ContainerQueryReply{Data: results, Freshness: rpcSchema.ComputeFreshness(latestTimestamp)}, nil
}

func (c *ContainerRepo) Usage(ctx context.Context, args rpcSchema.ContainerUsageArgs) (rpcSchema.ContainerUsageReply, error) {
	var containers []model.MonitorContainer
	db, err := c.agentDB(ctx)
	if err != nil {
		return rpcSchema.ContainerUsageReply{}, err
	}
	if err := db.Model(&model.MonitorContainer{}).
		Order("timestamp asc").
		Where("timestamp > ?", time.Unix(args.StartTime, 0)).Find(&containers).Error; err != nil {
		return rpcSchema.ContainerUsageReply{}, err
	}
	reply := rpcSchema.ContainerUsageReply{
		Names:    make([]string, 0),
		CPUUsage: make(map[string][]rpcSchema.Usage),
		MemUsage: make(map[string][]rpcSchema.Usage),
	}
	var latestTimestamp time.Time
	for _, item := range containers {
		if item.Timestamp.After(latestTimestamp) {
			latestTimestamp = item.Timestamp
		}
		if _, ok := reply.CPUUsage[item.Name]; !ok {
			reply.Names = append(reply.Names, item.Name)
			reply.CPUUsage[item.Name] = make([]rpcSchema.Usage, 0)
			reply.MemUsage[item.Name] = make([]rpcSchema.Usage, 0)
		}
		reply.CPUUsage[item.Name] = append(reply.CPUUsage[item.Name], rpcSchema.Usage{
			Timestamp: item.Timestamp.Unix(), Value: item.CPUPercent,
		})
		reply.MemUsage[item.Name] = append(reply.MemUsage[item.Name], rpcSchema.Usage{
			Timestamp: item.Timestamp.Unix(), Value: item.MemUsage,
		})
	}
	reply.Freshness = rpcSchema.ComputeFreshness(latestTimestamp)
	return reply, nil
}

func (c *ContainerRepo) ContainersByImage(ctx context.Context, image string) (num int, err error) {
	var count int64
	agentID, err := contextx.RequireAgentID(ctx)
	if err != nil {
		return 0, err
	}
	distinctContainers := c.DB.DB.Model(&model.MonitorContainer{}).
		Where("image = ?", image).
		Where("agent_id = ?", agentID).
		Select("agent_id, name").
		Group("agent_id, name")
	if err := c.DB.Table("(?) as containers", distinctContainers).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (c *ContainerRepo) ContainerCount(ctx context.Context, args rpcSchema.ContainerCountArgs) (rpcSchema.ContainerCountReply, error) {
	var count int64
	agentID, err := contextx.RequireAgentID(ctx)
	if err != nil {
		return rpcSchema.ContainerCountReply{}, err
	}
	distinctContainers := c.DB.DB.Model(&model.MonitorContainer{}).
		Where("ports != ?", "").
		Where("agent_id = ?", agentID).
		Select("agent_id, name").
		Group("agent_id, name")
	if err := c.DB.Table("(?) as containers", distinctContainers).Count(&count).Error; err != nil {
		return rpcSchema.ContainerCountReply{}, err
	}
	return rpcSchema.ContainerCountReply{Count: int(count)}, nil
}

func (c *ContainerRepo) ImageList(ctx context.Context, args rpcSchema.ImageQueryArgs) (rpcSchema.ImageQueryReply, error) {
	var results []model.MonitorImage
	db, err := c.agentDB(ctx)
	if err != nil {
		return rpcSchema.ImageQueryReply{}, err
	}
	if err := db.Model(&model.MonitorImage{}).
		Order("created_at desc").Offset((args.Page - 1) * args.Size).Limit(args.Size).Find(&results).Error; err != nil {
		return rpcSchema.ImageQueryReply{}, err
	}
	var list []rpcSchema.Image
	var latestTimestamp time.Time
	for _, result := range results {
		if result.Timestamp.After(latestTimestamp) {
			latestTimestamp = result.Timestamp
		}
		list = append(list, rpcSchema.Image{
			Timestamp: result.Timestamp,
			ImageID:   result.ImageID,
			Name:      result.Name,
			Tag:       result.Tag,
			Created:   result.Created,
			Size:      result.Size,
		})
	}
	return rpcSchema.ImageQueryReply{Data: list, Freshness: rpcSchema.ComputeFreshness(latestTimestamp)}, nil
}

func (c *ContainerRepo) ImageCount(ctx context.Context, args rpcSchema.ImageCountArgs) (rpcSchema.ImageCountReply, error) {
	var total int64
	db, err := c.agentDB(ctx)
	if err != nil {
		return rpcSchema.ImageCountReply{}, err
	}
	if err := db.Model(&model.MonitorImage{}).Order("created_at desc").Count(&total).Error; err != nil {
		return rpcSchema.ImageCountReply{}, err
	}
	return rpcSchema.ImageCountReply{Count: int(total)}, nil
}

func (c *ContainerRepo) NetworkList(ctx context.Context, args rpcSchema.NetworkQueryArgs) (rpcSchema.NetworkQueryReply, error) {
	var networks []model.MonitorNetwork
	db, err := c.agentDB(ctx)
	if err != nil {
		return rpcSchema.NetworkQueryReply{}, err
	}
	if err := db.Model(&model.MonitorNetwork{}).
		Order("created_at desc").Offset((args.Page - 1) * args.Size).Limit(args.Size).Find(&networks).Error; err != nil {
		return rpcSchema.NetworkQueryReply{}, err
	}
	var list []rpcSchema.Network
	var latestTimestamp time.Time
	for _, n := range networks {
		if n.Timestamp.After(latestTimestamp) {
			latestTimestamp = n.Timestamp
		}
		list = append(list, rpcSchema.Network{
			Timestamp: n.Timestamp,
			NetworkID: n.NetworkID,
			Name:      n.Name,
			Driver:    n.Driver,
			Scope:     n.Scope,
			Created:   n.Created,
			Internal:  n.Internal,
			Subnet:    n.Subnet,
			Gateway:   n.Gateway,
			Labels:    n.Labels,
		})
	}
	return rpcSchema.NetworkQueryReply{Data: list, Freshness: rpcSchema.ComputeFreshness(latestTimestamp)}, nil
}

func (c *ContainerRepo) NetworkCount(ctx context.Context, args rpcSchema.NetworkCountArgs) (rpcSchema.NetworkCountReply, error) {
	var total int64
	db, err := c.agentDB(ctx)
	if err != nil {
		return rpcSchema.NetworkCountReply{}, err
	}
	if err := db.Model(&model.MonitorNetwork{}).Order("created_at desc").Count(&total).Error; err != nil {
		return rpcSchema.NetworkCountReply{}, err
	}
	return rpcSchema.NetworkCountReply{Count: int(total)}, nil
}

func (c *ContainerRepo) GetDockerRegistryMirrors(ctx context.Context, args rpcSchema.GetDockerRegistryMirrorsArgs) (rpcSchema.GetDockerRegistryMirrorsReply, error) {
	// This still requires Agent access
	var reply rpcSchema.GetDockerRegistryMirrorsReply
	err := c.RPCClient.Call(ctx, "GetDockerRegistryMirrors", args, &reply)
	return reply, err
}

// ── Operational commands (RPC to Agent) ──

func (c *ContainerRepo) ContainerCreate(ctx context.Context, args rpcSchema.ContainerCreateArgs) (rpcSchema.ContainerCreateReply, error) {
	var reply rpcSchema.ContainerCreateReply
	err := c.RPCClient.Call(ctx, "ContainerCreate", args, &reply)
	return reply, err
}

func (c *ContainerRepo) ContainerUpdate(ctx context.Context, args rpcSchema.ContainerUpdateArgs) (rpcSchema.ContainerUpdateReply, error) {
	var reply rpcSchema.ContainerUpdateReply
	err := c.RPCClient.Call(ctx, "ContainerUpdate", args, &reply)
	return reply, err
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

func (c *ContainerRepo) ContainerLogs(ctx context.Context, args rpcSchema.ContainerLogsArgs) (rpcSchema.ContainerLogsReply, error) {
	var reply rpcSchema.ContainerLogsReply
	err := c.RPCClient.Call(ctx, "ContainerLogs", args, &reply)
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

func (c *ContainerRepo) ImageExport(ctx context.Context, args rpcSchema.ImageExportArgs) (rpcSchema.ImageExportReply, error) {
	var reply rpcSchema.ImageExportReply
	err := c.RPCClient.Call(ctx, "ImageExport", args, &reply)
	return reply, err
}

func (c *ContainerRepo) ImageDelete(ctx context.Context, args rpcSchema.ImageDeleteArgs) error {
	var reply rpcSchema.ImageDeleteReply
	return c.RPCClient.Call(ctx, "ImageDelete", args, &reply)
}

func (c *ContainerRepo) ImagesPrune(ctx context.Context) error {
	return c.RPCClient.Call(ctx, "ImagesPrune", nil, nil)
}

func (c *ContainerRepo) NetworkCreate(ctx context.Context, args rpcSchema.NetworkCreateArgs) (rpcSchema.NetworkCreateReply, error) {
	var reply rpcSchema.NetworkCreateReply
	err := c.RPCClient.Call(ctx, "NetworkCreate", args, &reply)
	return reply, err
}

func (c *ContainerRepo) NetworkDelete(ctx context.Context, args rpcSchema.NetworkDeleteArgs) error {
	var reply rpcSchema.NetworkDeleteReply
	return c.RPCClient.Call(ctx, "NetworkDelete", args, &reply)
}

func (c *ContainerRepo) SetDockerRegistryMirrors(ctx context.Context, args rpcSchema.SetDockerRegistryMirrorsArgs) error {
	var reply rpcSchema.SetDockerRegistryMirrorsReply
	return c.RPCClient.Call(ctx, "SetDockerRegistryMirrors", args, &reply)
}
