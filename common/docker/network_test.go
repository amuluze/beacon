// Package docker
// Date: 2024/07/09 14:14:21
// Author: Amu
// Description:
package docker

import (
	"context"
	"testing"
)

func TestListNetwork(t *testing.T) {
	manager, _ := NewManager()
	nets, _ := manager.ListNetwork(context.Background())
	for _, net := range nets {
		t.Logf("network: %#v\n", net)
	}
}

func TestCreateNetwork(t *testing.T) {
	manager, _ := NewManager()
	networkID, err := manager.CreateNetwork(context.Background(), "test", "bridge", "172.20.0.0/24", "172.20.0.1", map[string]string{CreatedByProbe: "true"})
	if err != nil {
		t.Fatalf("create network failed: %v\n", err)
	}
	t.Logf("network id: %#v\n", networkID)
}

func TestQueryNetwork(t *testing.T) {
	manager, _ := NewManager()
	net, _ := manager.GetNetworkByID(context.Background(), "7be8e024bcb58caff65d38b39e42dff05e292e3f2f30963ae51732250b45a33f")
	t.Logf("network detail: %#v\n", net)
}

func TestDeleteNetwork(t *testing.T) {
	manager, _ := NewManager()
	err := manager.DeleteNetwork(context.Background(), "6185489062d75740df5edffb2e4f282399e19d73d0740aadc911c81548ac7b7b")
	if err != nil {
		t.Errorf("delete network failed: %v\n", err)
	}
	t.Log("delete network success")
}

func TestPruneNetwork(t *testing.T) {
	manager, _ := NewManager()
	err := manager.PruneNetwork(context.Background())
	if err != nil {
		t.Errorf("prune network failed: %v\n", err)
	}
	t.Log("prune network success")
}
