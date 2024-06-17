// Package psutil
// Date: 2024/06/10 18:58:51
// Author: Amu
// Description:
package psutil

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/common"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type DiskInfo struct {
	Total   uint64
	Percent float64
	Used    uint64
}

type NetIO struct {
	Recv uint64 `json:"recv"`
	Send uint64 `json:"send"`
}

type DiskIO struct {
	Read  uint64 `json:"read"`
	Write uint64 `json:"write"`
}

type SystemInfo struct {
	Uptime          string
	Hostname        string
	Os              string
	Platform        string
	PlatformVersion string
	KernelVersion   string
	KernelArch      string
}

func GetMemInfo() (float64, uint64, uint64, error) {
	ctx := context.WithValue(context.Background(), common.EnvKey, common.EnvMap{})
	v, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return 0.0, 0, 0, err
	}
	return v.UsedPercent, v.Total, v.Used, nil
}

func GetCPUPercent() (float64, error) {
	ctx := context.WithValue(context.Background(), common.EnvKey, common.EnvMap{})
	totalPercent, err := cpu.PercentWithContext(ctx, 3*time.Second, false)
	if err != nil {
		return 0.0, err
	}
	return totalPercent[0], nil
}

func GetDiskInfo(devices map[string]struct{}) (map[string]DiskInfo, error) {
	ctx := context.WithValue(context.Background(), common.EnvKey, common.EnvMap{})
	diskMap := make(map[string]DiskInfo)
	infos, _ := disk.PartitionsWithContext(ctx, false)
	slog.Info("disk infos", "infos", infos)
	slog.Info("devices", "devices", devices)
	for _, info := range infos {
		usedInfo, err := disk.UsageWithContext(ctx, info.Mountpoint)
		if usedInfo == nil {
			slog.Error("get disk usage err", "err", err)
			continue
		}
		slog.Info("used info", "usedInfo", usedInfo)
		dev := strings.Split(info.Device, "/")[len(strings.Split(info.Device, "/"))-1]
		slog.Info(">>>>>>>>>>>>>>>>>>>> dev", "dev", dev)
		if _, ok := devices[dev]; !ok {
			continue
		}
		diskMap[dev] = DiskInfo{
			Total:   usedInfo.Total,
			Percent: usedInfo.UsedPercent,
			Used:    usedInfo.Used,
		}
	}
	return diskMap, nil
}

func GetDiskIO(devices map[string]struct{}) (map[string]DiskIO, error) {
	ctx := context.WithValue(context.Background(), common.EnvKey, common.EnvMap{})
	diskMap := make(map[string]DiskIO)
	// 实现磁盘IO的获取逻辑
	var names []string
	for k := range devices {
		names = append(names, k)
	}
	slog.Error("device names", "names", names)
	stat, err := disk.IOCountersWithContext(ctx, names...)
	if err != nil {
		slog.Error("get disk io err", "err", err)
		return diskMap, err
	}
	// 将获取到的磁盘IO信息填充到map中
	for k, v := range stat {
		slog.Info("disk stat value", "value", v, "key", k)
		if _, ok := devices[k]; !ok {
			continue
		}
		diskMap[k] = DiskIO{
			Read:  v.ReadBytes,
			Write: v.WriteBytes,
		}
	}
	return diskMap, nil
}

func GetNetworkIO(eth map[string]struct{}) (map[string]NetIO, error) {
	ctx := context.WithValue(context.Background(), common.EnvKey, common.EnvMap{})

	netMap := make(map[string]NetIO)
	IOCountersStat, err := net.IOCountersWithContext(ctx, true)
	if err != nil {
		slog.Error("get network io err", "err", err)
		return netMap, err
	}
	slog.Error("eth", "eth", eth)
	slog.Error("IOCountersStat", "IOCountersStat", IOCountersStat)
	for _, stat := range IOCountersStat {
		if _, ok := eth[stat.Name]; !ok {
			continue
		}
		netMap[stat.Name] = NetIO{
			Recv: stat.BytesRecv,
			Send: stat.BytesSent,
		}
	}
	return netMap, nil
}

// GetSystemInfo 获取系统信息
func GetSystemInfo() (*SystemInfo, error) {
	ctx := context.WithValue(context.Background(), common.EnvKey, common.EnvMap{})

	info, err := host.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}
	systemInfo := &SystemInfo{
		Uptime:          time.Unix(int64(info.BootTime), 0).Local().Format("2006-01-02 15:04:05"),
		Hostname:        info.HostID,
		Os:              info.OS,
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		KernelArch:      info.KernelArch,
	}
	return systemInfo, nil
}

// 字节单位转换
func fmtByte(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%.2fB", float64(size)/float64(1))
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.2fKB", float64(size)/float64(1024))
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2fMB", float64(size)/float64(1024*1024))
	} else if size < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2fGB", float64(size)/float64(1024*1024*1024))
	} else if size < 1024*1024*1024*1024*1024 {
		return fmt.Sprintf("%.2fTB", float64(size)/float64(1024*1024*1024*1024))
	} else {
		return fmt.Sprintf("%.2fEB", float64(size)/float64(1024*1024*1024*1024*2014))
	}
}
