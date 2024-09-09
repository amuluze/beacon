// Package compose
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package compose

import (
	"os"

	"github.com/mcuadros/go-defaults"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Subnet  string `yaml:"subnet" default:"172.18.0.0/24"`
	Gateway string `yaml:"gateway" default:"172.18.0.1"`
}

type DockerCompose struct {
	Services struct {
		Amprobe struct {
			Image         string `yaml:"image" default:"amuluze/amprobe:latest"`
			ContainerName string `yaml:"container_name" default:"amprobe"`
			Restart       string `yaml:"restart" default:"always"`
			Logging       struct {
				Options struct {
					MaxSize string `yaml:"max-size" default:"10m"`
					MaxFile string `yaml:"max-file" default:"10"`
				} `yaml:"options"`
			} `yaml:"logging"`
			Ports       []string `yaml:"ports" default:"[80:80,443:443]"`
			Environment struct {
				CreatedByProbe string `yaml:"CREATED_BY_PROBE" default:"true"`
			} `yaml:"environment"`
			Volumes  []string `yaml:"volumes" default:"[/var/run/docker.sock:/var/run/docker.sock,/tmp:/tmp]"`
			Networks struct {
				AppNet struct {
					Ipv4Address string `yaml:"ipv4_address" default:"172.18.0.2"`
				} `yaml:"amprobe"`
			} `yaml:"networks"`
		} `yaml:"amprobe"`
	} `yaml:"services"`
	Networks struct {
		Amprobe struct {
			Name   string `yaml:"name" default:"amprobe"`
			Driver string `yaml:"driver" default:"bridge"`
			Ipam   struct {
				Driver string   `yaml:"driver" default:"default"`
				Config []Config `yaml:"config"`
			} `yaml:"ipam"`
		} `yaml:"amprobe"`
	} `yaml:"networks"`
}

func DockerComposeConfig() DockerCompose {
	dockerCompose := DockerCompose{}
	defaults.SetDefaults(&dockerCompose)

	var subnetConfig Config
	defaults.SetDefaults(&subnetConfig)
	subnet := []Config{
		subnetConfig,
	}

	dockerCompose.Networks.Amprobe.Ipam.Config = subnet

	// 动态参数替换
	volumes := []string{
		"/data/amprobe/resources/amprobe/configs:/app/configs",
		"/data/amprobe/resources/amprobe/nginx/nginx.conf:/etc/nginx/nginx.conf",
		"/data/amprobe/resources/amprobe/nginx/conf.d:/etc/nginx/conf.d",
		"/data/amprobe/resources/amvector/socks/vector.sock:/app/vector.sock",
		"/data/amprobe/logs/amprobe:/app/logs",
	}
	dockerCompose.Services.Amprobe.Volumes = append(
		volumes,
		dockerCompose.Services.Amprobe.Volumes...,
	)

	return dockerCompose
}

func GenerateDockerCompose(dockerComposeFilePath string) error {
	conf := DockerComposeConfig()
	out, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	if err := os.WriteFile(dockerComposeFilePath, out, 0644); err != nil {
		return err
	}
	return nil
}
