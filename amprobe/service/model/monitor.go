// Package model
// Monitoring data GORM models stored on the Server side.
package model

import (
	"time"

	"gorm.io/gorm"
)

// ── Host ──

type MonitorHost struct {
	gorm.Model
	AgentID         string    `gorm:"index;index:idx_m_host_agent_time,priority:1"`
	Timestamp       time.Time `gorm:"index:idx_m_host_agent_time,priority:2"`
	Uptime          string
	Hostname        string
	Os              string
	Platform        string
	PlatformVersion string
	KernelVersion   string
	KernelArch      string
}

func (MonitorHost) TableName() string { return "m_host" }

// ── CPU ──

type MonitorCPU struct {
	gorm.Model
	AgentID    string    `gorm:"index;index:idx_m_cpu_agent_time,priority:1"`
	Timestamp  time.Time `gorm:"index:idx_m_cpu_agent_time,priority:2"`
	CPUPercent float64
}

func (MonitorCPU) TableName() string { return "m_cpu" }

// ── Memory ──

type MonitorMemory struct {
	gorm.Model
	AgentID    string    `gorm:"index;index:idx_m_memory_agent_time,priority:1"`
	Timestamp  time.Time `gorm:"index:idx_m_memory_agent_time,priority:2"`
	MemPercent float64
	MemTotal   float64
	MemUsed    float64
}

func (MonitorMemory) TableName() string { return "m_memory" }

// ── Disk ──

type MonitorDisk struct {
	gorm.Model
	AgentID     string    `gorm:"index;index:idx_m_disk_agent_time,priority:1"`
	Timestamp   time.Time `gorm:"index:idx_m_disk_agent_time,priority:2"`
	Device      string
	DiskPercent float64
	DiskTotal   float64
	DiskUsed    float64
	DiskRead    float64
	DiskWrite   float64
}

func (MonitorDisk) TableName() string { return "m_disk" }

// ── Network ──

type MonitorNet struct {
	gorm.Model
	AgentID   string    `gorm:"index;index:idx_m_net_agent_time,priority:1"`
	Timestamp time.Time `gorm:"index:idx_m_net_agent_time,priority:2"`
	Ethernet  string
	NetRecv   float64
	NetSend   float64
}

func (MonitorNet) TableName() string { return "m_net" }

// ── Docker ──

type MonitorDocker struct {
	gorm.Model
	AgentID       string    `gorm:"index;index:idx_m_docker_agent_time,priority:1"`
	Timestamp     time.Time `gorm:"index:idx_m_docker_agent_time,priority:2"`
	DockerVersion string
	APIVersion    string
	MinAPIVersion string
	GitCommit     string
	GoVersion     string
	Os            string
	Arch          string
}

func (MonitorDocker) TableName() string { return "m_docker" }

// ── Container ──

type MonitorContainer struct {
	gorm.Model
	AgentID     string    `gorm:"index;index:idx_m_container_agent_time,priority:1"`
	Timestamp   time.Time `gorm:"index:idx_m_container_agent_time,priority:2"`
	ContainerID string
	Name        string
	Image       string
	IP          string
	Ports       string
	State       string
	Uptime      string
	CPUPercent  float64
	MemPercent  float64
	MemUsage    float64
	MemLimit    float64
	Labels      string
}

func (MonitorContainer) TableName() string { return "m_container" }

// ── Image ──

type MonitorImage struct {
	gorm.Model
	AgentID   string    `gorm:"index;index:idx_m_image_agent_time,priority:1"`
	Timestamp time.Time `gorm:"index:idx_m_image_agent_time,priority:2"`
	ImageID   string
	Name      string
	Tag       string
	Created   string
	Size      string
	Number    int
}

func (MonitorImage) TableName() string { return "m_image" }

// ── Network ──

type MonitorNetwork struct {
	gorm.Model
	AgentID   string    `gorm:"index;index:idx_m_network_agent_time,priority:1"`
	Timestamp time.Time `gorm:"index:idx_m_network_agent_time,priority:2"`
	NetworkID string
	Name      string
	Driver    string
	Scope     string
	Created   string
	Internal  bool
	Subnet    string
	Gateway   string
	Labels    string
}

func (MonitorNetwork) TableName() string { return "m_network" }
