// Package rpc
// Date: 2024/06/11 19:27:16
// Author: Amu
// Description:
package rpc

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	HostServiceSet,
)
