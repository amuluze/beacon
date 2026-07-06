// Package service
// Date       : 2024/9/4 14:56
// Author     : Amu
// Description:
package service

import "github.com/google/wire"

var Set = wire.NewSet(
	AccountServiceSet,
)
