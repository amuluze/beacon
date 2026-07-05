// Package service
// Date: 2024/3/27 16:39
// Author: Amu
// Description:
package service

import "github.com/google/wire"

var Set = wire.NewSet(
	AuthServiceSet,
)
