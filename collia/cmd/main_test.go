package main

import "testing"

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
