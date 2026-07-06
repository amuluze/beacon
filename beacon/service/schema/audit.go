// Package schema
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package schema

type Audit struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Operate  string `json:"operate"`
	Created  string `json:"created"`
}

type AuditQueryArgs struct {
	Type string `json:"type,omitempty" validate:"lte=64"`
	Page int    `json:"page" validate:"required,gte=1"`
	Size int    `json:"size" validate:"required,gt=0,lte=100"`
}

type AuditQueryReply struct {
	Data  []Audit `json:"data"`
	Total int     `json:"total"`
	Page  int     `json:"page"`
	Size  int     `json:"size"`
}
