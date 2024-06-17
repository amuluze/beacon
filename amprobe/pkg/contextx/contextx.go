// Package contextx
// Date: 2024/3/27 16:28
// Author: Amu
// Description:
package contextx

import "context"

type (
	userIDCtx   struct{}
	usernameCtx struct{}
)

func NewUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDCtx{}, userID)
}

func FromUserID(ctx context.Context) string {
	v := ctx.Value(userIDCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameCtx{}, username)
}

func FromUsername(ctx context.Context) string {
	v := ctx.Value(usernameCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
