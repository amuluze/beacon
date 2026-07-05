// Package model
// Date: 2023/6/7 13:00
// Author: Amu
// Description:
package model

import (
	"gorm.io/gorm"
	"time"
)

type Licenses []*License

type License struct {
	gorm.Model
	DeletedAt *time.Time `sql:"index"`
	Content   string     `gorm:"not null"`
}

func (a License) TableName() string {
	return "license"
}
