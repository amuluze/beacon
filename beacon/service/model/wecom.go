package model

import "gorm.io/gorm"

const DefaultWeComConfigKey = "default"

// WeCom stores the singleton WeCom (企业微信) group robot configuration.
// Webhook is a credential and must never be returned directly by APIs.
type WeCom struct {
	gorm.Model
	Key     string `gorm:"type:varchar(64);uniqueIndex;not null;comment:配置键"`
	Enabled bool   `gorm:"comment:是否启用"`
	Webhook string `gorm:"type:text;comment:群机器人 Webhook"`
	AtAll   bool   `gorm:"comment:是否提醒所有人"`
}

func (w *WeCom) TableName() string {
	return "s_wecom"
}
