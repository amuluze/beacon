// Package schema
// Date: 2024/3/27 16:36
// Author: Amu
// Description:
package schema

type LoginArgs struct {
	Username string `json:"username" validate:"required,alphanum,gte=1,lte=64" description:"用户名"`
	Password string `json:"password" validate:"required,gte=1,lte=128" description:"用户密码"`
}

type LoginResult struct {
	MustChangePassword bool `json:"must_change_password" description:"是否必须修改密码"`
}

type PasswordUpdateArgs struct {
	Username    string `json:"username" validate:"required,gte=1,lte=64" description:"用户名"`
	OldPassword string `json:"old_password" validate:"required,gte=1,lte=128"  description:"旧密码"`
	NewPassword string `json:"new_password" validate:"required,gte=12,lte=128" description:"新密码至少12字符，且Unicode小写、大写、数字、其他字符至少三类；必须与旧密码不同"`
}

type UserInfo struct {
	ID                 string `json:"id" description:"用户 id"`
	Username           string `json:"username" description:"用户名"`
	Status             int    `json:"status" description:"状态"`
	IsAdmin            int    `json:"is_admin" description:"是否是管理员"`
	MustChangePassword bool   `json:"must_change_password" description:"是否必须修改密码"`
}
