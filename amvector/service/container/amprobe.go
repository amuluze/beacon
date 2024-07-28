// Package container
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package container

import (
	"context"
	"fmt"
	"github.com/amuluze/amprobe/amvector/pkg/profile"
)

func (m *Manager) CreateAmprobe(profile *profile.AmprobeProfile) error {
	// 同名 container 存在则跳过
	container, err := m.manager.HasSameNameContainer(context.TODO(), profile.ContainerName)
	if err != nil {
		return err
	} else if container {
		return nil
	}

	// 创建 container
	var vols []string
	for _, volume := range profile.Volumes.Volumes {
		var vol string
		if volume.AccessMode == "ro" {
			vol = fmt.Sprintf("%s:%s:%s", volume.Source, volume.Destination, volume.AccessMode)
		} else {
			vol = fmt.Sprintf("%s:%s:rw", volume.Source, volume.Destination)
		}
		vols = append(vols, vol)
	}
	cID, err := m.manager.CreateContainer(
		context.TODO(),
		profile.ContainerName,
		profile.Image,
		profile.Networks.Networks[0].Name,
		profile.Ports,
		vols,
		nil,
	)
	if err != nil {
		return err
	}
	// 启动容器
	err = m.manager.StartContainer(context.TODO(), cID)
	if err != nil {
		return err
	}
	return nil
}
