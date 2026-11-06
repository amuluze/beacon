package rpc

import (
	"context"
	"errors"
	"io"
	"os"
	"syscall"
	"testing"
	"time"

	rpcSchema "common/rpc/schema"
	"common/rpc/tunnel"
)

// TestIsPTYSessionEnd 固化 pty 会话结束判定：Linux 的 ptmx 在从属端关闭后
// 返回 EIO（而非 EOF），两者都必须按会话正常结束处理，否则 Web 终端在
// Linux 上会以 "pty read error" 结束而不是干净的 EOS。
func TestIsPTYSessionEnd(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"eof", io.EOF, true},
		{"eio", syscall.EIO, true},
		{"wrapped eio", &os.PathError{Op: "read", Path: "/dev/ptmx", Err: syscall.EIO}, true},
		{"closed", os.ErrClosed, false},
		{"other", errors.New("boom"), false},
	}
	for _, tc := range cases {
		if got := isPTYSessionEnd(tc.err); got != tc.want {
			t.Errorf("%s: isPTYSessionEnd = %v, want %v", tc.name, got, tc.want)
		}
	}
}

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

func TestTerminalSessionStream_RejectsUnsupportedShell(t *testing.T) {
	svc := &Service{}
	err := svc.TerminalSessionStream(context.Background(), rpcSchema.TerminalSessionArgs{
		SessionID: "test-shell",
		Shell:     "/tmp/untrusted-shell",
		Rows:      24,
		Cols:      80,
	}, func(*tunnel.Frame) {})
	if err == nil {
		t.Fatal("expected unsupported shell to be rejected")
	}
}

func TestPTYSessionResizeRejectsOutOfRangeSize(t *testing.T) {
	session := &ptySession{id: "test-size"}
	for _, size := range []struct {
		rows int
		cols int
	}{{rows: 0, cols: 80}, {rows: 24, cols: 0}, {rows: 65536, cols: 80}, {rows: 24, cols: 65536}} {
		if err := session.Resize(size.rows, size.cols); err == nil {
			t.Fatalf("expected size rows=%d cols=%d to be rejected", size.rows, size.cols)
		}
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
