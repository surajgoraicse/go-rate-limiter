# Rate Limiter – Iteration 2 (Concurrency-safe, per-IP)

This directory contains **Iteration 2** of a rate limiter built using the **Token Bucket algorithm** in Go.

The focus of this iteration is to move from a single-bucket, single-threaded design to something that is **safe under concurrency** and **usable in a real Go service**.

---

## Goal of Iteration 2

Build a rate limiter that:

- Is safe under concurrent access
- Supports multiple independent clients
- Uses **IP-based buckets**
- Preserves the exact refill correctness from Iteration 1

At this point, the limiter should be realistic enough to embed inside an API server.

---

## What This Iteration Adds

### New capabilities

- Thread-safe implementation
- One token bucket per IP address
- Lazy bucket creation
- Proper lock granularity
- Continuous (time-based) token refill

### Still intentionally excluded

- No bucket eviction or TTL
- No persistence
- No HTTP middleware
- No distributed coordination

These are deferred to later iterations.

---

## Design Overview

This iteration introduces two layers:

### RateLimiter (manager)

Responsible for:

- Mapping IP → bucket
- Synchronizing access to the buckets map
- Creating buckets lazily

```go
type RateLimiter struct {
	mu      sync.Mutex
	cap     int
	rate    float64
	buckets map[string]*bucket
}
```

## 2. bucket (per-IP token bucket)

Responsible for:

Token accounting

Time-based refill

Per-IP synchronization

Each bucket has its own mutex, allowing different IPs to proceed concurrently.

Refill Strategy

No background goroutines

No timers or tickers

Tokens are refilled only when Allow is called

Refill amount is calculated using elapsed time

Tokens are capped at capacity

This preserves correctness and avoids timing drift.

## Testing Expectations

Iteration 2 should pass tests that verify:

Correct throttling per IP

Independence across different IPs

Correct refill behavior over time

Safe concurrent access

No data races (go test -race)

## Example usage :

```go
rl := New(10, 5) // 10 tokens max, 5 tokens/sec

allowed, err := rl.Allow("192.168.1.10")
if err != nil {
	// invalid IP
}

if allowed {
	// request allowed
} else {
	// rate limited
}
```
