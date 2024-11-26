package proxy

import (
	"sync"
	"time"
)

type status string

var (
	statusOpen     status = "open"
	statusHalfOpen status = "half_open"
	statusClosed   status = "closed"
)

func (s status) String() string {
	return string(s)
}

func (s status) isOpen() bool {
	return s == statusOpen
}

func (s status) isHalfOpen() bool {
	return s == statusHalfOpen
}

func (s status) isClosed() bool {
	return s == statusClosed
}

type circuitBreaker struct {
	status status // status is the current state of the circuit breaker
	failureCount int // failureCount is the number of consecutive failures
	successCount int // successCount is the number of consecutive successes
	successThresholdRatio float64 // thresholdRatio is the ratio of failures to successes required to trip the circuit breaker
	initializeCountDuration time.Duration // initializeCountDuration is the duration to wait before resetting the failure count
	recoveryTimeout time.Duration // recoveryTimeout is the duration to wait before transitioning from open to half-open
	lastSuccessTime time.Time // lastSuccessTime is the time of the last success
	lastFailureTime time.Time // lastFailureTime is the time of the last failure
	mutex sync.Mutex
}

func newCircuitBreaker(successThresholdRatio float64, initializeCountDuration, recoveryTimeout time.Duration) *circuitBreaker {
	return &circuitBreaker{
		status: statusClosed,
		failureCount: 0,
		successCount: 0,
		successThresholdRatio: successThresholdRatio,
		initializeCountDuration: initializeCountDuration,
		recoveryTimeout: recoveryTimeout,
		mutex: sync.Mutex{},
	}
}

func (cb *circuitBreaker) isAllowed() bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.status {
	case statusOpen:
		if time.Since(cb.lastFailureTime) > cb.recoveryTimeout {
			cb.status = statusHalfOpen
			return true
		}

		return false
	case statusHalfOpen:
		successRatio := float64(cb.successCount) / float64(cb.failureCount)
		if successRatio >= cb.successThresholdRatio {
			cb.status = statusClosed
			return true
		}
		
		return false
	case statusClosed:
		return true
	}

	return false
}

func (cb *circuitBreaker) onFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if time.Since(cb.lastFailureTime) > cb.initializeCountDuration {
		cb.failureCount = 0
		cb.successCount = 0
	}

	cb.failureCount++
	cb.lastFailureTime = time.Now()
}

func (cb *circuitBreaker) onSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if time.Since(cb.lastSuccessTime) > cb.initializeCountDuration {
		cb.failureCount = 0
		cb.successCount = 0
	}
	
	cb.successCount++
	cb.lastSuccessTime = time.Now()
}
