// Package container
// Date: 2024/3/6 12:44
// Author: Amu
// Description:
package container

import (
	"github.com/google/wire"

	"amprobe/service/container/api"

	"amprobe/service/container/rpc"
)

var Set = wire.NewSet(
	api.Set,
	rpc.Set,
)
