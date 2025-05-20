package service

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func NewCache() *cache.Cache {
	return cache.New(5*time.Minute, 10*time.Minute)
}
