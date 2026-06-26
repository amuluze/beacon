package schema

import (
	"encoding/json"
	"testing"
)

func TestTerminalSessionArgs_Marshal(t *testing.T) {
	args := TerminalSessionArgs{
		Shell: "/bin/bash",
		Rows:  24,
		Cols:  80,
	}
	data, err := json.Marshal(args)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var decoded TerminalSessionArgs
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if decoded.Shell != "/bin/bash" || decoded.Rows != 24 || decoded.Cols != 80 {
		t.Errorf("decoded mismatch: %+v", decoded)
	}
}

func TestResizeTerminalArgs_Marshal(t *testing.T) {
	args := ResizeTerminalArgs{SessionID: "sess-1", Rows: 30, Cols: 120}
	data, err := json.Marshal(args)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var decoded ResizeTerminalArgs
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if decoded.SessionID != "sess-1" || decoded.Rows != 30 || decoded.Cols != 120 {
		t.Errorf("decoded mismatch: %+v", decoded)
	}
}

func TestTerminalInputArgs_Marshal(t *testing.T) {
	args := TerminalInputArgs{SessionID: "sess-1", Data: []byte("hello")}
	data, err := json.Marshal(args)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var decoded TerminalInputArgs
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if decoded.SessionID != "sess-1" || string(decoded.Data) != "hello" {
		t.Errorf("decoded mismatch: %+v", decoded)
	}
}
