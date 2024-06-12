// Package auth
// Date: 2024/3/27 16:38
// Author: Amu
// Description:
package auth

import (
	"github.com/amuluze/amprobe/service/auth/api"
	"github.com/amuluze/amprobe/service/auth/repository"
	"github.com/amuluze/amprobe/service/auth/service"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	api.Set,
	service.Set,
	repository.Set,
)
