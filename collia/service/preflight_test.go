package service

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// writeFakeSys 在临时目录构造 /sys 的子集，模拟块设备与网卡枚举。
func writeFakeSys(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	blockDir := filepath.Join(root, "block")
	netDir := filepath.Join(root, "class", "net")
	if err := os.MkdirAll(blockDir, 0o755); err != nil {
		t.Fatalf("mkdir block: %v", err)
	}
	if err := os.MkdirAll(netDir, 0o755); err != nil {
		t.Fatalf("mkdir net: %v", err)
	}
	for _, d := range []string{"sda", "sdb", "loop0", "ram0", "dm-0", "sr0"} {
		if err := os.MkdirAll(filepath.Join(blockDir, d), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", d, err)
		}
	}
	for _, n := range []string{"eth0", "eth1", "lo", "docker0", "veth1234", "br-abcdef"} {
		if err := os.MkdirAll(filepath.Join(netDir, n), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", n, err)
		}
	}
	return root
}

func TestProbeDiskDevicesFiltersVirtual(t *testing.T) {
	root := writeFakeSys(t)
	got := ProbeDiskDevices(root)
	want := []string{"sda", "sdb"}
	if !equalStrings(got, want) {
		t.Fatalf("ProbeDiskDevices = %v, want %v", got, want)
	}
}

// TestProbeDiskDevicesReturnsPartitions 验证采集目标是分区（vda3）而非整盘（vda）：
// 仅含 partition 文件的子目录被识别为分区；stat/holders 等属性条目被忽略。
func TestProbeDiskDevicesReturnsPartitions(t *testing.T) {
	root := t.TempDir()
	blockDir := filepath.Join(root, "block")
	if err := os.MkdirAll(filepath.Join(blockDir, "vda"), 0o755); err != nil {
		t.Fatalf("mkdir vda: %v", err)
	}
	for _, p := range []string{"vda1", "vda2", "vda3"} {
		dir := filepath.Join(blockDir, "vda", p)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", p, err)
		}
		if err := os.WriteFile(filepath.Join(dir, "partition"), []byte(p), 0o644); err != nil {
			t.Fatalf("write partition marker: %v", err)
		}
	}
	// 干扰项：同目录下的非分区条目应被忽略。
	if err := os.MkdirAll(filepath.Join(blockDir, "vda", "holders"), 0o755); err != nil {
		t.Fatalf("mkdir holders: %v", err)
	}
	if err := os.WriteFile(filepath.Join(blockDir, "vda", "stat"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write stat: %v", err)
	}

	got := ProbeDiskDevices(root)
	want := []string{"vda1", "vda2", "vda3"}
	if !equalStrings(got, want) {
		t.Fatalf("ProbeDiskDevices = %v, want %v", got, want)
	}
}

func TestProbeEthernetNamesFiltersVirtual(t *testing.T) {
	root := writeFakeSys(t)
	got := ProbeEthernetNames(root)
	want := []string{"eth0", "eth1"}
	if !equalStrings(got, want) {
		t.Fatalf("ProbeEthernetNames = %v, want %v", got, want)
	}
}

func TestProbeDiskDevicesMissingDir(t *testing.T) {
	if got := ProbeDiskDevices(t.TempDir()); len(got) != 0 {
		t.Fatalf("expected empty on missing /sys/block, got %v", got)
	}
}

func TestMergeProbeIntoConfigFillsEmpty(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yml")
	original := "task:\n  disk:\n    devices: []\n  ethernet:\n    names: []\n"
	if err := os.WriteFile(cfgPath, []byte(original), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	report := &PreflightReport{Disks: []string{"sda"}, Ethernets: []string{"eth0"}}
	if err := MergeProbeIntoConfig(cfgPath, report); err != nil {
		t.Fatalf("merge: %v", err)
	}
	if !report.ConfigUpdated {
		t.Fatal("ConfigUpdated should be true after merge")
	}
	// 再合并一次：已有非空列表，不应重复修改
	report2 := &PreflightReport{Disks: []string{"sda"}, Ethernets: []string{"eth0"}}
	if err := MergeProbeIntoConfig(cfgPath, report2); err != nil {
		t.Fatalf("merge2: %v", err)
	}
	if report2.ConfigUpdated {
		t.Fatal("ConfigUpdated should be false when lists already non-empty")
	}
}

func TestMergeProbeIntoConfigKeepsUserValues(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yml")
	original := "task:\n  disk:\n    devices:\n      - nvme0n1\n  ethernet:\n    names:\n      - ens33\n"
	if err := os.WriteFile(cfgPath, []byte(original), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	report := &PreflightReport{Disks: []string{"sda"}, Ethernets: []string{"eth0"}}
	if err := MergeProbeIntoConfig(cfgPath, report); err != nil {
		t.Fatalf("merge: %v", err)
	}
	if report.ConfigUpdated {
		t.Fatal("不应覆盖用户显式配置")
	}
}

func TestPreflightReportHasError(t *testing.T) {
	r := &PreflightReport{Errors: []string{"x"}}
	if !r.HasError() {
		t.Fatal("HasError should be true")
	}
	r2 := &PreflightReport{Warnings: []string{"y"}}
	if r2.HasError() {
		t.Fatal("HasError should be false with only warnings")
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	ac := append([]string(nil), a...)
	bc := append([]string(nil), b...)
	sort.Strings(ac)
	sort.Strings(bc)
	for i := range ac {
		if ac[i] != bc[i] {
			return false
		}
	}
	return true
}
