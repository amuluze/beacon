// Package rpc
// Date: 2026/6/26
// Author: Amu
// Description:
package rpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"sync"

	rpcSchema "common/rpc/schema"
	"common/rpc/tunnel"

	"github.com/creack/pty"
)

// ptySession holds an active PTY and its underlying process.
type ptySession struct {
	id      string
	ptyFile *os.File
	cmd     *exec.Cmd
	mu      sync.Mutex
	closed  bool
}

func (s *ptySession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	// Closing the PTY file first signals the process group.
	_ = s.ptyFile.Close()
	if s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
		_, _ = s.cmd.Process.Wait()
	}
	return nil
}

func (s *ptySession) Resize(rows, cols int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return fmt.Errorf("session %s is closed", s.id)
	}
	return pty.Setsize(s.ptyFile, &pty.Winsize{Rows: uint16(rows), Cols: uint16(cols)})
}

// terminalSessions tracks active PTY sessions keyed by Server-assigned session ID.
var (
	terminalSessions   = make(map[string]*ptySession)
	terminalSessionsMu sync.RWMutex
)

func getTerminalSession(id string) (*ptySession, bool) {
	terminalSessionsMu.RLock()
	defer terminalSessionsMu.RUnlock()
	s, ok := terminalSessions[id]
	return s, ok
}

func putTerminalSession(id string, s *ptySession) {
	terminalSessionsMu.Lock()
	defer terminalSessionsMu.Unlock()
	terminalSessions[id] = s
}

func deleteTerminalSession(id string) {
	terminalSessionsMu.Lock()
	defer terminalSessionsMu.Unlock()
	delete(terminalSessions, id)
}

// TerminalSessionStream starts a PTY shell session and streams output back to Server.
// It blocks until the shell process exits or the context is cancelled.
func (s *Service) TerminalSessionStream(ctx context.Context, args rpcSchema.TerminalSessionArgs, streamSender func(*tunnel.Frame)) error {
	shell := args.Shell
	if shell == "" {
		shell = "/bin/bash"
	}

	cmd := exec.CommandContext(ctx, shell)
	cmd.Env = os.Environ()

	ptty, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("failed to start pty: %w", err)
	}

	session := &ptySession{
		id:      args.SessionID,
		ptyFile: ptty,
		cmd:     cmd,
	}
	putTerminalSession(args.SessionID, session)
	defer func() {
		deleteTerminalSession(args.SessionID)
		_ = session.Close()
	}()

	if err := session.Resize(args.Rows, args.Cols); err != nil {
		slog.Warn("terminal: initial resize failed", "session_id", args.SessionID, "err", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	buf := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			streamSender(&tunnel.Frame{Eos: true})
			return ctx.Err()
		case err := <-done:
			_ = drainPTY(ptty, streamSender)
			streamSender(&tunnel.Frame{Eos: true})
			if err != nil && !errors.Is(err, context.Canceled) && err.Error() != "signal: killed" {
				return fmt.Errorf("shell exited: %w", err)
			}
			return nil
		default:
			n, err := ptty.Read(buf)
			if n > 0 {
				payload := make([]byte, n)
				copy(payload, buf[:n])
				streamSender(&tunnel.Frame{Payload: payload})
			}
			if err != nil {
				if errors.Is(err, io.EOF) {
					streamSender(&tunnel.Frame{Eos: true})
					return nil
				}
				streamSender(&tunnel.Frame{Eos: true})
				return fmt.Errorf("pty read error: %w", err)
			}
		}
	}
}

func drainPTY(ptty *os.File, streamSender func(*tunnel.Frame)) error {
	buf := make([]byte, 4096)
	for {
		n, err := ptty.Read(buf)
		if n > 0 {
			streamSender(&tunnel.Frame{Payload: append([]byte(nil), buf[:n]...)})
		}
		if err != nil {
			return err
		}
	}
}

// ResizeTerminal resizes an active PTY session.
func (s *Service) ResizeTerminal(ctx context.Context, args rpcSchema.ResizeTerminalArgs, reply *rpcSchema.ResizeTerminalReply) error {
	session, ok := getTerminalSession(args.SessionID)
	if !ok {
		return fmt.Errorf("session %s not found", args.SessionID)
	}
	if err := session.Resize(args.Rows, args.Cols); err != nil {
		return fmt.Errorf("resize failed: %w", err)
	}
	return nil
}
