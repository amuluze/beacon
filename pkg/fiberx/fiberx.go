// Package fiberx
// Date: 2024/3/6 13:15
// Author: Amu
// Description:
package fiberx

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/amuluze/amutool/errors"
	"github.com/gofiber/fiber/v2"
)

// GetToken Get jwt token from header (Authorization: Bearer xxx)
func GetToken(c *fiber.Ctx) string {
	var token string
	auth := c.GetReqHeaders()["Authorization"][0]
	fmt.Println(">>>>>>>>>", auth)
	prefix := "Bearer "
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token
}

// Success response.status = 200
func Success(c *fiber.Ctx, v interface{}) error {
	return ReturnJson(c, http.StatusOK, v)
}

type FailedResponse struct {
	Err string `json:"err"` // 响应错误，来自 greeter_service 层的错误信息
	Msg string `json:"msg"` // 错误消息，来自 controller 层的错误信息
}

// Failure response.status = 400
func Failure(c *fiber.Ctx, err error) error {
	var res *errors.Error
	if err != nil {
		if e, ok := err.(*errors.Error); ok {
			res = e
		} else {
			res = errors.UnWrapError(errors.ErrInternalServer)
		}
	} else {
		res = errors.UnWrapError(errors.ErrInternalServer)
	}

	if err := res.ERR; err != nil {
		if res.Message == "" {
			res.Message = err.Error()
		}
		// 正常来说api层返回给接口层的错误信息不用打印，在这里会进行统一打印，400-500之间会打印warn 日志
		// 500 会打印error 日志
		if status := res.Status; status >= 400 && status < 500 {
			fmt.Println("4xx")
		} else if status >= 500 {
			fmt.Println("5xx")
		}
	}

	return ReturnJson(c, res.Status, &FailedResponse{Err: res.Message, Msg: res.ERR.Error()})
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

// NotFound response.status = 404
func NotFound(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNotFound)
}

// MethodNotAllowed response.status = 405
func MethodNotAllowed(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusMethodNotAllowed)
}

func ReturnJson(c *fiber.Ctx, status int, v interface{}) error {
	c.Status(status)
	return c.JSON(v)
}

// ParseQuery Parse query parameter to struct
func ParseQuery(c *fiber.Ctx, obj interface{}) error {
	if err := c.QueryParser(obj); err != nil {
		return errors.New400Error(err.Error())
	}
	return nil
}

func ParseBody(c *fiber.Ctx, obj interface{}) error {
	if err := c.BodyParser(obj); err != nil {
		return errors.New400Error(err.Error())
	}
	return nil
}
