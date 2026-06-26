package terminal

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"amprobe/pkg/contextx"
	"amprobe/pkg/rpc"
	"amprobe/service/model"

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

// NewHandler creates a new terminal WebSocket handler.
func NewHandler(rpcClient rpc.Caller, db *database.DB, sessionDir string, recordingEnabled bool) *Handler {
	return &Handler{
		rpcClient:        rpcClient,
		db:               db,
		sessionDir:       sessionDir,
		recordingEnabled: recordingEnabled,
	}
}

// Handle processes a WebSocket terminal connection.
func (h *Handler) Handle(conn *websocket.Conn) {
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

	bridge := newBridge(h.rpcClient, conn, stream, recorder, sessionID, agentID)
	if err := bridge.run(ctx, rows, cols); err != nil {
		slog.Debug("terminal: bridge ended", "session_id", sessionID, "err", err)
	}

	_ = h.closeSession(sessionID, "closed")
	if recorder != nil {
		_ = recorder.Close()
	}
}

func (h *Handler) closeSession(sessionID, status string) error {
	now := time.Now()
	return h.db.Model(&model.Session{}).Where("session_id = ?", sessionID).Updates(map[string]interface{}{
		"status":   status,
		"ended_at": now,
	}).Error
}

func resolveAgentID(conn *websocket.Conn) string {
	agentID := conn.Headers("X-Agent-ID")
	if agentID == "" {
		agentID = conn.Query("agent_id")
	}
	return agentID
}

func defaultTerminalSize(conn *websocket.Conn) (int, int) {
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
	m := NewErrorMessage(msg)
	data, _ := json.Marshal(m)
	return conn.WriteMessage(websocket.TextMessage, data)
}
