package ratelimiter

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestConfig struct {
	t              *testing.T
	newWithLimit   func(limit Limit) RateLimiter
	newWithLimiter func(limiter Limiter) RateLimiter
}

func RateLimiterTest(config TestConfig) {
	t := config.t
	t.Run("Test_CreateWithLimit", func(t *testing.T) {
		config.newWithLimit(NewLimit(time.Second, 10))
	})
	t.Run("Test_CreateWithLimiter", func(t *testing.T) {
		config.newWithLimiter(func(_ string) Limit {
			return NewLimit(time.Second, 10)
		})
	})
	t.Run("Test_Take2onOneKey", func(t *testing.T) {
		limiter := config.newWithLimit(NewLimit(time.Second, 1))
		ok, err := limiter.Take(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, true, ok)
		ok, err = limiter.Take(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, false, ok)
	})
	t.Run("Test_TakeOnDifferentKeys", func(t *testing.T) {
		limiter := config.newWithLimit(NewLimit(time.Second, 1))
		ok, err := limiter.Take(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, true, ok)
		ok, err = limiter.Take(context.Background(), "2")
		assert.NoError(t, err)
		assert.Equal(t, true, ok)
	})
	t.Run("Test_TakeAfterExp", func(t *testing.T) {
		limiter := config.newWithLimit(NewLimit(time.Second, 1))
		ok, err := limiter.Take(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, true, ok)
		time.Sleep(time.Millisecond * 1001)
		ok, err = limiter.Take(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, true, ok)
	})
	t.Run("Test_Take5onOneKey", func(t *testing.T) {
		limiter := config.newWithLimit(NewLimit(time.Second, 5))
		for it := 0; it < 5; it++ {
			ok, err := limiter.Take(context.Background(), "1")
			assert.NoError(t, err)
			assert.Equal(t, true, ok)
		}
		ok, err := limiter.Take(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, false, ok)
	})
}
