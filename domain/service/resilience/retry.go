package resilience

import (
	loggerI "auth-service/domain/service/logger"
	"context"
	"math"
	"math/rand"
	"time"
)

// RetryConfig cấu hình cho retry mechanism
type RetryConfig struct {
	MaxRetries      int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	Multiplier      float64
	Jitter          bool
	RetryableErrors []error
	IsRetryable     func(error) bool
}

// RetryableOperation định nghĩa operation có thể retry
type RetryableOperation func(ctx context.Context, attempt int) (interface{}, error)

// Retryer interface cho retry mechanism
type Retryer interface {
	ExecuteWithRetry(ctx context.Context, operation RetryableOperation) (interface{}, error)
	ExecuteWithRetryAndCircuitBreaker(ctx context.Context, cb CircuitBreaker, operation RetryableOperation) (interface{}, error)
}

type retryer struct {
	config RetryConfig
	logger loggerI.Log
}

// NewRetryer tạo retryer mới
func NewRetryer(config RetryConfig, logger loggerI.Log) Retryer {
	if config.IsRetryable == nil {
		config.IsRetryable = defaultIsRetryable
	}

	if config.Multiplier <= 0 {
		config.Multiplier = 2.0
	}

	if config.InitialDelay <= 0 {
		config.InitialDelay = 100 * time.Millisecond
	}

	if config.MaxDelay <= 0 {
		config.MaxDelay = 30 * time.Second
	}

	return &retryer{
		config: config,
		logger: logger,
	}
}

// ExecuteWithRetry thực hiện operation với retry logic
func (r *retryer) ExecuteWithRetry(ctx context.Context, operation RetryableOperation) (interface{}, error) {
	var lastErr error

	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		result, err := operation(ctx, attempt)
		if err == nil {
			if attempt > 0 {
				r.logger.Info("Operation succeeded after retry",
					map[string]interface{}{
						"attempt": attempt,
					})
			}
			return result, nil
		}

		lastErr = err

		if !r.config.IsRetryable(err) {
			r.logger.Info("Error is not retryable",
				map[string]interface{}{
					"error": err.Error(),
				})
			return nil, err
		}

		if attempt == r.config.MaxRetries {
			r.logger.Error("Max retries exceeded",
				map[string]interface{}{
					"max_retries": r.config.MaxRetries,
					"last_error":  err.Error(),
				})
			break
		}

		delay := r.calculateDelay(attempt)
		r.logger.Warn("Operation failed, retrying",
			map[string]interface{}{
				"attempt":     attempt + 1,
				"max_retries": r.config.MaxRetries,
				"delay":       delay.String(),
				"error":       err.Error(),
			})

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}

	return nil, lastErr
}

// ExecuteWithRetryAndCircuitBreaker kết hợp retry với circuit breaker
func (r *retryer) ExecuteWithRetryAndCircuitBreaker(ctx context.Context, cb CircuitBreaker, operation RetryableOperation) (interface{}, error) {
	return r.ExecuteWithRetry(ctx, func(ctx context.Context, attempt int) (interface{}, error) {
		return cb.ExecuteWithContext(ctx, func(ctx context.Context) (interface{}, error) {
			return operation(ctx, attempt)
		})
	})
}

// calculateDelay tính toán delay cho retry với exponential backoff
func (r *retryer) calculateDelay(attempt int) time.Duration {
	delay := time.Duration(float64(r.config.InitialDelay) * math.Pow(r.config.Multiplier, float64(attempt)))

	if delay > r.config.MaxDelay {
		delay = r.config.MaxDelay
	}

	if r.config.Jitter {
		// Thêm jitter để tránh thundering herd
		jitter := time.Duration(rand.Float64() * float64(delay) * 0.1)
		delay += jitter
	}

	return delay
}

// ResilientClient kết hợp circuit breaker và retry mechanism
type ResilientClient struct {
	circuitBreaker CircuitBreaker
	retryer        Retryer
	logger         loggerI.Log
}

// NewResilientClient tạo resilient client
func NewResilientClient(cb CircuitBreaker, retryer Retryer, logger loggerI.Log) *ResilientClient {
	return &ResilientClient{
		circuitBreaker: cb,
		retryer:        retryer,
		logger:         logger,
	}
}

// Execute thực hiện operation với đầy đủ resilience patterns
func (rc *ResilientClient) Execute(ctx context.Context, operation RetryableOperation) (interface{}, error) {
	return rc.retryer.ExecuteWithRetryAndCircuitBreaker(ctx, rc.circuitBreaker, operation)
}

// GetCircuitBreakerState trả về trạng thái circuit breaker
func (rc *ResilientClient) GetCircuitBreakerState() CircuitBreakerState {
	return rc.circuitBreaker.State()
}

// GetCircuitBreakerCounts trả về số liệu thống kê
func (rc *ResilientClient) GetCircuitBreakerCounts() Counts {
	return rc.circuitBreaker.Counts()
}

// Default functions
func defaultIsRetryable(err error) bool {
	// Mặc định retry tất cả errors trừ context errors
	if err == context.Canceled || err == context.DeadlineExceeded {
		return false
	}
	return true
}

// Helper function để tạo RetryConfig with defaults
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:   3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
		IsRetryable:  defaultIsRetryable,
	}
}

// Helper function để tạo CircuitBreakerConfig with defaults
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		MaxRequests: 5,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		IsSuccessful: func(err error) bool {
			return err == nil
		},
	}
}
