package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CounterRepository struct {
	connRedisDb *redis.Client
}

func NewCounterRepository(connRedisDb *redis.Client) *CounterRepository {
	return &CounterRepository{connRedisDb}
}

func (c *CounterRepository) IncrementRedis(ctx context.Context, key string) (int64, error) {
	// Атомарное увеличение в Redis
	return c.connRedisDb.Incr(ctx, key).Result()
}

func (c *CounterRepository) Decrement(ctx context.Context, key string) (int64, error) {
	return c.connRedisDb.Decr(ctx, key).Result()
}

func (c *CounterRepository) Set(ctx context.Context, key string, limit int64, offset time.Duration) {
	c.connRedisDb.Set(ctx, key, limit, offset)
}

func (c *CounterRepository) Get(ctx context.Context, key string) (int64, error) {
	return c.connRedisDb.Get(ctx, key).Int64()
}

func (c *CounterRepository) Del(ctx context.Context, key string) error {
	return c.connRedisDb.Del(ctx, key).Err()
}
