// Package service
// Date: 2024/3/6 11:08
// Author: Amu
// Description:
package service

import (
	"github.com/spf13/viper"
)

type Config struct {
	Fiber  Fiber
	Gorm   Gorm
	DB     DB
	Log    Log
	Auth   Auth
	Casbin Casbin
	Task   Task
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
	Interval int      `yaml:"interval"`
	MaxAge   int      `yaml:"max_age"`
	Disk     Disk     `yaml:"disk"`
	Ethernet Ethernet `yaml:"ethernet"`
}

type Disk struct {
	Devices []string `yaml:"devices"`
}

type Ethernet struct {
	Names []string `yaml:"names"`
}

type Log struct {
	Output   string `yaml:"output"`
	Level    string `yaml:"level"`
	Rotation int    `yaml:"rotation"`
	MaxAge   int    `yaml:"max_age"`
}

type Auth struct {
	Enable         bool   `yaml:"enable"`
	SigningMethod  string `yaml:"signing_method"`
	SigningKey     string `yaml:"signing_key"`
	Expired        int    `yaml:"expired"`
	RefreshExpired int    `yaml:"refresh_expired"`
	Prefix         string `json:"prefix"`
}

type Casbin struct {
	Enable           bool `yaml:"enable"`
	Debug            bool `yaml:"debug"`
	AutoLoad         bool `yaml:"auto_load"`
	AutoLoadInternal int  `yaml:"auto_load_internal"`
}
