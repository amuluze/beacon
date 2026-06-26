// Package schema
// Date: 2026/6/26
// Author: Amu
// Description:
package schema

// TerminalSessionArgs is sent from Server to Agent to start a PTY shell session.
type TerminalSessionArgs struct {
	Shell string `json:"shell"`
	Rows  int    `json:"rows"`
	Cols  int    `json:"cols"`
}

// ResizeTerminalArgs is sent from Server to Agent to resize an active PTY.
type ResizeTerminalArgs struct {
	Rows int `json:"rows"`
	Cols int `json:"cols"`
}

// ResizeTerminalReply is the response from Agent after resizing a PTY.
type ResizeTerminalReply struct{}
