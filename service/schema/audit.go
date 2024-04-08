// Package schema
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package schema

import "time"

type Audit struct {
	ID       uint      `json:"id"`
	Username string    `json:"username"`
	Operate  string    `json:"operate"`
	Created  time.Time `json:"created"`
}

type AuditQueryArgs struct {
	Page int `json:"page" validate:"required"`
	Size int `json:"size" validate:"required,gt=0"`
}

type AuditQueryReply struct {
	Data  []Audit `json:"data"`
	Total int     `json:"total"`
	Page  int     `json:"page"`
	Size  int     `json:"size"`
}
