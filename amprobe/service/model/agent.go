// Package model
package model

import (
	"time"

	"gorm.io/gorm"
)

// Agent stores agent registration and status.
type Agent struct {
	gorm.Model
	AgentID  string    `gorm:"uniqueIndex;size:128"`
	Hostname string    `gorm:"size:128"`
	OS       string    `gorm:"size:64"`
	Arch     string    `gorm:"size:32"`
	Version  string    `gorm:"size:32"`
	LastSeen time.Time `gorm:"index"`
	Status   string    `gorm:"size:16;default:offline"` // online / offline
}

func (Agent) TableName() string { return "s_agent" }
