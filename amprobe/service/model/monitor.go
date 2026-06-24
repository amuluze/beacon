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
	AgentID         string `gorm:"index"`
	Timestamp       time.Time
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
	AgentID    string `gorm:"index"`
	Timestamp  time.Time
	CPUPercent float64
}

func (MonitorCPU) TableName() string { return "m_cpu" }

// ── Memory ──

type MonitorMemory struct {
	gorm.Model
	AgentID    string `gorm:"index"`
	Timestamp  time.Time
	MemPercent float64
	MemTotal   float64
	MemUsed    float64
}

func (MonitorMemory) TableName() string { return "m_memory" }

// ── Disk ──

type MonitorDisk struct {
	gorm.Model
	AgentID     string `gorm:"index"`
	Timestamp   time.Time
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
	AgentID   string `gorm:"index"`
	Timestamp time.Time
	Ethernet  string
	NetRecv   float64
	NetSend   float64
}

func (MonitorNet) TableName() string { return "m_net" }

// ── Docker ──

type MonitorDocker struct {
	gorm.Model
	AgentID       string `gorm:"index"`
	Timestamp     time.Time
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
	AgentID     string `gorm:"index"`
	Timestamp   time.Time
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
	AgentID   string `gorm:"index"`
	Timestamp time.Time
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
	AgentID   string `gorm:"index"`
	Timestamp time.Time
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
