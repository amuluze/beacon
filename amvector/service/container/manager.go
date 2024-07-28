// Package container
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package container

import "github.com/amuluze/docker"

type Manager struct {
	manager *docker.Manager
}

func NewContainerManager() *Manager {
	manager, err := docker.NewManager()
	if err != nil {
		panic(err)
	}
	return &Manager{manager: manager}
}
