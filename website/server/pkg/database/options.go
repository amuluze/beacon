// Package database
// Date       : 2024/8/22 14:25
// Author     : Amu
// Description:
package database

type Option func(*option)

type option struct {
	Debug        bool
	Type         string
	Host         string
	Port         string
	UserName     string
	Password     string
	DBName       string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
}

func WithDebug(debug bool) Option {
	return func(o *option) {
		o.Debug = debug
	}
}

func WithType(t string) Option {
	return func(o *option) {
		o.Type = t
	}
}

func WithHost(host string) Option {
	return func(o *option) {
		o.Host = host
	}
}

func WithPort(port string) Option {
	return func(o *option) {
		o.Port = port
	}
}

func WithUsername(username string) Option {
	return func(o *option) {
		o.UserName = username
	}
}

func WithPassword(password string) Option {
	return func(o *option) {
		o.Password = password
	}
}

func WithDBName(name string) Option {
	return func(o *option) {
		o.DBName = name
	}
}

func WithMaxLifetime(lifetime int) Option {
	return func(o *option) {
		o.MaxLifetime = lifetime
	}
}

func WithMaxOpenConns(maxOpen int) Option {
	return func(o *option) {
		o.MaxOpenConns = maxOpen
	}
}

func WithMaxIdleConns(maxIdle int) Option {
	return func(o *option) {
		o.MaxIdleConns = maxIdle
	}
}
