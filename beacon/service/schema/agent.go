// Package schema
package schema

import "time"

// ── Agent management ──

type AgentItem struct {
	AgentID  string    `json:"agent_id"`
	Hostname string    `json:"hostname"`
	OS       string    `json:"os"`
	Arch     string    `json:"arch"`
	Version  string    `json:"version"`
	LastSeen time.Time `json:"last_seen"`
	Status   string    `json:"status"`
}

type AgentQueryReply struct {
	Data  []AgentItem `json:"data"`
	Total int64       `json:"total"`
}

// ── Join Token management ──

type JoinTokenCreateArgs struct {
	Description string `json:"description"`
}

type JoinTokenCreateReply struct {
	Token       string    `json:"token"`
	Description string    `json:"description"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type JoinTokenDeleteArgs struct {
	Token string `json:"token" validate:"required"`
}
