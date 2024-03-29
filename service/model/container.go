// Package model
// Date: 2024/3/6 11:06
// Author: Amu
// Description:
package model

import (
	"time"

	"gorm.io/gorm"
)

type Containers []Container

type Container struct {
	gorm.Model
	Timestamp   time.Time
	ContainerID string
	Name        string
	Image       string
	IP          string
	State       string
	Uptime      string
	CPUPercent  float64
	MemPercent  float64
	MemUsage    float64
	MemLimit    float64
}

func (d *Container) TableName() string {
	return "s_container"
}

// func (d *Container) AfterCreate(tx *gorm.DB) error {
// 	// 数据清理
// 	tx.Unscoped().Where("timestamp < ?", time.Now().Add(-time.Minute*5)).Delete(&Container{})
// 	return nil
// }

type Docker struct {
	gorm.Model
	Timestamp     time.Time
	DockerVersion string
	APIVersion    string
	MinAPIVersion string
	GitCommit     string
	GoVersion     string
	Os            string
	Arch          string
}

func (d *Docker) TableName() string {
	return "s_docker"
}

// func (d *Docker) AfterCreate(tx *gorm.DB) error {
// 	tx.Unscoped().Where("timestamp < ?", time.Now().Add(-time.Minute*5)).Delete(&Docker{})
// 	return nil
// }

type Images []Image

type Image struct {
	gorm.Model
	Timestamp time.Time
	ImageID   string
	Name      string
	Tag       string
	Created   string
	Size      string
}

func (i *Image) TableName() string { return "s_image" }

// func (i *Image) AfterCreate(tx *gorm.DB) error {
// 	tx.Unscoped().Where("timestamp < ?", time.Now().Add(-time.Minute*5)).Delete(&Image{})
// 	return nil
// }
