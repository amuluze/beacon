// Package account
// Date       : 2024/9/4 14:55
// Author     : Amu
// Description:
package account

import (
	"amprobe/service/account/api"
	"amprobe/service/account/repository"
	"amprobe/service/account/service"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	api.Set,
	service.Set,
	repository.Set,
)
