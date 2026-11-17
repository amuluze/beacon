package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsColliaDaemonCommand(t *testing.T) {
	tests := []struct {
		name string
		argv []string
		want bool
	}{
		{name: "systemd daemon", argv: []string{"/usr/sbin/collia", "daemon", "-c", "/etc/collia/config.yml"}, want: true},
		{name: "relative daemon", argv: []string{"collia", "daemon"}, want: true},
		{name: "run subcommand", argv: []string{"/usr/sbin/collia", "run"}, want: false},
		{name: "remove subcommand", argv: []string{"/usr/sbin/collia", "remove"}, want: false},
		{name: "shell wrapper", argv: []string{"bash", "-c", "/usr/sbin/collia daemon -c /etc/collia/config.yml"}, want: false},
		{name: "single arg", argv: []string{"collia"}, want: false},
		{name: "empty", argv: nil, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isColliaDaemonCommand(tt.argv); got != tt.want {
				t.Fatalf("isColliaDaemonCommand(%q) = %v, want %v", tt.argv, got, tt.want)
			}
		})
	}
}

// 重复安装时以 systemd 单元文件是否存在判定「已安装」，据此走幂等分支。
func TestUnitFileExists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "collia.service")
	if unitFileExists(path) {
		t.Fatal("missing unit file reported as existing")
	}
	if err := os.WriteFile(path, []byte("[Unit]\n"), 0o600); err != nil {
		t.Fatalf("write unit file: %v", err)
	}
	if !unitFileExists(path) {
		t.Fatal("existing unit file not detected")
	}
}
