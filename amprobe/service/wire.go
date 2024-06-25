//go:build wireinject
// +build wireinject

package service

import (
	"github.com/amuluze/amprobe/service/audit"
	"github.com/amuluze/amprobe/service/auth"
	"github.com/amuluze/amprobe/service/container"
	"github.com/amuluze/amprobe/service/host"
	"github.com/amuluze/amprobe/service/model"
	"github.com/google/wire"
)

func BuildInjector(configFile string) (*Injector, func(), error) {
	wire.Build(
		NewConfig,
		NewLogger,
		NewDB,
		NewRPCClient,
		InitAuthStore,
		InitAuth,
		container.Set,
		host.Set,
		model.Set,
		auth.Set,
		audit.Set,
		NewLoggerHandler,
		NewTermHandler,
		RouterSet,
		NewFiberApp,
		//NewTimedTask,
		PrepareSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
