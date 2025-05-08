// Package model
// Date: 2024/06/10 18:44:53
// Author: Amu
// Description:
package model

import (
	"time"
)

type Container struct {
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

type Docker struct {
	Timestamp     time.Time
	DockerVersion string
	APIVersion    string
	MinAPIVersion string
	GitCommit     string
	GoVersion     string
	Os            string
	Arch          string
}

type Image struct {
	Timestamp time.Time
	ImageID   string
	Name      string
	Tag       string
	Created   string
	Size      string
	Number    int
}

type Network struct {
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
