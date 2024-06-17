// Package service
// Date: 2024/3/6 11:08
// Author: Amu
// Description:
package service

import (
	"github.com/spf13/viper"
)

type Config struct {
	Fiber    Fiber
	Rpc      Rpc
	Gorm     Gorm
	DB       DB
	Disk     Disk
	Task     Task
	Ethernet Ethernet
	Logger   Logger
	Auth     Auth
	InitData InitData
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

type Rpc struct {
	Address string
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

type Disk struct {
	Devices []string
}

type Task struct {
	Interval int
}

type Ethernet struct {
	Names []string
}

type Logger struct {
	File         string
	Level        string
	RotationTime int
	MaxAge       int
}

type Auth struct {
	Enable         bool
	SigningMethod  string
	SigningKey     string
	Expired        int
	RefreshExpired int
	Prefix         string
}

type InitData struct {
	Enable         bool
	InitConfigFile string
}
