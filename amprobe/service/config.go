// Package service
// Date: 2024/3/6 11:08
// Author: Amu
// Description:
package service

import (
	"log/slog"

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
}

// NewConfig Load config file (toml/json/yaml)
func NewConfig(configFile string) (*Config, error) {
	config := &Config{}

	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	warnInsecureDefaults(config)

	return config, nil
}

func warnInsecureDefaults(config *Config) {
	if config.Auth.Enable && config.Auth.SigningKey == "amprobe" {
		slog.Warn("auth signing key uses default development value")
	}
	if config.AgentInstall.Enable && config.AgentInstall.Token == "change-me" {
		slog.Warn("agent install token uses default development value")
	}
	if config.Control.Enable && config.Control.JoinToken == "" && config.AgentInstall.Token == "" {
		slog.Warn("control join token is empty")
	}
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
