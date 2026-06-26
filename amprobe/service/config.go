// Package service
// Date: 2024/3/6 11:08
// Author: Amu
// Description:
package service

import (
	"github.com/spf13/viper"
)

type Config struct {
	Fiber        Fiber
	Control      Control
	Gorm         Gorm
	DB           DB
	Log          Log
	Auth         Auth
	Casbin       Casbin
	Task         Task
	AgentInstall AgentInstall
	Session      Session
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

type Control struct {
	Enable         bool
	Address        string
	DefaultAgentID string
	TLSEnable      bool
	TLSCertDir     string
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
