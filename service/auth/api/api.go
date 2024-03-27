// Package api
// Date: 2024/3/27 16:38
// Author: Amu
// Description:
package api

import "github.com/google/wire"

var Set = wire.NewSet(
	NewLoginAPI,
)
