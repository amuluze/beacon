// Package api
// Date:   2024/10/14 17:31
// Author: Amu
// Description:
package api

import "github.com/google/wire"

var Set = wire.NewSet(
	NewAlarmAPI,
)
