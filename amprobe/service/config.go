// Package service
// Date: 2024/3/6 11:08
// Author: Amu
// Description:
package service

import (
	"strings"

	hostservice "amprobe/service/host/service"

	"github.com/spf13/viper"
)

type Config struct {
	App          App
	Fiber        Fiber
	Control      Control
	Gorm         Gorm
	DB           DB
	Log          Log
	Auth         Auth
	Casbin       Casbin
	Task         Task
	Retention    Retention
	AgentInstall AgentInstall
	Session      Session
	CORS         CORS
	RateLimit    RateLimit
}

// NewConfig Load config file (toml/json/yaml)
func NewConfig(configFile string) (*Config, error) {
	config := &Config{}

	viper.SetConfigFile(configFile)

	// 允许通过环境变量覆盖配置（12-factor）。
	// 敏感字段显式 BindEnv，确保 Unmarshal 能拾取（viper AutomaticEnv 对 Unmarshal 的已知坑）。
	// 环境变量命名：前缀 AMPROBE_ + 结构体路径（. 替换为 _），如 AMPROBE_AUTH_SIGNINGKEY。
	viper.SetEnvPrefix("AMPROBE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.BindEnv("App.Env")
	viper.BindEnv("Auth.SigningKey")
	viper.BindEnv("DB.Password")
	viper.BindEnv("Control.JoinToken")
	viper.BindEnv("AgentInstall.Token")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

type Fiber struct {
	Host            string
	Port            int
	ShutdownTimeout int
	SeverHeader     string
	AppName         string
	Prefork         bool
}

// App 承载进程级运行配置。
// Env 标识运行模式：development（默认，宽松）或 production（严格，敏感字段缺失即拒绝启动）。
type App struct {
	Env string
}

// CORS 控制跨域资源共享策略。
// AllowOrigins 为空时回退到本地开发域；生产部署应显式配置白名单。
type CORS struct {
	Enable       bool
	AllowOrigins []string
}

// RateLimit 控制分层速率限制阈值（按 IP，每分钟）。
// GlobalMax 是全路由兜底上限；LoginMax 针对 login/token_update 等敏感端点防爆破。
// Max <= 0 时由中间件回退到内置默认值。
type RateLimit struct {
	Enable    bool
	GlobalMax int
	LoginMax  int
}

type Control struct {
	Enable         bool
	Address        string
	DefaultAgentID string
	TLSEnable      bool
	TLSCertDir     string
	JoinToken      string
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

type Log struct {
	Output   string
	Level    string
	Rotation int
	MaxAge   int
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

type Session struct {
	Enabled   bool
	Directory string
}

type Retention struct {
	Days      int
	Staleness int
}

// NewStalenessMinutes extracts the staleness threshold as hostservice.StalenessMinutes.
func NewStalenessMinutes(config *Config) hostservice.StalenessMinutes {
	v := config.Retention.Staleness
	if v <= 0 {
		v = 5
	}
	return hostservice.StalenessMinutes(v)
}
