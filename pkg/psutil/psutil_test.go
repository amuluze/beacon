// Package psutil
// Date: 2024/3/6 10:55
// Author: Amu
// Description:
package psutil

import (
	"context"
	"testing"

	"github.com/shirou/gopsutil/v3/common"
	"github.com/shirou/gopsutil/v3/disk"
)

func TestGetCPUPercent(t *testing.T) {
	cpuPercent, err := GetCPUPercent()
	t.Log(cpuPercent, err)
}

func TestGetMemInfo(t *testing.T) {
	memPercent, total, used, err := GetMemInfo()
	t.Log(memPercent, total, used, err)
}

func TestGetDiskPercent(t *testing.T) {
	devices := map[string]struct{}{"vda3": {}}
	diskMap, err := GetDiskInfo(devices)
	t.Log(diskMap, err)
}

func TestGetNetworkPercent(t *testing.T) {
	eth := map[string]struct{}{"eth0": {}}
	netMap, err := GetNetworkIO(eth)
	t.Log(netMap, err)
}

func TestGetDiskIO(t *testing.T) {
	devices := map[string]struct{}{"vda3": {}}
	diskMap, err := GetDiskIO(devices)
	t.Log(diskMap, err)
}

func TestGetAllDisk(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.EnvKey, common.EnvMap{})
	infos, _ := disk.PartitionsWithContext(ctx, false) // false表示只获取物理分区
	for _, info := range infos {
		t.Logf("%#v\n", info)
	}
}
