// Package container
// Date: 2024/3/6 12:44
// Author: Amu
// Description:
package container

import (
	"github.com/google/wire"

	"github.com/amuluze/amprobe/service/container/api"
	"github.com/amuluze/amprobe/service/container/rpc"
)

var Set = wire.NewSet(
	api.Set,
	rpc.Set,
)
