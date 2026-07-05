// Package api
// Date       : 2024/9/4 14:55
// Author     : Amu
// Description:
package api

import "github.com/google/wire"

var Set = wire.NewSet(
	NewAccountAPI,
)
