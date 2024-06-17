// Package rpc
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package rpc

import (
	"context"

	"github.com/amuluze/amvector/service/model"

	"github.com/amuluze/amvector/service/schema"
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
	return nil
}

func (s *Service) ContainerDelete(ctx context.Context, args schema.ContainerDeleteArgs, reply *schema.ContainerDeleteReply) error {
	if err := s.Manager.DeleteContainer(ctx, args.ContainerID); err != nil {
		return err
	}
	if err := s.DB.Model(&model.Container{}).Delete(&model.Container{ContainerID: args.ContainerID}).Error; err != nil {
		return err
	}
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
	return nil
}

func (s *Service) ImageTag(ctx context.Context, args schema.ImageTagArgs, reply *schema.ImageTagReply) error {
	if err := s.Manager.TagImage(ctx, args.OldTag, args.NewTag); err != nil {
		return err
	}
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
	if err := s.Manager.RemoveImage(ctx, args.ImageID); err != nil {
		return err
	}
	if err := s.DB.Where("image_id = ?", args.ImageID).Delete(&model.Image{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) ImagesPrune(ctx context.Context) error {
	return s.Manager.PruneImages(ctx)
}
