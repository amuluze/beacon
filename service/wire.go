//go:build wireinject
// +build wireinject

package service

import (
	"amprobe/service/account"
	"amprobe/service/alarm"
	"amprobe/service/audit"
	"amprobe/service/auth"
	"amprobe/service/container"
	"amprobe/service/host"
	"amprobe/service/mail"
	"amprobe/service/model"
	"amprobe/service/ws"

	"github.com/google/wire"
)

func BuildInjector(configFile string, modelFile ModelConfig) (*Injector, func(), error) {
	wire.Build(
		NewConfig,
		NewLogger,
		NewDB,
		InitAuthStore,
		InitAuth,
		InitAdapter,
		InitCasbin,
		container.Set,
		host.Set,
		model.Set,
		auth.Set,
		audit.Set,
		account.Set,
		alarm.Set,
		mail.Set,
		ws.NewLoggerHandler,
		ws.NewTermHandler,
		NewTimedTask,
		RouterSet,
		NewFiberApp,
		PrepareSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
