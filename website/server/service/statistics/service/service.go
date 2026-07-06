// Package service
// Date:   2025/2/12 15:29
// Author: Amu
// Description:
package service

import "github.com/google/wire"

var Set = wire.NewSet(
	StatisticsServiceSet,
)
