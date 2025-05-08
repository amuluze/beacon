// Package main
// Date       : 2025/5/7 14:53
// Author     : Amu
// Description:
package main

import (
	"amprobe/assets"
	"amprobe/pkg/resources"
	"amprobe/pkg/utils"
)

func runSetup() error {
	if err := setupConfig(prefix, configFile); err != nil {
		return err
	}

	if err := setupResources(prefix); err != nil {
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
