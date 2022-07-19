package ratelimiter

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type RedisRateLimiter struct {
	client    *redis.Client
	keyPrefix string
	limiter   Limiter
}

func (rateLimit *RedisRateLimiter) Take(ctx context.Context, key string) (bool, error) {
	limit := rateLimit.limiter(key)
	key = rateLimit.keyPrefix + key + strconv.FormatInt(time.Now().UnixNano()/limit.interval.Nanoseconds(), 10)
	pipe := rateLimit.client.TxPipeline()
	totalCountCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, limit.interval)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	totalCount, err := totalCountCmd.Uint64()
	return err != nil || totalCount <= limit.max, err
}

func NewRedis(client *redis.Client, limit Limit, keyPrefix ...string) RateLimiter {
	prefix := ""
	if len(keyPrefix) != 0 {
		prefix = keyPrefix[0]
	}
	return &RedisRateLimiter{
		client: client,
		limiter: func(string) Limit {
			return limit
		},
		keyPrefix: prefix,
	}
}

func NewRedisWithLimiter(client *redis.Client, limiter Limiter, storeKeyPrefix ...string) RateLimiter {
	prefix := ""
	if len(storeKeyPrefix) != 0 {
		prefix = storeKeyPrefix[0]
	}
	return &RedisRateLimiter{
		client:    client,
		limiter:   limiter,
		keyPrefix: prefix,
	}
}
