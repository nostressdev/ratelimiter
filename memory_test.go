package ratelimiter

import (
	"testing"
	"time"
)

func TestMemoryRateLimiter(t *testing.T) {
	RateLimiterTest(TestConfig{
		t: t,
		newWithLimit: func(limit Limit) RateLimiter {
			return NewMemory(limit, time.Millisecond*200)
		},
		newWithLimiter: func(limiter Limiter) RateLimiter {
			return NewMemoryWithLimiter(limiter, time.Millisecond*200)
		},
	})
}
