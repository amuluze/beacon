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
	AuditCount(ctx context.Context) (int, error)
}

type AuditRepo struct {
	DB *database.DB
}

func NewAuditRepo(db *database.DB) *AuditRepo {
	return &AuditRepo{DB: db}
}

func (a *AuditRepo) AuditQuery(ctx context.Context, args schema.AuditQueryArgs) (model.Audits, error) {
	var audits model.Audits

	// Build the common conditions once and reuse across Type branches.
	// The AgentID filter is optional; only applied when non-empty.
	applyFilters := func(q *gorm.DB) *gorm.DB {
		if args.AgentID != "" {
			q = q.Where("agent_id = ?", args.AgentID)
		}
		return q
	}

	if args.Type == "system" {
		q := a.DB.Model(&model.Audit{}).Where("username = ?", "system")
		q = applyFilters(q)
		if err := q.Order("created_at DESC").
			Offset((args.Page - 1) * args.Size).
			Limit(args.Size).
			Find(&audits).Error; err != nil {
			return audits, nil
		}
	} else {
		q := a.DB.Model(&model.Audit{}).Where("username != ?", "system")
		q = applyFilters(q)
		if err := q.Order("created_at DESC").
			Offset((args.Page - 1) * args.Size).
			Limit(args.Size).
			Find(&audits).Error; err != nil {
			return audits, nil
		}
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
