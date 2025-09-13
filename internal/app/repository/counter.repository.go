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

func (c *CounterRepository) IncrementCounter(ctx context.Context, key string) (int64, error) {
	return c.connRedisDb.Incr(ctx, key).Result()
}

func (c *CounterRepository) DecrementCounter(ctx context.Context, key string) (int64, error) {
	return c.connRedisDb.Decr(ctx, key).Result()
}

func (c *CounterRepository) SetCounter(ctx context.Context, key string, limit int64, offset time.Duration) {
	c.connRedisDb.Set(ctx, key, limit, offset)
}

func (c *CounterRepository) GetCounter(ctx context.Context, key string) (int64, error) {
	return c.connRedisDb.Get(ctx, key).Int64()
}

func (c *CounterRepository) DelCounter(ctx context.Context, key string) error {
	return c.connRedisDb.Del(ctx, key).Err()
}
