// Package main
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package main

import (
	"path/filepath"

	"amvector/assets"
	"amvector/pkg/compose"
	"amvector/pkg/config"
	"amvector/pkg/resources"
	"amvector/pkg/utils"
)

func runSetup() error {
	if err := setupConfig(prefix, configFile); err != nil {
		return err
	}

	if err := setupResources(prefix); err != nil {
		return err
	}

	if err := setupComposeConfig(prefix); err != nil {
		return err
	}

	return nil
}

func setupConfig(prefix, configFile string) error {
	cfg, err := config.Create(prefix, configFile)
	if err != nil {
		return err
	}
	if err := cfg.Save(); err != nil {
		return err
	}
	return nil
}

func setupResources(prefix string) error {
	amprobeConfigDir := filepath.Join(prefix, resources.RootPath, resources.AmprobeConfigFolder)
	if err := utils.EnsureDirExists(amprobeConfigDir); err != nil {
		return err
	}
	if err := assets.CopyDir(
		filepath.Join(assets.ResourcesDir, resources.AmprobeConfigFolder),
		amprobeConfigDir,
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

func setupComposeConfig(prefix string) error {
	filePath := filepath.Join(prefix, "docker-compose.yml")
	err := compose.GenerateDockerCompose(filePath)
	if err != nil {
		return err
	}
	return nil
}
