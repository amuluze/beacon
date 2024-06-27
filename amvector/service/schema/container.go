// Package schema
// Date: 2024/06/11 00:14:42
// Author: Amu
// Description:
package schema

type VersionArgs struct {
}

type ContainerQueryArgs struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type QueryCountArgs struct{}

type QueryCountReply struct {
	Count int
}

type ContainerCreateArgs struct {
}

type ContainerCreateReply struct {
}

type ContainerDeleteArgs struct {
	ContainerID string
}

type ContainerDeleteReply struct {
}

type ContainerStartArgs struct {
	ContainerID string
}

type ContainerStartReply struct {
}

type ContainerStopArgs struct {
	ContainerID string
}

type ContainerStopReply struct {
}

type ContainerRestartArgs struct {
	ContainerID string
}

type ContainerRestartReply struct {
}

type ImageQueryArgs struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type ImagePullArgs struct {
	ImageName string
}

type ImagePullReply struct {
}

type ImageTagArgs struct {
	OldTag string
	NewTag string
}

type ImageTagReply struct {
}

type ImageCountArgs struct{}

type ImageCountReply struct {
	Count int
}

type ImageDeleteArgs struct {
	ImageID string
}

type ImageDeleteReply struct {
}

type ImageImportArgs struct {
	SourceFile string
}

type ImageImportReply struct{}

type ImageExportArgs struct {
	ImageIDs   []string
	TargetFile string
}

type ImageExportReply struct{}

type NetworkCreateArgs struct {
	Name           string
	Driver         string
	NetworkSegment string
	Labels         map[string]string
}

type NetworkCreateReply struct {
	NetworkID string
}

type NetworkQueryArgs struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type NetworkCountArgs struct{}

type NetworkCountReply struct {
	Count int
}

type NetworkDeleteArgs struct {
	NetworkID string
}

type NetworkDeleteReply struct {
}
