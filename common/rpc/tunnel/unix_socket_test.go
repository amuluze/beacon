package tunnel

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestTunnelOverUnixSocket 验证单机模式：Server 在 Unix Socket 上监听，
// Agent 经 unix:// 目标注册并完成一次双向 RPC。
func TestTunnelOverUnixSocket(t *testing.T) {
	// macOS 的 sun_path 上限为 104 字节，t.TempDir() 在 macOS 上过长；
	// 显式使用 /tmp 下的短目录，保证跨平台可监听。
	dir, err := os.MkdirTemp("/tmp", "beacon-uds-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	sock := filepath.Join(dir, "control", "collia.sock")

	server := NewServerTunnel()
	server.SetJoinToken("join-token")
	started := make(chan error, 1)
	go func() { started <- server.Start("unix://" + sock) }()
	t.Cleanup(server.Stop)

	// 等待监听就绪（socket 文件出现且可被拨号）；启动失败立即报错。
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		select {
		case err := <-started:
			t.Fatalf("server start: %v", err)
		default:
		}
		if _, err := os.Stat(sock); err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if _, err := os.Stat(sock); err != nil {
		t.Fatalf("unix socket not created: %v", err)
	}
	if info, err := os.Stat(sock); err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("unix socket permissions = %v, want 0600", info.Mode().Perm())
	}

	agent := NewAgentTunnel("unix://"+sock, "agent-uds")
	agent.SetJoinToken("join-token")
	agent.SetHandler(func(context.Context, string, []byte, func(*Frame)) ([]byte, error) {
		return []byte(`{"ok":true}`), nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
		_ = agent.Close()
	})
	go func() { _ = agent.Start(ctx) }()

	for time.Now().Before(deadline) && server.AgentCount() == 0 {
		time.Sleep(20 * time.Millisecond)
	}
	if server.AgentCount() != 1 {
		t.Fatal("agent did not register over the unix socket")
	}

	callCtx, callCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer callCancel()
	var reply struct {
		OK bool `json:"ok"`
	}
	if err := server.Call(callCtx, "agent-uds", "Ping", map[string]string{}, &reply); err != nil {
		t.Fatalf("rpc over unix socket failed: %v", err)
	}
	if !reply.OK {
		t.Fatal("unexpected rpc reply over unix socket")
	}

	// 先结束 Agent 连接再停止服务端：GracefulStop 会等待活动流结束，
	// 顺序颠倒会让测试挂起。
	cancel()
	_ = agent.Close()
	server.Stop()
	if _, err := os.Stat(sock); !os.IsNotExist(err) {
		t.Fatalf("unix socket must be removed on stop, stat err = %v", err)
	}
}

// TestAgentDialTarget 验证 unix 目标被改写并注入拨号器，非 unix 目标直接报错。
func TestAgentDialTarget(t *testing.T) {
	if _, _, err := agentDialTarget("beacon.example.com:17000"); err == nil {
		t.Fatal("tcp target must be rejected in single-machine mode")
	}
	target, opts, err := agentDialTarget("unix:///run/beacon.sock")
	if err != nil {
		t.Fatalf("unix target rejected: %v", err)
	}
	if target != "passthrough:///collia-unix" || len(opts) == 0 {
		t.Fatalf("unix target not adapted: %q", target)
	}
}

// TestUnixSocketPath 覆盖前缀解析边界。
func TestUnixSocketPath(t *testing.T) {
	cases := map[string]string{
		"unix:///run/beacon.sock": "/run/beacon.sock",
		":17000":                  "",
		"127.0.0.1:17000":         "",
	}
	for input, want := range cases {
		if got := unixSocketPath(input); got != want {
			t.Fatalf("unixSocketPath(%q) = %q, want %q", input, got, want)
		}
	}
}

var _ = json.Marshal // keep encoding/json referenced when reply decoding changes
