// Package schema
// Date:   2024/10/14 16:20
// Author: Amu
// Description:
package schema

type Mail struct {
	ID       uint   `json:"id"`
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}

type MailCreateArgs struct {
	Server   string `json:"server" validate:"required,gte=1,lte=256"`
	Port     int    `json:"port" validate:"required,gte=1,lte=65535"`
	Sender   string `json:"sender" validate:"required,email,lte=256"`
	Password string `json:"password" validate:"required,gte=1,lte=128"`
	Receiver string `json:"receiver" validate:"omitempty,email,lte=256"`
}

type MailUpdateArgs struct {
	ID       uint   `json:"id" validate:"required"`
	Server   string `json:"server" validate:"lte=256"`
	Port     int    `json:"port" validate:"omitempty,gte=1,lte=65535"`
	Sender   string `json:"sender" validate:"omitempty,email,lte=256"`
	Password string `json:"password" validate:"lte=128"`
	Receiver string `json:"receiver" validate:"omitempty,email,lte=256"`
}

type MailDeleteArgs struct {
	ID uint `json:"id" validate:"required"`
}

type MailTestArgs struct {
	Receiver string `json:"receiver" validate:"required"`
}
