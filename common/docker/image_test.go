// Package docker
// Date: 2024/07/09 14:14:04
// Author: Amu
// Description:
package docker

import (
	"context"
	"fmt"
	"testing"
)

func TestListImage(t *testing.T) {
	manager, _ := NewManager()
	images, _ := manager.ListImage(context.Background())
	for _, image := range images {
		fmt.Printf("image: %#v\n", image)
	}
}

func TestSearchImage(t *testing.T) {
	term := "ubuntu"
	manager, _ := NewManager()
	result, err := manager.SearchImage(context.Background(), term)
	if err != nil {
		fmt.Printf("err: %#v", err)
		return
	}
	for index, item := range result {
		fmt.Printf("index: %#v, item name: %#v\n", index, item.Name)
	}
}

func TestPullImage(t *testing.T) {
	imageName := "ubuntu:20.04"
	manager, _ := NewManager()
	err := manager.PullImage(context.Background(), imageName)
	fmt.Printf("pull error: %v", err)
}

func TestTagImage(t *testing.T) {
	oldTag := "ubuntu:latest"
	newTag := "ubuntu:22.04"
	manager, _ := NewManager()
	err := manager.TagImage(context.Background(), oldTag, newTag)
	t.Log("image tag error: ", err)
}

func TestExportImage(t *testing.T) {
	imageIDs := []string{"ubuntu:latest"}
	targetFile := "/Users/amu/Desktop/ubuntu.tar"
	manager, _ := NewManager()
	err := manager.ExportImage(context.Background(), imageIDs, targetFile)
	t.Log("export image error: ", err)
}

func TestImportImage(t *testing.T) {
	sourceFile := "/Users/amu/Desktop/ubuntu.tar"
	manager, _ := NewManager()
	err := manager.ImportImage(context.Background(), sourceFile)
	t.Log("import image error: ", err)
}
