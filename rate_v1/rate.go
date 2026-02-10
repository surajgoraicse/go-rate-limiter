package ratev1

import (
	"errors"
	"time"
)

var RateLimitted = errors.New("please try later")

type RateLimiter struct {
	token      float64
	refillRate float64
	capacity   int64
	lastRefill time.Time
}

func New(cap int64, refillRate int64) *RateLimiter {
	return &RateLimiter{
		token:      float64(cap),
		capacity:   cap,
		refillRate: float64(refillRate),
		lastRefill: time.Now(),
	}
}
func (r *RateLimiter) Allow() bool {

	elapsed := time.Since(r.lastRefill)

	if elapsed > 0 {
		r.token = r.token + float64(elapsed.Seconds())*r.refillRate
		if r.token > float64(r.capacity) {
			r.token = float64(r.capacity)
		}
		r.lastRefill = time.Now()
	}
	if r.token >= 1 {
		r.token--
		return true
	} else {
		return false
	}
}
