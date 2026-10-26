package rpc

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	rpcSchema "common/rpc/schema"
)

func withDaemonPath(t *testing.T, path string) {
	t.Helper()
	origin := DaemonJsonPath
	DaemonJsonPath = path
	t.Cleanup(func() { DaemonJsonPath = origin })
}

// 主机没有 /etc/docker/daemon.json 是常态：应返回空列表而不是报错。
func TestGetDockerRegistryMirrorsMissingFileReturnsEmpty(t *testing.T) {
	withDaemonPath(t, filepath.Join(t.TempDir(), "absent", "daemon.json"))

	reply := rpcSchema.GetDockerRegistryMirrorsReply{}
	svc := &Service{}
	if err := svc.GetDockerRegistryMirrors(context.Background(), rpcSchema.GetDockerRegistryMirrorsArgs{}, &reply); err != nil {
		t.Fatalf("expected no error for missing daemon.json, got %v", err)
	}
	if len(reply.Mirrors) != 0 {
		t.Fatalf("expected empty mirrors, got %v", reply.Mirrors)
	}
}

func TestGetDockerRegistryMirrorsReadsConfiguredMirrors(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "daemon.json")
	if err := os.WriteFile(path, []byte(`{"registry-mirrors":["https://mirror.example.com"]}`), 0600); err != nil {
		t.Fatalf("write daemon.json: %v", err)
	}
	withDaemonPath(t, path)

	reply := rpcSchema.GetDockerRegistryMirrorsReply{}
	svc := &Service{}
	if err := svc.GetDockerRegistryMirrors(context.Background(), rpcSchema.GetDockerRegistryMirrorsArgs{}, &reply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reply.Mirrors) != 1 || reply.Mirrors[0] != "https://mirror.example.com" {
		t.Fatalf("unexpected mirrors: %v", reply.Mirrors)
	}
}

// 配置文件损坏或字段类型异常时必须返回错误，而不是 panic。
func TestGetDockerRegistryMirrorsInvalidContentReturnsError(t *testing.T) {
	for name, content := range map[string]string{
		"broken json":  "{not-json",
		"not an array": `{"registry-mirrors":"https://mirror.example.com"}`,
		"not strings":  `{"registry-mirrors":[1,2]}`,
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "daemon.json")
			if err := os.WriteFile(path, []byte(content), 0600); err != nil {
				t.Fatalf("write daemon.json: %v", err)
			}
			withDaemonPath(t, path)

			reply := rpcSchema.GetDockerRegistryMirrorsReply{}
			svc := &Service{}
			if err := svc.GetDockerRegistryMirrors(context.Background(), rpcSchema.GetDockerRegistryMirrorsArgs{}, &reply); err == nil {
				t.Fatal("expected error for invalid daemon.json")
			}
		})
	}
}

// 写入时文件与目录都不存在：应创建并落盘；没有 systemd 时不应判定为失败。
func TestSetDockerRegistryMirrorsCreatesFileAndSurvivesNoSystemd(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "etc", "docker", "daemon.json")
	withDaemonPath(t, path)

	svc := &Service{}
	if err := svc.SetDockerRegistryMirrors(context.Background(), rpcSchema.SetDockerRegistryMirrorsArgs{
		Mirrors: []string{"https://mirror.example.com"},
	}, &rpcSchema.SetDockerRegistryMirrorsReply{}); err != nil {
		t.Fatalf("expected write to succeed without systemd, got %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read daemon.json: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected daemon.json to be written")
	}
}

// 已有配置时只应更新 registry-mirrors 字段，保留其它键。
func TestSetDockerRegistryMirrorsPreservesOtherKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "daemon.json")
	if err := os.WriteFile(path, []byte(`{"log-level":"debug"}`), 0600); err != nil {
		t.Fatalf("write daemon.json: %v", err)
	}
	withDaemonPath(t, path)

	svc := &Service{}
	if err := svc.SetDockerRegistryMirrors(context.Background(), rpcSchema.SetDockerRegistryMirrorsArgs{
		Mirrors: []string{"https://mirror.example.com"},
	}, &rpcSchema.SetDockerRegistryMirrorsReply{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read daemon.json: %v", err)
	}
	if len(data) == 0 || !containsAll(string(data), "log-level", "registry-mirrors") {
		t.Fatalf("expected both keys preserved, got %s", string(data))
	}
}

func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !contains(s, sub) {
			return false
		}
	}
	return true
}

func contains(s, sub string) bool {
	return len(sub) == 0 || len(s) >= len(sub) && (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})()
}
