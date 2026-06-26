// Package utils
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// SanitizePath validates that the path is absolute and does not contain
// directory traversal sequences. It returns the cleaned absolute path or an error.
func SanitizePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path")
	}
	cleaned := filepath.Clean(path)
	if !filepath.IsAbs(cleaned) {
		return "", fmt.Errorf("path must be absolute: %s", path)
	}
	if strings.Contains(cleaned, "..") {
		return "", fmt.Errorf("path contains traversal sequence: %s", path)
	}
	return cleaned, nil
}

func FileExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		if stat.IsDir() {
			return false, fmt.Errorf("[%s] is a directory", path)
		} else {
			return true, nil
		}
	}
}

func CopyFile(src, dst string) (int64, error) {
	src, err := SanitizePath(src)
	if err != nil {
		return 0, err
	}
	dst, err = SanitizePath(dst)
	if err != nil {
		return 0, err
	}

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src) //#nosec G304 -- path sanitized by SanitizePath
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst) //#nosec G304 -- path sanitized by SanitizePath
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
