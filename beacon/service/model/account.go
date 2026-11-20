// Package model
// Date: 2024/3/27 16:49
// Author: Amu
// Description:
package model

import (
	"beacon/pkg/utils/uuid"
	"time"
)

const (
	// DefaultAdminRoleName 与 Prepare 初始化创建的管理员角色名保持一致；
	// 用户绑定该角色时被判定为管理员。
	DefaultAdminRoleName = "管理员"

	// IsAdmin 状态取值与 IsAdmin 字段注释一致：1 表示是管理员，2 表示否。
	IsAdminYes = 1
	IsAdminNo  = 2
)

type Users []*User

type User struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;comment:唯一标识"`
	CreatedAt          time.Time `gorm:"comment:创建时间"`
	UpdatedAt          time.Time `gorm:"comment:更新时间"`
	Username           string    `gorm:"size:255;uniqueIndex;not null;comment:用户名"`
	Password           string    `gorm:"size:128;not null;comment:密码"`
	MustChangePassword bool      `gorm:"not null;default:false;comment:登录后必须修改密码"`
	Remark             string    `gorm:"size:200;comment:备注"`
	IsAdmin            int       `gorm:"default:0;comment:是否是管理员(1:是 2:否)"`
	Status             int       `gorm:"index;comment:状态(1:启用 2:停用)"`
	Roles              Roles     `gorm:"many2many:sys_user_roles;"`
}

func (a User) TableName() string {
	return "sys_user"
}

type Roles []*Role

type Role struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;comment:唯一标识"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
	Name      string    `gorm:"size:64;uniqueIndex;not null;comment:角色名称"`
	Status    int       `gorm:"index;comment:状态(1:启用 2:停用)"`
	Remark    string    `gorm:"size:200;comment:备注"`
	Resources Resources `gorm:"many2many:sys_role_resources;"`
}

func (a Role) TableName() string {
	return "sys_role"
}

type Resources []*Resource

type Resource struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;comment:唯一标识"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
	Name      string    `gorm:"size:64;uniqueIndex;not null;comment:资源名称"`
	Path      string    `gorm:"size:512;comment:资源路径"`
	Method    string    `gorm:"comment:资源请求方法"`
	Status    int       `gorm:"status;not null;comment:资源状态(1: 启用 2: 禁用)"`
}

func (a Resource) TableName() string {
	return "sys_resource"
}
