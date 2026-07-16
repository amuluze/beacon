// Package fiberx
// Date: 2024/3/6 13:15
// Author: Amu
// Description:
package fiberx

import (
	"net/http"
	"server/pkg/errors"
	"server/pkg/validatex"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetToken Get jwt token from header (Authorization: Bearer xxx)
func GetToken(c *fiber.Ctx) string {
	auth := c.Get("Authorization")
	prefix := "Bearer "
	if strings.HasPrefix(auth, prefix) {
		return auth[len(prefix):]
	}
	return ""
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

// ParseQuery Parse query parameter to struct
func ParseQuery(c *fiber.Ctx, obj interface{}) error {
	return c.QueryParser(obj)
}

func ParseBody(c *fiber.Ctx, obj interface{}) error {
	return c.BodyParser(obj)
}

// ParseBodyValidate 解析请求体后立即做 schema 校验，在系统边界快速失败。
func ParseBodyValidate(c *fiber.Ctx, obj interface{}) error {
	if err := c.BodyParser(obj); err != nil {
		return err
	}
	return validatex.ValidateStruct(obj)
}
