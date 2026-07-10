package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"common/rpc/tunnel"
)

// TestDispatcher_UnknownMethod 验证未注册方法返回可区分的 UnknownMethodError。
// 控制通道要求调用失败语义可区分（见 common/rpc/tunnel 约束）。
func TestDispatcher_UnknownMethod(t *testing.T) {
	d := NewDispatcher(&Service{})
	_, err := d.Call(context.Background(), "Nonexistent", nil, nil)
	if err == nil {
		t.Fatal("expected error for unknown method")
	}
	var unknownErr *UnknownMethodError
	if !errors.As(err, &unknownErr) {
		t.Fatalf("expected UnknownMethodError, got %T: %v", err, err)
	}
	if unknownErr.Method != "Nonexistent" {
		t.Fatalf("expected method name in error, got %q", unknownErr.Method)
	}
	if err.Error() != "unknown rpc method: Nonexistent" {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}

// echoArgs/echoReply 用于 RegisterUnary 测试。
type echoArgs struct {
	Name string `json:"name"`
}
type echoReply struct {
	Greeting string `json:"greeting"`
}

// TestRegisterUnary 验证 unary handler 的 args 反序列化、调用与 reply 序列化全链路。
func TestRegisterUnary(t *testing.T) {
	d := NewDispatcher(&Service{})
	RegisterUnary(d, "Echo", func(ctx context.Context, args echoArgs, reply *echoReply) error {
		reply.Greeting = "hello " + args.Name
		return nil
	})

	out, err := d.Call(context.Background(), "Echo", []byte(`{"name":"world"}`), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got echoReply
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("unmarshal reply: %v", err)
	}
	if got.Greeting != "hello world" {
		t.Fatalf("expected greeting 'hello world', got %q", got.Greeting)
	}
}

// TestRegisterUnary_InvalidPayload 验证非法 JSON args 返回带方法名的包装错误。
func TestRegisterUnary_InvalidPayload(t *testing.T) {
	d := NewDispatcher(&Service{})
	RegisterUnary(d, "Echo", func(ctx context.Context, args echoArgs, reply *echoReply) error {
		return nil
	})
	_, err := d.Call(context.Background(), "Echo", []byte(`{not-json`), nil)
	if err == nil {
		t.Fatal("expected error for invalid payload")
	}
}

// TestRegisterUnary_HandlerError 验证 handler 返回的错误原样透传，不被吞掉。
func TestRegisterUnary_HandlerError(t *testing.T) {
	d := NewDispatcher(&Service{})
	sentinel := errors.New("boom")
	RegisterUnary(d, "Fail", func(ctx context.Context, args echoArgs, reply *echoReply) error {
		return sentinel
	})
	_, err := d.Call(context.Background(), "Fail", []byte(`{}`), nil)
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error to propagate, got %v", err)
	}
}

// TestRegisterStream 验证 stream handler 接收 streamSender 并能发送多帧，且回复 payload 为 nil。
func TestRegisterStream(t *testing.T) {
	d := NewDispatcher(&Service{})
	RegisterStream(d, "Stream", func(ctx context.Context, args echoArgs, streamSender func(*tunnel.Frame)) error {
		streamSender(&tunnel.Frame{Payload: []byte("chunk1")})
		streamSender(&tunnel.Frame{Payload: []byte("chunk2")})
		return nil
	})

	var frames []*tunnel.Frame
	sender := func(f *tunnel.Frame) { frames = append(frames, f) }

	out, err := d.Call(context.Background(), "Stream", []byte(`{"name":"s"}`), sender)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Fatalf("streaming methods must return nil payload, got %v", out)
	}
	if len(frames) != 2 {
		t.Fatalf("expected 2 streamed frames, got %d", len(frames))
	}
	if string(frames[1].Payload) != "chunk2" {
		t.Fatalf("unexpected second frame payload: %q", string(frames[1].Payload))
	}
}

// TestRegisterStream_InvalidPayload 验证 stream handler 同样校验 args。
func TestRegisterStream_InvalidPayload(t *testing.T) {
	d := NewDispatcher(&Service{})
	RegisterStream(d, "Stream", func(ctx context.Context, args echoArgs, streamSender func(*tunnel.Frame)) error {
		return nil
	})
	_, err := d.Call(context.Background(), "Stream", []byte(`broken`), nil)
	if err == nil {
		t.Fatal("expected error for invalid stream payload")
	}
}

// TestRegister_Override 验证同名方法后注册覆盖先注册，符合 map 语义。
func TestRegister_Override(t *testing.T) {
	d := NewDispatcher(&Service{})
	RegisterUnary(d, "Override", func(ctx context.Context, args echoArgs, reply *echoReply) error {
		reply.Greeting = "first"
		return nil
	})
	RegisterUnary(d, "Override", func(ctx context.Context, args echoArgs, reply *echoReply) error {
		reply.Greeting = "second"
		return nil
	})
	out, err := d.Call(context.Background(), "Override", []byte(`{}`), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got echoReply
	_ = json.Unmarshal(out, &got)
	if got.Greeting != "second" {
		t.Fatalf("expected override to win, got %q", got.Greeting)
	}
}

// TestNewDispatcher_RegistersStandardMethods 验证四类标准控制方法均已注册，
// 避免回归导致控制通道方法丢失（Agent 控制能力契约）。
// 直接检查 handlers map，避免在空 Service 上触发真实 handler。
func TestNewDispatcher_RegistersStandardMethods(t *testing.T) {
	d := NewDispatcher(&Service{})
	cases := []string{
		// container
		"ContainerCreate", "ContainerUpdate", "ContainerStart", "ContainerStop", "ContainerLogs",
		"ImagePull", "NetworkCreate",
		// file
		"FilesSearch", "FileDownload", "FolderCreate",
		// system
		"Reboot", "GetSystemTime", "SetDNS",
		// terminal
		"TerminalSession", "ResizeTerminal", "TerminalClose",
	}
	for _, method := range cases {
		if _, ok := d.handlers[method]; !ok {
			t.Errorf("expected method %q to be registered", method)
		}
	}
}
