// Package api
// Date: 2024/3/6 12:44
// Author: Amu
// Description:
package api

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewContainerAPI,
)
