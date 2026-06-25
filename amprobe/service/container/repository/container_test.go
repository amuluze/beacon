package repository

import (
	"context"
	"testing"

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
