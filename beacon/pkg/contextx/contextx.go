// Package contextx
// Date: 2024/3/27 16:28
// Author: Amu
// Description:
package contextx

import (
	"context"
	"errors"
	"regexp"
)

var (
	ErrAgentIDRequired = errors.New("agent id is required")
	ErrMissingAgentID  = errors.New("agent id is missing")
	ErrInvalidAgentID  = errors.New("invalid agent id")
)

const maxAgentIDLen = 128

var validAgentIDRe = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// IsValidAgentID validates agent ID format: alphanumeric, dots, dashes, underscores only.
func IsValidAgentID(id string) bool {
	if id == "" || len(id) > maxAgentIDLen {
		return false
	}
	return validAgentIDRe.MatchString(id)
}

// ResolveAgentID extracts and validates agent ID from context.
// Returns ErrMissingAgentID if not present, ErrInvalidAgentID if format is invalid.
func ResolveAgentID(ctx context.Context) (string, error) {
	agentID := FromAgentID(ctx)
	if agentID == "" {
		return "", ErrMissingAgentID
	}
	if !IsValidAgentID(agentID) {
		return "", ErrInvalidAgentID
	}
	return agentID, nil
}

type (
	userIDCtx   struct{}
	usernameCtx struct{}
	agentIDCtx  struct{}
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

func NewAgentID(ctx context.Context, agentID string) context.Context {
	return context.WithValue(ctx, agentIDCtx{}, agentID)
}

func FromAgentID(ctx context.Context) string {
	v := ctx.Value(agentIDCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func RequireAgentID(ctx context.Context) (string, error) {
	agentID := FromAgentID(ctx)
	if agentID == "" {
		return "", ErrAgentIDRequired
	}
	return agentID, nil
}
