// Package mail
// Date:   2024/10/14 16:08
// Author: Amu
// Description:
package mail

import (
	"beacon/service/mail/api"
	"beacon/service/mail/repository"
	"beacon/service/mail/service"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	repository.Set,
	service.Set,
	api.Set,
)
