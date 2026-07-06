// Package host
// Date: 2024/3/6 12:43
// Author: Amu
// Description:
package host

import (
	"github.com/google/wire"

	"beacon/service/host/api"
	"beacon/service/host/repository"
	"beacon/service/host/service"
)

var Set = wire.NewSet(
	api.Set,
	service.Set,
	repository.Set,
)
