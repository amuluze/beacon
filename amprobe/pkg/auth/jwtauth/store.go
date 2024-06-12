// Package jwtauth
// Date: 2024/3/27 15:03
// Author: Amu
// Description:
package jwtauth

import (
	"fmt"
	"github.com/amuluze/amutool/errors"
	"github.com/patrickmn/go-cache"
	"time"
)

type Storer interface {
	// Set 存储令牌数据，并指定过期时间
	Set(tokenString string, expiration time.Duration) error
	// Check 检查令牌是否存在
	Check(tokenString string) (bool, error)
	// Close 关闭存储
	Close() error
}

type Store struct {
	Storage *cache.Cache
	Prefix  string
}

func (s *Store) wrapperKey(key string) string {
	return fmt.Sprintf("%s%s", s.Prefix, key)
}

func (s *Store) Close() error {
	return nil
}

func (s *Store) Set(tokenString string, expiration time.Duration) error {
	s.Storage.Set(s.wrapperKey(tokenString), "1", expiration)
	return nil
}

func (s *Store) Check(key string) (bool, error) {
	_, found := s.Storage.Get(s.wrapperKey(key))
	if found {
		return found, nil
	}
	return false, errors.New("key not found")
}
