package model

import (
	"testing"
)

func TestSession_TableName(t *testing.T) {
	s := Session{}
	if got := s.TableName(); got != "s_session" {
		t.Errorf("TableName() = %v, want %v", got, "s_session")
	}
}

func TestSession_Fields(t *testing.T) {
	s := Session{
		SessionID: "sess-001",
		AgentID:   "agent-01",
		UserID:    "user-01",
		Status:    "active",
		Width:     80,
		Height:    24,
	}
	if s.SessionID != "sess-001" {
		t.Errorf("unexpected SessionID: %v", s.SessionID)
	}
	if s.Width != 80 || s.Height != 24 {
		t.Errorf("unexpected terminal size: %dx%d", s.Width, s.Height)
	}
}
