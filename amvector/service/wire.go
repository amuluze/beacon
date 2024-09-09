//go:build wireinject
// +build wireinject

package service

import (
	"amvector/service/model"

	"github.com/google/wire"
)

func BuildInjector(configFile string, prefix Prefix) (*Injector, func(), error) {
	wire.Build(
		NewConfig,
		NewLogger,
		NewDB,
		model.Set,
		NewRPCServer,
		NewTimedTask,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
