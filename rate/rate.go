package rate

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
)

var (
	ErrInvalidIP   = errors.New("invalid IP address")
	ErrInvalidCap  = errors.New("invalid capacity : cap should be greater than zero")
	ErrInvalidRate = errors.New("invalid rate : rate should be greater than zero")
)

type bucket struct {
	mu           sync.Mutex
	token        float64
	cap          int
	lastRefilled time.Time
	rate         float64
}

type RateLimiter struct {
	mu      sync.Mutex         // global mutex
	cap     int                // capacity of the bucket
	rate    float64            // rate of refill
	buckets map[string]*bucket // key : ip , value : &bucket
}

type RateConfig struct {
	Cap  int
	Rate float64
}

// returns a new rate limiter instance
func new(cap int, rate float64) (*RateLimiter, error) {
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

// returns true if the request is allowed
func (r *RateLimiter) allow(ip string, c *echo.Context) (bool, error) {
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

	c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", r.cap))
	c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%.0f", b.token))
	// send res
	if b.token >= 1 {
		b.token--
		return true, nil
	}
	return false, nil
}


// echo rate limiter middleware
func NewRateLimiter(config RateConfig) echo.MiddlewareFunc {
	r, err := new(config.Cap, config.Rate)
	if err != nil {
		panic("rate limiter initialization failed : " + err.Error())
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			allow, err := r.allow(c.RealIP(), c)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": err.Error(),
				})
			}
			if !allow {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error":   "rate_limit_exceeded",
					"message": "Too many requests. Please try again later.",
				})
			}
			return next(c)
		}
	}
}
