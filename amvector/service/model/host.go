// Package model
// Date: 2024/06/10 18:44:35
// Author: Amu
// Description:
package model

import (
	"time"

	"gorm.io/gorm"
)

type Host struct {
	gorm.Model
	Timestamp       time.Time
	Uptime          string
	Hostname        string
	Os              string
	Platform        string
	PlatformVersion string
	KernelVersion   string
	KernelArch      string
}

func (d *Host) TableName() string {
	return "s_host"
}

type CPU struct {
	gorm.Model
	Timestamp  time.Time
	CPUPercent float64
}

func (d *CPU) TableName() string {
	return "s_cpu"
}

type Memory struct {
	gorm.Model
	Timestamp  time.Time
	MemPercent float64
	MemTotal   float64
	MemUsed    float64
}

func (d *Memory) TableName() string {
	return "s_memory"
}

type Disk struct {
	gorm.Model
	Timestamp   time.Time
	Device      string
	DiskPercent float64
	DiskTotal   float64
	DiskUsed    float64
	DiskRead    float64
	DiskWrite   float64
}

func (d *Disk) TableName() string {
	return "s_disk"
}

type Net struct {
	gorm.Model
	Timestamp time.Time
	Ethernet  string
	NetRecv   float64
	NetSend   float64
}

func (d *Net) TableName() string {
	return "s_net"
}
