# Rate Limiter (Token Bucket) – Go

This project incrementally builds a **production-grade rate limiter** using the **Token Bucket algorithm** in Go.

The implementation is developed in **fixed iterations**, each iteration strengthening correctness, concurrency, usability, and production readiness.  
Each phase builds on the previous one without changing the algorithmic foundation.

---

## Algorithm Choice

**Token Bucket** was chosen because it:
- Supports bursts naturally
- Is easy to reason about mathematically
- Maps well to distributed systems
- Is widely used in real-world APIs

---

## Iteration Roadmap

### Version 1: Minimal In-Memory Token Bucket

#### Goal
Build the smallest possible **correct** implementation of a token bucket.

This version focuses purely on algorithm correctness and time-based refill logic.

#### Scope
- Single process
- Single bucket
- In-memory only
- No concurrency
- No HTTP integration

#### Requirements
- Fixed bucket capacity
- Configurable refill rate (tokens per second)
- Lazy refill based on elapsed time
- Token consumption on request
- Deterministic behavior

#### Deliverables
- `TokenBucket` struct
- `Allow()` method
- Unit tests validating refill and exhaustion behavior

---

### Version 2: Concurrency-Safe Multi-Key Limiter

#### Goal
Support **multiple independent rate limits** safely under concurrent access.

This iteration introduces synchronization and bucket management.

#### Scope
- Single process
- Multiple buckets keyed by identifier
- Fully concurrency-safe
- In-memory only

#### Requirements
- Thread-safe token consumption
- Lazy bucket creation per key
- Independent rate limits per key
- Correct behavior under high concurrency
- No data races

#### Deliverables
- `Limiter` abstraction
- Safe bucket map
- Concurrent tests using `go test -race`

---

### Version 3: HTTP Middleware Integration

#### Goal
Expose the rate limiter as **HTTP middleware** suitable for APIs.

This iteration turns the limiter into a real application component.

#### Scope
- HTTP server integration
- Per-request rate limiting
- Client-visible behavior

#### Requirements
- Middleware-compatible design
- Proper HTTP status codes (`429 Too Many Requests`)
- Rate limit headers:
  - `X-RateLimit-Limit`
  - `X-RateLimit-Remaining`
  - `Retry-After`
- Configurable key extraction strategy (IP, header, token)

#### Deliverables
- HTTP middleware
- Integration tests
- Example server usage

---

### Version 4: Efficiency, Cleanup, and Time Correctness

#### Goal
Make the limiter **safe for long-running services**.

This iteration focuses on memory control, time handling, and efficiency.

#### Scope
- In-memory only
- High-QPS readiness
- Long uptime support

#### Requirements
- Bucket eviction strategy (TTL or LRU)
- Clock abstraction for testability
- Reduced lock contention
- Stable behavior under clock drift
- Optional cleanup goroutine


#### Deliverables
- Eviction mechanism
- Benchmarks
- Long-running stress tests

---

### Version 5: Production-Grade Architecture

#### Goal
Make the limiter **deployable in production** with confidence.

This iteration introduces extensibility, observability, and operational safety.

#### Scope
- Clean architecture
- Extensible design
- Production observability

#### Requirements
- Interface-based design
- Pluggable storage backends (in-memory first)
- Metrics (allowed, denied, latency)
- Graceful shutdown
- Fail-open / fail-closed configuration
- Clear documentation of guarantees and limits


#### Deliverables
- Clean package structure
- Metrics integration
- Benchmarks
- Production-ready README

---

## Design Principles

- Correctness before performance
- Lazy computation over background work
- Explicit trade-offs over hidden behavior
- Test-driven development
- Production concerns treated as first-class

---

## Final Outcome

By the end of this project, this repository will contain:
- A deeply understood token bucket implementation
- A production-ready Go rate limiter
- Clear documentation of guarantees and limitations
- Benchmarks and tests suitable for real-world usage

---

## License
MIT
