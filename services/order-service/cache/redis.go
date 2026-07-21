package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() *RedisCache {
	client := redis.NewClient(&redis.Options{Addr: "localhost:16379"})
	return &RedisCache{client: client}
}
func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisCache) SetOrder(ctx context.Context, orderID string, data []byte) error {
	key := "order:" + orderID
	value := data

	return r.client.Set(ctx, key, value, 10*time.Minute).Err()

}
func (r *RedisCache) GetOrder(ctx context.Context, orderID string) ([]byte, error) {
	key := "order:" + orderID
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(data), nil
}

func (r *RedisCache) DeleteOrder(ctx context.Context, orderID string) error {
	key := "order:" + orderID
	return r.client.Del(ctx, key).Err()
}
