// Package model
// Date:   2024/10/14 15:54
// Author: Amu
// Description:
package model

import "gorm.io/gorm"

type Mails []Mail

type Mail struct {
	gorm.Model
	Server   string `gorm:"type:varchar(255);comment:邮箱服务器"`
	Port     int    `gorm:"comment:端口"`
	Sender   string `gorm:"type:varchar(255);comment:发信账号"`
	Password string `gorm:"type:varchar(255);comment:密码"`
	Receiver string `gorm:"type:varchar(255);comment:邮箱服务器"`
}

func (d *Mail) TableName() string {
	return "s_mail"
}
