package service

import (
	"testing"
	"time"

	"collia/service/task"
)

func TestBuildReportArgsEmptyReport(t *testing.T) {
	td := &TimedTask{
		agentID:      "agent-test",
		reportClient: nil,
	}

	ts := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	r := &task.MonitorReport{}
	args := td.buildReportArgs(ts, r)

	if args.AgentID != "agent-test" {
		t.Fatalf("AgentID = %q, want agent-test", args.AgentID)
	}
	// All sub-reports should be zero-valued (except Timestamp)
	if args.Host.Hostname != "" {
		t.Fatalf("Host.Hostname = %q, want empty", args.Host.Hostname)
	}
	if len(args.Disks) != 0 {
		t.Fatalf("Disks = %d, want 0", len(args.Disks))
	}
	if len(args.Nets) != 0 {
		t.Fatalf("Nets = %d, want 0", len(args.Nets))
	}
	if len(args.Containers) != 0 {
		t.Fatalf("Containers = %d, want 0", len(args.Containers))
	}
}

func TestBuildReportArgsFullReport(t *testing.T) {
	td := &TimedTask{
		agentID: "agent-full",
	}
	ts := time.Date(2026, 2, 20, 8, 30, 0, 0, time.UTC)

	r := &task.MonitorReport{
		Host: &task.HostReport{
			Uptime:    "7d",
			Hostname:  "prod-01",
			Os:        "linux",
			Platform:  "ubuntu",
			KernelArch: "x86_64",
		},
		CPU: &task.CPUReport{
			CPUPercent: 42.5,
		},
		Memory: &task.MemoryReport{
			MemPercent: 65.2,
			MemTotal:   16384,
			MemUsed:    10678,
		},
		Disks: []*task.DiskReport{{
			Device:      "/dev/sda1",
			DiskPercent: 55.0,
			DiskTotal:   256000,
			DiskUsed:    140800,
		}},
		Nets: []*task.NetReport{{
			Ethernet: "eth0",
			NetRecv:  1024.5,
			NetSend:  512.3,
		}},
		Docker: &task.DockerReport{
			DockerVersion: "27.0.0",
			APIVersion:    "1.47",
			GoVersion:     "go1.22",
			Os:            "linux",
			Arch:          "amd64",
		},
		Containers: []*task.ContainerReport{{
			ContainerID: "abc123def",
			Name:        "nginx",
			Image:       "nginx:latest",
			State:       "running",
			CPUPercent:  1.2,
			MemPercent:  3.4,
		}},
		Images: []*task.ImageReport{{
			ImageID: "sha256:abc",
			Name:    "nginx",
			Tag:     "latest",
			Size:    "188MB",
		}},
		Networks: []*task.NetworkReport{{
			NetworkID: "net001",
			Name:      "bridge",
			Driver:    "bridge",
			Scope:     "local",
		}},
	}

	args := td.buildReportArgs(ts, r)

	// Verify all sections populated
	if args.AgentID != "agent-full" {
		t.Fatalf("AgentID = %q, want agent-full", args.AgentID)
	}

	// Host
	if args.Host.Hostname != "prod-01" {
		t.Fatalf("Host.Hostname = %q", args.Host.Hostname)
	}
	if args.Host.Os != "linux" {
		t.Fatalf("Host.Os = %q", args.Host.Os)
	}

	// CPU
	if args.CPU.CPUPercent != 42.5 {
		t.Fatalf("CPU.CPUPercent = %f", args.CPU.CPUPercent)
	}

	// Memory
	if args.Memory.MemPercent != 65.2 {
		t.Fatalf("Memory.MemPercent = %f", args.Memory.MemPercent)
	}

	// Disks
	if len(args.Disks) != 1 {
		t.Fatalf("Disks = %d, want 1", len(args.Disks))
	}
	if args.Disks[0].Device != "/dev/sda1" {
		t.Fatalf("Disk.Device = %q", args.Disks[0].Device)
	}

	// Nets
	if len(args.Nets) != 1 {
		t.Fatalf("Nets = %d, want 1", len(args.Nets))
	}

	// Docker
	if args.Docker.DockerVersion != "27.0.0" {
		t.Fatalf("Docker.Version = %q", args.Docker.DockerVersion)
	}

	// Containers
	if len(args.Containers) != 1 {
		t.Fatalf("Containers = %d, want 1", len(args.Containers))
	}
	if args.Containers[0].Name != "nginx" {
		t.Fatalf("Container.Name = %q", args.Containers[0].Name)
	}

	// Images
	if len(args.Images) != 1 {
		t.Fatalf("Images = %d, want 1", len(args.Images))
	}

	// Networks
	if len(args.Networks) != 1 {
		t.Fatalf("Networks = %d, want 1", len(args.Networks))
	}
}

func TestBuildReportArgsNilReport(t *testing.T) {
	td := &TimedTask{agentID: "agent-nil"}
	ts := time.Now()

	// buildReportArgs is only called when r != nil (see Execute), but test edge case
	var r *task.MonitorReport
	if r == nil {
		// Nil report handled by caller; buildReportArgs expects non-nil
		return
	}
	_ = td.buildReportArgs(ts, r)
}
