package ratelimit

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimit struct {
	client *redis.Client
}

func NewRateLimit() *RateLimit {
	client := redis.NewClient(&redis.Options{Addr: "localhost:16379"})
	return &RateLimit{client: client}
}
func (r *RateLimit) Allow(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		err = r.client.Expire(ctx, key, time.Minute).Err()
		if err != nil {
			return false, err
		}
	}
	if count > 10 {
		return false, nil
	}
	return true, nil
}
func (r *RateLimit) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		host, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			host = req.RemoteAddr
		}

		key := "rate_limit:" + host + ":" + req.URL.Path
		allowed, err := r.Allow(req.Context(), key)
		if err != nil {
			http.Error(w, "Rate limit error", http.StatusInternalServerError)
		}
		if !allowed {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next(w, req)
	}

}
