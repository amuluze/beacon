// Package service
// Description: collia 安装/启动前的环境与磁盘预检。
//
// 探测真实块设备与网卡（取代 config.yml 里硬编码的 vda2/eth0），
// 校验关键目录可写性与控制端口可达性，并能把探测结果合并回写 config。
package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// 默认探测根与超时。抽出常量便于测试注入。
const (
	defaultSysRoot      = "/sys"
	defaultDialTimeout  = 3 * time.Second
	preflightProbeRetry = 1 // 写入 config 前不重试，保证快速失败
)

// 虚拟/无关块设备前缀，探测时排除。
var excludedBlockPrefixes = []string{"loop", "ram", "dm-", "md", "sr", "fd", "zram"}

// 虚拟/桥接网卡前缀与精确名，探测时排除。
var (
	excludedNetExact    = map[string]struct{}{"lo": {}}
	excludedNetPrefixes = []string{"docker", "veth", "br-", "virbr", "tap", "tun"}
)

// PreflightReport 汇总预检结果。Errors 非空表示阻断性问题；Warnings 为非阻断提示。
type PreflightReport struct {
	Disks         []string
	Ethernets     []string
	Errors        []string
	Warnings      []string
	ConfigUpdated bool
}

// HasError 报告是否存在阻断性错误。
func (r *PreflightReport) HasError() bool { return len(r.Errors) > 0 }

// ProbeDiskDevices 枚举 sysRoot/block 下非虚拟块设备的「分区」名（vda3、nvme0n1p1 等）。
// 采集目标必须是分区而非整盘：disk usage 基于 mountpoint（分区级），
// 整盘 vda 没有挂载点，拿不到使用率，会导致 m_disk 始终为空。
// 整盘没有分区表（裸盘直接挂载，如某些 sda）时回退为整盘名。
// sysRoot 为空时使用默认 /sys。非 Linux 环境或目录不存在时返回空切片。
func ProbeDiskDevices(sysRoot string) []string {
	if sysRoot == "" {
		sysRoot = defaultSysRoot
	}
	if runtime.GOOS != "linux" {
		if _, err := os.Stat(filepath.Join(sysRoot, "block")); err != nil {
			return nil
		}
	}
	entries, err := os.ReadDir(filepath.Join(sysRoot, "block"))
	if err != nil {
		return nil
	}
	var devices []string
	for _, e := range entries {
		diskName := e.Name()
		if isExcludedName(diskName, excludedBlockPrefixes, nil) {
			continue
		}
		if parts := probePartitions(sysRoot, diskName); len(parts) > 0 {
			devices = append(devices, parts...)
		} else {
			// 无分区子设备（裸盘挂载场景），回退整盘。
			devices = append(devices, diskName)
		}
	}
	sort.Strings(devices)
	return devices
}

// probePartitions 枚举 /sys/block/<disk>/ 下属于分区子设备的目录名。
// 判据：子目录下存在 partition 文件（内容为分区号），据此与 stat/holders/size 等
// 属性文件/目录区分，兼容 vda3 与 nvme0n1p1 两种命名。
func probePartitions(sysRoot, diskName string) []string {
	entries, err := os.ReadDir(filepath.Join(sysRoot, "block", diskName))
	if err != nil {
		return nil
	}
	var parts []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		partName := e.Name()
		if _, err := os.Stat(filepath.Join(sysRoot, "block", diskName, partName, "partition")); err != nil {
			continue
		}
		parts = append(parts, partName)
	}
	sort.Strings(parts)
	return parts
}

// ProbeEthernetNames 枚举 sysRoot/class/net 下的物理/业务网卡名。
func ProbeEthernetNames(sysRoot string) []string {
	if sysRoot == "" {
		sysRoot = defaultSysRoot
	}
	if runtime.GOOS != "linux" {
		if _, err := os.Stat(filepath.Join(sysRoot, "class", "net")); err != nil {
			return nil
		}
	}
	entries, err := os.ReadDir(filepath.Join(sysRoot, "class", "net"))
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		name := e.Name()
		if isExcludedName(name, excludedNetPrefixes, excludedNetExact) {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func isExcludedName(name string, prefixes []string, exact map[string]struct{}) bool {
	if _, ok := exact[name]; ok {
		return true
	}
	for _, p := range prefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

// dirWritable 检查目录存在且可写（创建一个临时文件验证）。
func dirWritable(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("stat %s: %w", dir, err)
		}
		// 父目录可写即可创建
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}
	} else if !info.IsDir() {
		return fmt.Errorf("%s 不是目录", dir)
	}
	f, err := os.CreateTemp(dir, ".collia-preflight-*")
	if err != nil {
		return fmt.Errorf("%s 不可写: %w", dir, err)
	}
	_ = f.Close()
	_ = os.Remove(f.Name())
	return nil
}

