// Package repository
// Date: 2024/3/6 12:46
// Author: Amu
// Description:
package repository

import (
	"context"
	"github.com/amuluze/amprobe/service/schema"

	"github.com/amuluze/amprobe/service/model"
	"github.com/amuluze/amutool/database"
	"github.com/google/wire"
)

var ContainerRepoSet = wire.NewSet(NewContainerRepo, wire.Bind(new(IContainerRepo), new(*ContainerRepo)))

type IContainerRepo interface {
	ContainerList(ctx context.Context, args *schema.ContainerQueryArgs) (model.Containers, error)
	ContainerCount(ctx context.Context) (int, error)
	ImageList(ctx context.Context) (model.Images, error)
	Version(ctx context.Context) (model.Docker, error)
}

type ContainerRepo struct {
	DB *database.DB
}

func NewContainerRepo(db *database.DB) *ContainerRepo {
	return &ContainerRepo{DB: db}
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

func (a *ContainerRepo) ImageList(ctx context.Context) (model.Images, error) {
	var images model.Images
	if err := a.DB.Model(&model.Image{}).Group("name").Order("created_at desc").Find(&images).Error; err != nil {
		return images, err
	}
	return images, nil
}

func (a *ContainerRepo) Version(ctx context.Context) (model.Docker, error) {
	var docker model.Docker
	if err := a.DB.Model(&model.Docker{}).First(&docker).Error; err != nil {
		return docker, err
	}
	return docker, nil
}
