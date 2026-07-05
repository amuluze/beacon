// Package schema
// Date: 2024/3/27 16:36
// Author: Amu
// Description:
package schema

type LoginArgs struct {
	Username string `json:"username" validate:"required,gte=1" description:"用户名"`
	Password string `json:"password" validate:"required,gte=1" description:"用户密码"`
}

type LoginResult struct {
	AccessToken  string `json:"access_token" description:"访问令牌"`
	RefreshToken string `json:"refresh_token" description:"刷新令牌"`
}

type PasswordUpdateArgs struct {
	Username    string `json:"username" validate:"required,gte=1" description:"用户名"`
	OldPassword string `json:"old_password" validate:"required,gte=1"  description:"旧密码"`
	NewPassword string `json:"new_password" validate:"required,gte=1"  description:"新密码"`
}

type UserInfo struct {
	ID       string `json:"id" description:"用户 id"`
	Username string `json:"username" description:"用户名"`
	Status   int    `json:"status" description:"状态"`
	IsAdmin  int    `json:"is_admin" description:"是否是管理员"`
}
