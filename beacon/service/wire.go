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
	"beacon/service/dingtalk"
	"beacon/service/feishu"
	"beacon/service/host"
	"beacon/service/mail"
	"beacon/service/model"
	systemapi "beacon/service/system/api"
	"beacon/service/terminal"
	"beacon/service/wecom"

	"github.com/google/wire"
)

// NewStalenessMinutes provides the staleness threshold for host data freshness checks.
// Default is 300 seconds (5 minutes); can be configured via retention.
func NewStalenessMinutes() []int64 { return nil }

func BuildInjector(configFile string, modelFile ModeConf) (*Injector, func(), error) {
	wire.Build(
		NewStalenessMinutes,
		NewHealthProbe,
		NewConfig,
		InitAuthCookieManager,
		NewDataCipher,
		NewLogger,
		NewDB,
		NewRPCClient,
		NewReportService,
		InitAuthStore,
		InitAuth,
		InitAdapter,
		InitCasbin,
		container.Set,
		dingtalk.Set,
		feishu.Set,
		wecom.Set,
		host.Set,
		agent.Set,
		model.Set,
		auth.Set,
		audit.Set,
		account.Set,
		alarm.Set,
		mail.Set,
		NewLoggerHandler,
		terminal.NewHandler,
		NewTimedTask,
		ProvideSystemVersion,
		systemapi.NewSystemAPI,
		RouterSet,
		NewFiberApp,
		PrepareSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
