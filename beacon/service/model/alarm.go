// Package model
// Date:   2024/10/14 17:21
// Author: Amu
// Description:
package model

import "gorm.io/gorm"

const (
	DefaultAlarmMessageSettingKey = "default"
	DefaultAlarmMessageTitle      = "Beacon 主机告警"
)

type AlarmThreshold struct {
	gorm.Model
	Type      string `gorm:"type:varchar(255);comment:告警类型"`
	Duration  int    `gorm:"comment:取值间隔"`
	Threshold int    `gorm:"comment:阈值"`
}

func (d *AlarmThreshold) TableName() string {
	return "s_alarm_threshold"
}

type AlarmMessageSetting struct {
	gorm.Model
	Key          string `gorm:"type:varchar(64);uniqueIndex;not null;comment:配置键"`
	MessageTitle string `gorm:"type:varchar(128);not null;comment:告警消息标题"`
}

func (d *AlarmMessageSetting) TableName() string {
	return "s_alarm_message_setting"
}
