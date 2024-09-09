// Package rpc
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package schema

import "time"

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
	Timestamp   time.Time `json:"timestamp"`
	Device      string    `json:"device"`
	DiskPercent float64   `json:"disk_percent"`
	DiskTotal   float64   `json:"disk_total"`
	DiskUsed    float64   `json:"disk_used"`
	DiskRead    float64   `json:"disk_read"`
	DiskWrite   float64   `json:"disk_write"`
}

type Net struct {
	Timestamp time.Time `json:"timestamp"`
	Ethernet  string    `json:"ethernet"`
	NetRecv   float64   `json:"net_recv"`
	NetSend   float64   `json:"net_send"`
}

type HostSummaryArgs struct {
	Timestamp time.Time `json:"timestamp"`
}

type HostSummaryReply struct {
	Data Host `json:"data"`
}

type CPUSummaryArgs struct {
	Timestamp time.Time `json:"timestamp"`
}

type CPUSummaryReply struct {
	Data CPU `json:"data"`
}

type MemorySummaryArgs struct {
	Timestamp time.Time `json:"timestamp"`
}

type MemorySummaryReply struct {
	Data Memory `json:"data"`
}

type DiskSummaryArgs struct {
	Timestamp time.Time           `json:"timestamp"`
	Devices   map[string]struct{} `json:"devices"`
	Interval  int                 `json:"interval"`
}

type DiskSummaryReply struct {
	Data []Disk `json:"data"`
}

type NetSummaryArgs struct {
	Timestamp time.Time           `json:"timestamp"`
	Ethernets map[string]struct{} `json:"ethernets"`
	Interval  int                 `json:"interval"`
}

type NetSummaryReply struct {
	Data []Net `json:"data"`
}
