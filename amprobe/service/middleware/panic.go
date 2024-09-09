// Package middleware
// Date: 2024/3/27 16:29
// Author: Amu
// Description:
package middleware

import (
	"fmt"
	"log/slog"
	"runtime"

	"amprobe/pkg/fiberx"

	"github.com/amuluze/amutool/errors"
	"github.com/gofiber/fiber/v2"
)

var defaultStackTraceBufSize = 2048

func PanicMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, defaultStackTraceBufSize)

				buf = buf[:runtime.Stack(buf, false)]
				data := fmt.Sprintf("panic: %v\n%s\n", r, buf)
				slog.Error("panic", "panic", slog.String("data", data))
				_ = fiberx.Failure(ctx, errors.New500Error(data))
				return
			}
		}()
		return ctx.Next()
	}
}
