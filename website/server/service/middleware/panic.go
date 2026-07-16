// Package middleware
// Date: 2024/3/27 16:29
// Author: Amu
// Description:
package middleware

import (
	"fmt"
	"log/slog"
	"runtime"
	"server/pkg/errors"
	"server/pkg/fiberx"

	"github.com/gofiber/fiber/v2"
)

var defaultStackTraceBufSize = 2048

// PanicMiddleware 捕获下游 handler 抛出的 panic，
// 完整栈仅写入服务端日志，响应体只返回固定文案，避免泄露内部路径/符号。
func PanicMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, defaultStackTraceBufSize)
				buf = buf[:runtime.Stack(buf, false)]
				slog.Error("panic recovered",
					slog.Any("panic", r),
					slog.String("stack", fmt.Sprintf("%s", buf)),
				)
				_ = fiberx.Failure(ctx, errors.New500Error(""))
			}
		}()
		return ctx.Next()
	}
}
