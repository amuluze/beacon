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

// AgentAutoUpdateSetting 是控制台可见的 Agent 自动更新状态。
// deployment_enabled 表示部署级门禁（AgentUpdate.Enable），关闭时控制台开关不可用。
type AgentAutoUpdateSetting struct {
	AutoUpdate        bool `json:"auto_update"`
	DeploymentEnabled bool `json:"deployment_enabled"`
}

// AgentAutoUpdateSettingArgs 是自动更新开关的写入请求；指针用于区分未传字段与显式 false。
type AgentAutoUpdateSettingArgs struct {
	AutoUpdate *bool `json:"auto_update" validate:"required"`
}
