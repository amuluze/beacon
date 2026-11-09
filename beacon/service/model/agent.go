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

// DefaultAgentAutoUpdateSettingKey 是 Agent 自动更新开关的固定配置键。
const DefaultAgentAutoUpdateSettingKey = "default"

// AgentSetting 保存管理员在控制台维护的 Agent 级设置。
// 记录缺失表示使用产品默认值（自动更新开启）。
type AgentSetting struct {
	gorm.Model
	Key   string `gorm:"type:varchar(64);uniqueIndex;not null;comment:配置键"`
	Value bool   `gorm:"not null;comment:开关值"`
}

func (AgentSetting) TableName() string { return "s_agent_setting" }
