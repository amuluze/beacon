// Package jwtauth
// Date: 2026/6/26
// Author: Amu
// Description: unit tests for JWT auth options
package jwtauth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestDefaultOptions(t *testing.T) {
	o := defaultOptions

	if o.expired != 21600 {
		t.Errorf("default expired = %d, want 21600", o.expired)
	}
	if o.tokenType != "Bearer" {
		t.Errorf("default tokenType = %q, want %q", o.tokenType, "Bearer")
	}
}

func TestSetExpired(t *testing.T) {
	tests := []struct {
		name string
		opt  Option
		want int
	}{
		{"default", nil, 21600},
		{"override", SetExpired(3600), 3600},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := defaultOptions
			if tt.opt != nil {
				tt.opt(&o)
			}
			if o.expired != tt.want {
				t.Errorf("expired = %d, want %d", o.expired, tt.want)
			}
		})
	}
}

func TestSetRefreshExpired(t *testing.T) {
	o := defaultOptions
	SetRefreshExpired(86400)(&o)
	if o.refreshExpired != 86400 {
		t.Errorf("refreshExpired = %d, want 86400", o.refreshExpired)
	}
}

func TestSetSigningMethod(t *testing.T) {
	o := defaultOptions
	SetSigningMethod(jwt.SigningMethodHS256)(&o)
	if o.signingMethod != jwt.SigningMethodHS256 {
		t.Errorf("signingMethod not updated")
	}
}

func TestSetSigningKey(t *testing.T) {
	o := defaultOptions
	key := []byte("my-secret")
	SetSigningKey(key)(&o)
	if o.signingKey == nil {
		t.Error("signingKey is nil after SetSigningKey")
	}
}

func TestSetKeyfunc(t *testing.T) {
	o := defaultOptions
	called := false
	fn := func(t *jwt.Token) (interface{}, error) {
		called = true
		return []byte("test"), nil
	}
	SetKeyfunc(fn)(&o)
	if o.keyfunc == nil {
		t.Error("keyfunc is nil after SetKeyfunc")
	}
	// Verify the function is usable
	_, _ = o.keyfunc(&jwt.Token{})
	if !called {
		t.Error("keyfunc was not called")
	}
}
