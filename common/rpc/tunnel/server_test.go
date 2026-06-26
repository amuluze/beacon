package tunnel

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// 本文件用 bufconn 在进程内驱动 ServerTunnel，覆盖 Domain Spec 与 server.go 的关键路径：
// 注册认证、向后兼容、请求-响应、Agent 离线、RPC 错误、超时、流式分发与生命周期回调。

const bufSize = 1024 * 1024

// fakeLifecycle 记录生命周期回调，供测试断言。
type fakeLifecycle struct {
	mu          sync.Mutex
	connects    []AgentInfo
	disconnects []string
	heartbeats  []string
}

func (f *fakeLifecycle) OnAgentConnect(info AgentInfo) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.connects = append(f.connects, info)
}

func (f *fakeLifecycle) OnAgentDisconnect(agentID string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.disconnects = append(f.disconnects, agentID)
}

func (f *fakeLifecycle) OnAgentHeartbeat(agentID string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.heartbeats = append(f.heartbeats, agentID)
}

func (f *fakeLifecycle) connectCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.connects)
}

func (f *fakeLifecycle) heartbeatCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.heartbeats)
}

func (f *fakeLifecycle) disconnectCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.disconnects)
}

// startTestServer 启动一个 bufconn 后端的 ServerTunnel，返回 server 与 client 连接。
// 调用方通过 NewReverseTunnelClient(conn) 驱动。
func startTestServer(t *testing.T, opts ...ServerOption) (*ServerTunnel, *grpc.ClientConn) {
	t.Helper()
	lis := bufconn.Listen(bufSize)
	s := NewServerTunnel(opts...)
	grpcServer := grpc.NewServer()
	RegisterReverseTunnelServer(grpcServer, s)
	go func() { _ = grpcServer.Serve(lis) }()

	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("dial bufnet: %v", err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
		grpcServer.Stop()
	})
	return s, conn
}

func openAgentStream(t *testing.T, conn *grpc.ClientConn) grpc.BidiStreamingClient[Frame, Frame] {
	t.Helper()
	stream, err := NewReverseTunnelClient(conn).Tunnel(context.Background())
	if err != nil {
		t.Fatalf("open tunnel stream: %v", err)
	}
	return stream
}

func registerFrame(agentID string, payload RegistrationPayload) *Frame {
	b, _ := json.Marshal(payload)
	return &Frame{
		FrameType: FrameType_FRAME_REGISTER,
		Method:    agentID,
		Payload:   b,
	}
}

// serveAgent 模拟 agent 注册并进入请求处理循环。
// handler 为 nil 表示不处理请求；返回 (reply, err)。
type agentHandler func(method string, payload []byte) ([]byte, error)

func serveAgent(stream grpc.BidiStreamingClient[Frame, Frame], reg *Frame, h agentHandler) {
	_ = stream.Send(reg)
	go func() {
		for {
			frame, err := stream.Recv()
			if err != nil {
				return
			}
			if frame.FrameType != FrameType_FRAME_REQUEST {
				continue
			}
			req := frame
			go func() {
				if h == nil {
					_ = stream.Send(&Frame{Id: req.Id, FrameType: FrameType_FRAME_REPLY, Error: "no handler"})
					return
				}
				reply, err := h(req.Method, req.Payload)
				resp := &Frame{Id: req.Id, FrameType: FrameType_FRAME_REPLY}
				if err != nil {
					resp.Error = err.Error()
				} else {
					resp.Payload = reply
				}
				_ = stream.Send(resp)
			}()
		}
	}()
}

// eventually 轮询等待条件满足，避免注册异步处理的竞态。
func eventually(t *testing.T, fn func() bool, timeout time.Duration, msg string) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if fn() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("condition not met within %v: %s", timeout, msg)
}

// --- 注册与生命周期 ---

