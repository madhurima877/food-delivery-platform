package lock

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLock struct {
	client *redis.Client
}

func NewRedisLock() *RedisLock {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:16379",
	})
	return &RedisLock{client: client}
}
func (r *RedisLock) Acquire(ctx context.Context, orderID string) (bool, error) {
	key := "payment:lock:" + orderID
	acquired, err := r.client.SetNX(ctx, key, "locked", 30*time.Second).Result()
	if err != nil {
		return false, err
	}
	return acquired, nil
}

func (r *RedisLock) Release(ctx context.Context, orderID string) error {
	key := "payment:lock:" + orderID
	return r.client.Del(ctx, key).Err()
}
