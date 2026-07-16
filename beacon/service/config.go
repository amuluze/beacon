// Package service
// Date: 2024/3/6 11:08
// Author: Amu
// Description:
package service

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Fiber         Fiber
	Control       Control
	Gorm          Gorm
	DB            DB
	Log           Log
	Auth          Auth
	Casbin        Casbin
	Task          Task
	AgentInstall  AgentInstall
	InstallReport InstallReport
	CORS          CORS
	RateLimit     RateLimit
	App           App
	Retention     Retention
}

// NewConfig Load config file (toml/json/yaml)
func NewConfig(configFile string) (*Config, error) {
	config := &Config{}

	viper.SetConfigFile(configFile)

	// 注册部署期 env 覆盖，替代容器启动时 sed 改写 config.toml。
	// 必须在 ReadInConfig/Unmarshal 之前注册，Unmarshal 才会拾取 env 值。
	bindEnvs()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	// 历史环境变量别名兼容（AMPROBE_* 等多命名），优先级高于 BindEnv。
	overrideFromEnv(config)
	warnInsecureDefaults(config)

	return config, nil
}

// bindEnvs 将部署期可变字段绑定到 BEACON_* 环境变量。
// viper.BindEnv 单次只支持一个 env key；历史多命名兼容由 overrideFromEnv 处理。
func bindEnvs() {
	bindings := map[string]string{
		"db.dbname":                  "BEACON_DB_NAME",
		"auth.signingkey":            "BEACON_AUTH_SIGNING_KEY",
		"control.jointoken":          "BEACON_CONTROL_JOIN_TOKEN",
		"control.address":            "BEACON_CONTROL_ADDRESS",
		"agentinstall.token":         "BEACON_AGENT_INSTALL_TOKEN",
		"agentinstall.publicbaseurl": "BEACON_PUBLIC_BASE_URL",
		"agentinstall.controlport":   "BEACON_CONTROL_PORT",
	}
	for key, env := range bindings {
		if err := viper.BindEnv(key, env); err != nil {
			slog.Warn("failed to bind env override", "key", key, "env", env, "err", err)
		}
	}
}

// overrideFromEnv overrides sensitive config values from environment variables
// to prevent hard-coded secrets in production deployments.
func overrideFromEnv(config *Config) {
	if v := firstEnv("BEACON_AUTH_SIGNING_KEY", "BEACON_AUTH_SIGNINGKEY", "AMPROBE_AUTH_SIGNING_KEY", "AMPROBE_AUTH_SIGNINGKEY"); v != "" {
		config.Auth.SigningKey = v
	}
	if v := firstEnv("BEACON_AGENT_INSTALL_TOKEN", "BEACON_AGENT_INSTALLTOKEN", "AMPROBE_AGENT_INSTALL_TOKEN", "AMPROBE_AGENT_INSTALLTOKEN"); v != "" {
		config.AgentInstall.Token = v
	}
	if v := firstEnv("BEACON_CONTROL_JOIN_TOKEN", "BEACON_CONTROL_JOINTOKEN", "AMPROBE_CONTROL_JOIN_TOKEN", "AMPROBE_CONTROL_JOINTOKEN"); v != "" {
		config.Control.JoinToken = v
	}
}

func warnInsecureDefaults(config *Config) {
	if config.Auth.Enable && config.Auth.SigningKey == "beacon" {
		slog.Warn("auth signing key uses default development value; set BEACON_AUTH_SIGNING_KEY to override")
	}
	if config.AgentInstall.Enable && config.AgentInstall.Token == "change-me" {
		slog.Warn("agent install token uses default development value; set BEACON_AGENT_INSTALL_TOKEN to override")
	}
	if config.Control.Enable && config.Control.JoinToken == "" && config.AgentInstall.Token == "" {
		slog.Warn("control join token is empty; set BEACON_CONTROL_JOIN_TOKEN or AgentInstall.Token")
	}
}

func firstEnv(names ...string) string {
	for _, name := range names {
		if v := os.Getenv(name); v != "" {
			return v
		}
	}
	return ""
}

type Fiber struct {
	Host            string
	Port            int
	ShutdownTimeout int
	SeverHeader     string
	AppName         string
	Prefork         bool
}

type Control struct {
	Enable         bool
	Address        string
	// Deprecated: DefaultAgentID is no longer used. Agent selection must be explicit
	// via X-Agent-ID header or agent_id query parameter in every request.
	DefaultAgentID string
	JoinToken      string
	TLS            ControlTLS
}

type ControlTLS struct {
	Enable      bool
	CertDir     string
	ClientNames []string
}

type Gorm struct {
	GenDoc            bool
	Debug             bool
	DBType            string
	MaxLifetime       int
	MaxOpenConns      int
	MaxIdleConns      int
	TablePrefix       string
	EnableAutoMigrate bool
}

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Task struct {
	Interval int
}

type AgentInstall struct {
	Enable        bool
	Token         string
	PublicBaseURL string
	PackageDir    string
	ControlPort   int
	TLSEnable     bool
	CertDir       string
}

type InstallReport struct {
	Enable     bool
	URL        string
	InstallDir string
	IDFile     string
	Timeout    int
}

type Log struct {
	Output   string
	Level    string
	Rotation int
	MaxAge   int
}

type CORS struct {
	Enable       bool
	AllowOrigins []string
}

type RateLimit struct {
	Enable    bool
	GlobalMax int
}

type App struct {
	Env string
}

type Retention struct {
	Days int
}

type Auth struct {
	Enable         bool
	SigningMethod  string
	SigningKey     string
	Expired        int
	RefreshExpired int
	Prefix         string
}

type Casbin struct {
	Enable           bool
	Debug            bool
	AutoLoad         bool
	AutoLoadInternal int
}
