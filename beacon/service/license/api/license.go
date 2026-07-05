// Package api
// Date: 2023/6/7 12:59
// Author: Amu
// Description:
package api

import (
	"beacon/service/license/service"

	"github.com/gofiber/fiber/v2"
)

type LicenseAPI struct {
	LicenseService service.ILicenseService
}

func NewLicenseAPI(licenseService service.ILicenseService) *LicenseAPI {
	return &LicenseAPI{LicenseService: licenseService}
}

func (a *LicenseAPI) Query(c *fiber.Ctx) error {
	return nil
}

func (a *LicenseAPI) Upload(c *fiber.Ctx) error {
	return nil
}