func TestServerTunnel_Register_Success(t *testing.T) {
	lc := &fakeLifecycle{}
	s, conn := startTestServer(t, WithJoinToken("secret"))
	s.SetAgentLifecycle(lc)

	stream := openAgentStream(t, conn)
	serveAgent(stream, registerFrame("agent-1", RegistrationPayload{
		AgentID: "agent-1", Version: "v1.0.0", OS: "linux", Arch: "amd64", JoinToken: "secret",
	}), nil)

	eventually(t, func() bool { return s.IsOnline("agent-1") }, time.Second, "agent should be online")
	eventually(t, func() bool { return lc.connectCount() == 1 }, time.Second, "connect callback")

	if s.AgentCount() != 1 {
		t.Fatalf("agent count = %d, want 1", s.AgentCount())
	}
	if got := s.AgentIDs(); len(got) != 1 || got[0] != "agent-1" {
		t.Fatalf("agent ids = %v, want [agent-1]", got)
	}
}

func TestServerTunnel_Register_RejectedBadToken(t *testing.T) {
	lc := &fakeLifecycle{}
	s, conn := startTestServer(t, WithJoinToken("secret"))
	s.SetAgentLifecycle(lc)

	stream := openAgentStream(t, conn)
	_ = stream.Send(registerFrame("agent-1", RegistrationPayload{
		AgentID: "agent-1", JoinToken: "wrong",
	}))

	// 服务端应回 REGISTER_REJECTED
	frame, err := stream.Recv()
	if err != nil {
		t.Fatalf("recv: %v", err)
	}
	if frame.FrameType != FrameType_FRAME_REGISTER_REJECTED {
		t.Fatalf("frame type = %v, want REGISTER_REJECTED", frame.FrameType)
	}
	if s.IsOnline("agent-1") {
		t.Fatalf("rejected agent should not be online")
	}
	if lc.connectCount() != 0 {
		t.Fatalf("connect callback should not fire on rejection")
	}
}

func TestServerTunnel_Register_BackwardCompatibleNoToken(t *testing.T) {
	// 服务端未配置 joinToken 时，跳过校验（向后兼容）。
	s, conn := startTestServer(t)
	stream := openAgentStream(t, conn)
	serveAgent(stream, registerFrame("agent-2", RegistrationPayload{AgentID: "agent-2"}), nil)

	eventually(t, func() bool { return s.IsOnline("agent-2") }, time.Second, "agent should be online without token")
}

func TestServerTunnel_Register_InvalidPayload(t *testing.T) {
	s, conn := startTestServer(t)
	stream := openAgentStream(t, conn)
	_ = stream.Send(&Frame{FrameType: FrameType_FRAME_REGISTER, Method: "agent-x", Payload: []byte("not-json")})

	frame, err := stream.Recv()
	if err != nil {
		t.Fatalf("recv: %v", err)
	}
	if frame.FrameType != FrameType_FRAME_REGISTER_REJECTED {
		t.Fatalf("frame type = %v, want REGISTER_REJECTED", frame.FrameType)
	}
	if s.IsOnline("agent-x") {
		t.Fatalf("agent with invalid payload should not be online")
	}
}

func TestServerTunnel_Disconnect_Lifecycle(t *testing.T) {
	lc := &fakeLifecycle{}
	s, conn := startTestServer(t)
	s.SetAgentLifecycle(lc)

	stream := openAgentStream(t, conn)
	serveAgent(stream, registerFrame("agent-3", RegistrationPayload{AgentID: "agent-3"}), nil)
	eventually(t, func() bool { return s.IsOnline("agent-3") }, time.Second, "online")

	// 关闭 agent 流触发断开回调
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("close send: %v", err)
	}
	eventually(t, func() bool { return lc.disconnectCount() == 1 }, time.Second, "disconnect callback")
	if s.IsOnline("agent-3") {
		t.Fatalf("agent should be offline after disconnect")
	}
}