// serverReachable 对控制端口做一次 TCP 探测；不可达只产生警告，不阻断启动。
func serverReachable(server string, timeout time.Duration) bool {
	if server == "" {
		return false
	}
	conn, err := net.DialTimeout("tcp", server, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// Preflight 执行完整环境预检。需要校验的路径与端口从 config 派生。
// 返回的 report 同时携带探测到的设备/网卡，供调用方回写 config。
func Preflight(cfg *Config, sysRoot string) *PreflightReport {
	report := &PreflightReport{
		Disks:     ProbeDiskDevices(sysRoot),
		Ethernets: ProbeEthernetNames(sysRoot),
	}

	// 关键目录可写性
	logDir := cfg.Log.Output
	if logDir != "" {
		logDir = filepath.Dir(logDir)
		if err := dirWritable(logDir); err != nil {
			report.Errors = append(report.Errors, err.Error())
		}
	}
	if cfg.DB.DBType == "sqlite" || cfg.DB.DBType == "" {
		if cfg.DB.DBName != "" {
			dbDir := filepath.Dir(cfg.DB.DBName)
			if err := dirWritable(dbDir); err != nil {
				report.Errors = append(report.Errors, err.Error())
			}
		}
	}

	// 控制端口可达性（仅警告）
	if cfg.Control.Server != "" {
		if !serverReachable(cfg.Control.Server, defaultDialTimeout) {
			report.Warnings = append(report.Warnings,
				fmt.Sprintf("控制端口 %s 暂不可达，将在后台重试连接", cfg.Control.Server))
		}
	}

	// 探测结果充分性（仅警告）
	if len(report.Disks) == 0 {
		report.Warnings = append(report.Warnings, "未探测到块设备，磁盘监控可能无数据")
	}
	if len(report.Ethernets) == 0 {
		report.Warnings = append(report.Warnings, "未探测到网卡，网络监控可能无数据")
	}
	return report
}

// MergeProbeIntoConfig 把探测到的块设备/网卡合并进 configPath 指向的 YAML：
// 仅当 task.disk.devices 或 task.ethernet.names 为空或缺省时填充，不覆盖用户显式配置。
// 成功修改并落盘时返回 ConfigUpdated=true。
func MergeProbeIntoConfig(configPath string, report *PreflightReport) error {
	if configPath == "" {
		return errors.New("config path 为空")
	}
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	root := make(map[string]any)
	if err := yaml.Unmarshal(raw, &root); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	changed := false
	taskNode, ok := root["task"].(map[string]any)
	if !ok {
		taskNode = make(map[string]any)
		root["task"] = taskNode
	}

	if len(report.Disks) > 0 {
		diskNode, ok := taskNode["disk"].(map[string]any)
		if !ok {
			diskNode = make(map[string]any)
			taskNode["disk"] = diskNode
		}
		if !hasNonEmptyList(diskNode, "devices") {
			diskNode["devices"] = report.Disks
			changed = true
		}
	}
	if len(report.Ethernets) > 0 {
		ethNode, ok := taskNode["ethernet"].(map[string]any)
		if !ok {
			ethNode = make(map[string]any)
			taskNode["ethernet"] = ethNode
		}
		if !hasNonEmptyList(ethNode, "names") {
			ethNode["names"] = report.Ethernets
			changed = true
		}
	}

	if !changed {
		return nil
	}
	out, err := yaml.Marshal(root)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	tmp := configPath + ".tmp"
	if err := os.WriteFile(tmp, out, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	if err := os.Rename(tmp, configPath); err != nil {
		return fmt.Errorf("rename config: %w", err)
	}
	report.ConfigUpdated = true
	slog.Info("preflight 已将探测结果合并进 config", "config", configPath,
		"disks", report.Disks, "ethernets", report.Ethernets)
	return nil
}

// hasNonEmptyList 报告 node[key] 是否为非空列表。
func hasNonEmptyList(node map[string]any, key string) bool {
	v, ok := node[key]
	if !ok {
		return false
	}
	list, ok := v.([]any)
	if !ok {
		return false
	}
	return len(list) > 0
}

// LoadConfig 以默认 prefix 读取配置，供 cmd/main 在 preflight 时使用。
func LoadConfig(configFile string) (*Config, error) {
	return NewConfig(configFile, Prefix(defaultPrefix))
}

// PrintPreflight 把预检报告输出到 stderr，供 `collia probe` 等命令展示。
func PrintPreflight(r *PreflightReport) {
	fmt.Fprintf(os.Stderr, "preflight:\n")
	fmt.Fprintf(os.Stderr, "  disks:     %v\n", r.Disks)
	fmt.Fprintf(os.Stderr, "  ethernets: %v\n", r.Ethernets)
	for _, w := range r.Warnings {
		fmt.Fprintf(os.Stderr, "  warn: %s\n", w)
	}
	for _, e := range r.Errors {
		fmt.Fprintf(os.Stderr, "  ERROR: %s\n", e)
	}
	if r.ConfigUpdated {
		fmt.Fprintf(os.Stderr, "  config: 已合并探测结果\n")
	}
}
