// Package profile
// Date: 2024/7/3 09:34
// Author: Amu
// Description:
package profile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/amuluze/amprobe/amvector/pkg/utils"

	"github.com/amuluze/amprobe/amvector/pkg/resources"
	"gopkg.in/yaml.v2"
)

type Profile struct {
	Services        *Services                 `yaml:"services"`
	ServiceBuilders map[string]ServiceBuilder `yaml:"-"`
}

type Services struct {
	Network *NetworkProfile `yaml:"network"`
	Amprobe *AmprobeProfile `yaml:"amprobe"`
}

func (p *Profile) Save(prefix string) error {
	profileContent, err := yaml.Marshal(p)
	if err != nil {
		return err
	}
	// folder for save service_profile
	serviceProfileFolder := filepath.Join(prefix, resources.RootPath, resources.AmvectorFolder)
	if err := utils.EnsureDirExists(serviceProfileFolder); err != nil {
		return err
	}
	serviceProfile := filepath.Join(prefix, resources.RootPath, resources.AmvectorServiceProfilePath)
	fmt.Printf("server profile: %v\n", serviceProfile)
	if err := os.WriteFile(serviceProfile, profileContent, 0644); err != nil {
		return err
	}
	return nil
}

func (p *Profile) InitResources(prefix string) error {
	if err := utils.EnsureDirExists(prefix); err != nil {
		return err
	}

	for _, serviceBuilder := range p.ServiceBuilders {
		if err := serviceBuilder.InitResources(prefix); err != nil {
			return err
		}
	}
	return nil
}

func (p *Profile) BuildProfile(prefix string) error {
	for serviceName, serviceBuilder := range p.ServiceBuilders {
		if err := serviceBuilder.BuildProfile(prefix); err != nil {
			continue
		}
		switch serviceName {
		case "network":
			p.Services.Network = serviceBuilder.Entry().(*NetworkProfile)
		case "amprobe":
			p.Services.Amprobe = serviceBuilder.Entry().(*AmprobeProfile)
		}
	}
	return nil
}

func NewProfileBuilder(serviceBuilders map[string]ServiceBuilder) Builder {
	return &Profile{
		Services:        new(Services),
		ServiceBuilders: serviceBuilders,
	}
}

type Builder interface {
	Save(prefix string) error
	InitResources(prefix string) error
	BuildProfile(prefix string) error
}

func ReadProfile(profilePath string) (*Profile, error) {
	profileContent, err := os.ReadFile(profilePath)
	if err != nil {
		return nil, err
	}
	profile := &Profile{}
	if err := yaml.Unmarshal(profileContent, profile); err != nil {
		return nil, err
	}
	return profile, nil
}

type ServiceEntry struct {
	Name    string
	Desc    string
	Builder func() (ServiceBuilder, error)
}

type ServiceBuilder interface {
	Entry() interface{}
	InitResources(prefix string) error
	BuildProfile(prefix string) error
}
