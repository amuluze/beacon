// Package errors
// Date       : 2024/8/27 10:20
// Author     : Amu
// Description:
package errors

import "github.com/pkg/errors"

const (
	InternalServerError = "Internal server error"
	InvalidToken        = "invalid token"
	MethodNotAllow      = "method not allowed"
	NotFound            = "not found"
	TooManyRequests     = "too many requests"
	Forbidden           = "forbidden"
	BadRequest          = "bad request"
	Conflict            = "conflict"
	ServiceUnavailable  = "service unavailable"
	GatewayTimeout      = "gateway timeout"
)

var (
	Is          = errors.Is
	New         = errors.New
	Wrap        = errors.Wrap
	WithStack   = errors.WithStack
	WithMessage = errors.WithMessage
)

var (
	UnauthorizedError    = newError(401, InvalidToken)
	ForbiddenError       = newError(403, Forbidden)
	NotFoundError        = newError(404, NotFound)
	MethodNotAllowError  = newError(405, MethodNotAllow)
	TooManyRequestsError = newError(429, TooManyRequests)
)

type Error struct {
	Err    string // service 层错误消息
	Msg    string // api 层错误（可读）
	Status int    // 响应状态码
}

func (e Error) Error() string {
	if e.Err != "" {
		return e.Err
	}
	return e.Msg
}

func newError(status int, message string) Error {
	return Error{
		Msg:    message,
		Status: status,
	}
}

func New400Error(error string) Error {
	err := newError(400, BadRequest)
	err.Err = error
	return err
}

func New500Error(error string) Error {
	err := newError(500, InternalServerError)
	err.Err = error
	return err
}

func New401Error(msg string) Error {
	err := newError(401, InvalidToken)
	err.Err = msg
	return err
}

func New404Error(msg string) Error {
	err := newError(404, NotFound)
	err.Err = msg
	return err
}

func New409Error(msg string) Error {
	err := newError(409, Conflict)
	err.Err = msg
	return err
}

// New503Error 表示服务暂不可用（如目标 Agent 离线），用于把领域错误映射为可区分状态码。
func New503Error(msg string) Error {
	err := newError(503, ServiceUnavailable)
	err.Err = msg
	return err
}

// New504Error 表示网关超时（如 RPC 调用超时）。
func New504Error(msg string) Error {
	err := newError(504, GatewayTimeout)
	err.Err = msg
	return err
}
