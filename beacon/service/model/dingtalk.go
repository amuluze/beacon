package model

import "gorm.io/gorm"

const DefaultDingTalkConfigKey = "default"

// DingTalk stores the singleton DingTalk group robot configuration.
// Webhook and Secret are credentials and must never be returned directly by APIs.
type DingTalk struct {
	gorm.Model
	Key     string `gorm:"type:varchar(64);uniqueIndex;not null;comment:配置键"`
	Enabled bool   `gorm:"comment:是否启用"`
	Webhook string `gorm:"type:text;comment:群机器人 Webhook"`
	Secret  string `gorm:"type:varchar(255);comment:加签密钥"`
	AtAll   bool   `gorm:"comment:是否提醒所有人"`
}

func (d *DingTalk) TableName() string {
	return "s_dingtalk"
}
