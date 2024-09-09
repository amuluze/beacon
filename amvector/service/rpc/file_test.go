// Package rpc
// Date: 2024/06/25 00:25:29
// Author: Amu
// Description:
package rpc

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"amvector/service/schema"
)

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func TestFilesSearch(t *testing.T) {

	args := &schema.FilesSearchArgs{
		Path: "/Users/amu",
	}
	files, err := os.ReadDir(args.Path)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			size, _ := DirSize(filepath.Join(args.Path, file.Name()))
			fmt.Println(size)
		}
	}
}
