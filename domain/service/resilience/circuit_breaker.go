package resilience

import (
	loggerI "auth-service/domain/service/logger"
	"context"
	"errors"
	"sync"
	"time"
)

// CircuitBreakerState định nghĩa trạng thái của circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateOpen
	StateHalfOpen
)

func (s CircuitBreakerState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreakerConfig cấu hình cho circuit breaker
type CircuitBreakerConfig struct {
	MaxRequests   uint32                   // Số lượng requests tối đa trong half-open state
	Interval      time.Duration            // Thời gian để reset counter
	Timeout       time.Duration            // Thời gian chờ trước khi chuyển từ open sang half-open
	ReadyToTrip   func(counts Counts) bool // Function để quyết định khi nào trip
	OnStateChange func(name string, from CircuitBreakerState, to CircuitBreakerState)
	IsSuccessful  func(err error) bool // Function để xác định request có thành công không
}

// Counts lưu trữ số liệu thống kê
type Counts struct {
	Requests             uint32
	TotalSuccesses       uint32
	TotalFailures        uint32
	ConsecutiveSuccesses uint32
	ConsecutiveFailures  uint32
}

// CircuitBreaker interface
type CircuitBreaker interface {
	Execute(req func() (interface{}, error)) (interface{}, error)
	ExecuteWithContext(ctx context.Context, req func(ctx context.Context) (interface{}, error)) (interface{}, error)
	Name() string
	State() CircuitBreakerState
	Counts() Counts
}

// circuitBreaker implementation
type circuitBreaker struct {
	name          string
	maxRequests   uint32
	interval      time.Duration
	timeout       time.Duration
	readyToTrip   func(counts Counts) bool
	isSuccessful  func(err error) bool
	onStateChange func(name string, from CircuitBreakerState, to CircuitBreakerState)

	mutex      sync.Mutex
	state      CircuitBreakerState
	generation uint64
	counts     Counts
	expiry     time.Time
	logger     loggerI.Log
}

var (
	ErrTooManyRequests = errors.New("circuit breaker is open")
	ErrOpenState       = errors.New("circuit breaker is in open state")
)

// NewCircuitBreaker tạo circuit breaker mới
func NewCircuitBreaker(name string, config CircuitBreakerConfig, logger loggerI.Log) CircuitBreaker {
	cb := &circuitBreaker{
		name:          name,
		maxRequests:   config.MaxRequests,
		interval:      config.Interval,
		timeout:       config.Timeout,
		readyToTrip:   config.ReadyToTrip,
		isSuccessful:  config.IsSuccessful,
		onStateChange: config.OnStateChange,
		state:         StateClosed,
		expiry:        time.Now().Add(config.Interval),
		logger:        logger,
	}

	if cb.readyToTrip == nil {
		cb.readyToTrip = defaultReadyToTrip
	}

	if cb.isSuccessful == nil {
		cb.isSuccessful = defaultIsSuccessful
	}

	return cb
}

// Execute thực hiện request với circuit breaker protection
func (cb *circuitBreaker) Execute(req func() (interface{}, error)) (interface{}, error) {
	return cb.ExecuteWithContext(context.Background(), func(ctx context.Context) (interface{}, error) {
		return req()
	})
}

// ExecuteWithContext thực hiện request với context
func (cb *circuitBreaker) ExecuteWithContext(ctx context.Context, req func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	generation, err := cb.beforeRequest()
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			cb.afterRequest(generation, false)
			panic(r)
		}
	}()

	result, err := req(ctx)
	cb.afterRequest(generation, cb.isSuccessful(err))
	return result, err
}

func (cb *circuitBreaker) beforeRequest() (uint64, error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)

	if state == StateOpen {
		return generation, ErrOpenState
	} else if state == StateHalfOpen && cb.counts.Requests >= cb.maxRequests {
		return generation, ErrTooManyRequests
	}

	cb.counts.onRequest()
	return generation, nil
}

func (cb *circuitBreaker) afterRequest(before uint64, success bool) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)
	if generation != before {
		return
	}

	if success {
		cb.onSuccess(state, now)
	} else {
		cb.onFailure(state, now)
	}
}

func (cb *circuitBreaker) onSuccess(state CircuitBreakerState, now time.Time) {
	cb.counts.onSuccess()

	if state == StateHalfOpen {
		cb.counts.onSuccess()
		if cb.counts.ConsecutiveSuccesses >= cb.maxRequests {
			cb.setState(StateClosed, now)
		}
	}
}

func (cb *circuitBreaker) onFailure(state CircuitBreakerState, now time.Time) {
	cb.counts.onFailure()

	switch state {
	case StateClosed:
		if cb.readyToTrip(cb.counts) {
			cb.setState(StateOpen, now)
		}
	case StateHalfOpen:
		cb.setState(StateOpen, now)
	}
}

func (cb *circuitBreaker) currentState(now time.Time) (CircuitBreakerState, uint64) {
	switch cb.state {
	case StateClosed:
		if !cb.expiry.IsZero() && cb.expiry.Before(now) {
			cb.toNewGeneration(now)
		}
	case StateOpen:
		if cb.expiry.Before(now) {
			cb.setState(StateHalfOpen, now)
		}
	}
	return cb.state, cb.generation
}

func (cb *circuitBreaker) setState(state CircuitBreakerState, now time.Time) {
	if cb.state == state {
		return
	}

	prev := cb.state
	cb.state = state

	cb.toNewGeneration(now)

	if cb.onStateChange != nil {
		cb.onStateChange(cb.name, prev, state)
	}

	cb.logger.Info("Circuit breaker state changed",
		map[string]interface{}{
			"name": cb.name,
			"from": prev.String(),
			"to":   state.String(),
		})
}

func (cb *circuitBreaker) toNewGeneration(now time.Time) {
	cb.generation++
	cb.counts.clear()

	var zero time.Time
	switch cb.state {
	case StateClosed:
		if cb.interval == 0 {
			cb.expiry = zero
		} else {
			cb.expiry = now.Add(cb.interval)
		}
	case StateOpen:
		cb.expiry = now.Add(cb.timeout)
	default: // StateHalfOpen
		cb.expiry = zero
	}
}

func (cb *circuitBreaker) Name() string {
	return cb.name
}

func (cb *circuitBreaker) State() CircuitBreakerState {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	state, _ := cb.currentState(time.Now())
	return state
}

func (cb *circuitBreaker) Counts() Counts {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	return cb.counts
}

// Helper functions
func (c *Counts) onRequest() {
	c.Requests++
}

func (c *Counts) onSuccess() {
	c.TotalSuccesses++
	c.ConsecutiveSuccesses++
	c.ConsecutiveFailures = 0
}

func (c *Counts) onFailure() {
	c.TotalFailures++
	c.ConsecutiveFailures++
	c.ConsecutiveSuccesses = 0
}

func (c *Counts) clear() {
	c.Requests = 0
	c.TotalSuccesses = 0
	c.TotalFailures = 0
	c.ConsecutiveSuccesses = 0
	c.ConsecutiveFailures = 0
}

// Default functions
func defaultReadyToTrip(counts Counts) bool {
	return counts.ConsecutiveFailures > 5
}

func defaultIsSuccessful(err error) bool {
	return err == nil
}
