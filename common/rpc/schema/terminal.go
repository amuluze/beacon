// Package schema
// Date: 2026/6/26
// Author: Amu
// Description:
package schema

// TerminalSessionArgs is sent from Server to Agent to start a PTY shell session.
type TerminalSessionArgs struct {
	SessionID string `json:"session_id"`
	Shell     string `json:"shell"`
	Rows      int    `json:"rows"`
	Cols      int    `json:"cols"`
}

// ResizeTerminalArgs is sent from Server to Agent to resize an active PTY.
type ResizeTerminalArgs struct {
	SessionID string `json:"session_id"`
	Rows      int    `json:"rows"`
	Cols      int    `json:"cols"`
}

// TerminalInputArgs is sent from Server to Agent to deliver user input.
type TerminalInputArgs struct {
	SessionID string `json:"session_id"`
	Data      []byte `json:"data"`
}

// ResizeTerminalReply is the response from Agent after resizing a PTY.
type ResizeTerminalReply struct{}

// TerminalInputReply is the response from Agent after writing input.
type TerminalInputReply struct{}

// TerminalCloseArgs is sent from Server to Agent to close an active PTY session.
type TerminalCloseArgs struct {
	SessionID string `json:"session_id"`
}

// TerminalCloseReply is the response from Agent after closing a session.
type TerminalCloseReply struct{}
