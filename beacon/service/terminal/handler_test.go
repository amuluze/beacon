package terminal

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"beacon/pkg/contextx"
	"beacon/service/model"
	"common/database"
	rpcSchema "common/rpc/schema"
)

type handlerConn struct {
	*fakeConn
	headers map[string]string
	queries map[string]string
	locals  map[string]interface{}
}

func newHandlerConn() *handlerConn {
	return &handlerConn{
		fakeConn: newFakeConn(),
		headers:  make(map[string]string),
		queries:  make(map[string]string),
		locals:   make(map[string]interface{}),
	}
}

func (c *handlerConn) Headers(key string, defaultValue ...string) string {
	if value, ok := c.headers[key]; ok {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (c *handlerConn) Query(key string, defaultValue ...string) string {
	if value, ok := c.queries[key]; ok {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (c *handlerConn) Locals(key string, value ...interface{}) interface{} {
	if len(value) > 0 {
		c.locals[key] = value[0]
		return value[0]
	}
	return c.locals[key]
}

type handlerRPC struct {
	mu              sync.Mutex
	stream          chan []byte
	streamMethod    string
	streamAgentID   string
	streamArgs      rpcSchema.TerminalSessionArgs
	resizeCalls     int
	closedSessionID string
	closedAgentID   string
}

func (r *handlerRPC) StreamCall(ctx context.Context, method string, args interface{}) (<-chan []byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.streamMethod = method
	r.streamAgentID = contextx.FromAgentID(ctx)
	r.streamArgs = args.(rpcSchema.TerminalSessionArgs)
	return r.stream, nil
}

func (r *handlerRPC) Call(ctx context.Context, method string, args interface{}, reply interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	switch method {
	case "ResizeTerminal":
		r.resizeCalls++
		if r.resizeCalls == 1 {
			return errors.New("session not found")
		}
	case "TerminalClose":
		r.closedSessionID = args.(rpcSchema.TerminalCloseArgs).SessionID
		r.closedAgentID = contextx.FromAgentID(ctx)
	}
	return nil
}

func (r *handlerRPC) Close() error { return nil }

func TestHandlerStartsAgentPTYAndSignalsReady(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "terminal")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.Session)); err != nil {
		t.Fatalf("auto migrate session: %v", err)
	}

	rpcClient := &handlerRPC{stream: make(chan []byte, 1)}
	conn := newHandlerConn()
	conn.queries["agent_id"] = "node-01"
	conn.queries["rows"] = "21"
	conn.queries["cols"] = "139"
	conn.locals["user_id"] = "user-1"

	handler := NewHandler(rpcClient, db)
	done := make(chan struct{})
	go func() {
		handler.handle(conn)
		close(done)
	}()

	deadline := time.After(2 * time.Second)
	for {
		messages := conn.sentMessages()
		if len(messages) > 0 {
			if messages[0].Type != string(MessageTypeReady) {
				t.Fatalf("first websocket message type = %q, want ready", messages[0].Type)
			}
			break
		}
		select {
		case <-deadline:
			t.Fatal("timeout waiting for terminal ready message")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	rpcClient.mu.Lock()
	streamMethod := rpcClient.streamMethod
	streamAgentID := rpcClient.streamAgentID
	streamArgs := rpcClient.streamArgs
	resizeCalls := rpcClient.resizeCalls
	rpcClient.mu.Unlock()
	if streamMethod != "TerminalSession" {
		t.Fatalf("stream method = %q, want TerminalSession", streamMethod)
	}
	if streamAgentID != "node-01" {
		t.Fatalf("stream agent id = %q, want node-01", streamAgentID)
	}
	if streamArgs.Rows != 21 || streamArgs.Cols != 139 {
		t.Fatalf("terminal size = %dx%d, want 139x21", streamArgs.Cols, streamArgs.Rows)
	}
	if resizeCalls < 2 {
		t.Fatalf("resize calls = %d, want readiness retry", resizeCalls)
	}
	sessionID := streamArgs.SessionID

	close(conn.input)
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for terminal handler shutdown")
	}

	var session model.Session
	if err := db.Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		t.Fatalf("load session: %v", err)
	}
	if session.UserID != "user-1" || session.AgentID != "node-01" {
		t.Fatalf("session identity = user:%q agent:%q", session.UserID, session.AgentID)
	}
	if session.Status != "closed" || session.EndedAt == nil {
		t.Fatalf("session status = %q ended_at=%v, want closed with end time", session.Status, session.EndedAt)
	}

	rpcClient.mu.Lock()
	defer rpcClient.mu.Unlock()
	if rpcClient.closedSessionID != sessionID {
		t.Fatalf("closed session id = %q, want %q", rpcClient.closedSessionID, sessionID)
	}
	if rpcClient.closedAgentID != "node-01" {
		t.Fatalf("closed agent id = %q, want node-01", rpcClient.closedAgentID)
	}
}
