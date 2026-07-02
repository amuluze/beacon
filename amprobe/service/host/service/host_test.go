package service

import (
	"context"
	"testing"
	"time"

	"amprobe/service/testutil"
	rpcSchema "common/rpc/schema"
)

func TestHostServiceComputesInfoStalenessFromTimestamp(t *testing.T) {
	now := time.Now()
	staleTime := now.Add(-10 * time.Minute)
	repo := &testutil.FakeHostRepo{
		CPUInfoFn: func(context.Context, rpcSchema.CPUInfoArgs) (rpcSchema.CPUInfoReply, error) {
			return rpcSchema.CPUInfoReply{Timestamp: now.Unix(), Percent: 0.42}, nil
		},
		MemInfoFn: func(context.Context, rpcSchema.MemoryInfoArgs) (rpcSchema.MemoryInfoReply, error) {
			return rpcSchema.MemoryInfoReply{Timestamp: staleTime.Unix(), Percent: 0.70, Total: 100, Used: 70}, nil
		},
		DiskInfoFn: func(context.Context, rpcSchema.DiskInfoArgs) (rpcSchema.DiskInfoReply, error) {
			return rpcSchema.DiskInfoReply{Info: []rpcSchema.Disk{
				{Timestamp: now, Device: "/dev/disk1", DiskPercent: 0.30, DiskTotal: 100, DiskUsed: 30},
				{Timestamp: staleTime, Device: "/dev/disk2", DiskPercent: 0.60, DiskTotal: 100, DiskUsed: 60},
			}}, nil
		},
	}
	svc := NewHostService(repo, 5)

	cpu, err := svc.CPUInfo(context.Background())
	if err != nil {
		t.Fatalf("CPUInfo: %v", err)
	}
	if cpu.Timestamp != now.Unix() || cpu.Stale {
		t.Fatalf("CPUInfo timestamp/stale = (%d, %v), want (%d, false)", cpu.Timestamp, cpu.Stale, now.Unix())
	}

	mem, err := svc.MemInfo(context.Background())
	if err != nil {
		t.Fatalf("MemInfo: %v", err)
	}
	if mem.Timestamp != staleTime.Unix() || !mem.Stale {
		t.Fatalf("MemInfo timestamp/stale = (%d, %v), want (%d, true)", mem.Timestamp, mem.Stale, staleTime.Unix())
	}

	disk, err := svc.DiskInfo(context.Background())
	if err != nil {
		t.Fatalf("DiskInfo: %v", err)
	}
	if len(disk.Info) != 2 {
		t.Fatalf("disk info count = %d, want 2", len(disk.Info))
	}
	if disk.Info[0].Timestamp != now.Unix() || disk.Info[0].Stale {
		t.Fatalf("fresh disk timestamp/stale = (%d, %v), want (%d, false)", disk.Info[0].Timestamp, disk.Info[0].Stale, now.Unix())
	}
	if disk.Info[1].Timestamp != staleTime.Unix() || !disk.Info[1].Stale {
		t.Fatalf("stale disk timestamp/stale = (%d, %v), want (%d, true)", disk.Info[1].Timestamp, disk.Info[1].Stale, staleTime.Unix())
	}
}
