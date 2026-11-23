// Package service
// Date: 2024/3/6 11:08
// Author: Amu
// Description:
package service

import (
	"log/slog"
	"os"
	"strings"

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
	AgentUpdate   AgentUpdate
	InstallReport InstallReport
	CORS          CORS
	RateLimit     RateLimit
	App           App
	Retention     Retention
	Update        Update
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

	// 部署期 env 覆盖，优先级高于 BindEnv。
	overrideFromEnv(config)
	warnInsecureDefaults(config)
	if err := validateTransportSecurity(config); err != nil {
		return nil, err
	}

	return config, nil
}

// bindEnvs 将部署期可变字段绑定到 BEACON_* 环境变量。
// viper.BindEnv 单次只支持一个 env key；历史多命名兼容由 overrideFromEnv 处理。
func bindEnvs() {
	bindings := map[string]string{
		"app.env":                    "BEACON_APP_ENV",
		"app.readonly":               "BEACON_READ_ONLY",
		"app.noauth":                 "BEACON_NO_AUTH",
		"db.dbname":                  "BEACON_DB_NAME",
		"auth.signingkey":            "BEACON_AUTH_SIGNING_KEY",
		"control.jointoken":          "BEACON_CONTROL_JOIN_TOKEN",
		"control.address":            "BEACON_CONTROL_ADDRESS",
		"control.unixsocket":         "BEACON_CONTROL_SOCKET",
		"control.tls.enable":         "BEACON_CONTROL_TLS_ENABLE",
		"agentinstall.token":         "BEACON_AGENT_INSTALL_TOKEN",
		"agentinstall.tokenttl":      "BEACON_AGENT_INSTALL_TOKEN_TTL",
		"agentinstall.publicbaseurl": "BEACON_PUBLIC_BASE_URL",
		"agentinstall.controlsocket": "BEACON_AGENT_INSTALL_CONTROL_SOCKET",
		"agentinstall.tlsenable":     "BEACON_AGENT_INSTALL_TLS_ENABLE",
		"agentinstall.controlport":   "BEACON_CONTROL_PORT",
		"agentupdate.enable":         "BEACON_AGENT_UPDATE_ENABLE",
		"agentupdate.targetversion":  "BEACON_COLLIA_VERSION",
		"agentupdate.publickey":      "BEACON_AGENT_UPDATE_PUBLIC_KEY",
		"agentupdate.timeout":        "BEACON_AGENT_UPDATE_TIMEOUT",
		"agentupdate.maxconcurrent":  "BEACON_AGENT_UPDATE_MAX_CONCURRENT",
		"agentupdate.retryinterval":  "BEACON_AGENT_UPDATE_RETRY_INTERVAL",
		"update.enable":              "BEACON_UPDATE_ENABLE",
		"update.url":                 "BEACON_UPDATE_URL",
		"update.checkinterval":       "BEACON_UPDATE_CHECK_INTERVAL",
		"retention.days":             "BEACON_RETENTION_DAYS",
		"retention.auditdays":        "BEACON_RETENTION_AUDIT_DAYS",
		"cors.enable":                "BEACON_CORS_ENABLE",
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
	if v := firstEnv("BEACON_APP_ENV"); v != "" {
		config.App.Env = v
	}
	if v := firstEnv("BEACON_AUTH_SIGNING_KEY", "BEACON_AUTH_SIGNINGKEY"); v != "" {
		config.Auth.SigningKey = v
	}
	if v := firstEnv("BEACON_AGENT_INSTALL_TOKEN", "BEACON_AGENT_INSTALLTOKEN"); v != "" {
		config.AgentInstall.Token = v
	}
	if v := firstEnv("BEACON_CONTROL_JOIN_TOKEN", "BEACON_CONTROL_JOINTOKEN"); v != "" {
		config.Control.JoinToken = v
	}
	if v := firstEnv("BEACON_CORS_ALLOW_ORIGINS"); v != "" {
		origins := strings.Split(v, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
		config.CORS.AllowOrigins = origins
	}
}

func warnInsecureDefaults(config *Config) {
	if config.Auth.Enable && isWeakDefaultSigningKey(config.Auth.SigningKey) {
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
	ServerHeader    string
	AppName         string
	Prefork         bool
	// ReadTimeout / WriteTimeout / IdleTimeout 控制 fasthttp 连接的读、写、空闲超时（秒）。
	// 缺失或 <=0 时 NewFiberApp 兜底为下方默认值，避免慢速/半开连接无限占用内存。
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

type Control struct {
	Enable  bool
	Address string
	// UnixSocket 非空时反向隧道在 Unix Socket 上监听（单机部署），
	// 此时忽略 Address 与 TLS 强制要求；路径为进程可见路径。
	UnixSocket string
	JoinToken  string
	TLS        ControlTLS
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

// Update 控制 beacon 对 beacon-help 的版本轮询。Enable=false 完全关闭；
// CheckInterval 为秒，最低 minCheckInterval 以防风暴。
type Update struct {
	Enable        bool
	URL           string
	CheckInterval int
}

type AgentInstall struct {
	Enable bool
	Token  string
	// TokenTTL 控制安装令牌有效期（秒）：0 未配置时默认 1800（短时令牌），
	// 负值显式回退为返回静态令牌的兼容模式。
	TokenTTL      int
	PublicBaseURL string
	PackageDir    string
	ControlPort   int
	// ControlSocket 非空时把 Agent 的下发 config 指向 `unix://<路径>`（宿主视角），
	// 与 Control.UnixSocket（容器视角）配合实现单机部署；空值保持 TCP 行为。
	ControlSocket string
	TLSEnable     bool
	CertDir       string
}

// AgentUpdate controls Server-coordinated Collia updates. Durations are in
// seconds; Enable is the deployment-level gate for automatic scheduling: when
// false no automatic update is scheduled at all. When true, the console switch
// (persisted in s_agent_setting, missing row means enabled) can still disable
// automatic scheduling. Manual retries use the same target package and safety
// limits and ignore the console switch.
type AgentUpdate struct {
	Enable        bool
	TargetVersion string
	// PublicKey 是 Collia 发布签名的 Ed25519 公钥（base64）。配置后随安装
	// config 下发给 Agent 固定；侧车签名缺失时 Agent 拒绝升级。
	PublicKey     string
	Timeout       int
	MaxConcurrent int
	RetryInterval int
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
	Env      string
	ReadOnly bool
	NoAuth   bool
}

type Retention struct {
	// Days 是监控时间序列数据的保留天数；AuditDays 是审计日志的保留天数。
	Days      int
	AuditDays int
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
