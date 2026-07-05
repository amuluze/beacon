// Package greeter_service
// Date: 2023/6/7 12:59
// Author: Amu
// Description:
package service

import (
	"beacon/service/license/repository"
	"beacon/service/schema"
	"context"

	"github.com/google/wire"
)

var LicenseServiceSet = wire.NewSet(NewLicenseService, wire.Bind(new(ILicenseService), new(*LicenseService)))

type ILicenseService interface {
	Query(ctx context.Context, args *schema.GetLicenseArgs) (*schema.GetLicenseReply, error)
	Upload(ctx context.Context) error
}

type LicenseService struct {
	LicenseRepo repository.ILicenseRepository
}

func NewLicenseService(licenseRepo repository.ILicenseRepository) *LicenseService {
	return &LicenseService{LicenseRepo: licenseRepo}
}

func (a *LicenseService) Query(ctx context.Context, args *schema.GetLicenseArgs) (*schema.GetLicenseReply, error) {
	return nil, nil
}

func (a *LicenseService) Upload(ctx context.Context) error {
	return nil
}
