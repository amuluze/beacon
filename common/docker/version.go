// Package docker
// Date: 2024/07/09 14:19:40
// Author: Amu
// Description:
package docker

import "context"

type Version struct {
	DockerVersion string `json:"docker_version"`
	APIVersion    string `json:"api_version"`
	MinAPIVersion string `json:"min_api_version"`
	GitCommit     string `json:"git_commit"`
	GoVersion     string `json:"go_version"`
	OS            string `json:"os"`
	Arch          string `json:"arch"`
}

func (m *Manager) Version(ctx context.Context) (*Version, error) {
	serverVersion, err := m.client.ServerVersion(ctx)
	if err != nil {
		return nil, err
	}

	return &Version{
		DockerVersion: serverVersion.Version,
		APIVersion:    serverVersion.APIVersion,
		MinAPIVersion: serverVersion.MinAPIVersion,
		GitCommit:     serverVersion.GitCommit,
		GoVersion:     serverVersion.GoVersion,
		OS:            serverVersion.Os,
		Arch:          serverVersion.Arch,
	}, nil
}
