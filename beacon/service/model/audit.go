// Package model
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package model

import (
	"gorm.io/gorm"
)

type Audits []Audit

// Audit represents an operator / system audit log entry.
//
// Domain Spec I004 requires the agent_id to be carried on alarm / state-change
// audit messages so that operators can grep alerts back to a specific Agent.
// AgentID is intentionally nullable to preserve existing rows that were
// recorded before the column was introduced.
type Audit struct {
	gorm.Model
	Username string `gorm:"type:varchar(255);not null"`
	AgentID  string `gorm:"type:varchar(128);index"`
	Operate  string `gorm:"type:varchar(255);not null"`
}

func (d *Audit) TableName() string {
	return "s_audit"
}
