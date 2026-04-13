package service

import (
	"testing"
	"time"
)

func TestCircuitBreaker_WhenClosed_ShouldAllow(t *testing.T) {
	// Arrange
	cb := NewCircuitBreaker(3, 1*time.Second)

	// Act
	err := cb.Allow()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCircuitBreaker_WhenMaxFailuresReached_ShouldOpen(t *testing.T) {
	// Arrange
	cb := NewCircuitBreaker(3, 1*time.Second)

	// Act
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordFailure()

	// Assert
	err := cb.Allow()
	if err == nil {
		t.Fatal("expected circuit open error")
	}
}

func TestCircuitBreaker_WhenResetTimeoutPasses_ShouldTransitionToHalfOpen(t *testing.T) {
	// Arrange
	cb := NewCircuitBreaker(3, 50*time.Millisecond)
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordFailure()

	// Act
	time.Sleep(60 * time.Millisecond)

	// Assert
	err := cb.Allow()
	if err != nil {
		t.Fatalf("expected no error after reset timeout, got %v", err)
	}
}

func TestCircuitBreaker_WhenHalfOpenSuccess_ShouldClose(t *testing.T) {
	// Arrange
	cb := NewCircuitBreaker(3, 50*time.Millisecond)
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordFailure()
	time.Sleep(60 * time.Millisecond)
	_ = cb.Allow()

	// Act
	cb.RecordSuccess()

	// Assert
	err := cb.Allow()
	if err != nil {
		t.Fatalf("expected no error after success in half-open, got %v", err)
	}
	if cb.State() != "closed" {
		t.Fatalf("expected closed state, got %s", cb.State())
	}
}

func TestCircuitBreaker_WhenHalfOpenFailure_ShouldReopen(t *testing.T) {
	// Arrange
	cb := NewCircuitBreaker(3, 50*time.Millisecond)
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordFailure()
	time.Sleep(60 * time.Millisecond)
	_ = cb.Allow()

	// Act
	cb.RecordFailure()

	// Assert
	err := cb.Allow()
	if err == nil {
		t.Fatal("expected circuit open error after half-open failure")
	}
}
