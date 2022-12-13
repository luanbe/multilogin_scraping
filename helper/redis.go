package helper

import (
	"context"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisCache interface {
	AddRedisCache(key string, obj *struct{}, ttl time.Duration)
	GetRedisCache(key string, obj *struct{})
}

type RedisHelper struct {
	Addresses map[string]string
	Ring      *redis.Ring
	MyCache   *cache.Cache
	Ctx       context.Context
}

func NewRedisCache(address map[string]string) RedisCache {
	r := &RedisHelper{Addresses: address}
	r.CreateRedisConnection()
	return r
}

func (r *RedisHelper) CreateRedisConnection() {
	r.Ring = redis.NewRing(&redis.RingOptions{
		Addrs: r.Addresses,
	})
	r.MyCache = cache.New(&cache.Options{
		Redis:      r.Ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	r.Ctx = context.TODO()
}

func (r *RedisHelper) AddRedisCache(key string, obj *struct{}, ttl time.Duration) {
	err := r.MyCache.Set(&cache.Item{
		Ctx:   r.Ctx,
		Key:   key,
		Value: obj,
		TTL:   ttl,
	})
	failOnError(err, "Failed on add Redis Cache")
}

func (r *RedisHelper) GetRedisCache(key string, obj *struct{}) {
	err := r.MyCache.Get(r.Ctx, key, obj)
	failOnError(err, "Failed on get Redis Cache")
}
