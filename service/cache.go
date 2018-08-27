package service

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

//CacheService 内存缓存服务
type CacheService struct{}

var builtinCache *cache.Cache

func init() {
	builtinCache = cache.New(5*time.Minute, 10*time.Minute)
}
