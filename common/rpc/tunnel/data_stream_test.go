package tunnel

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

// 数据流角色的关键契约：
//   - 与控制流共享连接与鉴权（join token）；
//   - 不写入 Agent 在线注册表，不触发生命周期回调（不得覆盖控制流会话）；
//   - 批次载荷路由到 DataHandler，失败只丢该批并回执错误。

func TestServerTunnel_DataStream_DoesNotRegisterAgent(t *testing.T) {
	server, conn := startTestServer(t, WithJoinToken("join-secret"))
	life := &fakeLifecycle{}
	server.SetAgentLifecycle(life)

	received := make(chan string, 4)
	server.SetDataHandler(func(_ context.Context, agentID string, method string, payload []byte) error {
		received <- agentID + "|" + method + "|" + string(payload)
		return nil
	})

	stream := openAgentStream(t, conn)
	reg := registerFrame("node-1", RegistrationPayload{AgentID: "node-1", JoinToken: "join-secret", Role: RoleData})
	if err := stream.Send(reg); err != nil {
		t.Fatalf("send data registration: %v", err)
	}
	if err := stream.Send(&Frame{
		Id:        "b1",
		FrameType: FrameType_FRAME_STREAM_DATA,
		Method:    "MonitorReport",
		Payload:   []byte(`{"agent_id":"node-1"}`),
	}); err != nil {
		t.Fatalf("send data frame: %v", err)
	}

	select {
	case got := <-received:
		if got != `node-1|MonitorReport|{"agent_id":"node-1"}` {
			t.Fatalf("handler got %q", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("data handler not invoked")
	}

	// 数据流不得注册在线 Agent，也不得触发生命周期回调
	if ids := server.AgentIDs(); len(ids) != 0 {
		t.Fatalf("data stream registered agent(s): %v", ids)
	}
	if life.connectCount() != 0 || life.disconnectCount() != 0 {
		t.Fatalf("data stream triggered lifecycle: connect=%d disconnect=%d", life.connectCount(), life.disconnectCount())
	}
}

func TestServerTunnel_DataStream_RejectsBadToken(t *testing.T) {
	server, conn := startTestServer(t, WithJoinToken("join-secret"))
	server.SetDataHandler(func(context.Context, string, string, []byte) error { return nil })

	stream := openAgentStream(t, conn)
	reg := registerFrame("node-1", RegistrationPayload{AgentID: "node-1", JoinToken: "wrong", Role: RoleData})
	if err := stream.Send(reg); err != nil {
		t.Fatalf("send data registration: %v", err)
	}
	frame, err := stream.Recv()
	if err != nil {
		t.Fatalf("recv reject frame: %v", err)
	}
	if frame.FrameType != FrameType_FRAME_REGISTER_REJECTED {
		t.Fatalf("frame type = %v, want REGISTER_REJECTED", frame.FrameType)
	}
	if len(server.AgentIDs()) != 0 {
		t.Fatalf("rejected data stream registered agent: %v", server.AgentIDs())
	}
}

func TestServerTunnel_DataStream_RejectedPayloadKeepsStreamOpen(t *testing.T) {
	server, conn := startTestServer(t, WithJoinToken("join-secret"))
	reject := errors.New("agent id mismatch")
	server.SetDataHandler(func(context.Context, string, string, []byte) error { return reject })

	stream := openAgentStream(t, conn)
	reg := registerFrame("node-1", RegistrationPayload{AgentID: "node-1", JoinToken: "join-secret", Role: RoleData})
	if err := stream.Send(reg); err != nil {
		t.Fatalf("send data registration: %v", err)
	}
	if err := stream.Send(&Frame{
		Id:        "b1",
		FrameType: FrameType_FRAME_STREAM_DATA,
		Method:    "MonitorReport",
		Payload:   []byte(`{"agent_id":"other"}`),
	}); err != nil {
		t.Fatalf("send data frame: %v", err)
	}

	frame, err := stream.Recv()
	if err != nil {
		t.Fatalf("recv ack: %v", err)
	}
	if frame.FrameType != FrameType_FRAME_REPLY || frame.Error == "" {
		t.Fatalf("ack = %+v, want REPLY with error", frame)
	}

	// 拒绝后流仍可继续推送
	calls := 0
	server.SetDataHandler(func(context.Context, string, string, []byte) error {
		calls++
		return nil
	})
	if err := stream.Send(&Frame{
		Id:        "b2",
		FrameType: FrameType_FRAME_STREAM_DATA,
		Method:    "MonitorReport",
		Payload:   []byte(`{"agent_id":"node-1"}`),
	}); err != nil {
		t.Fatalf("send second data frame: %v", err)
	}
	ack, err := stream.Recv()
	if err != nil {
		t.Fatalf("recv second ack: %v", err)
	}
	if ack.FrameType != FrameType_FRAME_REPLY || ack.Error != "" {
		t.Fatalf("second ack = %+v, want REPLY without error", ack)
	}
	if calls != 1 {
		t.Fatalf("handler calls = %d, want 1", calls)
	}
}

func TestServerTunnel_DataRegistrationPayloadCarriesRole(t *testing.T) {
	// 旧 Agent 不带 role 字段时必须保持控制流语义（向后兼容）。
	payload := RegistrationPayload{AgentID: "node-1", JoinToken: "t"}
	encoded, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded RegistrationPayload
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded.Role != "" {
		t.Fatalf("role = %q, want empty (control by default)", decoded.Role)
	}
	if decoded.Role != "" && decoded.Role != RoleControl {
		t.Fatalf("role = %q, want control or empty", decoded.Role)
	}
}