func TestServerTunnel_Heartbeat_Lifecycle(t *testing.T) {
	lc := &fakeLifecycle{}
	s, conn := startTestServer(t)
	s.SetAgentLifecycle(lc)

	stream := openAgentStream(t, conn)
	serveAgent(stream, registerFrame("agent-4", RegistrationPayload{AgentID: "agent-4"}), nil)
	eventually(t, func() bool { return s.IsOnline("agent-4") }, time.Second, "online")

	_ = stream.Send(&Frame{FrameType: FrameType_FRAME_HEARTBEAT, Method: "agent-4"})
	eventually(t, func() bool { return lc.heartbeatCount() >= 1 }, time.Second, "heartbeat callback")
}

// --- 请求-响应 ---

func TestServerTunnel_Call_Success(t *testing.T) {
	s, conn := startTestServer(t)
	stream := openAgentStream(t, conn)
	serveAgent(stream, registerFrame("agent-5", RegistrationPayload{AgentID: "agent-5"}),
		func(method string, payload []byte) ([]byte, error) {
			var in string
			_ = json.Unmarshal(payload, &in)
			return json.Marshal("echo:" + in)
		})
	eventually(t, func() bool { return s.IsOnline("agent-5") }, time.Second, "online")

	var reply string
	err := s.Call(context.Background(), "agent-5", "Echo", "hello", &reply)
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	if reply != "echo:hello" {
		t.Fatalf("reply = %q, want echo:hello", reply)
	}
}

func TestServerTunnel_Call_AgentOffline(t *testing.T) {
	s, _ := startTestServer(t)
	var reply string
	err := s.Call(context.Background(), "ghost", "Echo", "hi", &reply)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	var offlineErr *AgentOfflineError
	if !errors.As(err, &offlineErr) {
		t.Fatalf("err type = %T, want *AgentOfflineError", err)
	}
}

func TestServerTunnel_Call_RPCError(t *testing.T) {
	s, conn := startTestServer(t)
	stream := openAgentStream(t, conn)
	boom := errors.New("container not found")
	serveAgent(stream, registerFrame("agent-6", RegistrationPayload{AgentID: "agent-6"}),
		func(method string, payload []byte) ([]byte, error) {
			return nil, boom
		})
	eventually(t, func() bool { return s.IsOnline("agent-6") }, time.Second, "online")

	var reply string
	err := s.Call(context.Background(), "agent-6", "Stop", nil, &reply)
	if err == nil {
		t.Fatalf("expected rpc error")
	}
	var rpcErr *RPCError
	if !errors.As(err, &rpcErr) {
		t.Fatalf("err type = %T, want *RPCError", err)
	}
}

