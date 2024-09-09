// Package audit
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package audit

import (
	"amprobe/service/audit/api"

	"amprobe/service/audit/repository"

	"amprobe/service/audit/service"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	api.Set,
	service.Set,
	repository.Set,
)
