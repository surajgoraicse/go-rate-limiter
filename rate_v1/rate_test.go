package ratev1

import (
	"testing"
	"time"
)

func TestAllowUntilExaution(t *testing.T) {
	r := New(3, 1)

	for range 3 {
		if !r.Allow() {
			t.Fatalf("expected first 3 calls to allow")
		}
	}

	if r.Allow() {
		t.Fatalf("expected bucket to exhaust")
	}
}

func TestRefillToken(t *testing.T) {
	r := New(3, 1)
	for range 3 {
		if !r.Allow() {
			t.Fatalf("expected first 3 calls to allow")
		}
	}

	if r.Allow() {
		t.Fatalf("expected bucket to exhaust")
	}
	time.Sleep(3 * time.Second)
	for range 3 {
		if !r.Allow() {
			t.Fatalf("expected to accept 3 calls after refill")
		}
	}
}

func TestNoRefillWithoutTime(t *testing.T) {
	r := New(1, 1000)
	r.Allow()

	if r.Allow() {
		t.Fatalf("unexpected refill without time passed")
	}

}
