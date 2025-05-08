// Package jwtauth
// Date: 2024/3/27 15:07
// Author: Amu
// Description:
package jwtauth

import (
	"amprobe/pkg/auth"

	"github.com/golang-jwt/jwt"
)

const defaultKey = "amprobe"

var defaultOptions = options{
	tokenType:     "Bearer",
	expired:       21600,
	signingMethod: jwt.SigningMethodES512,
	signingKey:    []byte(defaultKey),
	keyfunc: func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidToken
		}
		return []byte(defaultKey), nil
	},
}

type options struct {
	signingMethod  jwt.SigningMethod
	signingKey     interface{}
	keyfunc        jwt.Keyfunc
	expired        int
	refreshExpired int
	tokenType      string
}

type Option func(*options)

func SetSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

func SetSigningKey(key interface{}) Option {
	return func(o *options) {
		o.signingKey = key
	}
}

func SetKeyfunc(keyFunc jwt.Keyfunc) Option {
	return func(o *options) {
		o.keyfunc = keyFunc
	}
}

func SetExpired(expired int) Option {
	return func(o *options) {
		o.expired = expired
	}
}

func SetRefreshExpired(expired int) Option {
	return func(o *options) {
		o.refreshExpired = expired
	}
}
