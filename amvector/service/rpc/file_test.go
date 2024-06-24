// Package rpc
// Date: 2024/06/25 00:25:29
// Author: Amu
// Description:
package rpc

import (
	"fmt"
	"os"
	"testing"

	"github.com/amuluze/amvector/service/schema"
)

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
		fmt.Println(file.Name())
		info, _ := file.Info()
		fmt.Printf("file type: %#v\n", info.IsDir())
		fmt.Printf("file info: %#v\n", info.Size())
		fmt.Println("==================")
	}
}
