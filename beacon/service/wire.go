//go:build wireinject
// +build wireinject

package service

import (
	"beacon/service/account"
	"beacon/service/agent"
	"beacon/service/alarm"
	"beacon/service/audit"
	"beacon/service/auth"
	"beacon/service/container"
	"beacon/service/host"
	"beacon/service/mail"
	"beacon/service/model"

	"github.com/google/wire"
)

func BuildInjector(configFile string, modelFile ModeConf) (*Injector, func(), error) {
	wire.Build(
		NewConfig,
		NewLogger,
		NewDB,
		NewRPCClient,
		NewRPCCaller,
		NewServerTunnelFromResult,
		NewStalenessMinutes,
		NewReportService,
		InitAuthStore,
		InitAuth,
		InitAdapter,
		InitCasbin,
		container.Set,
		host.Set,
		agent.Set,
		model.Set,
		auth.Set,
		audit.Set,
		account.Set,
		alarm.Set,
		mail.Set,
		NewLoggerHandler,
		NewTerminalHandler,
		NewTermHandler,
		NewTimedTask,
		RouterSet,
		NewFiberApp,
		PrepareSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
