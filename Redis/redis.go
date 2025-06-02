package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Hosts    []string
	Password string
	DB       int
}

type RedisClient struct {
	client *redis.ClusterClient
}

func NewRedisClient(cfg RedisConfig) *RedisClient {
	// return &RedisClient{
	// 	client: redis.NewClusterClient(&redis.ClusterOptions{
	// 		Addrs:    cfg.Hosts, // Use all cluster nodes
	// 		Password: cfg.Password,
	// 	}),
	// }
	return &RedisClient{
		client: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    cfg.Hosts, // Use all cluster nodes
			Password: cfg.Password,
		}),
	}
}

func (r *RedisClient) IncrementCounter(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	// Use pipeline for atomic operations
	pipe := r.client.Pipeline()

	// Increment the counter
	incr := pipe.Incr(ctx, key)

	// Set expiration if this is the first increment (EXPIRE only if key is new)
	pipe.Expire(ctx, key, expiration)

	// Execute both commands
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return incr.Val(), nil
}

func (r *RedisClient) GetCounter(ctx context.Context, key string) (int64, error) {
	return r.client.Get(ctx, key).Int64()
}
