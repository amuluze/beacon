// Package utils
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package utils

import (
	"fmt"
	"os"
)

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
