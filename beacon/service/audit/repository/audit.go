// Package repository
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package repository

import (
	"beacon/service/model"
	"context"

	"beacon/service/schema"

	"common/database"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var AuditRepoSet = wire.NewSet(NewAuditRepo, wire.Bind(new(IAuditRepo), new(*AuditRepo)))

var _ IAuditRepo = (*AuditRepo)(nil)

type IAuditRepo interface {
	AuditQuery(ctx context.Context, args schema.AuditQueryArgs) (model.Audits, error)
	AuditCount(ctx context.Context, args schema.AuditQueryArgs) (int, error)
}

type AuditRepo struct {
	DB *database.DB
}

func NewAuditRepo(db *database.DB) *AuditRepo {
	return &AuditRepo{DB: db}
}

func (a *AuditRepo) filteredQuery(ctx context.Context, args schema.AuditQueryArgs) *gorm.DB {
	q := a.DB.WithContext(ctx).Model(&model.Audit{})
	if args.Type == "system" {
		q = q.Where("username = ?", "system")
	} else {
		q = q.Where("username != ?", "system")
	}
	if args.AgentID != "" {
		q = q.Where("agent_id = ?", args.AgentID)
	}
	return q
}

func (a *AuditRepo) AuditQuery(ctx context.Context, args schema.AuditQueryArgs) (model.Audits, error) {
	var audits model.Audits
	if err := a.filteredQuery(ctx, args).
		Order("created_at DESC").
		Offset((args.Page - 1) * args.Size).
		Limit(args.Size).
		Find(&audits).Error; err != nil {
		return nil, err
	}
	return audits, nil
}

func (a *AuditRepo) AuditCount(ctx context.Context, args schema.AuditQueryArgs) (int, error) {
	var count int64
	if err := a.filteredQuery(ctx, args).Count(&count).Error; err != nil {
		return int(count), err
	}
	return int(count), nil
}
