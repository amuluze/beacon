// Package profile
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package profile

import (
	"path/filepath"

	"github.com/amuluze/amprobe/amvector/assets"
	"github.com/amuluze/amprobe/amvector/pkg/resources"
	"github.com/amuluze/amprobe/amvector/pkg/utils"
	"github.com/docker/libcompose/yaml"
	"github.com/mcuadros/go-defaults"
)

type AmprobeProfile struct {
	ContainerName string         `yaml:"container_name" default:"amprobe"`
	Restart       string         `yaml:"restart" default:"always"`
	Image         string         `yaml:"image" default:"amuluze/amprobe:v1.3.5"`
	Ports         []string       `yaml:"ports" default:"[80:80,8000:8000]"`
	Volumes       *yaml.Volumes  `yaml:"volumes"`
	Networks      *yaml.Networks `yaml:"networks"`
}

func (a *AmprobeProfile) Entry() interface{} {
	return a
}

func (a *AmprobeProfile) InitResources(prefix string) error {
	amprobeDir := filepath.Join(prefix, resources.RootPath, resources.AmprobeConfigFolder)
	if err := utils.EnsureDirExists(amprobeDir); err != nil {
		return err
	}
	if err := assets.CopyDir(
		filepath.Join(assets.ResourcesDir, resources.AmprobeConfigFolder),
		amprobeDir,
	); err != nil {
		return err
	}

	amprobeNginxDir := filepath.Join(prefix, resources.RootPath, resources.AmprobeNginxFolder)
	if err := utils.EnsureDirExists(amprobeNginxDir); err != nil {
		return err
	}
	if err := assets.CopyDir(
		filepath.Join(assets.ResourcesDir, resources.AmprobeNginxFolder),
		amprobeNginxDir,
	); err != nil {
		return err
	}
	return nil
}

func (a *AmprobeProfile) BuildProfile(prefix string) error {
	defaults.SetDefaults(a)

	amprobeDir := filepath.Join(prefix, resources.AmprobeConfigFolder)
	amprobeNginxConfig := filepath.Join(prefix, resources.AmprobeNginxConfig)
	a.Volumes.Volumes = append(a.Volumes.Volumes, []*yaml.Volume{
		{
			Source:      amprobeDir,
			Destination: "/app/configs",
		},
		{
			Source:      amprobeNginxConfig,
			Destination: "/etc/nginx/nginx.conf",
		},
		{
			Source:      resources.AmvectorSockFile,
			Destination: "/app/vector.sock",
		},
		{
			Source:      "/var/run/docker.sock",
			Destination: "/var/run/docker.sock",
		},
		{
			Source:      "/tmp",
			Destination: "/tmp",
		},
	}...)
	bridgeNetwork := &yaml.Network{
		Name:        BridgeName,
		IPv4Address: ServiceAmprobeAddr,
	}
	a.Networks = &yaml.Networks{Networks: []*yaml.Network{bridgeNetwork}}
	return nil
}

func NewAmprobeBuilder() (ServiceBuilder, error) {
	return &AmprobeProfile{
		Volumes: new(yaml.Volumes),
	}, nil
}
