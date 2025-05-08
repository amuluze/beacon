// Package main
// Date       : 2025/5/7 14:53
// Author     : Amu
// Description:
package main

import (
	"amprobe/assets"
	"amprobe/pkg/utils"
	"amprobe/service/constants"
	"path/filepath"
)

func runSetup() error {
	if err := setupResources(); err != nil {
		return err
	}

	return nil
}

func setupResources() error {
	amprobeConfigDir := filepath.Join(constants.AmprobeFolder, constants.AmprobeConfigFolder)
	if err := utils.EnsureDirExists(amprobeConfigDir); err != nil {
		return err
	}
	if err := assets.CopyDir(
		filepath.Join(assets.ResourcesDir, constants.AmprobeConfigFolder),
		amprobeConfigDir,
	); err != nil {
		return err
	}

	return nil
}
