package redis

import (
	"coffee-shop/config"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, func()) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Host,
		MinIdleConns: 200,
		PoolSize:     12000,
		PoolTimeout:  time.Duration(240) * time.Second,
		DB:           0, // use default DB
	})

	disconnect := func() {
		client.Close()
	}
	return client, disconnect
}
