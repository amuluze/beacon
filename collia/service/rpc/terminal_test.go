package rpc

import (
	"context"
	"testing"
	"time"

	rpcSchema "common/rpc/schema"
	"common/rpc/tunnel"
)

func TestResizeTerminal_SessionNotFound(t *testing.T) {
	svc := &Service{}
	var reply rpcSchema.ResizeTerminalReply
	err := svc.ResizeTerminal(context.Background(), rpcSchema.ResizeTerminalArgs{
		SessionID: "missing-session",
		Rows:      30,
		Cols:      120,
	}, &reply)
	if err == nil {
		t.Fatal("expected error for missing session")
	}
}

func TestTerminalSessionStream_ContextCancel(t *testing.T) {
	svc := &Service{}
	var frames []*tunnel.Frame
	sender := func(f *tunnel.Frame) {
		frames = append(frames, f)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := svc.TerminalSessionStream(ctx, rpcSchema.TerminalSessionArgs{
		SessionID: "test-cancel",
		Shell:     "/bin/sh",
		Rows:      24,
		Cols:      80,
	}, sender)

	if err != nil && err != context.DeadlineExceeded {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify at least one frame was sent and the last one is EOS.
	if len(frames) == 0 {
		t.Fatal("expected at least one frame")
	}
	last := frames[len(frames)-1]
	if !last.Eos {
		t.Fatal("expected last frame to be EOS")
	}

	// Session should be removed after stream ends.
	if _, ok := getTerminalSession("test-cancel"); ok {
		t.Fatal("session should be removed after stream ends")
	}
}

func TestTerminalInput_SessionNotFound(t *testing.T) {
	svc := &Service{}
	var reply rpcSchema.TerminalInputReply
	err := svc.TerminalInput(context.Background(), rpcSchema.TerminalInputArgs{
		SessionID: "missing-session",
		Data:      []byte("hello"),
	}, &reply)
	if err == nil {
		t.Fatal("expected error for missing session")
	}
}

func TestTerminalClose_SessionNotFound(t *testing.T) {
	svc := &Service{}
	var reply rpcSchema.TerminalCloseReply
	err := svc.TerminalClose(context.Background(), rpcSchema.TerminalCloseArgs{
		SessionID: "missing-session",
	}, &reply)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
