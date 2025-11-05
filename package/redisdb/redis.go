package redisdb

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func InitRedis(addr, password string, db int) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if pong, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("gagal konek ke Redis: %w", err)
	} else {
		fmt.Println("Redis connected:", pong)
	}

	return &RedisClient{
		Client: rdb,
		Ctx:    ctx,
	}, nil
}
