// Package service
// Date: 2024/06/10 18:26:22
// Author: Amu
// Description:
package service

import (
	"log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	Rpc      Rpc
	Gorm     Gorm
	DB       DB
	Disk     Disk
	Task     Task
	Ethernet Ethernet
	Logger   Logger
}

func NewConfig(configFile string) (*Config, error) {
	config := &Config{}

	viper.SetConfigFile(configFile)
	slog.Info("config file", "info", configFile)
	if err := viper.ReadInConfig(); err != nil {
		slog.Error("read config error", "err", err)
		return nil, err
	}

	if err := viper.Unmarshal(config); err != nil {
		slog.Error("parse config error", "error", err)
		return nil, err
	}

	return config, nil
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
