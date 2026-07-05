// Package container
// Date: 2024/3/6 12:44
// Author: Amu
// Description:
package container

import (
	"github.com/google/wire"

	"beacon/service/container/api"
	"beacon/service/container/repository"
	"beacon/service/container/service"
)

var Set = wire.NewSet(
	api.Set,
	service.Set,
	repository.Set,
)
