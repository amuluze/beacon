// Package model
// Date: 2026/6/26
// Author: Amu
// Description:
package model

import (
	"time"

	"gorm.io/gorm"
)

// Session stores terminal session metadata and recording file path.
type Session struct {
	gorm.Model
	SessionID string    `gorm:"uniqueIndex;size:128"`
	AgentID   string    `gorm:"index;size:128"`
	UserID    string    `gorm:"size:128"`
	StartedAt time.Time `gorm:"index"`
	EndedAt   *time.Time
	FilePath  string `gorm:"size:512"`
	Status    string `gorm:"size:16;default:active"` // active / closed / failed
	Width     int
	Height    int
}

func (Session) TableName() string { return "s_session" }
