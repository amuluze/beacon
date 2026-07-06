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
	Control   Control   `yaml:"control"`
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
	Interval int      `yaml:"interval" mapstructure:"interval"`
	MaxAge   int      `yaml:"max_age" mapstructure:"max_age"`
	Disk     Disk     `yaml:"disk" mapstructure:"disk"`
	Ethernet Ethernet `yaml:"ethernet" mapstructure:"ethernet"`
	Report   Report   `yaml:"report" mapstructure:"report"`
}

type Report struct {
	URL     string `yaml:"url" mapstructure:"url"`
	Token   string `yaml:"token" mapstructure:"token"`
	AgentID string `yaml:"agent_id" mapstructure:"agent_id"`
}

type Control struct {
	Server    string `yaml:"server" mapstructure:"server"`
	AgentID   string `yaml:"agent_id" mapstructure:"agent_id"`
	JoinToken string `yaml:"join_token" mapstructure:"join_token"`
	TLS       TLS    `yaml:"tls" mapstructure:"tls"`
}

type TLS struct {
	Enable      bool     `yaml:"enable" mapstructure:"enable"`
	CertDir     string   `yaml:"cert_dir" mapstructure:"cert_dir"`
	ServerName  string   `yaml:"server_name" mapstructure:"server_name"`
	ClientNames []string `yaml:"client_names" mapstructure:"client_names"`
}

type Disk struct {
	Devices []string `yaml:"devices" mapstructure:"devices"`
}

type Ethernet struct {
	Names []string `yaml:"names" mapstructure:"names"`
}

type Log struct {
	Output   string `yaml:"output" mapstructure:"output"`
	Level    string `yaml:"level" mapstructure:"level"`
	Rotation int    `yaml:"rotation" mapstructure:"rotation"`
	MaxAge   int    `yaml:"max_age" mapstructure:"max_age"`
}

type DB struct {
	DBType   string `yaml:"dbtype" mapstructure:"dbtype"`
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
	DBName   string `yaml:"dbname" mapstructure:"dbname"`
	SSLMode  string `yaml:"sslmode" mapstructure:"sslmode"`
}

type Variables struct {
	ImageTag        string `yaml:"image_tag" mapstructure:"image_tag"`
	HostPrefix      string `yaml:"host_prefix" mapstructure:"host_prefix"`
	ContainerPrefix string `yaml:"container_prefix" mapstructure:"container_prefix"`
}
