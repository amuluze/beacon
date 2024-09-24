// Package host
// Date: 2024/3/6 12:43
// Author: Amu
// Description:
package host

import (
	"github.com/google/wire"
	
	"amprobe/service/host/api"
	"amprobe/service/host/repository"
	"amprobe/service/host/service"
)

var Set = wire.NewSet(
	api.Set,
	service.Set,
	repository.Set,
)
