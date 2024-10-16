// Package model
// Date:   2024/10/14 17:21
// Author: Amu
// Description:
package model

import "gorm.io/gorm"

type AlarmThreshold struct {
	gorm.Model
	Type      string `gorm:"type:varchar(255)comment:告警类型"`
	Duration  int    `gorm:"comment:取值间隔"`
	Threshold int    `gorm:"comment:阈值"`
}

func (d *AlarmThreshold) TableName() string {
	return "s_alarm_threshold"
}
