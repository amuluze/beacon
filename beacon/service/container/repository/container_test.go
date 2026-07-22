package repository

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"beacon/pkg/contextx"
	beaconrpc "beacon/pkg/rpc"
	"beacon/service/model"
	"common/database"
	rpcSchema "common/rpc/schema"
)

type fakeCaller struct {
	method string
	err    error
}

func (f *fakeCaller) Call(ctx context.Context, method string, args interface{}, reply interface{}) error {
	f.method = method
	if f.err != nil {
		return f.err
	}
	if r, ok := reply.(*rpcSchema.ContainerUpdateReply); ok {
		r.ContainerID = "updated"
	}
	return nil
}

func (f *fakeCaller) StreamCall(ctx context.Context, method string, args interface{}) (<-chan []byte, error) {
	return nil, nil
}

func (f *fakeCaller) Close() error {
	return nil
}

func TestContainerUpdateCallsAgentRPC(t *testing.T) {
	caller := &fakeCaller{}
	repo := &ContainerRepo{RPCClient: caller}

	reply, err := repo.ContainerUpdate(context.Background(), rpcSchema.ContainerUpdateArgs{ContainerID: "abc123"})
	if err != nil {
		t.Fatalf("ContainerUpdate returned error: %v", err)
	}
	if caller.method != "ContainerUpdate" {
		t.Fatalf("called method %q, want ContainerUpdate", caller.method)
	}
	if reply.ContainerID != "updated" {
		t.Fatalf("reply container id = %q, want updated", reply.ContainerID)
	}
}

func TestContainerUpdateRejectsMissingAgentIDBeforeRPC(t *testing.T) {
	repo := &ContainerRepo{RPCClient: beaconrpc.NewTunnelClient(nil)}

	_, err := repo.ContainerUpdate(context.Background(), rpcSchema.ContainerUpdateArgs{ContainerID: "abc123"})
	if !errors.Is(err, contextx.ErrMissingAgentID) {
		t.Fatalf("error = %v, want ErrMissingAgentID", err)
	}
}

func TestContainerUpdateRejectsInvalidAgentIDBeforeRPC(t *testing.T) {
	repo := &ContainerRepo{RPCClient: beaconrpc.NewTunnelClient(nil)}
	ctx := contextx.NewAgentID(context.Background(), "agent/../../etc")

	_, err := repo.ContainerUpdate(ctx, rpcSchema.ContainerUpdateArgs{ContainerID: "abc123"})
	if !errors.Is(err, contextx.ErrInvalidAgentID) {
		t.Fatalf("error = %v, want ErrInvalidAgentID", err)
	}
}

func TestContainerListRequiresAgentID(t *testing.T) {
	repo := &ContainerRepo{}
	if _, err := repo.ContainerList(context.Background(), rpcSchema.ContainerQueryArgs{}); !errors.Is(err, contextx.ErrAgentIDRequired) {
		t.Fatalf("error = %v, want ErrAgentIDRequired", err)
	}
}

func TestContainerListReturnsLatestBatchForSelectedAgent(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.MonitorContainer)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	now := time.Now()
	rows := []model.MonitorContainer{
		{AgentID: "agent-a", Name: "app", ContainerID: "old", Timestamp: now.Add(-time.Minute), Ports: "80", State: "stopped"},
		{AgentID: "agent-a", Name: "removed", ContainerID: "removed", Timestamp: now.Add(-time.Minute), Image: "unused:test", State: "stopped"},
		{AgentID: "agent-a", Name: "app", ContainerID: "new", Timestamp: now, Ports: "80", State: "running"},
		{AgentID: "agent-b", Name: "app", ContainerID: "other-agent", Timestamp: now.Add(time.Minute), Ports: "80", State: "paused"},
	}
	if err := db.Create(&rows).Error; err != nil {
		t.Fatalf("create containers: %v", err)
	}

	repo := &ContainerRepo{DB: db}
	ctx := contextx.NewAgentID(context.Background(), "agent-a")
	reply, err := repo.ContainerList(ctx, rpcSchema.ContainerQueryArgs{Page: 1, Size: 10})
	if err != nil {
		t.Fatalf("ContainerList returned error: %v", err)
	}
	if len(reply.Data) != 1 {
		t.Fatalf("len(reply.Data) = %d, want 1", len(reply.Data))
	}
	if got := reply.Data[0].ContainerID; got != "new" {
		t.Fatalf("container id = %q, want latest selected agent container new", got)
	}
	if reply.Freshness.Stale {
		t.Fatalf("freshness = %+v, want non-stale latest data", reply.Freshness)
	}

	count, err := repo.ContainerCount(ctx, rpcSchema.ContainerCountArgs{})
	if err != nil {
		t.Fatalf("ContainerCount returned error: %v", err)
	}
	if count.Count != 1 {
		t.Fatalf("container count = %d, want 1", count.Count)
	}

	byRemovedImage, err := repo.ContainersByImage(ctx, "unused:test")
	if err != nil {
		t.Fatalf("ContainersByImage returned error: %v", err)
	}
	if byRemovedImage != 0 {
		t.Fatalf("containers by removed image = %d, want 0 outside latest report batch", byRemovedImage)
	}
}

