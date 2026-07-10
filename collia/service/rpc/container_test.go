package rpc

import (
	"context"
	"errors"
	"reflect"
	"testing"

	rpcSchema "common/rpc/schema"

	"github.com/amuluze/docker"
)

var errContainerMutation = errors.New("container mutation failed")

type fakeContainerMutationManager struct {
	containers []docker.ContainerSummary
	createID   string
	createErr  error
	operations []string
	createArgs rpcSchema.ContainerUpdateArgs
}

func (f *fakeContainerMutationManager) ListContainer(context.Context) ([]docker.ContainerSummary, error) {
	return f.containers, nil
}

func (f *fakeContainerMutationManager) CreateContainer(_ context.Context, name, image, network string, ports, volumes, environments, commands []string, labels map[string]string) (string, error) {
	f.operations = append(f.operations, "create:"+name)
	f.createArgs = rpcSchema.ContainerUpdateArgs{
		Name: name, Image: image, Network: network, Ports: ports, Volumes: volumes,
		Environments: environments, Commands: commands, Labels: labels,
	}
	if f.createErr != nil {
		return "", f.createErr
	}
	return f.createID, nil
}

func (f *fakeContainerMutationManager) StartContainer(_ context.Context, id string) error {
	f.operations = append(f.operations, "start:"+id)
	return nil
}

func (f *fakeContainerMutationManager) StopContainer(_ context.Context, id string) error {
	f.operations = append(f.operations, "stop:"+id)
	return nil
}

func (f *fakeContainerMutationManager) DeleteContainer(_ context.Context, id string) error {
	f.operations = append(f.operations, "delete:"+id)
	return nil
}

func (f *fakeContainerMutationManager) RenameContainer(_ context.Context, id, name string) error {
	f.operations = append(f.operations, "rename:"+id+":"+name)
	return nil
}

type fakeRestartPolicyUpdater struct {
	err        error
	operations *[]string
}

func (f *fakeRestartPolicyUpdater) UpdateRestartPolicy(_ context.Context, id, policy string) error {
	*f.operations = append(*f.operations, "policy:"+id+":"+policy)
	return f.err
}

func TestContainerUpdateRecreatesAndRestoresRunningState(t *testing.T) {
	manager := &fakeContainerMutationManager{
		containers: []docker.ContainerSummary{{ID: "old-id", Name: "old", State: "running"}},
		createID:   "new-id",
	}
	updater := &fakeRestartPolicyUpdater{operations: &manager.operations}
	args := rpcSchema.ContainerUpdateArgs{
		ContainerID: "old-id", Name: "new", Image: "nginx:latest", Network: "bridge",
		Ports: []string{"8080:80"}, RestartPolicy: "unless-stopped",
	}

	newID, err := recreateContainer(context.Background(), manager, updater, args)
	if err != nil {
		t.Fatalf("recreateContainer: %v", err)
	}
	if newID != "new-id" {
		t.Fatalf("new id = %q, want new-id", newID)
	}
	want := []string{
		"stop:old-id",
		"rename:old-id:old-beacon-backup",
		"create:new",
		"policy:new-id:unless-stopped",
		"start:new-id",
		"delete:old-id",
	}
	if !reflect.DeepEqual(manager.operations, want) {
		t.Fatalf("operations = %#v, want %#v", manager.operations, want)
	}
}

func TestContainerUpdateRollsBackWhenCreateFails(t *testing.T) {
	manager := &fakeContainerMutationManager{
		containers: []docker.ContainerSummary{{ID: "old-id", Name: "old", State: "running"}},
		createErr:  errContainerMutation,
	}
	updater := &fakeRestartPolicyUpdater{operations: &manager.operations}

	_, err := recreateContainer(context.Background(), manager, updater, rpcSchema.ContainerUpdateArgs{
		ContainerID: "old-id", Name: "new", Image: "nginx:latest", Network: "bridge",
	})
	if !errors.Is(err, errContainerMutation) {
		t.Fatalf("error = %v, want create error", err)
	}
	want := []string{
		"stop:old-id",
		"rename:old-id:old-beacon-backup",
		"create:new",
		"rename:old-id:old",
		"start:old-id",
	}
	if !reflect.DeepEqual(manager.operations, want) {
		t.Fatalf("operations = %#v, want %#v", manager.operations, want)
	}
}

func TestContainerUpdateRollsBackWhenRestartPolicyFails(t *testing.T) {
	manager := &fakeContainerMutationManager{
		containers: []docker.ContainerSummary{{ID: "old-id", Name: "old", State: "stopped"}},
		createID:   "new-id",
	}
	updater := &fakeRestartPolicyUpdater{err: errContainerMutation, operations: &manager.operations}

	_, err := recreateContainer(context.Background(), manager, updater, rpcSchema.ContainerUpdateArgs{
		ContainerID: "old-id", Name: "new", Image: "nginx:latest", Network: "bridge", RestartPolicy: "always",
	})
	if !errors.Is(err, errContainerMutation) {
		t.Fatalf("error = %v, want policy error", err)
	}
	want := []string{
		"rename:old-id:old-beacon-backup",
		"create:new",
		"policy:new-id:always",
		"delete:new-id",
		"rename:old-id:old",
	}
	if !reflect.DeepEqual(manager.operations, want) {
		t.Fatalf("operations = %#v, want %#v", manager.operations, want)
	}
}

func TestContainerUpdateRejectsUnknownContainer(t *testing.T) {
	manager := &fakeContainerMutationManager{}
	updater := &fakeRestartPolicyUpdater{operations: &manager.operations}

	_, err := recreateContainer(context.Background(), manager, updater, rpcSchema.ContainerUpdateArgs{ContainerID: "missing"})
	if err == nil {
		t.Fatal("expected missing container error")
	}
	if len(manager.operations) != 0 {
		t.Fatalf("unexpected operations: %#v", manager.operations)
	}
}

func TestContainerUpdatePreservesOmittedRuntimeConfiguration(t *testing.T) {
	old := docker.ContainerSummary{
		ID:           "old-id",
		Name:         "old",
		Image:        "old:image",
		Network:      "old-network",
		State:        "stopped",
		Ports:        []string{"8080:80"},
		Volumes:      []string{"/data:/app/data"},
		Environments: []string{"TZ=Asia/Shanghai"},
		Labels:       map[string]string{"owner": "platform"},
	}
	manager := &fakeContainerMutationManager{
		containers: []docker.ContainerSummary{old},
		createID:   "new-id",
	}
	updater := &fakeRestartPolicyUpdater{operations: &manager.operations}

	_, err := recreateContainer(context.Background(), manager, updater, rpcSchema.ContainerUpdateArgs{
		ContainerID: "old-id",
		Name:        "renamed",
	})
	if err != nil {
		t.Fatalf("recreateContainer: %v", err)
	}
	if !reflect.DeepEqual(manager.createArgs.Ports, old.Ports) {
		t.Fatalf("ports = %#v, want %#v", manager.createArgs.Ports, old.Ports)
	}
	if !reflect.DeepEqual(manager.createArgs.Volumes, old.Volumes) {
		t.Fatalf("volumes = %#v, want %#v", manager.createArgs.Volumes, old.Volumes)
	}
	if !reflect.DeepEqual(manager.createArgs.Environments, old.Environments) {
		t.Fatalf("environments = %#v, want %#v", manager.createArgs.Environments, old.Environments)
	}
	if !reflect.DeepEqual(manager.createArgs.Labels, old.Labels) {
		t.Fatalf("labels = %#v, want %#v", manager.createArgs.Labels, old.Labels)
	}
}
