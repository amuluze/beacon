// Package terminal provides WebSocket terminal bridge and asciinema session recording.
package terminal

// MessageType is the type of a WebSocket terminal message.
type MessageType string

const (
	MessageTypeInput  MessageType = "input"
	MessageTypeOutput MessageType = "output"
	MessageTypeResize MessageType = "resize"
	MessageTypeError  MessageType = "error"
)

// Message is the JSON envelope exchanged between browser and Server.
type Message struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
	Rows int    `json:"rows,omitempty"`
	Cols int    `json:"cols,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

// NewInputMessage creates an input message from base64 data.
func NewInputMessage(data string) Message {
	return Message{Type: string(MessageTypeInput), Data: data}
}

// NewOutputMessage creates an output message with base64 data.
func NewOutputMessage(data string) Message {
	return Message{Type: string(MessageTypeOutput), Data: data}
}

// NewResizeMessage creates a resize message.
func NewResizeMessage(rows, cols int) Message {
	return Message{Type: string(MessageTypeResize), Rows: rows, Cols: cols}
}

// NewErrorMessage creates an error message.
func NewErrorMessage(msg string) Message {
	return Message{Type: string(MessageTypeError), Msg: msg}
}
