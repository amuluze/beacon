package rpc

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"common/database"
	rpcSchema "common/rpc/schema"

	"collia/service/model"

	"github.com/amuluze/docker"
)

var errContainerMutation = errors.New("container mutation failed")

func TestImageImportRejectsEmptyData(t *testing.T) {
	err := (&Service{}).ImageImport(context.Background(), rpcSchema.ImageImportArgs{FileName: "image.tar"}, &rpcSchema.ImageImportReply{})
	if err == nil {
		t.Fatal("expected empty image data to be rejected")
	}
}

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

type fakeContainerRuntimeDeleter struct {
	err     error
	deleted []string
	// onDelete 在删除执行时回调，用于验证"记录此时尚未移除"。
	onDelete func()
}

func (f *fakeContainerRuntimeDeleter) DeleteContainer(_ context.Context, containerID string) error {
	f.deleted = append(f.deleted, containerID)
	if f.onDelete != nil {
		f.onDelete()
	}
	if f.err != nil {
		return f.err
	}
	return nil
}

func newContainerDeleteTestDB(t *testing.T) *database.DB {
	t.Helper()
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "container")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.Container)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func seedContainerRecord(t *testing.T, db *database.DB, containerID string) {
	t.Helper()
	if err := db.Create(&model.Container{ContainerID: containerID, Name: "seed"}).Error; err != nil {
		t.Fatalf("seed container: %v", err)
	}
}

func containerRecordCount(t *testing.T, db *database.DB, containerID string) int64 {
	t.Helper()
	var count int64
	if err := db.Model(&model.Container{}).Where("container_id = ?", containerID).Count(&count).Error; err != nil {
		t.Fatalf("count container: %v", err)
	}
	return count
}

// TestContainerDeleteKeepsRecordWhenRuntimeDeleteFails 覆盖运行时删除失败的场景：
// 本地记录必须保留，避免记录消失但容器仍在节点上。
func TestContainerDeleteKeepsRecordWhenRuntimeDeleteFails(t *testing.T) {
	db := newContainerDeleteTestDB(t)
	seedContainerRecord(t, db, "c-fail")

	deleter := &fakeContainerRuntimeDeleter{err: errContainerMutation}
	err := deleteContainerThenRecord(context.Background(), deleter, db, "c-fail")
	if err == nil {
		t.Fatal("expected runtime delete failure to be propagated")
	}
	if got := containerRecordCount(t, db, "c-fail"); got != 1 {
		t.Fatalf("record count = %d, want 1 (record must survive failed delete)", got)
	}
}

// TestContainerDeleteRemovesRecordAfterRuntimeDelete 覆盖正常删除：
// 运行时删除成功后本地记录必须被移除。
func TestContainerDeleteRemovesRecordAfterRuntimeDelete(t *testing.T) {
	db := newContainerDeleteTestDB(t)
	seedContainerRecord(t, db, "c-ok")

	deleter := &fakeContainerRuntimeDeleter{}
	if err := deleteContainerThenRecord(context.Background(), deleter, db, "c-ok"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if got := containerRecordCount(t, db, "c-ok"); got != 0 {
		t.Fatalf("record count = %d, want 0", got)
	}
	if len(deleter.deleted) != 1 || deleter.deleted[0] != "c-ok" {
		t.Fatalf("runtime delete calls = %v, want [c-ok]", deleter.deleted)
	}
}

// TestContainerDeleteRunsRuntimeBeforeRecordRemoval 锁定删除顺序：
// 运行时删除执行时本地记录必须仍然在库。
func TestContainerDeleteRunsRuntimeBeforeRecordRemoval(t *testing.T) {
	db := newContainerDeleteTestDB(t)
	seedContainerRecord(t, db, "c-order")

	deleter := &fakeContainerRuntimeDeleter{}
	deleter.onDelete = func() {
		if got := containerRecordCount(t, db, "c-order"); got != 1 {
			t.Fatalf("record count during runtime delete = %d, want 1 (runtime must be deleted first)", got)
		}
	}

	if err := deleteContainerThenRecord(context.Background(), deleter, db, "c-order"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
