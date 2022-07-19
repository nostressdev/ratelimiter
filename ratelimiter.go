package ratelimiter

import "context"

type Limiter func(string) Limit

type RateLimiter interface {
	Take(ctx context.Context, key string) (bool, error)
}
