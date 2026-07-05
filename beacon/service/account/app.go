// Package account
// Date       : 2024/9/4 14:55
// Author     : Amu
// Description:
package account

import (
	"beacon/service/account/api"
	"beacon/service/account/repository"
	"beacon/service/account/service"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	api.Set,
	service.Set,
	repository.Set,
)
