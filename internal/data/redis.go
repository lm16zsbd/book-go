package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"serica-go/internal/conf"
)

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient(c *conf.Bootstrap) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Data.Redis.Addr,
		Password: c.Data.Redis.Password,
		DB:       c.Data.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &RedisClient{rdb: rdb}, nil
}

func (r *RedisClient) Close() error {
	return r.rdb.Close()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.rdb.Set(ctx, key, value, ttl).Err()
}

func (r *RedisClient) Del(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}

func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.rdb.Incr(ctx, key).Result()
}

func (r *RedisClient) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.rdb.Expire(ctx, key, ttl).Err()
}

func (r *RedisClient) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisClient) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, key, data, ttl).Err()
}

func (r *RedisClient) Pipeline() redis.Pipeliner {
	return r.rdb.Pipeline()
}
