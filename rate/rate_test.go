package ratev2

import (
	"errors"
	"sync"
	"testing"
)

func TestInvalidParameter(t *testing.T) {
	t.Run("check invalid capacity", func(t *testing.T) {
		_, err := New(0, 1)
		if !errors.Is(err, ErrInvalidCap) {
			t.Errorf("invalid error returned")
		}
	})
	t.Run("check invalid rate", func(t *testing.T) {
		_, err := New(1, 0)
		if !errors.Is(err, ErrInvalidRate) {
			t.Errorf("invalid error returned")
		}
	})

}

func TestPerKeyIsolation(t *testing.T) {
	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"

	n := 10
	limiter, _ := New(n, 1)
	for range n {
		v1, _ := limiter.Allow(ip1)
		v2, _ := limiter.Allow(ip2)
		if !v1 || !v2 {
			t.Errorf("key isolation failed")
		}
	}
}

func TestRejectInvalidIP(t *testing.T) {
	ip := "invalid ip"
	r, _ := New(10, 1)
	_, err := r.Allow(ip)
	if !errors.Is(err, ErrInvalidIP) {
		t.Errorf("failed to reject invalid IP")
	}
}

func TestConcurrentCalls(t *testing.T) {
	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"

	n := 100
	limiter, _ := New(n, 1)
	var wg sync.WaitGroup
	for range 500 {
		wg.Go(func() {
			limiter.Allow(ip1)
			limiter.Allow(ip2)
		})
	}
	wg.Wait()
}

func TestIPNormalization(t *testing.T) {
	ip1 := "::ffff:127.0.0.1"
	ip2 := "127.0.0.1"
	n := 2
	r, _ := New(n, 1)
	r.Allow(ip1)
	r.Allow(ip2)
	if v, _ := r.Allow(ip1); v {
		t.Error("IP normalization failed")
	}
}
