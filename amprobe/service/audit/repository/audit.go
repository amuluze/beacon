// Package repository
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package repository

import (
	"amprobe/service/model"
	"context"

	"amprobe/service/schema"

	"github.com/amuluze/amutool/database"
	"github.com/google/wire"
)

var AuditRepoSet = wire.NewSet(NewAuditRepo, wire.Bind(new(IAuditRepo), new(*AuditRepo)))

type IAuditRepo interface {
	AuditQuery(ctx context.Context, args *schema.AuditQueryArgs) (model.Audits, error)
	AuditCount(ctx context.Context) (int, error)
}

type AuditRepo struct {
	DB *database.DB
}

func NewAuditRepo(db *database.DB) *AuditRepo {
	return &AuditRepo{DB: db}
}

func (a *AuditRepo) AuditQuery(ctx context.Context, args *schema.AuditQueryArgs) (model.Audits, error) {
	var audits model.Audits
	if err := a.DB.Model(&model.Audit{}).Order("created_at DESC").Offset((args.Page - 1) * args.Size).Limit(args.Size).Find(&audits).Error; err != nil {
		return audits, err
	}
	return audits, nil
}

func (a *AuditRepo) AuditCount(ctx context.Context) (int, error) {
	var count int64
	if err := a.DB.Model(&model.Audit{}).Count(&count).Error; err != nil {
		return int(count), err
	}
	return int(count), nil
}
