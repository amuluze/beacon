// Package fiberx
// Date: 2024/3/6 13:15
// Author: Amu
// Description:
package fiberx

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"amprobe/pkg/contextx"
	pkgerrors "amprobe/pkg/errors"
	tunnelpkg "common/rpc/tunnel"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

// Success response.status = 200
func Success(c *fiber.Ctx, v interface{}) error {
	return ReturnJson(c, http.StatusOK, v)
}

type FailedResponse struct {
	Err string `json:"err"` // 响应错误，来自 service 层的错误信息
	Msg string `json:"msg"` // 错误消息，来自 api 层的错误信息
}

// Failure response.status = 400
func Failure(c *fiber.Ctx, err pkgerrors.Error) error {
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

// ServiceError wraps a service-layer error into an appropriate HTTP error.
// 优先识别领域错误并映射为语义化状态码（满足 Domain R001：Agent 不可达/未实现
// 必须返回可区分错误，禁止降级为成功空结果或统一 500）：
//   - 已是 errors.Error：原样返回；
//   - Agent 离线（*tunnel.AgentOfflineError）：503；
//   - agent_id 缺失或格式非法：400；
//   - 记录不存在（gorm.ErrRecordNotFound）：404；
//   - 调用超时（context.DeadlineExceeded）：504；
//   - 其他：500。
func ServiceError(err error) pkgerrors.Error {
	var e pkgerrors.Error
	if errors.As(err, &e) {
		return e
	}
	// 目标 Agent 离线：503，便于前端区分"服务故障"与"Agent 不可达"。
	var offlineErr *tunnelpkg.AgentOfflineError
	if errors.As(err, &offlineErr) {
		return pkgerrors.New503Error(offlineErr.Error())
	}
	// agent_id 缺失或格式非法：400（客户端输入问题）。
	if errors.Is(err, contextx.ErrMissingAgentID) || errors.Is(err, contextx.ErrInvalidAgentID) {
		return pkgerrors.New400Error(err.Error())
	}
	// 记录不存在：404（如某 Agent 从未上报数据）。
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return pkgerrors.New404Error(err.Error())
	}
	// RPC/调用超时：504。
	if errors.Is(err, context.DeadlineExceeded) {
		return pkgerrors.New504Error(err.Error())
	}
	return pkgerrors.New500Error(err.Error())
}

// ParseQuery Parse query parameter to struct
func ParseQuery(c *fiber.Ctx, obj interface{}) error {
	return c.QueryParser(obj)
}

func ParseBody(c *fiber.Ctx, obj interface{}) error {
	return c.BodyParser(obj)
}
