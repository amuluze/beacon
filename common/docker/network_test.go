// Package docker
// Date: 2024/07/09 14:14:21
// Author: Amu
// Description:
package docker

import (
	"context"
	"testing"
	"time"
)

func TestNetworkLifecycle(t *testing.T) {
	manager := testManager(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	name := uniqueDockerName("network-only")

	networkID, err := manager.CreateNetwork(ctx, name, "bridge", "", "", map[string]string{CreatedByProbe: "integration-test"})
	if err != nil {
		t.Fatalf("create network: %v", err)
	}
	cleanupNetwork(t, manager, networkID)
	if exists, err := manager.HasSameNameNetwork(ctx, name); err != nil || !exists {
		t.Fatalf("network exists = %v, err = %v", exists, err)
	}
	byName, err := manager.GetNetworkByName(ctx, name)
	if err != nil {
		t.Fatalf("get network by name: %v", err)
	}
	if byName.ID != networkID {
		t.Fatalf("network ID by name = %q, want %q", byName.ID, networkID)
	}
	byID, err := manager.GetNetworkByID(ctx, networkID)
	if err != nil {
		t.Fatalf("get network by ID: %v", err)
	}
	if byID.Name != name {
		t.Fatalf("network name by ID = %q, want %q", byID.Name, name)
	}
	networks, err := manager.ListNetwork(ctx)
	if err != nil {
		t.Fatalf("list networks: %v", err)
	}
	found := false
	for _, item := range networks {
		if item.ID == networkID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("created network %s missing from network list", networkID)
	}
}
