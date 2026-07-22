package schema

import (
	"testing"

	"beacon/pkg/validatex"
)

func TestContainerCreateArgsAcceptsReportedDockerNetworkID(t *testing.T) {
	args := ContainerCreateArgs{
		ContainerName: "codex-e2e-container",
		ImageName:     "codex-e2e-image:latest",
		NetworkID:     "112dbc",
		NetworkMode:   "bridge",
		NetworkName:   "codex-e2e-network",
		RestartPolicy: "always",
	}

	if err := validatex.ValidateStruct(&args); err != nil {
		t.Fatalf("reported Docker network ID should be valid: %v", err)
	}
}
