package ratelimiter

import "time"

type Limit struct {
	interval time.Duration
	max      uint64
}

func NewLimit(interval time.Duration, max uint64) Limit {
	return Limit{
		interval: interval,
		max:      max,
	}
}
