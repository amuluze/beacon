// Package api
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package api

import "github.com/google/wire"

var Set = wire.NewSet(
	NewAuditAPI,
)
