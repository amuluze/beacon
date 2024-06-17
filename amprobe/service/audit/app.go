// Package audit
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package audit

import (
	"github.com/amuluze/amprobe/service/audit/api"
	"github.com/amuluze/amprobe/service/audit/repository"
	"github.com/amuluze/amprobe/service/audit/service"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	api.Set,
	service.Set,
	repository.Set,
)