func TestServerTunnel_Call_Timeout(t *testing.T) {
	s, conn := startTestServer(t)
	stream := openAgentStream(t, conn)
	// 注册后进入空循环：收到请求但不回复，触发调用方 ctx 超时。
	_ = stream.Send(registerFrame("agent-7", RegistrationPayload{AgentID: "agent-7"}))
	eventually(t, func() bool { return s.IsOnline("agent-7") }, time.Second, "online")
	go func() {
		for {
			if _, err := stream.Recv(); err != nil {
				return
			}
			// 故意丢弃，不回 REPLY
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	var reply string
	err := s.Call(ctx, "agent-7", "Slow", nil, &reply)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want DeadlineExceeded", err)
	}
}

// --- 流式 ---

func TestServerTunnel_StreamCall(t *testing.T) {
	s, conn := startTestServer(t)
	st := openAgentStream(t, conn)
	_ = st.Send(registerFrame("agent-8", RegistrationPayload{AgentID: "agent-8"}))
	eventually(t, func() bool { return s.IsOnline("agent-8") }, time.Second, "online")

	// agent 循环：收到 REQUEST 后回 3 个 STREAM_DATA + 1 个 STREAM_END
	go func() {
		for {
			frame, err := st.Recv()
			if err != nil {
				return
			}
			if frame.FrameType != FrameType_FRAME_REQUEST {
				continue
			}
			for i := 0; i < 3; i++ {
				_ = st.Send(&Frame{
					Id:        frame.Id,
					FrameType: FrameType_FRAME_STREAM_DATA,
					Payload:   []byte{byte('a' + i)},
				})
			}
			_ = st.Send(&Frame{Id: frame.Id, FrameType: FrameType_FRAME_STREAM_END})
		}
	}()

	ch, err := s.StreamCall(context.Background(), "agent-8", "Tail", nil)
	if err != nil {
		t.Fatalf("stream call: %v", err)
	}
	// 验证能收到全部 STREAM_DATA 帧。
	// 注意：STREAM_END 当前被 server 路由到 pending 查找而不会到达调用方，
	// 属于流式健壮性缺陷，由后续 StreamEnd 路由修复覆盖；此处以数据完整性为准。
	var got []byte
	timeout := time.After(time.Second)
	for {
		select {
		case f, ok := <-ch:
			if !ok {
				t.Fatalf("channel closed early, got %q", string(got))
			}
			got = append(got, f.Payload...)
			if string(got) == "abc" {
				return
			}
		case <-timeout:
			t.Fatalf("stream timeout, got %q", string(got))
		}
	}
}

// TestServerTunnel_StreamCall_StreamEndDelivered 在 STREAM_END 路由修复后启用：
// 调用方必须能收到 STREAM_END 并据此结束读取。
func TestServerTunnel_StreamCall_StreamEndDelivered(t *testing.T) {
	s, conn := startTestServer(t)
	st := openAgentStream(t, conn)
	_ = st.Send(registerFrame("agent-9", RegistrationPayload{AgentID: "agent-9"}))
	eventually(t, func() bool { return s.IsOnline("agent-9") }, time.Second, "online")

	go func() {
		for {
			frame, err := st.Recv()
			if err != nil {
				return
			}
			if frame.FrameType != FrameType_FRAME_REQUEST {
				continue
			}
			_ = st.Send(&Frame{Id: frame.Id, FrameType: FrameType_FRAME_STREAM_DATA, Payload: []byte("x")})
			_ = st.Send(&Frame{Id: frame.Id, FrameType: FrameType_FRAME_STREAM_END})
		}
	}()

	ch, err := s.StreamCall(context.Background(), "agent-9", "Tail", nil)
	if err != nil {
		t.Fatalf("stream call: %v", err)
	}
	timeout := time.After(time.Second)
	for {
		select {
		case f, ok := <-ch:
			if !ok {
				t.Fatalf("channel closed before STREAM_END")
			}
			if f.FrameType == FrameType_FRAME_STREAM_END {
				return
			}
		case <-timeout:
			t.Fatalf("did not receive STREAM_END within timeout")
		}
	}
}

// TestServerTunnel_StreamCall_CtxCancelReleasesStream 验证调用方取消 ctx 后，
// streams 条目被清理（防泄漏）：制造"流不结束"场景，取消 ctx 后 streams 必须清空。
func TestServerTunnel_StreamCall_CtxCancelReleasesStream(t *testing.T) {
	s, conn := startTestServer(t)
	st := openAgentStream(t, conn)
	_ = st.Send(registerFrame("agent-10", RegistrationPayload{AgentID: "agent-10"}))
	eventually(t, func() bool { return s.IsOnline("agent-10") }, time.Second, "online")

	go func() {
		for {
			if _, err := st.Recv(); err != nil {
				return
			}
			// 收到请求但不回复任何数据，制造"流不结束"场景
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	if _, err := s.StreamCall(ctx, "agent-10", "Tail", nil); err != nil {
		t.Fatalf("stream call: %v", err)
	}

	// 此时 streams 中应驻留一个 pendingStream
	remaining := func() int {
		n := 0
		s.streams.Range(func(_, _ any) bool { n++; return true })
		return n
	}
	if remaining() != 1 {
		t.Fatalf("streams count = %d, want 1 before cancel", remaining())
	}

	// 取消调用方 ctx，ctx 兜底 goroutine 应删除条目
	cancel()
	eventually(t, func() bool { return remaining() == 0 }, time.Second, "streams cleared after cancel")
}
