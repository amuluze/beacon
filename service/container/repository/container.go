// Package repository
// Date: 2024/3/6 12:46
// Author: Amu
// Description:
package repository

import (
	"context"

	"github.com/amuluze/amprobe/service/model"
	"github.com/amuluze/amutool/database"
	"github.com/google/wire"
)

var ContainerRepoSet = wire.NewSet(NewContainerRepo, wire.Bind(new(IContainerRepo), new(*ContainerRepo)))

type IContainerRepo interface {
	ContainerList(ctx context.Context) (model.Containers, error)
}

type ContainerRepo struct {
	DB *database.DB
}

func NewContainerRepo(db *database.DB) *ContainerRepo {
	return &ContainerRepo{DB: db}
}

func (a *ContainerRepo) ContainerList(ctx context.Context) (model.Containers, error) {
	var containers model.Containers
	if err := a.DB.Model(&model.Container{}).Group("name").Order("created_at desc").Find(&containers).Error; err != nil {
		return containers, err
	}
	return containers, nil
}
