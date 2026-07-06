// Package errors
// Date       : 2024/8/27 10:20
// Author     : Amu
// Description:
package errors

import (
	"context"
	stderrors "errors"
	"net/http"

	"beacon/pkg/contextx"
	tunnel "common/rpc/tunnel"

	pkgerrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

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
	Is          = pkgerrors.Is
	New         = pkgerrors.New
	Wrap        = pkgerrors.Wrap
	WithStack   = pkgerrors.WithStack
	WithMessage = pkgerrors.WithMessage
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

// New401Error creates a 401 Unauthorized error with the given message.
func New401Error(error string) Error {
	err := newError(401, InvalidToken)
	err.Err = error
	return err
}

// New409Error creates a 409 Conflict error with the given message.
func New409Error(error string) Error {
	err := newError(409, "conflict")
	err.Err = error
	return err
}

func New500Error(error string) Error {
	err := newError(500, InternalServerError)
	err.Err = error
	return err
}

func FromError(err error) Error {
	if err == nil {
		return Error{}
	}

	if stderrors.Is(err, contextx.ErrAgentIDRequired) {
		e := newError(http.StatusBadRequest, "agent id is required")
		e.Err = err.Error()
		return e
	}

	if stderrors.Is(err, contextx.ErrMissingAgentID) {
		e := newError(http.StatusBadRequest, "agent id is missing")
		e.Err = err.Error()
		return e
	}

	if stderrors.Is(err, contextx.ErrInvalidAgentID) {
		e := newError(http.StatusBadRequest, "invalid agent id")
		e.Err = err.Error()
		return e
	}

	var offline *tunnel.AgentOfflineError
	if stderrors.As(err, &offline) {
		e := newError(http.StatusServiceUnavailable, "agent offline")
		e.Err = err.Error()
		return e
	}

	var unauthorized *tunnel.AgentUnauthorizedError
	if stderrors.As(err, &unauthorized) {
		e := newError(http.StatusUnauthorized, "agent unauthorized")
		e.Err = err.Error()
		return e
	}

	var duplicate *tunnel.DuplicateAgentError
	if stderrors.As(err, &duplicate) {
		e := newError(http.StatusConflict, "agent already connected")
		e.Err = err.Error()
		return e
	}

	var invalidAgent *tunnel.InvalidAgentIDError
	if stderrors.As(err, &invalidAgent) {
		e := newError(http.StatusBadRequest, "invalid agent")
		e.Err = err.Error()
		return e
	}

	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		e := newError(http.StatusNotFound, "not found")
		e.Err = err.Error()
		return e
	}

	if stderrors.Is(err, context.DeadlineExceeded) {
		e := newError(http.StatusGatewayTimeout, "upstream timeout")
		e.Err = err.Error()
		return e
	}

	return New500Error(err.Error())
}
