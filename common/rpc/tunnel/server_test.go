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

	// Server should reject with an error, which closes the stream
	// Use a goroutine + timeout to avoid hanging on Recv
	type result struct {
		frame *Frame
		err   error
	}
	ch := make(chan result, 1)
	go func() {
		f, e := stream.Recv()
		ch <- result{f, e}
	}()
	var err error
	select {
	case r := <-ch:
		_, err = r.frame, r.err
	case <-time.After(5 * time.Second):
		t.Fatalf("recv timeout: expected stream to close with error")
	}
	if err == nil {
		t.Fatalf("expected error from closed stream, got nil")
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
	// Verify all STREAM_DATA frames are received.
	// Stream_END delivery is verified separately by TestServerTunnel_StreamCall_StreamEndDelivered.
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

// ── E2E Integration Tests ──
// 全链路覆盖 Domain I001 / I003 / R001 / R006：
// Agent 注册 → Server 派发 RPC → Agent 处理 → Server 接收响应 → Agent 断开

// TestE2E_RegisterCallDisconnect 覆盖完整的 Agent 生命周期链路。
func TestE2E_RegisterCallDisconnect(t *testing.T) {
	lc := &fakeLifecycle{}
	s, conn := startTestServer(t, WithJoinToken("e2e-token"))
	s.SetAgentLifecycle(lc)

	// Phase 1: Agent registers with metadata
	stream := openAgentStream(t, conn)
	payload, _ := json.Marshal(RegistrationPayload{
		AgentID: "e2e-agent", Version: "v2.0.0", OS: "linux", Arch: "arm64", JoinToken: "e2e-token",
	})
	registerFrame := &Frame{
		FrameType: FrameType_FRAME_REGISTER,
		Method:    "e2e-agent",
		Payload:   payload,
	}
	serveAgent(stream, registerFrame, func(method string, payload []byte) ([]byte, error) {
		// Agent handler: echo the payload back
		return payload, nil
	})
	eventually(t, func() bool { return s.IsOnline("e2e-agent") }, time.Second, "agent should be online")
	eventually(t, func() bool { return lc.connectCount() == 1 }, time.Second, "connect callback should fire")

	// Phase 2: Server dispatches RPC call
	var reply map[string]interface{}
	err := s.Call(context.Background(), "e2e-agent", "GetInfo", map[string]string{"key": "value"}, &reply)
	if err != nil {
		t.Fatalf("E2E Call: %v", err)
	}
	if reply["key"] != "value" {
		t.Fatalf("E2E reply = %v, want key=value", reply)
	}

	// Phase 3: Agent disconnects
	err = stream.CloseSend()
	if err != nil {
		t.Fatalf("close send: %v", err)
	}
	eventually(t, func() bool { return !s.IsOnline("e2e-agent") }, time.Second, "agent should be offline")
	eventually(t, func() bool { return lc.disconnectCount() == 1 }, time.Second, "disconnect callback should fire")
}

// TestE2E_MultiAgentIsolation 验证多 Agent 场景下调用正确隔离。
// Domain I003: 容器运行时操作必须只影响目标节点。
// 注：bufconn 测试中，两个 Agent 通过同一 gRPC 连接的不同 stream 注册，
// 这模拟了 NAT/防火墙后的真实场景（Server 看到同一 IP 但不同 Agent）。
func TestE2E_MultiAgentIsolation(t *testing.T) {
	lc := &fakeLifecycle{}
	s, conn := startTestServer(t)
	s.SetAgentLifecycle(lc)

	// Register agent-a on first stream
	streamA := openAgentStream(t, conn)
	serveAgent(streamA, registerFrame("agent-a", RegistrationPayload{AgentID: "agent-a"}),
		func(method string, payload []byte) ([]byte, error) {
			return json.Marshal(map[string]string{"handler": "a"})
		})
	eventually(t, func() bool { return s.IsOnline("agent-a") }, time.Second, "agent-a online")

	// Register agent-b on second stream (same connection, different stream)
	streamB := openAgentStream(t, conn)
	serveAgent(streamB, registerFrame("agent-b", RegistrationPayload{AgentID: "agent-b"}),
		func(method string, payload []byte) ([]byte, error) {
			return json.Marshal(map[string]string{"handler": "b"})
		})
	eventually(t, func() bool { return s.IsOnline("agent-b") }, time.Second, "agent-b online")

	// Calls to each agent return distinct results
	var replyA, replyB map[string]string
	if err := s.Call(context.Background(), "agent-a", "Do", nil, &replyA); err != nil {
		t.Fatalf("Call agent-a: %v", err)
	}
	if err := s.Call(context.Background(), "agent-b", "Do", nil, &replyB); err != nil {
		t.Fatalf("Call agent-b: %v", err)
	}
	if replyA["handler"] != "a" || replyB["handler"] != "b" {
		t.Fatalf("multi-agent isolation failed: a=%q b=%q", replyA["handler"], replyB["handler"])
	}

	// Validate AgentIDs and AgentCount
	ids := s.AgentIDs()
	if len(ids) != 2 {
		t.Fatalf("AgentIDs = %v, want 2 agents", ids)
	}
	if s.AgentCount() != 2 {
		t.Fatalf("AgentCount = %d, want 2", s.AgentCount())
	}
}

// TestE2E_RegistrationRejectsUnauthorized 验证未授权 Agent 注册被拒绝。
// Domain R005: 无有效认证的上报请求被拒绝。
func TestE2E_RegistrationRejectsUnauthorized(t *testing.T) {
	s, conn := startTestServer(t, WithJoinToken("secure-token"))
	stream := openAgentStream(t, conn)

	// Send registration with wrong token
	payload, _ := json.Marshal(RegistrationPayload{
		AgentID: "eve", JoinToken: "wrong-token",
	})
	_ = stream.Send(&Frame{
		FrameType: FrameType_FRAME_REGISTER,
		Method:    "eve",
		Payload:   payload,
	})

	// Agent should receive REGISTER_REJECTED
	frame, err := stream.Recv()
	if err != nil {
		t.Fatalf("recv: %v", err)
	}
	if frame.FrameType != FrameType_FRAME_REGISTER_REJECTED {
		t.Fatalf("frame type = %v, want REGISTER_REJECTED", frame.FrameType)
	}
	if s.IsOnline("eve") {
		t.Fatalf("unauthorized agent eve should not be online")
	}
}

// TestE2E_HeartbeatKeepsAgentAlive 验证心跳保持 Agent 在线。
func TestE2E_HeartbeatKeepsAgentAlive(t *testing.T) {
	lc := &fakeLifecycle{}
	s, conn := startTestServer(t)
	s.SetAgentLifecycle(lc)

	stream := openAgentStream(t, conn)
	serveAgent(stream, registerFrame("alive-agent", RegistrationPayload{AgentID: "alive-agent"}), nil)
	eventually(t, func() bool { return s.IsOnline("alive-agent") }, time.Second, "online")

	// Send heartbeats
	for i := 0; i < 3; i++ {
		_ = stream.Send(&Frame{
			FrameType: FrameType_FRAME_HEARTBEAT,
			Method:    "alive-agent",
		})
		time.Sleep(10 * time.Millisecond)
	}

	eventually(t, func() bool { return lc.heartbeatCount() >= 3 }, 2*time.Second, "3+ heartbeats")
	if !s.IsOnline("alive-agent") {
		t.Fatalf("alive-agent should still be online after heartbeats")
	}
}

// TestUnmarshalReplyNilReply 验证当 reply=nil 时（即 fire-and-forget RPC），
// 即使 agent 返回了非空 payload，也不应触发 json.Unmarshal 错误。
//
// 背景：beacon 端 ContainerRepo.ImagesPrune 调用 RPCClient.Call(ctx, "ImagesPrune", nil, nil)
// 把 reply 传为 nil；agent 端 d.Register("ImagesPrune", ...) 返回 json.Marshal(struct{}{}) = "{}"。
// 修复前：unmarshalReply([]byte("{}"), nil) → json: Unmarshal(nil) 错误。
// 修复后：reply == nil 时直接返回 nil，payload 被有意丢弃。
func TestUnmarshalReplyNilReply(t *testing.T) {
	cases := []struct {
		name string
		data []byte
	}{
		{"nil payload", nil},
		{"empty payload", []byte{}},
		{"empty object payload (ImagesPrune shape)", []byte("{}")},
		{"non-empty payload", []byte(`{"foo":"bar"}`)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := unmarshalReply(tc.data, nil); err != nil {
				t.Fatalf("unmarshalReply(_, nil) returned %v, want nil (caller does not care about payload)", err)
			}
		})
	}
}

// TestUnmarshalReplyNonNilReply 验证当 reply 非 nil 时，正常解析 JSON。
func TestUnmarshalReplyNonNilReply(t *testing.T) {
	type reply struct {
		Foo string `json:"foo"`
	}
	var r reply
	if err := unmarshalReply([]byte(`{"foo":"bar"}`), &r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Foo != "bar" {
		t.Fatalf("got %q, want %q", r.Foo, "bar")
	}

	// 空 payload 写到非 nil reply 时，不报错也不修改字段（保持零值）
	r = reply{}
	if err := unmarshalReply(nil, &r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Foo != "" {
		t.Fatalf("empty payload should leave reply at zero value, got %q", r.Foo)
	}
}
