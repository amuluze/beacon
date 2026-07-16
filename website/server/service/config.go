// Package service
// Date: 2025/02/12 15:03:22
// Author: Amu
// Description:
package service

import "github.com/spf13/viper"

type Config struct {
	App   App
	Fiber Fiber
	Gorm  Gorm
	DB    DB
	Log   Log
}

// NewConfig Load config file (toml/json/yaml)
func NewConfig(configFile string) (*Config, error) {
	config := &Config{}

	viper.SetConfigFile(configFile)

	// 允许通过环境变量覆盖运行环境，便于容器化部署
	viper.SetEnvPrefix("APP")
	viper.BindEnv("App.Env", "APP_ENV")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	if config.App.Env == "" {
		config.App.Env = "production"
	}

	return config, nil
}

// App 全局应用配置
type App struct {
	// Env 运行环境：dev / production；production 下关闭 pprof 等调试面
	Env string
	// CORSAllowOrigins 允许的跨域来源；为空时回退为同源不开放跨域
	CORSAllowOrigins []string
}

// IsProduction 是否为生产环境
func (a *App) IsProduction() bool {
	return a.Env == "production"
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

type Log struct {
	Output   string
	Level    string
	Rotation int
	MaxAge   int
}
