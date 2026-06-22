// Package psutil
// Date: 2024/06/10 18:59:19
// Author: Amu
// Description:
package psutil

import "testing"

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
