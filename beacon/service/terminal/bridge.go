package terminal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"beacon/pkg/contextx"
	"beacon/pkg/rpc"
	rpcSchema "common/rpc/schema"

	"github.com/gofiber/contrib/websocket"
)

// Connection abstracts the WebSocket methods the bridge needs, enabling tests
// to inject a fake connection instead of relying on *websocket.Conn.
type Connection interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	WriteControl(messageType int, data []byte, deadline time.Time) error
	// SetReadDeadline 设置底层读超时；用于兜底回收半开 TCP 上的阻塞读 goroutine。
	SetReadDeadline(t time.Time) error
	// SetWriteDeadline 设置底层写超时；用于兜底回收半开 TCP 上的阻塞写 goroutine。
	SetWriteDeadline(t time.Time) error
}

// bridge forwards data between a WebSocket connection and an Agent PTY stream.
type bridge struct {
	rpcClient   rpc.Caller
	conn        Connection
	stream      <-chan []byte
	recorder    *Recorder
	sessionID   string
	agentID     string
	pingEvery   time.Duration
	writeTimeout time.Duration
}

const (
	terminalPingInterval = 30 * time.Second
	// terminalReadDeadline 是 WebSocket 读超时兜底。半开 TCP（连接还在但客户端不读不写）
	// 会让 ReadMessage 永久阻塞；这里设一个足够长、覆盖正常 idle 终端的阈值，
	// 超时即视为连接死亡并 cancel 整个 bridge，释放 goroutine 与 server 端 stream 资源。
	// 正常交互中每次读到消息都会重置该 deadline，不会误杀活跃会话。
	terminalReadDeadline = 10 * time.Minute
	// terminalWriteDeadline 是单次 WebSocket 写超时，防止半开 TCP 写永久阻塞。
	terminalWriteDeadline = 10 * time.Second
)

func newBridge(rpcClient rpc.Caller, conn Connection, stream <-chan []byte, recorder *Recorder, sessionID, agentID string) *bridge {
	return &bridge{
		rpcClient:    rpcClient,
		conn:         conn,
		stream:       stream,
		recorder:     recorder,
		sessionID:    sessionID,
		agentID:      agentID,
		pingEvery:    terminalPingInterval,
		writeTimeout: terminalWriteDeadline,
	}
}

func (b *bridge) run(ctx context.Context, rows, cols int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		// Notify Agent to release PTY resources.
		closeCtx := contextx.NewAgentID(context.Background(), b.agentID)
		closeCtx, closeCancel := context.WithTimeout(closeCtx, 5*time.Second)
		defer closeCancel()
		_ = b.rpcClient.Call(closeCtx, "TerminalClose", rpcSchema.TerminalCloseArgs{
			SessionID: b.sessionID,
		}, &rpcSchema.TerminalCloseReply{})
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	// 兜底读超时：防止半开 TCP 导致 readWebSocket 永久阻塞。
	// readWebSocket 每次成功读取会重置该 deadline。
	if err := b.conn.SetReadDeadline(time.Now().Add(terminalReadDeadline)); err != nil {
		slog.Debug("terminal: set initial read deadline failed", "session_id", b.sessionID, "err", err)
	}

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
		// 收到任意消息即重置读超时，活跃会话不会被误杀。
		if err := b.conn.SetReadDeadline(time.Now().Add(terminalReadDeadline)); err != nil {
			slog.Debug("terminal: reset read deadline failed", "session_id", b.sessionID, "err", err)
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
	pingTicker := time.NewTicker(b.pingEvery)
	defer pingTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pingTicker.C:
			// Keep idle terminal sessions alive through reverse proxies. Browsers
			// answer WebSocket ping frames automatically, so this remains outside
			// the application-level terminal message protocol.
			if err := b.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second)); err != nil {
				slog.Debug("terminal: websocket ping failed", "session_id", b.sessionID, "err", err)
				return
			}
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
	// 写超时兜底：半开 TCP 下 WriteMessage 会永久阻塞，导致 readTunnel goroutine 卡死，
	// server 端 stream 与缓冲无法释放。
	if err := b.conn.SetWriteDeadline(time.Now().Add(b.writeTimeout)); err != nil {
		return fmt.Errorf("set write deadline: %w", err)
	}
	return b.conn.WriteMessage(websocket.TextMessage, data)
}
