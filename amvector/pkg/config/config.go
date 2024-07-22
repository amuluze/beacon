// Package config
// Date: 2024/7/2 14:38
// Author: Amu
// Description: amprobe/amvector 的配置文件 config.yml 的生成
package config

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/amuluze/amprobe/amvector/pkg/resources"
	"github.com/amuluze/amprobe/amvector/pkg/utils"

	"gopkg.in/yaml.v2"
)

func Create(prefix, filename string) (*Config, error) {
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	cfg := new(Config)
	cfg.filename = filename

	if err := cfg.loadDefault(prefix); err != nil {
		return nil, err
	}

	if data, err := io.ReadAll(fp); err != nil {
		return nil, err
	} else if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	cfg.loadVariables()
	return cfg, nil
}

type Config struct {
	filename string
	// 配置项
	Log struct {
		Output   string `yaml:"output"`
		Level    string `yaml:"level"`
		Rotation int    `yaml:"rotation"`
		MaxAge   int    `yaml:"max_age"`
	} `yaml:"log"`

	Task struct {
		Interval int `yaml:"interval"`
		Disk     struct {
			Devices []string `yaml:"devices"`
		} `yaml:"disk"`
		Ethernet struct {
			Names []string `yaml:"names"`
		} `yaml:"ethernet"`
	} `yaml:"task"`
	DB struct {
		DBType   string `yaml:"dbtype"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"db"`
	Profile   string `yaml:"profile"`
	Variables struct {
		ImageTag        string `yaml:"image_tag"`
		HostPrefix      string `yaml:"host_prefix"`
		ContainerPrefix string `yaml:"container_prefix"`
		// generated variables
		HostResourceDir      string `yaml:"-"`
		HostLogsDir          string `yaml:"-"`
		ContainerResourceDir string `yaml:"-"`
		ContainerLogsDir     string `yaml:"-"`
	} `yaml:"variables"`
}

func (c *Config) loadDefault(prefix string) error {
	c.Variables.ImageTag = "latest"
	c.Variables.HostPrefix = prefix
	c.Variables.ContainerPrefix = "/"
	c.loadVariables()

	c.Log.Output = filepath.Join(c.Variables.HostLogsDir, "vector.log")
	c.Log.Level = "info"
	c.Log.Rotation = 1
	c.Log.MaxAge = 7

	c.Task.Interval = 30
	c.Task.Disk.Devices = []string{"vda2"}
	c.Task.Ethernet.Names = []string{"eth0"}

	c.DB.DBType = "sqlite"
	c.DB.DBName = filepath.Join(c.Variables.HostResourceDir, resources.AmvectorStorageConfigDBPath)
	c.Profile = filepath.Join(c.Variables.HostResourceDir, resources.AmvectorServiceProfilePath)
	return nil
}

func (c *Config) loadVariables() {
	c.Variables.HostResourceDir = filepath.Join(c.Variables.HostPrefix, resources.RootPath)
	c.Variables.HostLogsDir = filepath.Join(c.Variables.HostPrefix, "logs")

	_ = utils.EnsureDirExists(c.Variables.HostResourceDir)
	_ = utils.EnsureDirExists(c.Variables.HostLogsDir)
	_ = utils.EnsureDirExists(filepath.Join(c.Variables.HostResourceDir, resources.AmvectorSockFolder))
	_ = utils.EnsureDirExists(filepath.Join(c.Variables.HostResourceDir, resources.AmvectorStorageFolder))

	c.Variables.ContainerResourceDir = filepath.Join(c.Variables.ContainerPrefix, resources.RootPath)
	c.Variables.ContainerLogsDir = filepath.Join(c.Variables.ContainerPrefix, "logs")
}

func (c *Config) Save() error {
	if c.filename == "" {
		return errors.New("empty filename")
	}
	out, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	if err := os.WriteFile(c.filename, out, 0644); err != nil {
		return err
	}
	return nil
}
