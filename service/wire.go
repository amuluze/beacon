//go:build wireinject
// +build wireinject

package service

import (
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
		container.Set,
		host.Set,
		model.Set,
		NewLoggerHandler,
		RouterSet,
		NewFiberApp,
		NewTimedTask,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
