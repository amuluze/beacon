//go:build wireinject
// +build wireinject

package service

import (
	"amprobe/service/audit"
	"amprobe/service/auth"
	"amprobe/service/container"
	"amprobe/service/host"
	"amprobe/service/model"

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
