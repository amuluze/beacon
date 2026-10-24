// Package docker
// Date: 2024/7/11 19:48
// Author: Amu
// Description:
package docker

import (
	"context"
	"testing"
	"time"
)

func TestVersion(t *testing.T) {
	manager := testManager(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	version, err := manager.Version(ctx)
	if err != nil {
		t.Fatalf("get Docker version: %v", err)
	}
	if version.DockerVersion == "" || version.APIVersion == "" || version.OS == "" || version.Arch == "" {
		t.Fatalf("incomplete Docker version: %#v", version)
	}
}
