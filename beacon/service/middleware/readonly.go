// Package middleware
// Date: 2026/11/17
// Description: read-only deployment guard.
package middleware

import (
	"beacon/pkg/errors"
	"beacon/pkg/fiberx"

	"github.com/gofiber/fiber/v2"
)

// readOnlyAllowedPaths 是只读模式下仍必须放行的非安全方法端点：
//   - 认证生命周期（登录、刷新、登出、改密）不修改监控业务数据；
//   - Agent 上报是服务端数据入口，若被拦截会导致只读展示随时间失真；
//   - 文件下载在 HTTP 语义上是 POST，但不会修改服务端状态。
var readOnlyAllowedPaths = map[string]struct{}{
	"/api/v1/auth/csrf":          {},
	"/api/v1/auth/logout":        {},
	"/api/v1/auth/token_update":  {},
	"/api/v1/host/report":        {},
	"/api/v1/host/file_download": {},
}

// readOnlyAuthenticatedPaths 仅在保留登录的只读模式下放行；匿名会话不能改密。
var readOnlyAuthenticatedPaths = map[string]struct{}{
	"/api/v1/auth/login":       {},
	"/api/v1/auth/pass_update": {},
}

// ReadOnlyMiddleware rejects every mutating request not explicitly allowlisted.
// NoAuth deployments install it unconditionally because anonymous access must
// never imply write access.
func ReadOnlyMiddleware(noAuth bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if isSafeMethod(c.Method()) {
			return c.Next()
		}
		if _, ok := readOnlyAllowedPaths[c.Path()]; ok {
			return c.Next()
		}
		if !noAuth {
			if _, ok := readOnlyAuthenticatedPaths[c.Path()]; ok {
				return c.Next()
			}
		}
		return fiberx.Failure(c, errors.Error{Status: fiber.StatusForbidden, Msg: "read-only mode"})
	}
}
