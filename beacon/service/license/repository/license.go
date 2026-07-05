// Package repository
// Date: 2023/6/7 12:59
// Author: Amu
// Description:
package repository

import (
	"beacon/service/model"
	"beacon/service/schema"
	"common/database"
	"context"

	"github.com/google/wire"
)

var LicenseRepoSet = wire.NewSet(NewLicenseRepository, wire.Bind(new(ILicenseRepository), new(*LicenseRepository)))

type ILicenseRepository interface {
	Query(ctx context.Context, args *schema.GetLicenseArgs) (*model.License, error)
	Upload(ctx context.Context) error
}

type LicenseRepository struct {
	DB *database.DB
}

func NewLicenseRepository(db *database.DB) *LicenseRepository {
	return &LicenseRepository{DB: db}
}
func (a *LicenseRepository) Query(ctx context.Context, args *schema.GetLicenseArgs) (*model.License, error) {
	return nil, nil
}

func (a *LicenseRepository) Upload(ctx context.Context) error {
	return nil
}
