package feishu

import (
	"beacon/service/feishu/api"
	"beacon/service/feishu/client"
	"beacon/service/feishu/repository"
	"beacon/service/feishu/service"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	repository.Set,
	client.Set,
	service.Set,
	api.Set,
)
