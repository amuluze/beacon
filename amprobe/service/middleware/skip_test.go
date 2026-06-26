// Package middleware
// Date: 2026/6/26
// Author: Amu
// Description: unit tests for skipper middleware logic
package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAllowPathPrefixSkipper(t *testing.T) {
	tests := []struct {
		name     string
		prefixes []string
		path     string
		want     bool
	}{
		{
			name:     "matching prefix",
			prefixes: []string{"/api/v1/auth"},
			path:     "/api/v1/auth/login",
			want:     true,
		},
		{
			name:     "non-matching prefix",
			prefixes: []string{"/api/v1/auth"},
			path:     "/api/v1/host/info",
			want:     false,
		},
		{
			name:     "empty prefixes",
			prefixes: []string{},
			path:     "/any/path",
			want:     false,
		},
		{
			name:     "multiple prefixes first match",
			prefixes: []string{"/api/v1/auth", "/api/v1/index"},
			path:     "/api/v1/auth/login",
			want:     true,
		},
		{
			name:     "multiple prefixes second match",
			prefixes: []string{"/api/v1/auth", "/api/v1/index"},
			path:     "/api/v1/index/index",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skipper := AllowPathPrefixSkipper(tt.prefixes...)

			var result bool
			app := fiber.New()
			app.Use(func(c *fiber.Ctx) error {
				result = skipper(c)
				return c.SendStatus(200)
			})
			req := httptest.NewRequest("GET", tt.path, nil)
			_, _ = app.Test(req, -1)

			if result != tt.want {
				t.Errorf("AllowPathPrefixSkipper(%v) on %q = %v, want %v",
					tt.prefixes, tt.path, result, tt.want)
			}
		})
	}
}

func TestAllowPathPrefixNoSkipper(t *testing.T) {
	tests := []struct {
		name     string
		prefixes []string
		path     string
		want     bool
	}{
		{
			name:     "matching prefix returns false",
			prefixes: []string{"/api/v1/auth"},
			path:     "/api/v1/auth/login",
			want:     false,
		},
		{
			name:     "non-matching prefix returns true",
			prefixes: []string{"/api/v1/auth"},
			path:     "/api/v1/host/info",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skipper := AllowPathPrefixNoSkipper(tt.prefixes...)

			var result bool
			app := fiber.New()
			app.Use(func(c *fiber.Ctx) error {
				result = skipper(c)
				return c.SendStatus(200)
			})
			req := httptest.NewRequest("GET", tt.path, nil)
			_, _ = app.Test(req, -1)

			if result != tt.want {
				t.Errorf("AllowPathPrefixNoSkipper(%v) on %q = %v, want %v",
					tt.prefixes, tt.path, result, tt.want)
			}
		})
	}
}

func TestJoinRouter(t *testing.T) {
	tests := []struct {
		method string
		path   string
		want   string
	}{
		{"GET", "/api/v1/auth", "GET/api/v1/auth"},
		{"POST", "api/v1/host", "POST/api/v1/host"},
		{"get", "/test", "GET/test"},
	}

	for _, tt := range tests {
		t.Run(tt.method+tt.path, func(t *testing.T) {
			result := JoinRouter(tt.method, tt.path)
			if result != tt.want {
				t.Errorf("JoinRouter(%q, %q) = %q, want %q",
					tt.method, tt.path, result, tt.want)
			}
		})
	}
}

func TestSkipHandler(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		skippers []SkipperFunc
		want     bool
	}{
		{
			name:     "no skippers",
			path:     "/api/v1/auth/login",
			skippers: nil,
			want:     false,
		},
		{
			name:     "one matching skipper",
			path:     "/api/v1/auth/login",
			skippers: []SkipperFunc{AllowPathPrefixSkipper("/api/v1/auth")},
			want:     true,
		},
		{
			name: "all non-matching skippers",
			path: "/api/v1/auth/login",
			skippers: []SkipperFunc{
				AllowPathPrefixSkipper("/api/v1/host"),
				AllowPathPrefixSkipper("/api/v1/container"),
			},
			want: false,
		},
		{
			name: "first matches short-circuits",
			path: "/api/v1/auth/login",
			skippers: []SkipperFunc{
				AllowPathPrefixSkipper("/api/v1/auth"),
				AllowPathPrefixSkipper("/api/v1/host"),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			app := fiber.New()
			app.Use(func(c *fiber.Ctx) error {
				result = SkipHandler(c, tt.skippers...)
				return c.SendStatus(200)
			})
			req := httptest.NewRequest("GET", tt.path, nil)
			_, _ = app.Test(req, -1)

			if result != tt.want {
				t.Errorf("SkipHandler = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestAllowMethodAndPathPrefixSkipper(t *testing.T) {
	skipper := AllowMethodAndPathPrefixSkipper("POST/api/v1/auth")

	tests := []struct {
		method string
		path   string
		want   bool
	}{
		{"POST", "/api/v1/auth/login", true},
		{"GET", "/api/v1/auth/login", false},
		{"POST", "/api/v1/host/info", false},
	}

	for _, tt := range tests {
		t.Run(tt.method+tt.path, func(t *testing.T) {
			var result bool
			app := fiber.New()
			app.Use(func(c *fiber.Ctx) error {
				result = skipper(c)
				return c.SendStatus(200)
			})
			req := httptest.NewRequest(tt.method, tt.path, nil)
			_, _ = app.Test(req, -1)

			if result != tt.want {
				t.Errorf("AllowMethodAndPathPrefixSkipper on %s %s = %v, want %v",
					tt.method, tt.path, result, tt.want)
			}
		})
	}
}
