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
	SSLMode      string
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

// WithSSLMode 设置 postgres 的 sslmode；为空时回退 disable，便于生产开启 TLS 到数据库。
func WithSSLMode(mode string) Option {
	return func(o *option) {
		o.SSLMode = mode
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
