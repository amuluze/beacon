// Package schema
// Date: 2024/3/6 13:20
// Author: Amu
// Description:
package schema

import "time"

type Container struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	IP            string `json:"ip"`
	State         string `json:"state"`
	Uptime        string `json:"uptime"`
	CPUPercent    string `json:"cpu_percent"`
	MemoryPercent string `json:"memory_percent"`
	MemoryUsage   string `json:"memory_usage"`
	MemoryLimit   string `json:"memory_limit"`
}

type ContainerQueryArgs struct {
	Page int `json:"page" validate:"required"`
	Size int `json:"size" validate:"gte=0"`
}

type ContainerQueryRely struct {
	Data  []Container `json:"data"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

type Image struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Created string `json:"created"`
	Size    string `json:"size"`
}

type ImageQueryArgs struct {
	Page int `json:"page" validate:"required"`
	Size int `json:"size" validate:"gte=0"`
}

type ImageQueryReply struct {
	Data  []Image `json:"data"`
	Total int     `json:"total"`
	Page  int     `json:"page"`
	Size  int     `json:"size"`
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
