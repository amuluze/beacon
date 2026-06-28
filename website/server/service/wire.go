//go:build wireinject
// +build wireinject

package service

import (
	"server/service/model"
	"server/service/statistics"

	"github.com/google/wire"
)

func BuildInjector(configFile string) (*Injector, func(), error) {
	wire.Build(
		NewConfig,
		NewLogger,
		NewDB,
		model.Set,
		statistics.Set,
		RouterSet,
		NewFiberApp,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
