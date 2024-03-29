// Package service
// Date: 2024/3/27 16:06
// Author: Amu
// Description:
package service

import (
	"github.com/amuluze/amprobe/pkg/auth"
	"github.com/amuluze/amprobe/pkg/auth/jwtauth"
	"github.com/golang-jwt/jwt"
	"github.com/patrickmn/go-cache"
	"time"
)

func InitAuthStore(config *Config) (*jwtauth.Store, func(), error) {
	var err error
	authStore := &jwtauth.Store{
		Storage: cache.New(5*time.Minute, 60*time.Second),
		Prefix:  config.Auth.Prefix,
	}
	cleanFunc := func() { err = authStore.Close() }
	return authStore, cleanFunc, err
}

func InitAuth(config *Config, authStore *jwtauth.Store) (auth.Auther, func(), error) {
	var opts []jwtauth.Option
	opts = append(opts, jwtauth.SetExpired(config.Auth.Expired))
	opts = append(opts, jwtauth.SetSigningKey([]byte(config.Auth.SigningKey)))
	opts = append(opts, jwtauth.SetKeyfunc(func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidToken
		}
		return []byte(config.Auth.SigningKey), nil
	}))
	var method jwt.SigningMethod
	switch config.Auth.SigningMethod {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	default:
		method = jwt.SigningMethodHS512
	}
	opts = append(opts, jwtauth.SetSigningMethod(method))
	var err error
	jAuth := jwtauth.New(authStore, opts...)
	cleanFunc := func() {
		err = jAuth.Release()
	}
	return jAuth, cleanFunc, err
}
