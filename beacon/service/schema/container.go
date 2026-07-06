// Package schema
// Date: 2024/3/6 13:20
// Author: Amu
// Description:
package schema

import "time"

type Container struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Image         string            `json:"image"`
	IP            string            `json:"ip"`
	Ports         string            `json:"ports"`
	ServerType    string            `json:"server_type"`
	State         string            `json:"state"`
	Uptime        string            `json:"uptime"`
	CPUPercent    string            `json:"cpu_percent"`
	MemoryPercent string            `json:"memory_percent"`
	MemoryUsage   string            `json:"memory_usage"`
	MemoryLimit   string            `json:"memory_limit"`
	Labels        map[string]string `json:"labels"`
}

type ContainerQueryArgs struct {
	Page int `json:"page" validate:"required,gte=1"`
	Size int `json:"size" validate:"gte=0,lte=100"`
}

type QueryCountArgs struct{}

type QueryCountReply struct {
	Count int `json:"count"`
}

type ContainerCreateArgs struct {
	ContainerName string            `json:"container_name" validate:"required,gte=1,lte=128"`
	ImageName     string            `json:"image_name" validate:"required,gte=1,lte=256"`
	NetworkMode   string            `json:"network_mode" validate:"required,oneof=bridge host none"`
	NetworkID     string            `json:"network_id" validate:"omitempty,len=36"`
	NetworkName   string            `json:"network_name" validate:"lte=128"`
	Ports         []string          `json:"ports" validate:"max=64"`
	Volumes       []string          `json:"volumes" validate:"max=64"`
	Environments  []string          `json:"environments" validate:"max=64"`
	Labels        map[string]string `json:"labels" validate:"max=32"`
}

type ContainerCreateReply struct {
	ContainerID string `json:"container_id"`
}

type ContainerUpdateArgs struct {
	ContainerID   string            `json:"container_id" validate:"required,len=64"`
	ContainerName string            `json:"container_name,omitempty" validate:"lte=128"`
	ImageName     string            `json:"image_name,omitempty" validate:"lte=256"`
	NetworkName   string            `json:"network_name,omitempty" validate:"lte=128"`
	Ports         []string          `json:"ports,omitempty" validate:"max=64"`
	Volumes       []string          `json:"volumes,omitempty" validate:"max=64"`
	Environments  []string          `json:"environment,omitempty" validate:"max=64"`
	Labels        map[string]string `json:"labels,omitempty" validate:"max=32"`
}

type ContainerUpdateReply struct {
	ContainerID string `json:"container_id"`
}

type ContainerDeleteArgs struct {
	ContainerID string `json:"container_id" validate:"required"`
}

type ContainerStartArgs struct {
	ContainerID string `json:"container_id" validate:"required"`
}

type ContainerStartReply struct{}

type ContainerStopArgs struct {
	ContainerID string `json:"container_id" validate:"required"`
}

type ContainerStopReply struct{}

type ContainerRestartArgs struct {
	ContainerID string `json:"container_id" validate:"required"`
}

type ContainerRestartReply struct{}

type ContainerQueryRely struct {
	Data      []Container `json:"data"`
	Total     int         `json:"total"`
	Page      int         `json:"page"`
	Size      int         `json:"size"`
	Freshness Freshness   `json:"freshness"`
}

type Image struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Created string `json:"created"`
	Size    string `json:"size"`
	Number  int    `json:"number"`
}

type ImageQueryArgs struct {
	Page int `json:"page" validate:"required,gte=1"`
	Size int `json:"size" validate:"gt=0,lte=100"`
}

type ImageQueryReply struct {
	Data      []Image   `json:"data"`
	Total     int       `json:"total"`
	Page      int       `json:"page"`
	Size      int       `json:"size"`
	Freshness Freshness `json:"freshness"`
}

type ImageTagArgs struct {
	OldTag string `json:"old_tag"`
	NewTag string `json:"new_tag"`
}

type ImagePullArgs struct {
	ImageName string `json:"image_name" validate:"required,gte=1,lte=256"`
}

type ImagePullReply struct{}

type ImageImportArgs struct {
	SourceFile string `json:"source_file"`
	FileName   string `json:"file_name"`
	Data       []byte `json:"-"`
}

type ImageImportReply struct{}

type ImageExportArgs struct {
	ImageName string `json:"image_names" validate:"required"`
	ImageID   string `json:"image_ids" validate:"required"`
}

type ImageExportRPCArgs struct {
	ImageIDs   []string `json:"image_ids"`
	TargetFile string   `json:"target_file"`
}

type ImageExportReply struct {
	FileName string `json:"file_name"`
	Data     []byte `json:"-"`
}

type ImageDeleteArgs struct {
	ImageID string `json:"image_id" validate:"required"`
}

type ImageDeleteReply struct{}

type ImageCountArgs struct{}

type ImageCountReply struct {
	Count int `json:"count"`
}

type NetworkCreateArgs struct {
	Name    string            `json:"name" validate:"required,gte=1,lte=128"`
	Driver  string            `json:"driver" validate:"required,oneof=bridge host overlay macvlan none"`
	Subnet  string            `json:"subnet" validate:"required,gte=1,lte=64"`
	Gateway string            `json:"gateway" validate:"required,gte=1,lte=64"`
	Labels  map[string]string `json:"labels" validate:"max=32"`
}

type NetworkCreateReply struct {
	NetworkID string `json:"network_id"`
}

type NetworkQueryArgs struct {
	Page int `json:"page" validate:"required,gte=1"`
	Size int `json:"size" validate:"gte=0,lte=100"`
}

type Network struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Driver  string            `json:"driver"`
	Created string            `json:"created"`
	Subnet  string            `json:"subnet"`
	Gateway string            `json:"gateway"`
	Labels  map[string]string `json:"labels"`
}

type NetworkQueryReply struct {
	Data      []Network `json:"data"`
	Total     int       `json:"total"`
	Page      int       `json:"page"`
	Size      int       `json:"size"`
	Freshness Freshness `json:"freshness"`
}

type NetworkCountArgs struct{}

type NetworkCountReply struct {
	Count int `json:"count"`
}

type NetworkDeleteArgs struct {
	NetworkID string `json:"network_id" validate:"required"`
}

type NetworkDeleteReply struct {
}

type Docker struct {
	Timestamp     time.Time
	DockerVersion string    `json:"docker_version"`
	APIVersion    string    `json:"api_version"`
	MinAPIVersion string    `json:"min_api_version"`
	GitCommit     string    `json:"git_commit"`
	GoVersion     string    `json:"go_version"`
	Os            string    `json:"os"`
	Arch          string    `json:"arch"`
	Freshness     Freshness `json:"freshness"`
}

type GetDockerRegistryMirrorsArgs struct {
}

type GetDockerRegistryMirrorsReply struct {
	Mirrors []string `json:"registry_mirrors"`
}

type SetDockerRegistryMirrorsArgs struct {
	Mirrors []string `json:"registry_mirrors" validate:"required,max=8"`
}

type SetDockerRegistryMirrorsReply struct{}

type ContainerUsageArgs struct {
	StartTime int64 `query:"start_time"`
	EndTime   int64 `query:"end_time"`
}

type ContainerUsageReply struct {
	Names     []string           `json:"names"`
	CPUUsage  map[string][]Usage `json:"cpu_usage"`
	MemUsage  map[string][]Usage `json:"mem_usage"`
	Freshness Freshness          `json:"freshness"`
}
