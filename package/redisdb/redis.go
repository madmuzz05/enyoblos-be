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
	opt, _ := redis.ParseURL(fmt.Sprintf("rediss://default_ro:%s@%s", password, addr))
	rdb := redis.NewClient(opt)

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
