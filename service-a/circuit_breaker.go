package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	closed int32 = iota
	halfOpen
	open
)

type CircuitBreaker struct {
	failures         int64
	successes        int64
	state            int32
	lastFailure      int64
	failureThreshold int64
	successThreshold int64
	keepOpenFor      time.Duration
}

func NewCircuitBreaker(failureThreshold, successThreshold int64, keepOpenFor time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            closed,
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		keepOpenFor:      keepOpenFor,
	}
}

func (cb *CircuitBreaker) Call(fn func() (*http.Response, error)) (*http.Response, error) {
	if cb.keepOpen() {
		return nil, fmt.Errorf("circuit breaker is still open")
	}

	res, err := fn()

	if err != nil || res == nil || res.StatusCode >= 500 {
		cb.onError()
	} else {
		cb.onSuccess()
	}

	return res, err
}

func (cb *CircuitBreaker) keepOpen() bool {
	return atomic.LoadInt32(&cb.state) == open && time.Now().UnixNano()-atomic.LoadInt64(&cb.lastFailure) < int64(cb.keepOpenFor)
}

func (cb *CircuitBreaker) onSuccess() {
	atomic.AddInt64(&cb.successes, 1)
	atomic.StoreInt64(&cb.failures, 0)
	if atomic.LoadInt32(&cb.state) == halfOpen {
		if atomic.LoadInt64(&cb.successes) >= cb.successThreshold {
			log.Println("Circuit breaker is closed")
			atomic.StoreInt32(&cb.state, closed)
		}
	} else if atomic.LoadInt32(&cb.state) == open {
		if time.Now().UnixNano()-atomic.LoadInt64(&cb.lastFailure) >= int64(cb.keepOpenFor) {
			log.Println("Circuit breaker is half-open")
			atomic.StoreInt32(&cb.state, halfOpen)
		}
	}
}

func (cb *CircuitBreaker) onError() {
	atomic.AddInt64(&cb.failures, 1)
	atomic.StoreInt64(&cb.successes, 0)
	atomic.StoreInt64(&cb.lastFailure, time.Now().UnixNano())
	if atomic.LoadInt64(&cb.failures) >= cb.failureThreshold {
		log.Println("Circuit breaker is open")
		atomic.StoreInt32(&cb.state, open)
	}
}
