package terminal

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"beacon/pkg/contextx"
	"beacon/pkg/rpc"
	"beacon/service/model"

	"common/database"
	rpcSchema "common/rpc/schema"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

// Handler upgrades WebSocket connections and bridges them to Agent PTY sessions.
type Handler struct {
	rpcClient        rpc.Caller
	db               *database.DB
	sessionDir       string
	recordingEnabled bool
}

const (
	terminalReadyTimeout       = 3 * time.Second
	terminalReadyRetryInterval = 25 * time.Millisecond
)

// ClientConnection is the WebSocket contract used by the terminal handler.
// Fiber copies request metadata and locals to websocket.Conn during upgrade.
type ClientConnection interface {
	Connection
	Headers(key string, defaultValue ...string) string
	Query(key string, defaultValue ...string) string
	Locals(key string, value ...interface{}) interface{}
}

// NewHandler creates a new terminal WebSocket handler.
func NewHandler(rpcClient rpc.Caller, db *database.DB) *Handler {
	return &Handler{
		rpcClient: rpcClient,
		db:        db,
	}
}

// Handle processes a WebSocket terminal connection.
func (h *Handler) Handle(conn *websocket.Conn) {
	h.handle(conn)
}

func (h *Handler) handle(conn ClientConnection) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	agentID := resolveAgentID(conn)
	if agentID == "" {
		slog.Error("terminal: missing agent_id")
		_ = sendError(conn, "missing agent_id")
		_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "missing agent_id"), time.Now().Add(time.Second))
		return
	}
	ctx = contextx.NewAgentID(ctx, agentID)

	sessionID := uuid.NewString()
	rows, cols := defaultTerminalSize(conn)

	session := &model.Session{
		SessionID: sessionID,
		AgentID:   agentID,
		UserID:    resolveUserID(conn),
		StartedAt: time.Now(),
		Status:    "active",
		Width:     cols,
		Height:    rows,
	}
	if err := h.db.Create(session).Error; err != nil {
		slog.Error("terminal: create session record failed", "session_id", sessionID, "err", err)
		_ = sendError(conn, "failed to create session")
		return
	}

	var recorder *Recorder
	if h.recordingEnabled {
		path, err := CleanSessionPath(h.sessionDir, sessionID)
		if err != nil {
			slog.Error("terminal: invalid session path", "session_id", sessionID, "err", err)
			_ = sendError(conn, "invalid session configuration")
			_ = h.closeSession(sessionID, "failed")
			return
		}
		recorder, err = NewRecorder(path, cols, rows)
		if err != nil {
			slog.Error("terminal: create recorder failed", "session_id", sessionID, "err", err)
			_ = sendError(conn, "failed to start recording")
			_ = h.closeSession(sessionID, "failed")
			return
		}
		session.FilePath = path
		_ = h.db.Model(&model.Session{}).Where("session_id = ?", sessionID).Update("file_path", path).Error
	}

	stream, err := h.rpcClient.StreamCall(ctx, "TerminalSession", rpcSchema.TerminalSessionArgs{
		SessionID: sessionID,
		Shell:     "/bin/bash",
		Rows:      rows,
		Cols:      cols,
	})
	if err != nil {
		slog.Error("terminal: start agent session failed", "session_id", sessionID, "agent_id", agentID, "err", err)
		_ = sendError(conn, fmt.Sprintf("failed to start terminal: %v", err))
		_ = h.closeSession(sessionID, "failed")
		if recorder != nil {
			_ = recorder.Close()
		}
		return
	}
	if err := h.awaitAgentReady(ctx, sessionID, rows, cols); err != nil {
		slog.Error("terminal: agent session did not become ready", "session_id", sessionID, "agent_id", agentID, "err", err)
		_ = sendError(conn, fmt.Sprintf("failed to prepare terminal: %v", err))
		_ = h.closeAgentSession(agentID, sessionID)
		_ = h.closeSession(sessionID, "failed")
		if recorder != nil {
			_ = recorder.Close()
		}
		return
	}
	if err := sendMessage(conn, NewReadyMessage()); err != nil {
		slog.Debug("terminal: send ready message failed", "session_id", sessionID, "err", err)
		_ = h.closeAgentSession(agentID, sessionID)
		_ = h.closeSession(sessionID, "failed")
		if recorder != nil {
			_ = recorder.Close()
		}
		return
	}

	bridge := newBridge(h.rpcClient, conn, stream, recorder, sessionID, agentID)
	if err := bridge.run(ctx, rows, cols); err != nil {
		slog.Debug("terminal: bridge ended", "session_id", sessionID, "err", err)
	}

	_ = h.closeSession(sessionID, "closed")
	if recorder != nil {
		_ = recorder.Close()
	}
}

func (h *Handler) awaitAgentReady(ctx context.Context, sessionID string, rows, cols int) error {
	readyCtx, cancel := context.WithTimeout(ctx, terminalReadyTimeout)
	defer cancel()

	ticker := time.NewTicker(terminalReadyRetryInterval)
	defer ticker.Stop()

	var lastErr error
	for {
		lastErr = h.rpcClient.Call(readyCtx, "ResizeTerminal", rpcSchema.ResizeTerminalArgs{
			SessionID: sessionID,
			Rows:      rows,
			Cols:      cols,
		}, &rpcSchema.ResizeTerminalReply{})
		if lastErr == nil {
			return nil
		}

		select {
		case <-readyCtx.Done():
			return fmt.Errorf("agent PTY readiness timeout: %w", lastErr)
		case <-ticker.C:
		}
	}
}

func (h *Handler) closeAgentSession(agentID, sessionID string) error {
	ctx, cancel := context.WithTimeout(contextx.NewAgentID(context.Background(), agentID), 5*time.Second)
	defer cancel()
	return h.rpcClient.Call(ctx, "TerminalClose", rpcSchema.TerminalCloseArgs{SessionID: sessionID}, &rpcSchema.TerminalCloseReply{})
}

func (h *Handler) closeSession(sessionID, status string) error {
	now := time.Now()
	return h.db.Model(&model.Session{}).Where("session_id = ?", sessionID).Updates(map[string]interface{}{
		"status":   status,
		"ended_at": now,
	}).Error
}

func resolveAgentID(conn ClientConnection) string {
	agentID := conn.Headers("X-Agent-ID")
	if agentID == "" {
		agentID = conn.Query("agent_id")
	}
	return agentID
}

func resolveUserID(conn ClientConnection) string {
	userID, _ := conn.Locals("user_id").(string)
	return userID
}

func defaultTerminalSize(conn ClientConnection) (int, int) {
	rows := conn.Query("rows")
	cols := conn.Query("cols")
	r, _ := strconv.Atoi(rows)
	c, _ := strconv.Atoi(cols)
	if r == 0 {
		r = 24
	}
	if c == 0 {
		c = 80
	}
	return r, c
}

func sendError(conn Connection, msg string) error {
	return sendMessage(conn, NewErrorMessage(msg))
}

func sendMessage(conn Connection, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, data)
}
