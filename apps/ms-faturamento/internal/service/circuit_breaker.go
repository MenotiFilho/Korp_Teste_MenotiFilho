package service

import (
	"errors"
	"sync"
	"time"
)

const (
	circuitClosed   = 0
	circuitOpen     = 1
	circuitHalfOpen = 2
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

type CircuitBreaker struct {
	mu              sync.Mutex
	state           int
	failureCount    int
	successCount    int
	lastFailureTime time.Time

	maxFailures  int
	resetTimeout time.Duration
	halfOpenMax  int
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:        circuitClosed,
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		halfOpenMax:  1,
	}
}

func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case circuitClosed:
		return nil
	case circuitOpen:
		if time.Since(cb.lastFailureTime) >= cb.resetTimeout {
			cb.state = circuitHalfOpen
			cb.successCount = 0
			return nil
		}
		return ErrCircuitOpen
	case circuitHalfOpen:
		if cb.successCount < cb.halfOpenMax {
			return nil
		}
		return ErrCircuitOpen
	}
	return nil
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case circuitHalfOpen:
		cb.successCount++
		if cb.successCount >= cb.halfOpenMax {
			cb.state = circuitClosed
			cb.failureCount = 0
		}
	case circuitClosed:
		cb.failureCount = 0
	}
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount++
	cb.lastFailureTime = time.Now()

	if cb.state == circuitHalfOpen {
		cb.state = circuitOpen
		return
	}

	if cb.failureCount >= cb.maxFailures {
		cb.state = circuitOpen
	}
}

func (cb *CircuitBreaker) State() string {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case circuitClosed:
		return "closed"
	case circuitOpen:
		return "open"
	case circuitHalfOpen:
		return "half-open"
	}
	return "unknown"
}
