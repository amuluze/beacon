// Package model
// Date: 2024/3/6 11:05
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

// func (d *Host) AfterCreate(tx *gorm.DB) error {
// 	// 数据清理
// 	tx.Unscoped().Where("timestamp < ?", time.Now().Add(-time.Hour*24)).Delete(&Host{})
// 	return nil
// }

type CPU struct {
	gorm.Model
	Timestamp  time.Time
	CPUPercent float64
}

func (d *CPU) TableName() string {
	return "s_cpu"
}

// func (d *CPU) AfterCreate(tx *gorm.DB) error {
// 	tx.Unscoped().Where("timestamp < ?", time.Now().Add(-time.Hour*24)).Delete(&CPU{})
// 	return nil
// }

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

// func (d *Memory) AfterCreate(tx *gorm.DB) error {
// 	tx.Unscoped().Where("timestamp < ?", time.Now().Add(-time.Hour*24)).Delete(&Memory{})
// 	return nil
// }

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

// func (d *Disk) AfterCreate(tx *gorm.DB) error {
// 	tx.Unscoped().Where("timestamp < ?", time.Now().Add(-time.Hour*24)).Delete(&Disk{})
// 	return nil
// }

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

// func (d *Net) AfterCreate(tx *gorm.DB) error {
// 	tx.Unscoped().Where("timestamp < ?", time.Now().Add(-time.Hour*24)).Delete(&Net{})
// 	return nil
// }
