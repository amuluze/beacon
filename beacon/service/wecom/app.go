package wecom

import (
	"beacon/service/wecom/api"
	"beacon/service/wecom/client"
	"beacon/service/wecom/repository"
	"beacon/service/wecom/service"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	repository.Set,
	client.Set,
	service.Set,
	api.Set,
)
