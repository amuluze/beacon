// Package errors
// Date: 2026/6/26
// Author: Amu
// Description: unit tests for error types and factory functions
package errors

import (
	"testing"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  Error
		want string
	}{
		{"with Err field", Error{Err: "detail message", Msg: "user message", Status: 400}, "detail message"},
		{"without Err field", Error{Msg: "user message", Status: 404}, "user message"},
		{"empty Error", Error{Status: 500}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNew400Error(t *testing.T) {
	err := New400Error("bad input")
	if err.Status != 400 {
		t.Errorf("Status = %d, want 400", err.Status)
	}
	if err.Err != "bad input" {
		t.Errorf("Err = %q, want %q", err.Err, "bad input")
	}
	if err.Msg != BadRequest {
		t.Errorf("Msg = %q, want %q", err.Msg, BadRequest)
	}
}

func TestNew401Error(t *testing.T) {
	err := New401Error("invalid credentials")
	if err.Status != 401 {
		t.Errorf("Status = %d, want 401", err.Status)
	}
	if err.Err != "invalid credentials" {
		t.Errorf("Err = %q, want %q", err.Err, "invalid credentials")
	}
	if err.Msg != InvalidToken {
		t.Errorf("Msg = %q, want %q", err.Msg, InvalidToken)
	}
}

func TestNew404Error(t *testing.T) {
	err := New404Error("resource not found")
	if err.Status != 404 {
		t.Errorf("Status = %d, want 404", err.Status)
	}
	if err.Err != "resource not found" {
		t.Errorf("Err = %q, want %q", err.Err, "resource not found")
	}
	if err.Msg != NotFound {
		t.Errorf("Msg = %q, want %q", err.Msg, NotFound)
	}
}

func TestNew409Error(t *testing.T) {
	err := New409Error("duplicate entry")
	if err.Status != 409 {
		t.Errorf("Status = %d, want 409", err.Status)
	}
	if err.Err != "duplicate entry" {
		t.Errorf("Err = %q, want %q", err.Err, "duplicate entry")
	}
	if err.Msg != Conflict {
		t.Errorf("Msg = %q, want %q", err.Msg, Conflict)
	}
}

func TestNew500Error(t *testing.T) {
	err := New500Error("internal failure")
	if err.Status != 500 {
		t.Errorf("Status = %d, want 500", err.Status)
	}
	if err.Err != "internal failure" {
		t.Errorf("Err = %q, want %q", err.Err, "internal failure")
	}
	if err.Msg != InternalServerError {
		t.Errorf("Msg = %q, want %q", err.Msg, InternalServerError)
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name   string
		err    Error
		status int
		msg    string
	}{
		{"UnauthorizedError", UnauthorizedError, 401, InvalidToken},
		{"ForbiddenError", ForbiddenError, 403, Forbidden},
		{"NotFoundError", NotFoundError, 404, NotFound},
		{"MethodNotAllowError", MethodNotAllowError, 405, MethodNotAllow},
		{"TooManyRequestsError", TooManyRequestsError, 429, TooManyRequests},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Status != tt.status {
				t.Errorf("Status = %d, want %d", tt.err.Status, tt.status)
			}
			if tt.err.Msg != tt.msg {
				t.Errorf("Msg = %q, want %q", tt.err.Msg, tt.msg)
			}
			if tt.err.Err != "" {
				t.Errorf("Err = %q, want empty", tt.err.Err)
			}
		})
	}
}
