# RateLimiter

### Usage example

```golang
// Create insatnce
limiter := ratelimiter.NewMemory(ratelimiter.NewLimit(time.Second, 1))
// Or with limiter
limiter := ratelimiter.NewRedisWithLimiter(client, func (key string) *Limit {
	if key == "1" {
        return NewLimit(time.Second, 10)
    }
    return NewLimit(time.Second, 1000000000)
}, "ratelimiter_test_"+randString(16))
// Try to take
ok, err := limiter.Take(context.Background(), "1")
assert.NoError(t, err)
assert.Equal(t, true, ok)
ok, err = limiter.Take(context.Background(), "1")
assert.NoError(t, err)
assert.Equal(t, false, ok)
```

### Internal

You can use one of two implementations:

```golang
type MemoryRateLimiter struct
type RedisRateLimiter struct
```

First is just in memory implementation, second uses redis for ratelimiting. It depends on:

```
github.com/go-redis/redis/v8
```

You can specify limit for all keys or create function which gets a key and returns limit for it. It will be executed
while trying to take a value

Optionally for memory implementation you can pass GC time interval. For redis you may specify prefix, which will be used
to store keys

### Testing

Memory tests can be executed without eny enviroment

To run redis tests you need to specify enviroment variable REDIS_ADDR. It should be an addr of redis host without auth