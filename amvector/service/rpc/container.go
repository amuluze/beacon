// Package rpc
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package rpc

import (
	rpcSchema "common/rpc/schema"
	"context"
	"encoding/json"
	"gorm.io/gorm"
	"log/slog"
	"time"

	"amvector/service/model"
	"github.com/amuluze/docker"
)

func (s *Service) Version(ctx context.Context, args rpcSchema.DockerArgs, reply *rpcSchema.DockerReply) error {
	var result model.Docker
	if err := s.DB.Model(&model.Docker{}).First(&result).Error; err != nil {
		return err
	}
	reply.Data.Timestamp = result.Timestamp
	reply.Data.DockerVersion = result.DockerVersion
	reply.Data.APIVersion = result.APIVersion
	reply.Data.MinAPIVersion = result.MinAPIVersion
	reply.Data.GitCommit = result.GitCommit
	reply.Data.GoVersion = result.GoVersion
	reply.Data.Os = result.Os
	reply.Data.Arch = result.Arch
	return nil
}

func (s *Service) ContainerList(ctx context.Context, args rpcSchema.ContainerQueryArgs, reply *rpcSchema.ContainerQueryReply) error {
	var containers []model.Container
	if err := s.DB.Model(&model.Container{}).Order("created_at desc").Offset((args.Page - 1) * args.Size).Limit(args.Size).Find(&containers).Error; err != nil {
		return err
	}
	var results []rpcSchema.Container
	for _, container := range containers {
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
	reply.Data = results
	return nil
}

func (s *Service) ContainerCount(ctx context.Context, args rpcSchema.ContainerCountArgs, reply *rpcSchema.ContainerCountReply) error {
	var count int64
	if err := s.DB.Model(&model.Container{}).Count(&count).Error; err != nil {
		return err
	}
	reply.Count = int(count)
	return nil
}

func (s *Service) ContainerCreate(ctx context.Context, args rpcSchema.ContainerCreateArgs, reply *rpcSchema.ContainerCreateReply) error {
	var containerID string
	var err error
	if containerID, err = s.Manager.CreateContainer(
		ctx,
		args.Name,
		args.Image,
		args.Network,
		args.Ports,
		args.Volumes,
		args.Environments,
		nil,
		args.Labels,
	); err != nil {
		return err
	}
	s.containerTask()
	reply.ContainerID = containerID
	return nil
}

func (s *Service) ContainerUpdate(ctx context.Context, args rpcSchema.ContainerUpdateArgs, reply *rpcSchema.ContainerUpdateReply) error {
	return nil
}

func (s *Service) containerTask() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cs, err := s.Manager.ListContainer(ctx)
	if err != nil {
		slog.Error("failed to list containers", "error", err)
		return
	}
	var containers []model.Container
	for _, info := range cs {
		labels, _ := json.Marshal(info.Labels)
		var d model.Container
		d.Timestamp = time.Now()
		d.ContainerID = info.ID[:6]
		d.Name = info.Name
		d.State = info.State
		d.Image = info.Image
		d.Uptime = info.Uptime
		d.IP = info.IP
		d.Labels = string(labels)

		cpuPercent, err := s.Manager.GetContainerCpu(ctx, info.ID[:6])
		if err != nil {
			slog.Error("failed to get container cpu", "error", err)
		}
		d.CPUPercent = cpuPercent

		memPercent, used, limit, err := s.Manager.GetContainerMem(ctx, info.ID[:6])
		if err != nil {
			slog.Error("failed to get container mem", "error", err)
		}
		d.MemPercent = memPercent

		d.MemUsage = used
		d.MemLimit = limit
		if _, ok := s.cache.Get(info.Image); !ok {
			s.cache.Set(info.Image, 1, 2*time.Minute)
		} else {
			count, err := s.cache.IncrementInt(info.Image, 1)
			slog.Info("container image cache", "image", info.Image, "count", count, "error", err)
		}
		containers = append(containers, d)
	}
	if err := s.DB.Unscoped().Where("1 = 1").Delete(&model.Container{}).Error; err != nil {
		slog.Error("failed to delete container", "error", err)
	}
	s.DB.Model(&model.Container{}).Create(&containers)
}

