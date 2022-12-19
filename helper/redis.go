package helper

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

type RedisCache interface {
	SetRedis(key string, value interface{}, ttl time.Duration) error
	GetRedis(key string, dest interface{}) error
}

type RedisHelper struct {
	Client   *redis.Client
	Address  string
	Password string
	DB       int
	Context  context.Context
	Utils    *UtilHelper
}

func NewRedisCache(address, password string, db int, logger *zap.Logger) RedisCache {
	utils := &UtilHelper{}
	r := &RedisHelper{Address: address, Password: password, DB: db, Utils: utils}
	r.CreateRedisConnection()
	return r
}

func (r *RedisHelper) CreateRedisConnection() {
	r.Context = context.Background()
	r.Client = redis.NewClient(&redis.Options{
		Addr:     r.Address,
		Password: r.Password, // no password set
		DB:       r.DB,       // use default DB
	})
}

func (r *RedisHelper) SetRedis(key string, value interface{}, ttl time.Duration) error {
	valueByte, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err = r.Client.Set(r.Context, key, valueByte, ttl).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisHelper) GetRedis(key string, dest interface{}) error {
	valueByte, err := r.Client.Get(r.Context, key).Bytes()
	if err != nil {
		return err
	}

	if err = json.Unmarshal(valueByte, dest); err != nil {
		return err
	}
	return nil
}
