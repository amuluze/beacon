// Package schema
// Date       : 2024/9/4 15:14
// Author     : Amu
// Description:
package schema

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Remark    string `json:"remark"`
	IsAdmin   int    `json:"isAdmin"`
	Status    int    `json:"status"`
	Roles     []Role `json:"roles"`
	CreatedAt string `json:"created_at"`
}

type Role struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Remark    string     `json:"remark"`
	Status    int        `json:"status"`
	Resources []Resource `json:"resources"`
	CreatedAt string     `json:"created_at"`
}

type Resource struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Method string `json:"method"`
	Status int    `json:"status"`
}

type UserQueryArgs struct {
	ID        string   `json:"id,omitempty"`
	IDs       []string `json:"ids,omitempty"`
	Username  string   `json:"username,omitempty"`
	Usernames []string `json:"usernames,omitempty"`
	Status    int      `json:"status,omitempty"`
}

type UserQueryReply struct {
	Data  []User `json:"data"`
	Total int    `json:"total"`
}

type UserCreateArgs struct {
	Username string   `json:"username" validate:"required"`
	Password string   `json:"password" validate:"required"`
	Remark   string   `json:"remark,omitempty"`
	Status   int      `json:"status" validate:"required"`
	IsAdmin  int      `json:"isAdmin" validate:"required"`
	RoleIDs  []string `json:"role_ids" validate:"required"`
}

type UserUpdateArgs struct {
	ID       string   `json:"id" validate:"required"`
	Username string   `json:"username,omitempty"`
	Remark   string   `json:"remark,omitempty"`
	Status   int      `json:"status,omitempty"`
	IsAdmin  int      `json:"isAdmin,omitempty"`
	RoleIDs  []string `json:"role_ids,omitempty"`
}

type UserDeleteArgs struct {
	IDs []string `json:"ids,omitempty"`
}

type RoleQueryArgs struct {
	ID     string   `json:"id,omitempty"`
	IDs    []string `json:"ids,omitempty"`
	Name   string   `json:"name,omitempty"`
	Names  []string `json:"names,omitempty"`
	Status int      `json:"status,omitempty"`
}

type RoleQueryReply struct {
	Data  []Role `json:"data"`
	Total int    `json:"total"`
}

type RoleCreateArgs struct {
	Name        string   `json:"name" validate:"required"`
	Remark      string   `json:"remark,omitempty"`
	Status      int      `json:"status" validate:"required"`
	ResourceIDs []string `json:"resource_ids,omitempty"`
}

type RoleUpdateArgs struct {
	ID          string   `json:"id" validate:"required"`
	Name        string   `json:"name,omitempty"`
	Remark      string   `json:"remark,omitempty"`
	Status      int      `json:"status,omitempty"`
	ResourceIDs []string `json:"resource_ids,omitempty"`
}

type RoleDeleteArgs struct {
	IDs []string `json:"ids,omitempty"`
}

type ResourceQueryArgs struct {
	ID string `json:"id,omitempty"`
}

type ResourceQueryReply struct {
	Data  []Resource `json:"data"`
	Total int        `json:"total"`
}
