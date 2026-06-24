// Package schema
// Push report types for Agent → Server monitoring data transfer
package schema

import "time"

// ── Host report ──

type HostReport struct {
	Timestamp       time.Time `json:"timestamp"`
	Uptime          string    `json:"uptime"`
	Hostname        string    `json:"hostname"`
	Os              string    `json:"os"`
	Platform        string    `json:"platform"`
	PlatformVersion string    `json:"platform_version"`
	KernelArch      string    `json:"kernel_arch"`
	KernelVersion   string    `json:"kernel_version"`
}

// ── CPU report ──

type CPUReport struct {
	Timestamp  time.Time `json:"timestamp"`
	CPUPercent float64   `json:"cpu_percent"`
}

// ── Memory report ──

type MemoryReport struct {
	Timestamp  time.Time `json:"timestamp"`
	MemPercent float64   `json:"mem_percent"`
	MemTotal   float64   `json:"mem_total"`
	MemUsed    float64   `json:"mem_used"`
}

// ── Disk report ──

type DiskReport struct {
	Timestamp   time.Time `json:"timestamp"`
	Device      string    `json:"device"`
	DiskPercent float64   `json:"disk_percent"`
	DiskTotal   float64   `json:"disk_total"`
	DiskUsed    float64   `json:"disk_used"`
	DiskRead    float64   `json:"disk_read"`
	DiskWrite   float64   `json:"disk_write"`
}

// ── Net report ──

type NetReport struct {
	Timestamp time.Time `json:"timestamp"`
	Ethernet  string    `json:"ethernet"`
	NetRecv   float64   `json:"net_recv"`
	NetSend   float64   `json:"net_send"`
}

// ── Docker report ──

type DockerReport struct {
	Timestamp     time.Time `json:"timestamp"`
	DockerVersion string    `json:"docker_version"`
	APIVersion    string    `json:"api_version"`
	MinAPIVersion string    `json:"min_api_version"`
	GitCommit     string    `json:"git_commit"`
	GoVersion     string    `json:"go_version"`
	Os            string    `json:"os"`
	Arch          string    `json:"arch"`
}

// ── Container report ──

type ContainerReport struct {
	Timestamp   time.Time `json:"timestamp"`
	ContainerID string    `json:"container_id"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	IP          string    `json:"ip"`
	Ports       string    `json:"ports"`
	State       string    `json:"state"`
	Uptime      string    `json:"uptime"`
	CPUPercent  float64   `json:"cpu_percent"`
	MemPercent  float64   `json:"mem_percent"`
	MemUsage    float64   `json:"mem_usage"`
	MemLimit    float64   `json:"mem_limit"`
	Labels      string    `json:"labels"`
}

// ── Image report ──

type ImageReport struct {
	Timestamp time.Time `json:"timestamp"`
	ImageID   string    `json:"image_id"`
	Name      string    `json:"name"`
	Tag       string    `json:"tag"`
	Created   string    `json:"created"`
	Size      string    `json:"size"`
	Number    int       `json:"number"`
}

// ── Network report ──

type NetworkReport struct {
	Timestamp time.Time `json:"timestamp"`
	NetworkID string    `json:"network_id"`
	Name      string    `json:"name"`
	Driver    string    `json:"driver"`
	Scope     string    `json:"scope"`
	Created   string    `json:"created"`
	Internal  bool      `json:"internal"`
	Subnet    string    `json:"subnet"`
	Gateway   string    `json:"gateway"`
	Labels    string    `json:"labels"`
}

// ── Batch push args / reply (all at once for efficiency) ──

type MonitorReportArgs struct {
	AgentID    string            `json:"agent_id"`
	Host       HostReport        `json:"host"`
	CPU        CPUReport         `json:"cpu"`
	Memory     MemoryReport      `json:"memory"`
	Disks      []DiskReport      `json:"disks"`
	Nets       []NetReport       `json:"nets"`
	Docker     DockerReport      `json:"docker"`
	Containers []ContainerReport `json:"containers"`
	Images     []ImageReport     `json:"images"`
	Networks   []NetworkReport   `json:"networks"`
}

type MonitorReportReply struct{}
