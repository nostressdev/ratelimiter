package ratelimiter

import (
	"context"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type MemoryRateLimiter struct {
	limiter    Limiter
	storage    map[string]uint64
	expiration map[string]time.Time
	lock       sync.Mutex
	gcTicker   *time.Ticker
}

func (mem *MemoryRateLimiter) Take(_ context.Context, key string) (bool, error) {
	mem.lock.Lock()
	defer mem.lock.Unlock()
	limit := mem.limiter(key)
	key = key + "_" + strconv.FormatInt(time.Now().UnixNano()/limit.interval.Nanoseconds(), 10)
	now := mem.storage[key] + 1
	if now > limit.max {
		return false, nil
	}
	mem.storage[key] = now
	mem.expiration[key] = time.Now().Add(limit.interval)
	return true, nil
}

func (mem *MemoryRateLimiter) clean() {
	mem.lock.Lock()
	defer mem.lock.Unlock()
	now := time.Now().UnixNano()
	for key, _ := range mem.expiration {
		if mem.expiration[key].UnixNano() < now {
			delete(mem.expiration, key)
			delete(mem.storage, key)
		}
	}
}

func (mem *MemoryRateLimiter) gc() {
	for range mem.gcTicker.C {
		mem.clean()
	}
}

func NewMemory(limit Limit, gcInterval ...time.Duration) RateLimiter {
	result := &MemoryRateLimiter{
		storage:    make(map[string]uint64),
		expiration: make(map[string]time.Time),
		limiter: func(string) Limit {
			return limit
		},
	}
	if len(gcInterval) > 0 {
		result.gcTicker = time.NewTicker(gcInterval[0])
	} else {
		result.gcTicker = time.NewTicker(time.Second)
	}
	go result.gc()
	runtime.SetFinalizer(result, func(mem *MemoryRateLimiter) {
		mem.gcTicker.Stop()
	})
	return result
}

func NewMemoryWithLimiter(limiter Limiter, gcInterval ...time.Duration) RateLimiter {
	result := &MemoryRateLimiter{
		storage:    make(map[string]uint64),
		expiration: make(map[string]time.Time),
		limiter:    limiter,
	}
	if len(gcInterval) > 0 {
		result.gcTicker = time.NewTicker(gcInterval[0])
	} else {
		result.gcTicker = time.NewTicker(time.Second)
	}
	go result.gc()
	runtime.SetFinalizer(result, func(mem *MemoryRateLimiter) {
		mem.gcTicker.Stop()
	})
	return result
}
