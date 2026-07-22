package terminal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"beacon/pkg/contextx"
	rpcSchema "common/rpc/schema"

	"github.com/gofiber/contrib/websocket"
)

// fakeConn simulates a websocket.Conn for bridge testing.
type fakeConn struct {
	mu       sync.Mutex
	closed   bool
	sent     [][]byte
	controls []int
	input    chan []byte // messages to be read by the bridge
	writeErr error
}

func newFakeConn() *fakeConn {
	return &fakeConn{input: make(chan []byte, 16)}
}

func (f *fakeConn) pushMessage(msg Message) {
	data, _ := json.Marshal(msg)
	f.input <- data
}

func (f *fakeConn) ReadMessage() (int, []byte, error) {
	data, ok := <-f.input
	if !ok {
		return 0, nil, errors.New("closed")
	}
	return 1, data, nil
}

func (f *fakeConn) WriteMessage(messageType int, data []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.writeErr != nil {
		return f.writeErr
	}
	f.sent = append(f.sent, data)
	return nil
}

func (f *fakeConn) WriteControl(messageType int, data []byte, deadline time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.writeErr != nil {
		return f.writeErr
	}
	f.controls = append(f.controls, messageType)
	if messageType == websocket.CloseMessage {
		f.closed = true
	}
	return nil
}

func (f *fakeConn) controlMessages() []int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]int(nil), f.controls...)
}

func (f *fakeConn) sentMessages() []Message {
	f.mu.Lock()
	defer f.mu.Unlock()
	var msgs []Message
	for _, raw := range f.sent {
		var m Message
		_ = json.Unmarshal(raw, &m)
		msgs = append(msgs, m)
	}
	return msgs
}

// fakeRPC implements rpc.Caller for bridge testing.
type fakeRPC struct {
	mu           sync.Mutex
	inputs       [][]byte
	resizes      []rpcSchema.ResizeTerminalArgs
	closed       []rpcSchema.TerminalCloseArgs
	closeAgentID string
}

func (f *fakeRPC) Call(ctx context.Context, method string, args interface{}, reply interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	switch method {
	case "TerminalInput":
		a := args.(rpcSchema.TerminalInputArgs)
		f.inputs = append(f.inputs, a.Data)
	case "ResizeTerminal":
		a := args.(rpcSchema.ResizeTerminalArgs)
		f.resizes = append(f.resizes, a)
	case "TerminalClose":
		a := args.(rpcSchema.TerminalCloseArgs)
		f.closed = append(f.closed, a)
		f.closeAgentID = contextx.FromAgentID(ctx)
	}
	return nil
}

func (f *fakeRPC) StreamCall(ctx context.Context, method string, args interface{}) (<-chan []byte, error) {
	return nil, nil
}

func (f *fakeRPC) Close() error { return nil }

func TestBridge_ForwardsOutputAndInput(t *testing.T) {
	conn := newFakeConn()
	stream := make(chan []byte, 4)
	stream <- []byte("hello")
	close(stream)

	rpcClient := &fakeRPC{}
	recorder := &Recorder{} // closed recorder, no file writes
	recorder.closed = true

	b := newBridge(rpcClient, conn, stream, recorder, "sess-test", "agent-1")

	done := make(chan struct{})
	go func() {
		_ = b.run(context.Background(), 24, 80)
		close(done)
	}()

	// Wait briefly for the output frame to be written.
	deadline := time.After(2 * time.Second)
	for {
		select {
		case <-deadline:
			t.Fatal("timeout waiting for output message")
		default:
		}
		if msgs := conn.sentMessages(); len(msgs) > 0 {
			if msgs[0].Type != string(MessageTypeOutput) {
				t.Fatalf("expected output message, got %v", msgs[0].Type)
			}
			break
		}
		time.Sleep(time.Millisecond * 10)
	}

	// Close input channel to terminate readWebSocket.
	close(conn.input)
	<-done

	// TerminalClose should have been called on shutdown.
	rpcClient.mu.Lock()
	closedCount := len(rpcClient.closed)
	closeSessionID := ""
	if closedCount > 0 {
		closeSessionID = rpcClient.closed[0].SessionID
	}
	rpcClient.mu.Unlock()
	if closedCount == 0 {
		t.Fatal("expected TerminalClose call")
	}
	if closeSessionID != "sess-test" {
		t.Fatalf("unexpected close session id: %v", closeSessionID)
	}
	if rpcClient.closeAgentID != "agent-1" {
		t.Fatalf("close agent id = %q, want agent-1", rpcClient.closeAgentID)
	}
}

func TestBridge_ForwardsInputAndResize(t *testing.T) {
	conn := newFakeConn()
	stream := make(chan []byte, 4) // keep open

	rpcClient := &fakeRPC{}
	recorder := &Recorder{}
	recorder.closed = true

	b := newBridge(rpcClient, conn, stream, recorder, "sess-input", "agent-1")

	done := make(chan struct{})
	go func() {
		_ = b.run(context.Background(), 24, 80)
		close(done)
	}()

	encoded := base64.StdEncoding.EncodeToString([]byte("ls\n"))
	conn.pushMessage(NewInputMessage(encoded))
	conn.pushMessage(NewResizeMessage(30, 120))

	deadline := time.After(2 * time.Second)
	for {
		select {
		case <-deadline:
			t.Fatal("timeout waiting for input/resize forwarding")
		default:
		}
		rpcClient.mu.Lock()
		inputCount := len(rpcClient.inputs)
		resizeCount := len(rpcClient.resizes)
		rpcClient.mu.Unlock()
		if inputCount > 0 && resizeCount > 0 {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}

	rpcClient.mu.Lock()
	firstInput := rpcClient.inputs[0]
	firstResize := rpcClient.resizes[0]
	rpcClient.mu.Unlock()

	if string(firstInput) != "ls\n" {
		t.Fatalf("unexpected input: %q", firstInput)
	}
	if firstResize.Rows != 30 || firstResize.Cols != 120 {
		t.Fatalf("unexpected resize: %+v", firstResize)
	}

	close(conn.input)
	close(stream)
	<-done
}

func TestBridge_PingsIdleWebSocket(t *testing.T) {
	conn := newFakeConn()
	stream := make(chan []byte)
	rpcClient := &fakeRPC{}
	b := newBridge(rpcClient, conn, stream, nil, "sess-idle", "agent-1")
	b.pingEvery = 10 * time.Millisecond

	done := make(chan struct{})
	go func() {
		_ = b.run(context.Background(), 24, 80)
		close(done)
	}()

	deadline := time.After(2 * time.Second)
	for {
		select {
		case <-deadline:
			t.Fatal("timeout waiting for WebSocket ping")
		default:
		}
		if controls := conn.controlMessages(); len(controls) > 0 {
			if controls[0] != websocket.PingMessage {
				t.Fatalf("control message = %d, want WebSocket ping", controls[0])
			}
			break
		}
		time.Sleep(time.Millisecond)
	}

	close(conn.input)
	close(stream)
	<-done
}
