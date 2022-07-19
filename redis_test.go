package ratelimiter

import (
	"github.com/go-redis/redis/v8"
	"math/rand"
	"os"
	"testing"
)

const RedisAddr = "REDIS_ADDR"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestRedisRateLimiter(t *testing.T) {
	redisAddr := os.Getenv(RedisAddr)
	if redisAddr == "" {
		t.Skip()
	}
	rand.Seed(179)
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	RateLimiterTest(TestConfig{
		t: t,
		newWithLimit: func(limit Limit) RateLimiter {
			return NewRedis(client, limit, "ratelimiter_test_"+randString(16))
		},
		newWithLimiter: func(limiter Limiter) RateLimiter {
			return NewRedisWithLimiter(client, limiter, "ratelimiter_test_"+randString(16))
		},
	})
}
