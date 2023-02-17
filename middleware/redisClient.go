package middleware

import (
	"github.com/go-redis/redis/v8"
	"time"
)

func NewRedisClient(addr, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    password,
		DB:          db,
		PoolSize:    10,               // 连接池大小
		PoolTimeout: 30 * time.Second, // 连接池等待时间
	})
}
