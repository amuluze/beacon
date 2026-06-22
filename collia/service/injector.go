// Package service
// Date: 2024/06/10 18:20:47
// Author: Amu
// Description:
package service

import (
	"github.com/amuluze/amutool/logger"
	"github.com/google/wire"
)

var InjectorSet = wire.NewSet(NewInjector)

type Injector struct {
	Config    *Config
	Logger    *logger.Logger
	Task      *TimedTask
	RPCServer *Server
}

func NewInjector(config *Config, task *TimedTask, rpcServer *Server, logx *logger.Logger) (*Injector, error) {
	return &Injector{
		Config:    config,
		Task:      task,
		RPCServer: rpcServer,
		Logger:    logx,
	}, nil
}
