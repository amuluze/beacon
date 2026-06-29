package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSanitizePathWithin_ValidPath 验证沙箱内路径通过校验。
func TestSanitizePathWithin_ValidPath(t *testing.T) {
	root := t.TempDir()
	// 沙箱内子路径
	p := filepath.Join(root, "data", "file.txt")
	if got, err := SanitizePathWithin(p, root); err != nil {
		t.Fatalf("expected sandbox path to pass, got err: %v", err)
	} else if !strings.HasPrefix(got, root) {
		t.Fatalf("returned path %q not under root %q", got, root)
	}
	// 等于 root 本身也应通过
	if _, err := SanitizePathWithin(root, root); err != nil {
		t.Fatalf("root itself should pass: %v", err)
	}
}

// TestSanitizePathWithin_RejectsEscape 验证沙箱外路径被拒绝（防越界读写系统文件）。
func TestSanitizePathWithin_RejectsEscape(t *testing.T) {
	root := t.TempDir()
	cases := []struct {
		name string
		path string
	}{
		{"etc passwd", "/etc/passwd"},
		{"root ssh", "/root/.ssh/authorized_keys"},
		{"sibling", filepath.Join(filepath.Dir(root), "other")},
		{"absolute outside", "/var/log/syslog"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := SanitizePathWithin(tc.path, root); err == nil {
				t.Fatalf("path %q outside root %q should be rejected", tc.path, root)
			}
		})
	}
}

// TestSanitizePathWithin_RejectsTraversal 验证 .. 穿越被拒绝。
func TestSanitizePathWithin_RejectsTraversal(t *testing.T) {
	root := t.TempDir()
	// /root/../etc 风格的穿越
	escape := filepath.Join(root, "..", "..", "etc", "passwd")
	if _, err := SanitizePathWithin(escape, root); err == nil {
		t.Fatal("traversal path should be rejected")
	}
}

// TestSanitizePathWithin_RejectsSymlinkEscape 验证符号链接逃逸被拒绝。
// 在沙箱内建一个指向 /etc 的软链，校验应失败。
func TestSanitizePathWithin_RejectsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	link := filepath.Join(root, "evil")
	// 指向沙箱外的目标（/etc 在大多数系统存在）
	if err := os.Symlink("/etc", link); err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}
	if _, err := SanitizePathWithin(link, root); err == nil {
		t.Fatal("symlink escaping sandbox should be rejected")
	}
}

// TestSanitizePathWithin_NonExistentPath 验证不存在的路径（创建场景）仍受沙箱约束：
// 沙箱内的不存在路径通过，沙箱外的不存在路径被拒。
func TestSanitizePathWithin_NonExistentPath(t *testing.T) {
	root := t.TempDir()
	// 沙箱内、尚不存在
	inside := filepath.Join(root, "newdir", "newfile.txt")
	if _, err := SanitizePathWithin(inside, root); err != nil {
		t.Fatalf("non-existent path inside sandbox should pass: %v", err)
	}
	// 沙箱外、尚不存在
	outside := "/etc/newly_created_evil.txt"
	if _, err := SanitizePathWithin(outside, root); err == nil {
		t.Fatal("non-existent path outside sandbox should be rejected")
	}
}

// TestSanitizePathWithin_InvalidInputs 验证空值、相对路径、空 root 报错。
func TestSanitizePathWithin_InvalidInputs(t *testing.T) {
	root := t.TempDir()
	cases := []struct {
		name string
		path string
		root string
	}{
		{"empty path", "", root},
		{"empty root", "/some/abs", ""},
		{"relative path", "relative/path", root},
		{"relative root", "/abs", "relative"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := SanitizePathWithin(tc.path, tc.root); err == nil {
				t.Fatalf("expected error for %s", tc.name)
			}
		})
	}
}
