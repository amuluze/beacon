package dingtalk

import (
	"beacon/service/dingtalk/api"
	"beacon/service/dingtalk/client"
	"beacon/service/dingtalk/repository"
	"beacon/service/dingtalk/service"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	repository.Set,
	client.Set,
	service.Set,
	api.Set,
)
