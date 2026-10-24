// Package docker
// Date: 2024/07/09 14:15:12
// Author: Amu
// Description:
package docker

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

const defaultTestImage = "alpine:3.20"

func testManager(t *testing.T) *Manager {
	t.Helper()

	manager, err := NewManager()
	if err != nil {
		t.Fatalf("create Docker manager: %v", err)
	}
	t.Cleanup(func() {
		if err := manager.client.Close(); err != nil {
			t.Errorf("close Docker client: %v", err)
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := manager.Version(ctx); err != nil {
		t.Fatalf("connect to Docker daemon at %q: %v", manager.client.DaemonHost(), err)
	}
	return manager
}

func testImageName() string {
	if image := os.Getenv("BEACON_DOCKER_TEST_IMAGE"); image != "" {
		return image
	}
	return defaultTestImage
}

func ensureTestImage(t *testing.T, manager *Manager) string {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	imageName := testImageName()
	if _, err := manager.GetImageByName(ctx, imageName); err == nil {
		return imageName
	}
	if err := manager.PullImage(ctx, imageName); err != nil {
		t.Fatalf("pull integration image %q: %v", imageName, err)
	}
	if _, err := manager.GetImageByName(ctx, imageName); err != nil {
		t.Fatalf("integration image %q unavailable after pull: %v", imageName, err)
	}
	return imageName
}

func uniqueDockerName(kind string) string {
	return fmt.Sprintf("beacon-docker-it-%s-%d", kind, time.Now().UnixNano())
}

func cleanupContainer(t *testing.T, manager *Manager, containerID string) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := manager.DeleteContainer(ctx, containerID); err != nil {
			t.Logf("cleanup container %s: %v", containerID, err)
		}
	})
}

func cleanupNetwork(t *testing.T, manager *Manager, networkID string) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := manager.DeleteNetwork(ctx, networkID); err != nil {
			t.Logf("cleanup network %s: %v", networkID, err)
		}
	})
}

func TestNewManagerUsesDockerHost(t *testing.T) {
	const host = "unix:///tmp/beacon-docker-test.sock"
	t.Setenv("DOCKER_HOST", host)

	manager, err := NewManager()
	if err != nil {
		t.Fatalf("create Docker manager: %v", err)
	}
	defer manager.client.Close()
	if got := manager.client.DaemonHost(); got != host {
		t.Fatalf("daemon host = %q, want %q", got, host)
	}
}
