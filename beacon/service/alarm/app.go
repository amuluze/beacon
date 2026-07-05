// Package alarm
// Date:   2024/10/14 17:25
// Author: Amu
// Description:
package alarm

import (
	"beacon/service/alarm/api"
	"beacon/service/alarm/repository"
	"beacon/service/alarm/service"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	api.Set,
	service.Set,
	repository.Set,
)
