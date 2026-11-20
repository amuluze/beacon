// Package service
// Date: 2024/06/11 19:29:11
// Author: Amu
// Description:
package service

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"beacon/pkg/rpc"
	"beacon/service/agent"
	healthapi "beacon/service/health/api"
	"common/database"
	tunnelpkg "common/rpc/tunnel"
	transporttls "common/transport/tlsconfig"
)

// globalTunnel guards the singleton tunnel instance.
// Set once during initialization; read safely by health probes and lifecycle hooks.
var globalTunnel *tunnelpkg.ServerTunnel

// NewRPCClient creates the tunnel-based RPC client.
// Server listens for reverse connections from Agents.
func NewRPCClient(config *Config, agentService *agent.Service, certManager *CertManager, db *database.DB) (rpc.Caller, error) {
	if !config.Control.Enable {
		globalTunnel = nil
		if agentService != nil {
			agentService.SetUpdater(nil)
			agentService.SetTunnel(nil)
		}
		slog.Info("reverse tunnel server disabled")
		return rpc.NewDisabledClient(), nil
	}
	if strings.EqualFold(config.App.Env, productionEnv) && !config.Control.TLS.Enable && config.Control.UnixSocket == "" {
		return nil, ErrControlTLSRequired
	}

	joinToken, err := resolveControlToken(controlJoinToken(config), config.App.Env)
	if err != nil {
		return nil, err
	}

	slog.Info("reverse tunnel server prepared", "addr", controlListenAddress(config))

	if config.Control.TLS.Enable {
		if err := certManager.EnsureControlCerts(); err != nil {
			return nil, fmt.Errorf("ensure control certs: %w", err)
		}
	}

	tun := tunnelpkg.NewServerTunnel()
	tun.SetJoinToken(joinToken)
	if joinToken == "" {
		// 非生产且未配置 join token 时显式允许匿名注册以保留本地开发便利；
		// 生产模式 resolveControlToken 已拒绝空 token，不会走到这里。
		tun.SetAllowAnonymous(true)
	}
	if config.Control.TLS.Enable {
		tlsCfg, err := transporttls.ServerConfig(config.Control.TLS.CertDir, config.Control.TLS.ClientNames)
		if err != nil {
			return nil, err
		}
		tun.SetTLSConfig(tlsCfg)
		slog.Info("reverse tunnel tls enabled", "cert_dir", config.Control.TLS.CertDir)
	}
	client := rpc.NewTunnelClient(tun)
	updater, err := agent.NewUpdater(agent.UpdateConfig{
		Enable:        config.AgentUpdate.Enable,
		TargetVersion: config.AgentUpdate.TargetVersion,
		Timeout:       time.Duration(config.AgentUpdate.Timeout) * time.Second,
		MaxConcurrent: config.AgentUpdate.MaxConcurrent,
		RetryInterval: time.Duration(config.AgentUpdate.RetryInterval) * time.Second,
		PublicBaseURL: config.AgentInstall.PublicBaseURL,
		InstallToken:  config.AgentInstall.Token,
		PackageDir:    config.AgentInstall.PackageDir,
	}, client, db)
	if err != nil {
		return nil, fmt.Errorf("configure collia updater: %w", err)
	}
	agentService.SetUpdater(updater)
	agentService.SetTunnel(tun)
	globalTunnel = tun
	// 监听在 NewInjector 挂载数据处理器之后启动（见 injector.go），
	// 消除「连接已接受但 Store 未就绪」的启动窗口静默丢批。
	return client, nil
}

func controlJoinToken(config *Config) string {
	if config.Control.JoinToken != "" {
		return config.Control.JoinToken
	}
	return config.AgentInstall.Token
}

// ServerTunnelFromHolder returns the singleton tunnel instance held by globalTunnel.
// It exposes the live tunnel so other components (e.g. health.Probe) can wire
// into it without taking on the lifetime ownership.
func ServerTunnelFromHolder() *tunnelpkg.ServerTunnel {
	return globalTunnel
}

// controlListenAddress 返回控制通道监听目标：Unix Socket 模式优先，
// 否则回退 TCP 地址（空配置默认 17000）。
func controlListenAddress(config *Config) string {
	if socket := strings.TrimSpace(config.Control.UnixSocket); socket != "" {
		return "unix://" + socket
	}
	if addr := config.Control.Address; addr != "" {
		return addr
	}
	return ":17000"
}

// dbPinger 适配 *database.DB 到健康探针的 DBPinger 接口。
type dbPinger struct{ db *database.DB }

func (p dbPinger) Ping() error {
	sqlDB, err := p.db.DB.DB()
	if err != nil {
		return fmt.Errorf("open sql connection: %w", err)
	}
	return sqlDB.Ping()
}

// NewHealthProbe 构造真实接线的健康探针：数据库失联时 /ready 必须翻转为
// 503，tunnel 存在时参与就绪判定。caller 参数用于让 Wire 在 NewRPCClient
// （发布全局 tunnel）之后才构造探针，本身不参与实现。
func NewHealthProbe(caller rpc.Caller, db *database.DB) *healthapi.Probe {
	_ = caller
	p := healthapi.NewProbe()
	if db != nil {
		p.SetDB(dbPinger{db: db})
	}
	if tun := ServerTunnelFromHolder(); tun != nil {
		p.SetTunnel(tun)
	}
	return p
}
