// Package schema
// Date: 2024/3/27 16:36
// Author: Amu
// Description:
package schema

type User struct {
	ID       string `json:"id,string"`                             // 唯一标识
	Username string `json:"username"`                              // 用户名
	Email    string `json:"email"`                                 // 邮箱
	Status   int    `json:"status"  description:"用户状态(1:启用 2:停用)"` // 用户状态(1:启用 2:停用)
}

type Users []*User
