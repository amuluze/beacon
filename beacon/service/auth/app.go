// Package auth
// Date: 2024/3/27 16:38
// Author: Amu
// Description:
package auth

import (
	"beacon/service/auth/api"
	"beacon/service/auth/repository"
	"beacon/service/auth/service"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	api.Set,
	service.Set,
	repository.Set,
)
