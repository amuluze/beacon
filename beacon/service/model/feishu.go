package model

import "gorm.io/gorm"

const DefaultFeishuConfigKey = "default"

// Feishu stores the singleton Feishu group robot configuration.
// Webhook and Secret are credentials and must never be returned directly by APIs.
type Feishu struct {
	gorm.Model
	Key     string `gorm:"type:varchar(64);uniqueIndex;not null;comment:配置键"`
	Enabled bool   `gorm:"comment:是否启用"`
	Webhook string `gorm:"type:text;comment:群机器人 Webhook"`
	Secret  string `gorm:"type:text;comment:加密签名密钥"`
}

func (f *Feishu) TableName() string {
	return "s_feishu"
}
