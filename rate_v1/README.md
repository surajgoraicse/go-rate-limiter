# Token Bucket Rate Limiter (Iteration 1)

This repository contains **Iteration 1** of a rate limiter implemented using the **Token Bucket algorithm**, written in Go.

The goal of this iteration is **correctness and clarity**, not production readiness.

---


## Token Bucket Model

A token bucket has:

- **capacity**: maximum number of tokens the bucket can hold
- **tokens**: current number of available tokens
- **refill rate**: tokens added per second
- **last refill time**: timestamp of the last refill calculation

Each request consumes **1 token**.

If a token is available, the request is allowed.  
If not, it is rejected.

---

## Core Design Decisions

### Lazy Refill

Tokens are **not refilled continuously**.

Instead, tokens are recalculated **only when `Allow()` is called**, using elapsed time since the last refill.

This avoids:

- background goroutines
- timing drift
- test flakiness
- unnecessary complexity

---

### Floating-Point Tokens

Tokens are stored as `float64`.

Reason:  
Refill is continuous over time.

Example:

- Refill rate = 5 tokens/sec  
- Elapsed time = 100ms  
- Tokens added = 0.5  

Using integers would force rounding and introduce incorrect behavior.

---

### Single-Threaded Assumption

This iteration assumes:

- A single goroutine calling `Allow()`
- No concurrent access

Concurrency and atomicity are intentionally postponed to later iterations.

---

## Goal of Iteration 1

Build the **smallest possible correct token bucket** with the following constraints:

- Single process
- In-memory only
- Single bucket
- No persistence
- No HTTP or networking
- No background goroutines

This iteration focuses entirely on **time-based token refill correctness**.

---

