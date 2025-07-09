// Package assets
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package assets

import (
	"amprobe/pkg/utils"
	"fmt"
	"os"
	"path/filepath"
)

const ResourcesDir string = "assets/resources"

func CopyFile(srcFile, dstFile string) error {
	// skip if exist
	if exists, err := utils.FileExists(dstFile); err != nil {
		return err
	} else if exists {
		return nil
	}

	data, err := Asset(srcFile)
	if err != nil {
		return err
	}

	info, err := AssetInfo(srcFile)
	if err != nil {
		return err
	}

	if err := utils.EnsureFileDirExists(dstFile); err != nil {
		return err
	}

	fmt.Println("dist file: ", dstFile)
	if err := os.WriteFile(dstFile, data, info.Mode()); err != nil {
		fmt.Println("write file error")
		return err
	}

	return nil
}

func CopyDir(srcDir, dstDir string) error {
	children, err := AssetDir(srcDir)
	if err != nil {
		return CopyFile(srcDir, dstDir)
	}
	fmt.Println("children: ", children)
	for _, child := range children {
		if err := CopyDir(filepath.Join(srcDir, child), filepath.Join(dstDir, child)); err != nil {
			fmt.Println(filepath.Join(srcDir, child), filepath.Join(dstDir, child))
			return err
		}
	}
	return nil
}
