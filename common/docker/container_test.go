// Package docker
// Date: 2024/07/09 14:13:44
// Author: Amu
// Description:
package docker

import (
	"context"
	"math"
	"testing"
	"time"
)

func TestPercentageReturnsJSONSafeValue(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		total float64
		want  float64
	}{
		{name: "normal", value: 25, total: 100, want: 25},
		{name: "zero total", value: 0, total: 0, want: 0},
		{name: "negative delta", value: -1, total: 100, want: 0},
		{name: "infinite input", value: math.Inf(1), total: 100, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := percentage(tt.value, tt.total)
			if got != tt.want {
				t.Fatalf("percentage(%v, %v) = %v, want %v", tt.value, tt.total, got, tt.want)
			}
			if math.IsNaN(got) || math.IsInf(got, 0) {
				t.Fatalf("percentage(%v, %v) returned non-finite value %v", tt.value, tt.total, got)
			}
		})
	}
}

func TestContainerLifecycle(t *testing.T) {
	manager := testManager(t)
	imageName := ensureTestImage(t, manager)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	networkName := uniqueDockerName("network")
	networkID, err := manager.CreateNetwork(ctx, networkName, "bridge", "", "", map[string]string{CreatedByProbe: "integration-test"})
	if err != nil {
		t.Fatalf("create network: %v", err)
	}
	cleanupNetwork(t, manager, networkID)

	containerName := uniqueDockerName("container")
	containerID, err := manager.CreateContainer(
		ctx,
		containerName,
		imageName,
		networkName,
		nil,
		nil,
		[]string{"BEACON_DOCKER_INTEGRATION=1"},
		[]string{"sh", "-c", "echo beacon-docker-integration-ready; sleep 60"},
		map[string]string{CreatedByProbe: "integration-test", ServerTypeLabel: WebServer},
	)
	if err != nil {
		t.Fatalf("create container: %v", err)
	}
	cleanupContainer(t, manager, containerID)

	if err := manager.StartContainer(ctx, containerID); err != nil {
		t.Fatalf("start container: %v", err)
	}
	if exists, err := manager.ContainerExists(ctx, containerID); err != nil || !exists {
		t.Fatalf("container exists = %v, err = %v", exists, err)
	}
	if exists, err := manager.HasSameNameContainer(ctx, containerName); err != nil || !exists {
		t.Fatalf("container name exists = %v, err = %v", exists, err)
	}
	if got, err := manager.GetContainerIDByContainerName(ctx, containerName); err != nil || got != containerID {
		t.Fatalf("container ID by name = %q, err = %v, want %q", got, err, containerID)
	}

	containers, err := manager.ListContainer(ctx)
	if err != nil {
		t.Fatalf("list containers: %v", err)
	}
	found := false
	for _, container := range containers {
		if container.ID == containerID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("created container %s missing from container list", containerID)
	}

	memPercent, _, memLimit, err := manager.GetContainerMem(ctx, containerID)
	if err != nil {
		t.Fatalf("get container memory: %v", err)
	}
	if math.IsNaN(memPercent) || math.IsInf(memPercent, 0) || memLimit <= 0 {
		t.Fatalf("invalid memory stats: percent=%v limit=%v", memPercent, memLimit)
	}
	cpuPercent, err := manager.GetContainerCpu(ctx, containerID)
	if err != nil {
		t.Fatalf("get container CPU: %v", err)
	}
	if math.IsNaN(cpuPercent) || math.IsInf(cpuPercent, 0) {
		t.Fatalf("invalid CPU stats: %v", cpuPercent)
	}

	logs, err := manager.ContainerLogs(ctx, containerID)
	if err != nil {
		t.Fatalf("open container logs: %v", err)
	}
	if err := logs.Close(); err != nil {
		t.Fatalf("close container logs: %v", err)
	}

	renamed := uniqueDockerName("renamed")
	if err := manager.RenameContainer(ctx, containerID, renamed); err != nil {
		t.Fatalf("rename container: %v", err)
	}
	if exists, err := manager.HasSameNameContainer(ctx, renamed); err != nil || !exists {
		t.Fatalf("renamed container exists = %v, err = %v", exists, err)
	}
	if err := manager.RestartContainer(ctx, containerID); err != nil {
		t.Fatalf("restart container: %v", err)
	}
	if err := manager.StopContainer(ctx, containerID); err != nil {
		t.Fatalf("stop container: %v", err)
	}
}
