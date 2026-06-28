// Package task
package task

import (
	"errors"
	"time"

	"github.com/amuluze/docker"
	"github.com/patrickmn/go-cache"
)

const (
	LatestDiskReadKey   = "latest_disk_io_read_"
	LatestDisKWriteKey  = "latest_disk_io_write_"
	LatestNetReceiveKey = "latest_net_io_receive_"
	LatestNetSendKey    = "latest_net_io_send_"
)

var _ ITask = (*Task)(nil)

type ITask interface {
	Report(timestamp time.Time) (*MonitorReport, error)
}

type Task struct {
	interval int
	manager  *docker.Manager
	devices  map[string]struct{}
	ethernet map[string]struct{}
	cache    *cache.Cache
}

// MonitorReport holds all collected monitoring data for one tick.
type MonitorReport struct {
	Host       *HostReport
	CPU        *CPUReport
	Memory     *MemoryReport
	Disks      []*DiskReport
	Nets       []*NetReport
	Docker     *DockerReport
	Containers []*ContainerReport
	Images     []*ImageReport
	Networks   []*NetworkReport
	Error      error
}

type HostReport struct {
	Uptime          string
	Hostname        string
	Os              string
	Platform        string
	PlatformVersion string
	KernelVersion   string
	KernelArch      string
}

type CPUReport struct {
	CPUPercent float64
}

type MemoryReport struct {
	MemPercent float64
	MemTotal   float64
	MemUsed    float64
}

type DiskReport struct {
	Device      string
	DiskPercent float64
	DiskTotal   float64
	DiskUsed    float64
	DiskRead    float64
	DiskWrite   float64
}

type NetReport struct {
	Ethernet string
	NetRecv  float64
	NetSend  float64
}

type DockerReport struct {
	DockerVersion string
	APIVersion    string
	MinAPIVersion string
	GitCommit     string
	GoVersion     string
	Os            string
	Arch          string
}

type ContainerReport struct {
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

type ImageReport struct {
	ImageID string
	Name    string
	Tag     string
	Created string
	Size    string
	Number  int
}

type NetworkReport struct {
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

func NewTask(interval int, manager *docker.Manager, dev map[string]struct{}, eth map[string]struct{}) *Task {
	return &Task{
		interval: interval,
		manager:  manager,
		devices:  dev,
		ethernet: eth,
		cache:    cache.New(5*time.Minute, 60*time.Second),
	}
}

// Report collects all monitoring data and returns it as a MonitorReport.
func (a *Task) Report(timestamp time.Time) (*MonitorReport, error) {
	hostReport, hostErr := a.HostTask()
	cpuReport, cpuErr := a.CPUTask()
	memoryReport, memoryErr := a.MemoryTask()
	diskReports, diskErr := a.DiskTask()
	netReports, netErr := a.NetTask()

	report := &MonitorReport{
		Host:       hostReport,
		CPU:        cpuReport,
		Memory:     memoryReport,
		Disks:      diskReports,
		Nets:       netReports,
		Docker:     a.DockerTask(),
		Containers: a.ContainerTask(),
		Images:     a.ImageTask(),
		Networks:   a.NetworkTask(),
	}
	report.Error = errors.Join(hostErr, cpuErr, memoryErr, diskErr, netErr)
	if report.Error != nil {
		return report, report.Error
	}
	return report, nil
}
