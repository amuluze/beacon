// Package license
// Date: 2023/6/7 12:56
// Author: Amu
// Description:
package license

import (
	"beacon/service/license/api"
	"beacon/service/license/repository"
	"beacon/service/license/service"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	api.Set,
	service.Set,
	repository.Set,
)
