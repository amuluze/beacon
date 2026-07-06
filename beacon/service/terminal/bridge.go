package terminal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"amprobe/pkg/rpc"
	rpcSchema "common/rpc/schema"

	"github.com/gofiber/contrib/websocket"
)

// Connection abstracts the WebSocket methods the bridge needs, enabling tests
// to inject a fake connection instead of relying on *websocket.Conn.
type Connection interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	WriteControl(messageType int, data []byte, deadline time.Time) error
}

// bridge forwards data between a WebSocket connection and an Agent PTY stream.
type bridge struct {
	rpcClient rpc.Caller
	conn      Connection
	stream    <-chan []byte
	recorder  *Recorder
	sessionID string
	agentID   string
}

func newBridge(rpcClient rpc.Caller, conn Connection, stream <-chan []byte, recorder *Recorder, sessionID, agentID string) *bridge {
	return &bridge{
		rpcClient: rpcClient,
		conn:      conn,
		stream:    stream,
		recorder:  recorder,
		sessionID: sessionID,
		agentID:   agentID,
	}
}

func (b *bridge) run(ctx context.Context, rows, cols int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		// Notify Agent to release PTY resources.
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer closeCancel()
		_ = b.rpcClient.Call(closeCtx, "TerminalClose", rpcSchema.TerminalCloseArgs{
			SessionID: b.sessionID,
		}, &rpcSchema.TerminalCloseReply{})
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		b.readWebSocket(ctx, cancel)
	}()

	go func() {
		defer wg.Done()
		b.readTunnel(ctx, cancel)
	}()

	wg.Wait()
	return nil
}

// readWebSocket reads messages from browser and forwards to Agent.
func (b *bridge) readWebSocket(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, data, err := b.conn.ReadMessage()
		if err != nil {
			slog.Debug("terminal: websocket read failed", "session_id", b.sessionID, "err", err)
			return
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			slog.Warn("terminal: invalid websocket message", "session_id", b.sessionID, "err", err)
			continue
		}

		switch MessageType(msg.Type) {
		case MessageTypeInput:
			decoded, err := base64.StdEncoding.DecodeString(msg.Data)
			if err != nil {
				slog.Warn("terminal: invalid input base64", "session_id", b.sessionID, "err", err)
				continue
			}
			if err := b.rpcClient.Call(ctx, "TerminalInput", rpcSchema.TerminalInputArgs{
				SessionID: b.sessionID,
				Data:      decoded,
			}, &rpcSchema.TerminalInputReply{}); err != nil {
				slog.Debug("terminal: send input failed", "session_id", b.sessionID, "err", err)
				return
			}

		case MessageTypeResize:
			if err := b.rpcClient.Call(ctx, "ResizeTerminal", rpcSchema.ResizeTerminalArgs{
				SessionID: b.sessionID,
				Rows:      msg.Rows,
				Cols:      msg.Cols,
			}, &rpcSchema.ResizeTerminalReply{}); err != nil {
				slog.Debug("terminal: resize failed", "session_id", b.sessionID, "err", err)
				continue
			}
			if b.recorder != nil {
				_ = b.recorder.Resize(msg.Cols, msg.Rows)
			}

		default:
			slog.Debug("terminal: unknown message type", "session_id", b.sessionID, "type", msg.Type)
		}
	}
}

// readTunnel reads output from Agent and forwards to browser.
func (b *bridge) readTunnel(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case frame, ok := <-b.stream:
			if !ok {
				_ = sendError(b.conn, "agent stream closed")
				_ = b.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "agent stream closed"), time.Now().Add(time.Second))
				return
			}
			if len(frame) == 0 {
				continue
			}
			encoded := base64.StdEncoding.EncodeToString(frame)
			if err := b.writeMessage(NewOutputMessage(encoded)); err != nil {
				slog.Debug("terminal: write output failed", "session_id", b.sessionID, "err", err)
				return
			}
			if b.recorder != nil {
				if err := b.recorder.WriteOutput(frame); err != nil {
					slog.Error("terminal: record output failed", "session_id", b.sessionID, "err", err)
				}
			}
		}
	}
}

func (b *bridge) writeMessage(msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}
	return b.conn.WriteMessage(websocket.TextMessage, data)
}
