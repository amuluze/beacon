// Package model
// Date: 2024/3/27 16:49
// Author: Amu
// Description:
package model

import (
	"github.com/amuluze/amprobe/pkg/utils/uuid"
	"time"
)

type Users []*User

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;comment:唯一标识"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
	Username  string    `gorm:"size:255;uniqueIndex;not null;comment:用户名"`
	Password  string    `gorm:"size:128;not null;comment:密码"`
	Remark    *string   `gorm:"size:200;comment:备注"`
	Status    int       `gorm:"index;default:0;comment:状态(1:启用 0:停用)"`
}

func (a User) TableName() string {
	return "sys_user"
}
