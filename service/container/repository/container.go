// Package repository
// Date: 2024/3/6 12:46
// Author: Amu
// Description:
package repository

import (
	"context"
	"github.com/amuluze/amprobe/service/schema"
	"github.com/amuluze/amutool/docker"
	"github.com/amuluze/amutool/errors"

	"github.com/amuluze/amprobe/service/model"
	"github.com/amuluze/amutool/database"
	"github.com/google/wire"
)

var ContainerRepoSet = wire.NewSet(NewContainerRepo, wire.Bind(new(IContainerRepo), new(*ContainerRepo)))

type IContainerRepo interface {
	ContainerList(ctx context.Context, args *schema.ContainerQueryArgs) (model.Containers, error)
	ContainerCount(ctx context.Context) (int, error)
	ContainerStart(ctx context.Context, args *schema.ContainerStartArgs) error
	ContainerStop(ctx context.Context, args *schema.ContainerStopArgs) error
	ContainerRemove(ctx context.Context, args *schema.ContainerRemoveArgs) error
	ContainerRestart(ctx context.Context, args *schema.ContainerRestartArgs) error
	ImageList(ctx context.Context, args *schema.ImageQueryArgs) (model.Images, error)
	ImageRemove(ctx context.Context, args *schema.ImageRemoveArgs) error
	ImagesPrune(ctx context.Context) error
	ImageCount(ctx context.Context) (int, error)
	Version(ctx context.Context) (model.Docker, error)
}

type ContainerRepo struct {
	DB      *database.DB
	Manager *docker.Manager
}

func NewContainerRepo(db *database.DB) *ContainerRepo {
	manager, err := docker.NewManager()
	if err != nil {
		panic(err)
	}
	return &ContainerRepo{DB: db, Manager: manager}
}

func (a *ContainerRepo) ContainerList(ctx context.Context, args *schema.ContainerQueryArgs) (model.Containers, error) {
	var containers model.Containers
	if err := a.DB.Model(&model.Container{}).Group("name").Order("created_at desc").Offset((args.Page - 1) * args.Size).Limit(args.Size).Find(&containers).Error; err != nil {
		return containers, err
	}
	return containers, nil
}

func (a *ContainerRepo) ContainerCount(ctx context.Context) (int, error) {
	var total int64
	if err := a.DB.Model(&model.Container{}).Group("name").Order("created_at desc").Count(&total).Error; err != nil {
		return int(total), err
	}
	return int(total), nil
}

func (a *ContainerRepo) ImageList(ctx context.Context, args *schema.ImageQueryArgs) (model.Images, error) {
	var images model.Images
	if err := a.DB.Model(&model.Image{}).Group("name").Order("created_at desc").Offset((args.Page - 1) * args.Size).Limit(args.Size).Find(&images).Error; err != nil {
		return images, err
	}
	return images, nil
}

func (a *ContainerRepo) ImageCount(ctx context.Context) (int, error) {
	var total int64
	if err := a.DB.Model(&model.Images{}).Group("name").Order("created_at desc").Count(&total).Error; err != nil {
		return int(total), err
	}
	return int(total), nil
}

func (a *ContainerRepo) Version(ctx context.Context) (model.Docker, error) {
	var docker model.Docker
	if err := a.DB.Model(&model.Docker{}).First(&docker).Error; err != nil {
		return docker, err
	}
	return docker, nil
}

func (a *ContainerRepo) ContainerStart(ctx context.Context, args *schema.ContainerStartArgs) error {
	err := a.Manager.StartContainer(ctx, args.ContainerID)
	if err != nil {
		return errors.New400Error("failed start container")
	}
	a.DB.Model(&model.Container{}).Where("container_id = ?", args.ContainerID).Update("state", "running")
	return err
}

func (a *ContainerRepo) ContainerStop(ctx context.Context, args *schema.ContainerStopArgs) error {
	err := a.Manager.StopContainer(ctx, args.ContainerID)
	if err != nil {
		return errors.New400Error("failed to stop container")
	}
	a.DB.Model(&model.Container{}).Where("container_id = ?", args.ContainerID).Update("state", "exited")
	return err
}

func (a *ContainerRepo) ContainerRemove(ctx context.Context, args *schema.ContainerRemoveArgs) error {
	err := a.Manager.DeleteContainer(ctx, args.ContainerID)
	if err != nil {
		return errors.New400Error("failed to remove container")
	}
	a.DB.Model(&model.Container{}).Delete(&model.Container{ContainerID: args.ContainerID})
	return err
}

func (a *ContainerRepo) ContainerRestart(ctx context.Context, args *schema.ContainerRestartArgs) error {
	err := a.Manager.RestartContainer(ctx, args.ContainerID)
	if err != nil {
		return errors.New400Error("failed to restart container")
	}
	a.DB.Model(&model.Container{}).Where("container_id = ?", args.ContainerID).Update("state", "running")
	return err
}

func (a *ContainerRepo) ImageRemove(ctx context.Context, args *schema.ImageRemoveArgs) error {
	err := a.Manager.RemoveImage(ctx, args.ImageID)
	if err != nil {
		return errors.New400Error(err.Error())
	}
	a.DB.Where("image_id = ?", args.ImageID).Delete(&model.Image{})
	return err
}

func (a *ContainerRepo) ImagesPrune(ctx context.Context) error {
	return a.Manager.PruneImages(ctx)
}
