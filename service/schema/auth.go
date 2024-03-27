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
	Token     string `json:"token" description:"访问令牌"`
	TokenType string `json:"token_type" description:"令牌类型"`
	ExpiresAt string `json:"expires_at" description:"过期时间戳"`
}

type PasswordUpdateArgs struct {
	Username    string `json:"username" validate:"required,gte=1" description:"用户名"`
	OldPassword string `json:"old_password" validate:"required,gte=1"  description:"旧密码"`
	NewPassword string `json:"new_password" validate:"required,gte=1"  description:"新密码"`
}

type UserInfo struct {
	ID       string `json:"id" description:"用户 id"`
	Username string `json:"username" description:"用户名"`
	Email    string `json:"email" description:"邮箱"`
	Status   int    `json:"status" description:"状态"`
}
