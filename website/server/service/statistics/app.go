// Package statistics
// Date:   2025/2/12 15:32
// Author: Amu
// Description:
package statistics

import (
	"github.com/google/wire"
	"server/service/statistics/api"
	"server/service/statistics/repository"
	"server/service/statistics/service"
)

var Set = wire.NewSet(
	repository.Set,
	service.Set,
	api.Set,
)
