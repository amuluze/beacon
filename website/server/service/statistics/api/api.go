// Package api
// Date:   2025/2/12 15:32
// Author: Amu
// Description:
package api

import "github.com/google/wire"

var Set = wire.NewSet(
	NewStatisticsAPI,
)
