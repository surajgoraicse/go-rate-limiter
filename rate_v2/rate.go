package ratev2

import (
	"errors"
	"net"
	"sync"
	"time"
)

var ErrInvalidIP = errors.New("invalid IP address")
var ErrInvalidCap = errors.New("invalid capacity : cap should be greater than zero")
var ErrInvalidRate = errors.New("invalid rate : rate should be greater than zero")

type bucket struct {
	mu           sync.Mutex
	token        float64
	cap          int
	lastRefilled time.Time
	rate         float64
}

type RateLimiter struct {
	mu      sync.Mutex
	cap     int
	rate    float64
	buckets map[string]*bucket
}

func New(cap int, rate float64) (*RateLimiter, error) {
	if cap <= 0 {
		return nil, ErrInvalidCap
	}
	if rate <= 0 {
		return nil, ErrInvalidRate
	}
	return &RateLimiter{
		rate:    rate,
		cap:     cap,
		buckets: make(map[string]*bucket),
	}, nil
}

func (r *RateLimiter) Allow(ip string) (bool, error) {
	parse := net.ParseIP(ip)
	if parse == nil {
		return false, ErrInvalidIP
	}
	key := parse.String()

	r.mu.Lock()
	b, exist := r.buckets[key]
	if !exist {
		// will create a bucket and return response
		b = &bucket{
			token:        float64(r.cap),
			cap:          r.cap,
			rate:         r.rate,
			lastRefilled: time.Now(),
		}
		r.buckets[key] = b
	}
	r.mu.Unlock()

	b.mu.Lock()
	defer b.mu.Unlock()

	// refill the bucket
	timeElapsed := time.Since(b.lastRefilled).Seconds()
	if timeElapsed > 0 {
		b.token = min(float64(b.cap), b.token+(timeElapsed)*b.rate)
		b.lastRefilled = time.Now()
	}

	// send res
	if b.token >= 1 {
		b.token--
		return true, nil
	}
	return false, nil
}
