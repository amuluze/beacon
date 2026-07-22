package rpc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	rpcSchema "common/rpc/schema"

	"github.com/amuluze/docker"
	dockercontainer "github.com/docker/docker/api/types/container"
	dockerclient "github.com/docker/docker/client"
)

type containerMutationManager interface {
	ListContainer(ctx context.Context) ([]docker.ContainerSummary, error)
	CreateContainer(ctx context.Context, containerName, imageName, networkName string, ports, volumes, environments, commands []string, labels map[string]string) (string, error)
	StartContainer(ctx context.Context, containerID string) error
	StopContainer(ctx context.Context, containerID string) error
	DeleteContainer(ctx context.Context, containerID string) error
	RenameContainer(ctx context.Context, containerID, newName string) error
}

type restartPolicyUpdater interface {
	UpdateRestartPolicy(ctx context.Context, containerID, policy string) error
}

type dockerRestartPolicyUpdater struct{}

func (dockerRestartPolicyUpdater) UpdateRestartPolicy(ctx context.Context, containerID, policy string) error {
	client, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()

	_, err = client.ContainerUpdate(ctx, containerID, dockercontainer.UpdateConfig{
		RestartPolicy: dockercontainer.RestartPolicy{Name: dockercontainer.RestartPolicyMode(normalizeRestartPolicy(policy))},
	})
	return err
}

func normalizeRestartPolicy(policy string) string {
	if policy == "" {
		return "always"
	}
	return policy
}

func findContainer(containers []docker.ContainerSummary, containerID string) (docker.ContainerSummary, error) {
	var matches []docker.ContainerSummary
	for _, item := range containers {
		if item.ID == containerID || strings.HasPrefix(item.ID, containerID) {
			matches = append(matches, item)
		}
	}
	if len(matches) == 0 {
		return docker.ContainerSummary{}, fmt.Errorf("container %s not found", containerID)
	}
	if len(matches) > 1 {
		return docker.ContainerSummary{}, fmt.Errorf("container id %s is ambiguous", containerID)
	}
	return matches[0], nil
}

func rollbackContainer(ctx context.Context, manager containerMutationManager, old docker.ContainerSummary, backupName, newContainerID string, wasRunning bool) error {
	var rollbackErrors []error
	if newContainerID != "" {
		if err := manager.DeleteContainer(ctx, newContainerID); err != nil {
			rollbackErrors = append(rollbackErrors, fmt.Errorf("delete replacement: %w", err))
		}
	}
	if err := manager.RenameContainer(ctx, old.ID, old.Name); err != nil {
		rollbackErrors = append(rollbackErrors, fmt.Errorf("restore container name from %s: %w", backupName, err))
	}
	if wasRunning {
		if err := manager.StartContainer(ctx, old.ID); err != nil {
			rollbackErrors = append(rollbackErrors, fmt.Errorf("restart original container: %w", err))
		}
	}
	return errors.Join(rollbackErrors...)
}

func recreateContainer(ctx context.Context, manager containerMutationManager, updater restartPolicyUpdater, args rpcSchema.ContainerUpdateArgs) (string, error) {
	containers, err := manager.ListContainer(ctx)
	if err != nil {
		return "", err
	}
	old, err := findContainer(containers, args.ContainerID)
	if err != nil {
		return "", err
	}

	name := args.Name
	if name == "" {
		name = old.Name
	}
	image := args.Image
	if image == "" {
		image = old.Image
	}
	network := args.Network
	if network == "" {
		network = old.Network
	}
	ports := args.Ports
	if ports == nil {
		ports = old.Ports
	}
	volumes := args.Volumes
	if volumes == nil {
		volumes = old.Volumes
	}
	environments := args.Environments
	if environments == nil {
		environments = old.Environments
	}
	labels := args.Labels
	if labels == nil {
		labels = old.Labels
	}
	backupName := old.Name + "-beacon-backup"
	wasRunning := old.State == "running"

	if wasRunning {
		if err := manager.StopContainer(ctx, old.ID); err != nil {
			return "", fmt.Errorf("stop original container: %w", err)
		}
	}
	if err := manager.RenameContainer(ctx, old.ID, backupName); err != nil {
		if wasRunning {
			_ = manager.StartContainer(ctx, old.ID)
		}
		return "", fmt.Errorf("rename original container: %w", err)
	}

	newID, err := manager.CreateContainer(ctx, name, image, network, ports, volumes, environments, args.Commands, labels)
	if err != nil {
		rollbackErr := rollbackContainer(ctx, manager, old, backupName, "", wasRunning)
		return "", errors.Join(err, rollbackErr)
	}
	if err := updater.UpdateRestartPolicy(ctx, newID, normalizeRestartPolicy(args.RestartPolicy)); err != nil {
		rollbackErr := rollbackContainer(ctx, manager, old, backupName, newID, wasRunning)
		return "", errors.Join(err, rollbackErr)
	}
	if wasRunning {
		if err := manager.StartContainer(ctx, newID); err != nil {
			rollbackErr := rollbackContainer(ctx, manager, old, backupName, newID, true)
			return "", errors.Join(err, rollbackErr)
		}
	}
	if err := manager.DeleteContainer(ctx, old.ID); err != nil {
		return "", fmt.Errorf("delete original container backup %s: %w", backupName, err)
	}
	return newID, nil
}
