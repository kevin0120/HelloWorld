package utils

import (
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/pkg/errors"
)

var ENV_CACHE_SIZE_LIMIT = GetEnvInt("ENV_CACHE_SIZE_LIMIT", 1024*1024*1024*5) //默认5Mb

var ENV_CACHE_DEFAULT_TTL = GetEnvInt("ENV_CACHE_DEFAULT_TTL", 60*10) // 默认TTL 为10分钟

type LRUCache struct {
	cache *ttlcache.Cache
}

type FuncLRUCacheExpirationCallback = ttlcache.ExpireCallback

var glbLRUCache = LRUCache{nil}

func OpenNewLRUCache(ttl time.Duration) (cache *LRUCache, err error) {
	cache = &glbLRUCache
	if cache.cache != nil {
		return
	}
	cc := ttlcache.NewCache()
	cache.cache = cc
	if ttl == 0 {
		ttl = time.Second * time.Duration(ENV_CACHE_DEFAULT_TTL)
	}
	cc.SetCacheSizeLimit(ENV_CACHE_SIZE_LIMIT)
	err = cc.SetTTL(ttl)
	return
}

func (cache *LRUCache) Get(key string) (result interface{}, err error) {
	if cache == nil {
		err = errors.New("Cache do not init")
		return
	}
	return cache.cache.Get(key)
}

func (cache *LRUCache) Set(key string, val interface{}, ttl time.Duration) (err error) {
	if cache == nil {
		err = errors.New("Set Cache Is Empty")
		return
	}
	c := cache.cache
	if ttl == 0 {
		err = c.Set(key, val)

	} else {
		err = c.SetWithTTL(key, val, ttl)
	}
	return
}

func (cache *LRUCache) Del(key string) error {
	return cache.cache.Remove(key)
}

func (cache *LRUCache) SetExpirationCallback(callback FuncLRUCacheExpirationCallback) (err error) {
	if cache == nil {
		err = errors.New("SetExpirationCallback Cache Is Empty")
		return
	}
	cache.cache.SetExpirationCallback(callback)
	return
}

func (cache *LRUCache) Close() (err error) {
	if cache == nil {
		err = errors.New("Set Cache Is Empty")
		return
	}
	c := cache.cache
	err = c.Purge()
	if err != nil {
		return
	}
	err = c.Close()
	return
}
