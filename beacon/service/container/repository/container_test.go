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
}

func (f *fakeCaller) Call(ctx context.Context, method string, args interface{}, reply interface{}) error {
	f.method = method
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

func TestContainerListReturnsLatestPerNameForSelectedAgent(t *testing.T) {
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
}
