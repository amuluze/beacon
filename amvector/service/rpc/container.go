// Package rpc
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package rpc

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/amuluze/docker"

	"amvector/service/model"

	"amvector/service/schema"
)

func (s *Service) DockerVersion(ctx context.Context, args schema.VersionArgs, reply *model.Docker) error {
	if err := s.DB.Model(&model.Docker{}).First(reply).Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) ContainerList(ctx context.Context, args schema.ContainerQueryArgs, reply *model.Containers) error {
	if err := s.DB.Model(&model.Container{}).Order("created_at desc").Offset((args.Page - 1) * args.Size).Limit(args.Size).Find(reply).Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) ContainerCount(ctx context.Context, args schema.QueryCountArgs, reply *schema.QueryCountReply) error {
	var count int64
	if err := s.DB.Model(&model.Container{}).Count(&count).Error; err != nil {
		return err
	}
	reply.Count = int(count)
	return nil
}

func (s *Service) ContainerCreate(ctx context.Context, args schema.ContainerCreateArgs, reply *schema.ContainerCreateReply) error {
	var containerID string
	var err error
	if containerID, err = s.Manager.CreateContainer(
		ctx,
		args.ContainerName,
		args.ImageName,
		args.NetworkName,
		args.Ports,
		args.Volumes,
		args.Environments,
		args.Labels,
	); err != nil {
		return err
	}
	s.containerTask()
	s.imageTask()
	reply.ContainerID = containerID
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

func (s *Service) ContainerDelete(ctx context.Context, args schema.ContainerDeleteArgs, reply *schema.ContainerDeleteReply) error {
	if err := s.Manager.DeleteContainer(ctx, args.ContainerID); err != nil {
		return err
	}
	s.containerTask()
	return nil
}

func (s *Service) ContainerStart(ctx context.Context, args schema.ContainerStartArgs, reply *schema.ContainerStartReply) error {
	if err := s.Manager.StartContainer(ctx, args.ContainerID); err != nil {
		return err
	}
	if err := s.DB.Model(&model.Container{}).Where("container_id = ?", args.ContainerID).Update("state", "running").Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) ContainerStop(ctx context.Context, args schema.ContainerStopArgs, reply *schema.ContainerStopReply) error {
	if err := s.Manager.StopContainer(ctx, args.ContainerID); err != nil {
		return err
	}
	if err := s.DB.Model(&model.Container{}).Where("container_id = ?", args.ContainerID).Update("state", "stopped").Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) ContainerRestart(ctx context.Context, args schema.ContainerRestartArgs, reply *schema.ContainerRestartReply) error {
	if err := s.Manager.RestartContainer(ctx, args.ContainerID); err != nil {
		return err
	}
	if err := s.DB.Model(&model.Container{}).Where("container_id = ?", args.ContainerID).Update("state", "running").Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) ImageList(ctx context.Context, args schema.ImageQueryArgs, reply *model.Images) error {
	if err := s.DB.Model(&model.Image{}).Order("created_at desc").Offset((args.Page - 1) * args.Size).Limit(args.Size).Find(reply).Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) ImagePull(ctx context.Context, args schema.ImagePullArgs, reply *schema.ImagePullReply) error {
	if err := s.Manager.PullImage(ctx, args.ImageName); err != nil {
		return err
	}
	s.imageTask()
	return nil
}

func (s *Service) ImageTag(ctx context.Context, args schema.ImageTagArgs, reply *schema.ImageTagReply) error {
	if err := s.Manager.TagImage(ctx, args.OldTag, args.NewTag); err != nil {
		return err
	}
	s.imageTask()
	return nil
}

func (s *Service) ImageCount(ctx context.Context, args schema.ImageCountArgs, reply *schema.ImageCountReply) error {
	var total int64
	if err := s.DB.Model(&model.Images{}).Order("created_at desc").Count(&total).Error; err != nil {
		return err
	}
	reply.Count = int(total)
	return nil
}

func (s *Service) ImageDelete(ctx context.Context, args schema.ImageDeleteArgs, reply *schema.ImageDeleteReply) error {
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

func (s *Service) ImageImport(ctx context.Context, args schema.ImageImportArgs, reply *schema.ImageImportReply) error {
	if err := s.Manager.ImportImage(ctx, args.SourceFile); err != nil {
		return err
	}
	s.imageTask()
	return nil
}

func (s *Service) ImageExport(ctx context.Context, args schema.ImageExportArgs, reply *schema.ImageExportReply) error {
	if err := s.Manager.ExportImage(ctx, args.ImageIDs, args.TargetFile); err != nil {
		return err
	}
	return nil
}

func (s *Service) NetworkCreate(ctx context.Context, args schema.NetworkCreateArgs, reply *schema.NetworkCreateReply) error {
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

func (s *Service) NetworkList(ctx context.Context, args schema.NetworkQueryArgs, reply *model.Networks) error {
	if err := s.DB.Model(&model.Network{}).Order("created_at desc").Offset((args.Page - 1) * args.Size).Limit(args.Size).Find(reply).Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) NetworkCount(ctx context.Context, args schema.NetworkCountArgs, reply *schema.NetworkCountReply) error {
	var total int64
	if err := s.DB.Model(&model.Network{}).Order("created_at desc").Count(&total).Error; err != nil {
		return err
	}
	reply.Count = int(total)
	return nil
}

func (s *Service) NetworkDelete(ctx context.Context, args schema.NetworkDeleteArgs, reply *schema.NetworkDeleteReply) error {
	if err := s.Manager.DeleteNetwork(ctx, args.NetworkID); err != nil {
		return err
	}
	s.networkTask()
	return nil
}
