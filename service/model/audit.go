// Package model
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package model

import (
	"gorm.io/gorm"
)

type Audits []Audit

type Audit struct {
	gorm.Model
	Username string `gorm:"type:varchar(255);not null"`
	Operate  string `gorm:"type:varchar(255);not null"`
}

func (d *Audit) TableName() string {
	return "s_audit"
}
