// Package api
// Date: 2024/3/6 11:21
// Author: Amu
// Description:
package api

import (
	"beacon/pkg/fiberx"
	"beacon/service/schema"

	"github.com/gofiber/fiber/v2"
)

// SystemAPI exposes read-only information about the running Beacon Server.
type SystemAPI struct {
	version func() schema.SystemVersion
}

// NewSystemAPI builds the system information API around a version provider.
func NewSystemAPI(version func() schema.SystemVersion) *SystemAPI {
	return &SystemAPI{version: version}
}

// Version returns the current build version and the latest update-check result.
func (a *SystemAPI) Version(c *fiber.Ctx) error {
	current := a.version()
	return fiberx.Success(c, schema.SystemVersion{
		Version:         current.Version,
		UpdateAvailable: current.UpdateAvailable,
		LatestVersion:   current.LatestVersion,
		ReleaseNotes:    current.ReleaseNotes,
		CheckedAt:       current.CheckedAt,
	})
}
