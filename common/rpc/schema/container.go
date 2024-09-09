// Package rpc
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package schema

import "time"

type Docker struct {
	Timestamp     time.Time `json:"timestamp"`
	DockerVersion string    `json:"docker_version"`
	APIVersion    string    `json:"api_version"`
	MinAPIVersion string    `json:"min_api_version"`
	GitCommit     string    `json:"git_commit"`
	GoVersion     string    `json:"go_version"`
	Os            string    `json:"os"`
	Arch          string    `json:"arch"`
}

type Container struct {
	Timestamp    time.Time `json:"timestamp"`
	ContainerID  string    `json:"container_id"`
	Name         string    `json:"name"`
	Image        string    `json:"image"`
	Network      string    `json:"network"`
	IP           string    `json:"ip"`
	Ports        string    `json:"ports"`
	State        string    `json:"state"`
	Uptime       string    `json:"uptime"`
	Volumes      string    `json:"volumes"`
	Environments string    `json:"environments"`
	Commands     string    `json:"commands"`
	Labels       string    `json:"labels"`
	CPUPercent   float64   `json:"cpu_percent"`
	MemPercent   float64   `json:"mem_percent"`
	MemUsage     float64   `json:"mem_usage"`
	MemLimit     float64   `json:"mem_limit"`
}

type Image struct {
	Timestamp time.Time `json:"timestamp"`
	ImageID   string    `json:"image_id"`
	Name      string    `json:"name"`
	Tag       string    `json:"tag"`
	Created   string    `json:"created"`
	Size      string    `json:"size"`
}

type Network struct {
	Timestamp time.Time `json:"timestamp"`
	NetworkID string    `json:"network_id"`
	Name      string    `json:"name"`
	Driver    string    `json:"driver"`
	Scope     string    `json:"scope"`
	Created   string    `json:"created"`
	Internal  bool      `json:"internal"`
	Subnet    string    `json:"subnet"`
	Gateway   string    `json:"gateway"`
	Labels    string    `json:"labels"`
}

type DockerSummaryArgs struct {
	Timestamp time.Time `json:"timestamp"`
}

type DockerSummaryReply struct {
	Data Docker `json:"data"`
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

type ContainerDeleteArgs struct {
	ContainerID string `json:"container_id"`
}

type ContainerDeleteReply struct{}

type ContainerSummaryArgs struct {
	Timestamp time.Time `json:"timestamp"`
}

type ContainerSummaryReply struct {
	Data []Container `json:"data"`
}

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

type ImageSummaryArgs struct {
	Timestamp time.Time `json:"timestamp"`
}

type ImageSummaryReply struct {
	Data []Image `json:"data"`
}

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

type NetworkUpdateArgs struct {
	NetworkID string            `json:"network_id"`
	Name      string            `json:"name"`
	Driver    string            `json:"driver"`
	Subnet    string            `json:"subnet"`
	Gateway   string            `json:"gateway"`
	Labels    map[string]string `json:"labels"`
}

type NetworkUpdateReply struct{}

type NetworkDeleteArgs struct {
	NetworkID string `json:"network_id"`
}

type NetworkDeleteReply struct{}

type NetworkSummaryArgs struct {
	Timestamp time.Time `json:"timestamp"`
}

type NetworkSummaryReply struct {
	Data []Network `json:"data"`
}
