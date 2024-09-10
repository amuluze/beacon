// Package schema
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package schema

import "time"

type Usage struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

type Host struct {
	Timestamp       time.Time `json:"timestamp"`
	Uptime          string    `json:"uptime"`
	Hostname        string    `json:"hostname"`
	Os              string    `json:"os"`
	Platform        string    `json:"platform"`
	PlatformVersion string    `json:"platform_version"`
	KernelArch      string    `json:"kernel_arch"`
	KernelVersion   string    `json:"kernel_version"`
}

type CPU struct {
	Timestamp  time.Time `json:"timestamp"`
	CPUPercent float64   `json:"cpu_percent"`
}

type Memory struct {
	Timestamp  time.Time `json:"timestamp"`
	MemPercent float64   `json:"mem_percent"`
	MemTotal   float64   `json:"mem_total"`
	MemUsed    float64   `json:"mem_used"`
}

type Disk struct {
	Timestamp   time.Time `json:"timestamp,omitempty"`
	Device      string    `json:"device,omitempty"`
	DiskPercent float64   `json:"disk_percent,omitempty"`
	DiskTotal   float64   `json:"disk_total,omitempty"`
	DiskUsed    float64   `json:"disk_used,omitempty"`
	DiskRead    float64   `json:"disk_read,omitempty"`
	DiskWrite   float64   `json:"disk_write,omitempty"`
}

type DiskIO struct {
	Timestamp int64   `json:"timestamp"`
	IORead    float64 `json:"io_read"`
	IOWrite   float64 `json:"io_write"`
}

type DiskUsage struct {
	Device string   `json:"device"`
	Data   []DiskIO `json:"data"`
}

type Net struct {
	Timestamp time.Time `json:"timestamp"`
	Ethernet  string    `json:"ethernet"`
	NetRecv   float64   `json:"net_recv"`
	NetSend   float64   `json:"net_send"`
}

type NetIO struct {
	Timestamp int64   `json:"timestamp"`
	BytesSent float64 `json:"bytes_sent"`
	BytesRecv float64 `json:"bytes_recv"`
}

type NetUsage struct {
	Ethernet string  `json:"ethernet"`
	Data     []NetIO `json:"data"`
}

type HostInfoArgs struct{}

type HostInfoReply struct {
	Timestamp       int64  `json:"timestamp"`
	Uptime          string `json:"uptime"`
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platform_version"`
	KernelVersion   string `json:"kernel_version"`
	KernelArch      string `json:"kernel_arch"`
}

type CPUInfoArgs struct{}

type CPUInfoReply struct {
	Percent float64 `json:"percent"`
}

type CPUUsageArgs struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type CPUUsageReply struct {
	Data []Usage `json:"data"`
}

type MemoryInfoArgs struct{}

type MemoryInfoReply struct {
	Percent float64 `json:"percent"`
	Total   float64 `json:"total"`
	Used    float64 `json:"used"`
}

type MemoryUsageArgs struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type MemoryUsageReply struct {
	Data []Usage `json:"data"`
}

type DiskInfoArgs struct{}

type DiskInfoReply struct {
	Info []Disk `json:"info"`
}

type DiskUsageArgs struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type DiskUsageReply struct {
	Usage []DiskUsage `json:"usage"`
}

type NetUsageArgs struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type NetUsageReply struct {
	Usage []NetUsage `json:"usage"`
}
