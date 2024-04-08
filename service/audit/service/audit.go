// Package service
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package service

import (
	"context"
	"fmt"
	"github.com/amuluze/amprobe/service/audit/repository"
	"github.com/amuluze/amprobe/service/schema"
	"github.com/amuluze/amutool/errors"
	"github.com/google/wire"
)

var AuditServiceSet = wire.NewSet(NewAuditService, wire.Bind(new(IAuditService), new(*AuditService)))

type IAuditService interface {
	AuditQuery(ctx context.Context, args *schema.AuditQueryArgs) (*schema.AuditQueryReply, error)
}

type AuditService struct {
	AuditRepo repository.IAuditRepo
}

func NewAuditService(repo repository.IAuditRepo) *AuditService {
	return &AuditService{AuditRepo: repo}
}

func (a AuditService) AuditQuery(ctx context.Context, args *schema.AuditQueryArgs) (*schema.AuditQueryReply, error) {
	audits, err := a.AuditRepo.AuditQuery(ctx, args)
	if err != nil {
		return &schema.AuditQueryReply{}, errors.New400Error(err.Error())
	}
	fmt.Println("audits:", audits)

	var list []schema.Audit
	for _, audit := range audits {
		list = append(list, schema.Audit{
			ID:       audit.ID,
			Username: audit.Username,
			Operate:  audit.Operate,
			Created:  audit.CreatedAt,
		})
	}
	total, _ := a.AuditRepo.AuditCount(ctx)
	return &schema.AuditQueryReply{Data: list, Total: total, Page: args.Page, Size: args.Size}, nil
}