func TestContainerListIncludesContainersWithoutPublishedPorts(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.MonitorContainer)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	now := time.Now()
	rows := []model.MonitorContainer{
		{AgentID: "agent-a", Name: "created-app", ContainerID: "created", Timestamp: now, State: "created", Ports: ""},
		{AgentID: "agent-a", Name: "web-app", ContainerID: "running", Timestamp: now, State: "running", Ports: "8080"},
	}
	if err := db.Create(&rows).Error; err != nil {
		t.Fatalf("create containers: %v", err)
	}

	repo := &ContainerRepo{DB: db}
	ctx := contextx.NewAgentID(context.Background(), "agent-a")
	reply, err := repo.ContainerList(ctx, rpcSchema.ContainerQueryArgs{Page: 1, Size: 10})
	if err != nil {
		t.Fatalf("ContainerList returned error: %v", err)
	}
	if len(reply.Data) != 2 {
		t.Fatalf("len(reply.Data) = %d, want 2 including portless created container", len(reply.Data))
	}
	if reply.Data[0].ContainerID != "created" || reply.Data[0].Ports != "" {
		t.Fatalf("first container = %+v, want portless created container", reply.Data[0])
	}

	count, err := repo.ContainerCount(ctx, rpcSchema.ContainerCountArgs{})
	if err != nil {
		t.Fatalf("ContainerCount returned error: %v", err)
	}
	if count.Count != 2 {
		t.Fatalf("container count = %d, want 2 including portless created container", count.Count)
	}
}

func TestImageDeleteInvalidatesSelectedAgentCacheAfterRPCSuccess(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.MonitorImage)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	rows := []model.MonitorImage{
		{AgentID: "agent-a", ImageID: "delete-me", Name: "test", Tag: "delete"},
		{AgentID: "agent-a", ImageID: "keep-me", Name: "test", Tag: "keep"},
		{AgentID: "agent-b", ImageID: "delete-me", Name: "test", Tag: "other-agent"},
	}
	if err := db.Create(&rows).Error; err != nil {
		t.Fatalf("create images: %v", err)
	}

	caller := &fakeCaller{}
	repo := &ContainerRepo{RPCClient: caller, DB: db}
	ctx := contextx.NewAgentID(context.Background(), "agent-a")
	if err := repo.ImageDelete(ctx, rpcSchema.ImageDeleteArgs{ImageID: "delete-me"}); err != nil {
		t.Fatalf("ImageDelete returned error: %v", err)
	}
	if caller.method != "ImageDelete" {
		t.Fatalf("called method %q, want ImageDelete", caller.method)
	}

	assertImageCacheCount(t, db, "agent-a", "delete-me", 0)
	assertImageCacheCount(t, db, "agent-a", "keep-me", 1)
	assertImageCacheCount(t, db, "agent-b", "delete-me", 1)
}

func TestImageDeleteKeepsCacheWhenAgentRPCFails(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.MonitorImage)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	if err := db.Create(&model.MonitorImage{AgentID: "agent-a", ImageID: "keep-on-failure"}).Error; err != nil {
		t.Fatalf("create image: %v", err)
	}

	rpcErr := errors.New("agent image delete failed")
	repo := &ContainerRepo{RPCClient: &fakeCaller{err: rpcErr}, DB: db}
	ctx := contextx.NewAgentID(context.Background(), "agent-a")
	if err := repo.ImageDelete(ctx, rpcSchema.ImageDeleteArgs{ImageID: "keep-on-failure"}); !errors.Is(err, rpcErr) {
		t.Fatalf("ImageDelete error = %v, want %v", err, rpcErr)
	}

	assertImageCacheCount(t, db, "agent-a", "keep-on-failure", 1)
}

func assertImageCacheCount(t *testing.T, db *database.DB, agentID, imageID string, want int64) {
	t.Helper()
	var got int64
	if err := db.Model(&model.MonitorImage{}).
		Where("agent_id = ? AND image_id = ?", agentID, imageID).
		Count(&got).Error; err != nil {
		t.Fatalf("count image cache: %v", err)
	}
	if got != want {
		t.Fatalf("image cache count for %s/%s = %d, want %d", agentID, imageID, got, want)
	}
}

func TestContainerUsageRespectsClosedTimeWindow(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.MonitorContainer)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	start := time.Unix(1_000, 0)
	middle := time.Unix(1_100, 0)
	end := time.Unix(1_200, 0)
	after := time.Unix(1_300, 0)
	for index, timestamp := range []time.Time{start, middle, end, after} {
		if err := db.Create(&model.MonitorContainer{
			AgentID: "agent-a", Timestamp: timestamp, Name: "app",
			CPUPercent: float64(index), MemUsage: float64(index),
		}).Error; err != nil {
			t.Fatalf("create container sample: %v", err)
		}
	}

	repo := &ContainerRepo{DB: db}
	ctx := contextx.NewAgentID(context.Background(), "agent-a")
	reply, err := repo.Usage(ctx, rpcSchema.ContainerUsageArgs{StartTime: start.Unix(), EndTime: end.Unix()})
	if err != nil {
		t.Fatalf("Usage returned error: %v", err)
	}
	want := []int64{start.Unix(), middle.Unix(), end.Unix()}
	assertContainerUsageTimestamps(t, "CPU", want, reply.CPUUsage["app"])
	assertContainerUsageTimestamps(t, "memory", want, reply.MemUsage["app"])
}

func assertContainerUsageTimestamps(t *testing.T, metric string, want []int64, items []rpcSchema.Usage) {
	t.Helper()
	got := make([]int64, 0, len(items))
	for _, item := range items {
		got = append(got, item.Timestamp)
	}
	if len(got) != len(want) {
		t.Fatalf("%s timestamps = %v, want %v", metric, got, want)
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("%s timestamps = %v, want %v", metric, got, want)
		}
	}
}
