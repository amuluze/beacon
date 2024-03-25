// Package schema
// Date: 2024/3/6 13:20
// Author: Amu
// Description:
package schema

import "time"

type Container struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Image         string  `json:"image"`
	IP            string  `json:"ip"`
	State         string  `json:"state"`
	Uptime        string  `json:"uptime"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
	MemoryUsage   float64 `json:"memory_usage"`
	MemoryLimit   float64 `json:"memory_limit"`
}

type ContainerQueryRely struct {
	Containers []Container `json:"containers"`
}

type Image struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Created string `json:"created"`
	Size    string `json:"size"`
}

type ImageQueryReply struct {
	Images []Image `json:"images"`
}

type Docker struct {
	Timestamp     time.Time
	DockerVersion string `json:"docker_version"`
	APIVersion    string `json:"api_version"`
	MinAPIVersion string `json:"min_api_version"`
	GitCommit     string `json:"git_commit"`
	GoVersion     string `json:"go_version"`
	Os            string `json:"os"`
	Arch          string `json:"arch"`
}