func (s *Service) ContainerDelete(ctx context.Context, args rpcSchema.ContainerDeleteArgs, reply *rpcSchema.ContainerDeleteReply) error {
	if err := s.DB.RunInTransaction(func(tx *gorm.DB) error {
		if err := s.Manager.DeleteContainer(ctx, args.ContainerID); err != nil {
			return err
		}
		if err := tx.Model(&model.Container{}).Where("container_id = ?", args.ContainerID).Delete(&model.Container{}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *Service) ContainerStart(ctx context.Context, args rpcSchema.ContainerStartArgs, reply *rpcSchema.ContainerStartReply) error {
	if err := s.Manager.StartContainer(ctx, args.ContainerID); err != nil {
		return err
	}
	if err := s.DB.Model(&model.Container{}).Where("container_id = ?", args.ContainerID).Update("state", "running").Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) ContainerStop(ctx context.Context, args rpcSchema.ContainerStopArgs, reply *rpcSchema.ContainerStopReply) error {
	if err := s.Manager.StopContainer(ctx, args.ContainerID); err != nil {
		return err
	}
	if err := s.DB.Model(&model.Container{}).Where("container_id = ?", args.ContainerID).Update("state", "stopped").Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) ContainerRestart(ctx context.Context, args rpcSchema.ContainerRestartArgs, reply *rpcSchema.ContainerRestartReply) error {
	if err := s.Manager.RestartContainer(ctx, args.ContainerID); err != nil {
		return err
	}
	if err := s.DB.Model(&model.Container{}).Where("container_id = ?", args.ContainerID).Update("state", "running").Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) ImageList(ctx context.Context, args rpcSchema.ImageQueryArgs, reply *rpcSchema.ImageQueryReply) error {
	var results []model.Image
	if err := s.DB.Model(&model.Image{}).Order("created_at desc").Offset((args.Page - 1) * args.Size).Limit(args.Size).Find(&results).Error; err != nil {
		return err
	}
	var list []rpcSchema.Image
	for _, result := range results {
		list = append(list, rpcSchema.Image{
			Timestamp: result.Timestamp,
			ImageID:   result.ImageID,
			Name:      result.Name,
			Tag:       result.Tag,
			Created:   result.Created,
			Size:      result.Size,
		})
	}
	reply.Data = list
	return nil
}

func (s *Service) ImagePull(ctx context.Context, args rpcSchema.ImagePullArgs, reply *rpcSchema.ImagePullReply) error {
	if err := s.Manager.PullImage(ctx, args.ImageName); err != nil {
		return err
	}
	s.imageTask()
	return nil
}

func (s *Service) ImageTag(ctx context.Context, args rpcSchema.ImageTagArgs, reply *rpcSchema.ImageTagReply) error {
	if err := s.Manager.TagImage(ctx, args.OldTag, args.NewTag); err != nil {
		return err
	}
	s.imageTask()
	return nil
}

func (s *Service) ImageCount(ctx context.Context, args rpcSchema.ImageCountArgs, reply *rpcSchema.ImageCountReply) error {
	var total int64
	if err := s.DB.Model(&model.Images{}).Order("created_at desc").Count(&total).Error; err != nil {
		return err
	}
	reply.Count = int(total)
	return nil
}

func (s *Service) ImageDelete(ctx context.Context, args rpcSchema.ImageDeleteArgs, reply *rpcSchema.ImageDeleteReply) error {
	if err := s.Manager.DeleteImage(ctx, args.ImageID); err != nil {
		return err
	}
	s.imageTask()
	return nil
}

func (s *Service) imageTask() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	images, err := s.Manager.ListImage(ctx)
	if err != nil {
		slog.Error("failed to get version", "error", err)
		return
	}
	var list model.Images
	duplicateImage := make(map[string]struct{})
	for _, im := range images {
		val, ok := s.cache.Get(im.Name + ":" + im.Tag)
		if !ok {
			slog.Error("failed to get image cache", "error", err)
			val = 0
		}
		if _, ok := duplicateImage[im.ID]; !ok {
			duplicateImage[im.ID] = struct{}{}
		} else {
			if im.Tag != "latest" {
				continue
			}
		}
		list = append(list, model.Image{
			Timestamp: time.Now(),
			ImageID:   im.ID[7:19],
			Name:      im.Name,
			Number:    val.(int),
			Tag:       im.Tag,
			Created:   im.Created,
			Size:      im.Size,
		})
		s.cache.Delete(im.Name + ":" + im.Tag)
	}
	if err := s.DB.Unscoped().Where("1 = 1").Delete(&model.Image{}).Error; err != nil {
		slog.Error("failed to delete image", "error", err)
	}
	s.DB.Model(&model.Image{}).Create(&list)
}

func (s *Service) ImagesPrune(ctx context.Context) error {
	return s.Manager.PruneImages(ctx)
}

func (s *Service) ImageImport(ctx context.Context, args rpcSchema.ImageImportArgs, reply *rpcSchema.ImageImportReply) error {
	if err := s.Manager.ImportImage(ctx, args.SourceFile); err != nil {
		return err
	}
	s.imageTask()
	return nil
}

func (s *Service) ImageExport(ctx context.Context, args rpcSchema.ImageExportArgs, reply *rpcSchema.ImageExportReply) error {
	if err := s.Manager.ExportImage(ctx, args.ImageIDs, args.TargetFile); err != nil {
		return err
	}
	return nil
}

func (s *Service) NetworkCreate(ctx context.Context, args rpcSchema.NetworkCreateArgs, reply *rpcSchema.NetworkCreateReply) error {
	args.Labels[docker.CreatedByProbe] = "true"
	if networkID, err := s.Manager.CreateNetwork(ctx, args.Name, args.Driver, args.Subnet, args.Gateway, args.Labels); err != nil {
		return err
	} else {
		s.networkTask()
		reply.NetworkID = networkID
		return nil
	}
}

func (s *Service) networkTask() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	nets, err := s.Manager.ListNetwork(ctx)
	if err != nil {
		slog.Error("failed to get network", "error", err)
		return
	}
	var list model.Networks
	for _, net := range nets {
		subnet := ""
		gateway := ""
		labels, _ := json.Marshal(net.Labels)
		if len(net.SubNet) > 0 {
			subnet = net.SubNet[0].Subnet
			gateway = net.SubNet[0].Gateway
		}
		list = append(list, model.Network{
			Timestamp: time.Now(),
			NetworkID: net.ID,
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
	if err := s.DB.Unscoped().Where("1 = 1").Delete(&model.Network{}).Error; err != nil {
		slog.Error("failed to delete network", "error", err)
	}
	s.DB.Model(&model.Network{}).Create(&list)
}

func (s *Service) NetworkList(ctx context.Context, args rpcSchema.NetworkQueryArgs, reply *rpcSchema.NetworkQueryReply) error {
	if err := s.DB.Model(&model.Network{}).Order("created_at desc").Offset((args.Page - 1) * args.Size).Limit(args.Size).Find(reply).Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) NetworkCount(ctx context.Context, args rpcSchema.NetworkCountArgs, reply *rpcSchema.NetworkCountReply) error {
	var total int64
	if err := s.DB.Model(&model.Network{}).Order("created_at desc").Count(&total).Error; err != nil {
		return err
	}
	reply.Count = int(total)
	return nil
}

func (s *Service) NetworkDelete(ctx context.Context, args rpcSchema.NetworkDeleteArgs, reply *rpcSchema.NetworkDeleteReply) error {
	if err := s.Manager.DeleteNetwork(ctx, args.NetworkID); err != nil {
		return err
	}
	s.networkTask()
	return nil
}
