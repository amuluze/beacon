// Package service
// Date: 2024/06/10 18:26:22
// Author: Amu
// Description:
package service

import (
	"log/slog"

	"github.com/spf13/viper"
)

type Prefix string

type Config struct {
	prefix    Prefix    `yaml:"-"`
	Log       Log       `yaml:"log"`
	Task      Task      `yaml:"task"`
	DB        DB        `yaml:"db"`
	Variables Variables `yaml:"variables"`
}

func NewConfig(configFile string, prefix Prefix) (*Config, error) {
	config := &Config{}

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		slog.Error("read config error", "err", err)
		return nil, err
	}

	if err := viper.Unmarshal(config); err != nil {
		slog.Error("parse config error", "error", err)
		return nil, err
	}
	config.prefix = prefix
	return config, nil
}

type Task struct {
	Interval int      `yaml:"interval"`
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

type DB struct {
	DBType   string `yaml:"dbtype"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type Variables struct {
	ImageTag        string `yaml:"image_tag"`
	HostPrefix      string `yaml:"host_prefix"`
	ContainerPrefix string `yaml:"container_prefix"`
}
