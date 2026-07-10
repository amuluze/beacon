// Package schema
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package schema

// Audit is the wire representation of an operator / system audit row.
//
// Domain Spec I004 requires the AgentID to travel alongside the human-readable
// Operate message so audit consumers can filter alerts by source Agent.
type Audit struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	AgentID  string `json:"agent_id,omitempty"`
	Operate  string `json:"operate"`
	Created  string `json:"created"`
}

type AuditQueryArgs struct {
	// Type filter is the historical toggle: "system" returns only system
	// rows (Username = "system"), any other value returns operator rows.
	Type string `json:"type,omitempty" validate:"lte=64"`
	// AgentID is an optional filter that limits results to a single Agent.
	// Empty string disables the filter.
	AgentID string `json:"agent_id,omitempty"`
	Page    int    `json:"page" validate:"required,gte=1"`
	Size    int    `json:"size" validate:"required,gt=0,lte=100"`
}

type AuditQueryReply struct {
	Data  []Audit `json:"data"`
	Total int     `json:"total"`
	Page  int     `json:"page"`
	Size  int     `json:"size"`
}
