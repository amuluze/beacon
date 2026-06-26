//go:build wireinject
// +build wireinject

package service

import (
	"collia/service/model"

	"github.com/google/wire"
)

func BuildInjector(configFile string, prefix Prefix, version string) (*Injector, func(), error) {
	wire.Build(
		NewConfig,
		NewLogger,
		NewDB,
		model.Set,
		NewVersion,
		NewRPCServer,
		NewTimedTask,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
