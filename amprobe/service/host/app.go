// Package host
// Date: 2024/3/6 12:43
// Author: Amu
// Description:
package host

import (
	"github.com/google/wire"

	"amprobe/service/host/api"

	"amprobe/service/host/rpc"
)

var Set = wire.NewSet(
	api.Set,
	rpc.Set,
)
