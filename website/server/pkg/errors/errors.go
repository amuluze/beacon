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
