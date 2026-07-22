// Package docker
// Date: 2024/7/11 19:48
// Author: Amu
// Description:
package docker

import (
	"context"
	"testing"
)

func TestVersion(t *testing.T) {
	manager, _ := NewManager()
	version, _ := manager.Version(context.TODO())
	t.Logf("version: %#v", version)
}
