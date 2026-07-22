// Package fiberx
// Date: 2024/3/6 13:15
// Author: Amu
// Description:
package fiberx

import (
	"net/http"
	"strings"

	"beacon/pkg/errors"
	"github.com/gofiber/fiber/v2"
)

// GetToken Get jwt token from header (Authorization: Bearer xxx)
func GetToken(c *fiber.Ctx) string {
	var token string
	auth := c.Get("Authorization")
	prefix := "Bearer "
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token
}

// GetWebSocketToken returns the access token used during a WebSocket
// handshake. Browsers cannot attach an Authorization header to WebSocket
// requests, so the terminal client may use the token query parameter instead.
// Header authentication takes precedence when both forms are present.
func GetWebSocketToken(c *fiber.Ctx) string {
	if token := GetToken(c); token != "" {
		return token
	}
	return c.Query("token")
}

// Success response.status = 200
func Success(c *fiber.Ctx, v interface{}) error {
	return ReturnJson(c, http.StatusOK, v)
}

type FailedResponse struct {
	Err string `json:"err"` // 响应错误，来自 service 层的错误信息
	Msg string `json:"msg"` // 错误消息，来自 api 层的错误信息
}

// Failure response.status = 400
func Failure(c *fiber.Ctx, err errors.Error) error {
	return ReturnJson(c, err.Status, &FailedResponse{Err: err.Err, Msg: err.Msg})
}

// Unauthorized response.status = 401 when token error or token is fail
func Unauthorized(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusUnauthorized)
}

// NoContent response.status = 204
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNoContent)
}

// Forbidden response.status = 403 when permission error
func Forbidden(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusForbidden)
}

func ReturnJson(c *fiber.Ctx, status int, v interface{}) error {
	c.Status(status)
	return c.JSON(v)
}

// ServiceError converts a Go error to errors.Error using FromError.
// If the input is already an errors.Error, it is returned as-is.
func ServiceError(err error) errors.Error {
	if e, ok := err.(errors.Error); ok {
		return e
	}
	if e, ok := err.(*errors.Error); ok {
		return *e
	}
	return errors.FromError(err)
}

// ParseQuery Parse query parameter to struct
func ParseQuery(c *fiber.Ctx, obj interface{}) error {
	return c.QueryParser(obj)
}

func ParseBody(c *fiber.Ctx, obj interface{}) error {
	return c.BodyParser(obj)
}
