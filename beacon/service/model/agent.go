// Package model
package model

import (
	"time"

	"gorm.io/gorm"
)

// Agent stores agent registration and status.
type Agent struct {
	gorm.Model
	AgentID  string    `gorm:"uniqueIndex;size:128" json:"agent_id"`
	Hostname string    `gorm:"size:128" json:"hostname"`
	OS       string    `gorm:"size:64" json:"os"`
	Arch     string    `gorm:"size:32" json:"arch"`
	Version  string    `gorm:"size:32" json:"version"`
	LastSeen time.Time `gorm:"index" json:"last_seen"`
	Status   string    `gorm:"size:16;default:offline" json:"status"` // online / offline
}

func (Agent) TableName() string { return "s_agent" }
