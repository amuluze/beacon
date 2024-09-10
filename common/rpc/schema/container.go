// Package schema
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package schema

import "time"

type Docker struct {
	Timestamp     time.Time `json:"timestamp,omitempty"`
	DockerVersion string    `json:"docker_version,omitempty"`
	APIVersion    string    `json:"api_version,omitempty"`
	MinAPIVersion string    `json:"min_api_version,omitempty"`
	GitCommit     string    `json:"git_commit,omitempty"`
	GoVersion     string    `json:"go_version,omitempty"`
	Os            string    `json:"os,omitempty"`
	Arch          string    `json:"arch,omitempty"`
}

type Container struct {
	Timestamp    time.Time `json:"timestamp,omitempty"`
	ContainerID  string    `json:"container_id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Image        string    `json:"image,omitempty"`
	Network      string    `json:"network,omitempty"`
	IP           string    `json:"ip,omitempty"`
	Ports        string    `json:"ports,omitempty"`
	State        string    `json:"state,omitempty"`
	Uptime       string    `json:"uptime,omitempty"`
	Volumes      string    `json:"volumes,omitempty"`
	Environments string    `json:"environments,omitempty"`
	Commands     string    `json:"commands,omitempty"`
	Labels       string    `json:"labels,omitempty"`
	CPUPercent   float64   `json:"cpu_percent,omitempty"`
	MemPercent   float64   `json:"mem_percent,omitempty"`
	MemUsage     float64   `json:"mem_usage,omitempty"`
	MemLimit     float64   `json:"mem_limit,omitempty"`
}

type Image struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	ImageID   string    `json:"image_id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Tag       string    `json:"tag,omitempty"`
	Created   string    `json:"created,omitempty"`
	Size      string    `json:"size,omitempty"`
}

type Network struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	NetworkID string    `json:"network_id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Driver    string    `json:"driver,omitempty"`
	Scope     string    `json:"scope,omitempty"`
	Created   string    `json:"created,omitempty"`
	Internal  bool      `json:"internal,omitempty"`
	Subnet    string    `json:"subnet,omitempty"`
	Gateway   string    `json:"gateway,omitempty"`
	Labels    string    `json:"labels,omitempty"`
}

type DockerArgs struct {
	Timestamp time.Time `json:"timestamp"`
}

type DockerReply struct {
	Data Docker `json:"data"`
}

type ContainerQueryArgs struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type ContainerQueryReply struct {
	Data []Container `json:"data"`
}

type ContainerCountArgs struct{}

type ContainerCountReply struct {
	Count int `json:"count"`
}

type ImageQueryArgs struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type ImageQueryReply struct {
	Data []Image `json:"data"`
}

type ImageCountArgs struct{}

type ImageCountReply struct {
	Count int `json:"count"`
}

type NetworkQueryArgs struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type NetworkQueryReply struct {
	Data []Network `json:"data"`
}

type NetworkCountArgs struct{}

type NetworkCountReply struct {
	Count int `json:"count"`
}

type ContainerCreateArgs struct {
	Name         string            `json:"name"`
	Image        string            `json:"image"`
	Network      string            `json:"network"`
	Ports        []string          `json:"ports"`
	Volumes      []string          `json:"volumes"`
	Environments []string          `json:"environments"`
	Commands     []string          `json:"commands"`
	Labels       map[string]string `json:"labels"`
}

type ContainerCreateReply struct {
	ContainerID string `json:"container_id"`
}

type ContainerUpdateArgs struct {
	ContainerID  string            `json:"container_id"`
	Name         string            `json:"name"`
	Image        string            `json:"image"`
	Network      string            `json:"network"`
	Ports        []string          `json:"ports"`
	Volumes      []string          `json:"volumes"`
	Environments []string          `json:"environments"`
	Commands     []string          `json:"commands"`
	Labels       map[string]string `json:"labels"`
}

type ContainerUpdateReply struct {
	ContainerID string `json:"container_id"`
}

type ContainerDeleteArgs struct {
	ContainerID string `json:"container_id"`
}

type ContainerDeleteReply struct{}

type ContainerStartArgs struct {
	ContainerID string `json:"container_id"`
}

type ContainerStartReply struct{}

type ContainerStopArgs struct {
	ContainerID string `json:"container_id"`
}

type ContainerStopReply struct{}

type ContainerRestartArgs struct {
	ContainerID string `json:"container_id"`
}

type ContainerRestartReply struct{}

type ImagePullArgs struct {
	ImageName string `json:"image_name"`
}

type ImagePullReply struct{}

type ImageTagArgs struct {
	OldTag string `json:"old_tag"`
	NewTag string `json:"new_tag"`
}

type ImageTagReply struct{}

type ImageDeleteArgs struct {
	ImageID string `json:"image_id"`
}

type ImageDeleteReply struct{}

type ImageImportArgs struct {
	SourceFile string `json:"source_file"`
}

type ImageImportReply struct{}

type ImageExportArgs struct {
	ImageIDs   []string `json:"image_ids"`
	TargetFile string   `json:"target_file"`
}

type ImageExportReply struct{}

type NetworkCreateArgs struct {
	Name    string            `json:"name"`
	Driver  string            `json:"driver"`
	Subnet  string            `json:"subnet"`
	Gateway string            `json:"gateway"`
	Labels  map[string]string `json:"labels"`
}

type NetworkCreateReply struct {
	NetworkID string `json:"network_id"`
}

type NetworkDeleteArgs struct {
	NetworkID string `json:"network_id"`
}

type NetworkDeleteReply struct{}
