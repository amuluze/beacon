// Package utils
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureDirExists 确保 filePath 存在，filePath 是一个绝对路径
func EnsureDirExists(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("could not create directory %s: %v", filePath, err)
		}
	}
	return nil
}

func EnsureFileDirExists(path string) error {
	return EnsureDirExists(filepath.Dir(path))
}
