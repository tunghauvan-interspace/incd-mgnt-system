package retry

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"
)

// RetryPolicy defines the retry behavior
type RetryPolicy struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Multiplier  float64
}

// DefaultRetryPolicy returns a sensible default retry policy
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxAttempts: 3,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    5 * time.Second,
		Multiplier:  2.0,
	}
}

// RetryableFunc is a function that can be retried
type RetryableFunc func() error

// IsRetryable determines if an error should trigger a retry
type IsRetryable func(error) bool

// DefaultIsRetryable returns true for transient errors
func DefaultIsRetryable(err error) bool {
	// For now, retry most errors except validation errors
	// In a real implementation, you'd check for specific error types
	if err == nil {
		return false
	}
	
	errStr := err.Error()
	// Don't retry validation errors
	if strings.Contains(errStr, "validation") || strings.Contains(errStr, "invalid") || strings.Contains(errStr, "malformed") {
		return false
	}
	
	return true
}

// Retryer handles retry logic with exponential backoff
type Retryer struct {
	policy      *RetryPolicy
	isRetryable IsRetryable
}

// NewRetryer creates a new retryer with the given policy
func NewRetryer(policy *RetryPolicy, isRetryable IsRetryable) *Retryer {
	if policy == nil {
		policy = DefaultRetryPolicy()
	}
	if isRetryable == nil {
		isRetryable = DefaultIsRetryable
	}
	
	return &Retryer{
		policy:      policy,
		isRetryable: isRetryable,
	}
}

// Execute attempts to execute the function with retry logic
func (r *Retryer) Execute(ctx context.Context, fn RetryableFunc) error {
	var lastErr error
	
	for attempt := 0; attempt < r.policy.MaxAttempts; attempt++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation cancelled: %w", ctx.Err())
		default:
		}
		
		// Execute the function
		err := fn()
		if err == nil {
			return nil // Success
		}
		
		lastErr = err
		
		// Check if error is retryable
		if !r.isRetryable(err) {
			return fmt.Errorf("non-retryable error: %w", err)
		}
		
		// Don't sleep after the last attempt
		if attempt == r.policy.MaxAttempts-1 {
			break
		}
		
		// Calculate delay with exponential backoff
		delay := r.calculateDelay(attempt)
		
		// Sleep with context cancellation support
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation cancelled during retry: %w", ctx.Err())
		case <-time.After(delay):
			// Continue to next attempt
		}
	}
	
	return fmt.Errorf("max retry attempts (%d) exceeded, last error: %w", r.policy.MaxAttempts, lastErr)
}

// calculateDelay calculates the delay for the given attempt using exponential backoff
func (r *Retryer) calculateDelay(attempt int) time.Duration {
	// Calculate exponential backoff: base * multiplier^attempt
	delay := float64(r.policy.BaseDelay) * math.Pow(r.policy.Multiplier, float64(attempt))
	
	// Cap at maximum delay
	if delay > float64(r.policy.MaxDelay) {
		delay = float64(r.policy.MaxDelay)
	}
	
	return time.Duration(delay)
}

// ExecuteWithCustomIsRetryable executes with a custom retry check for this specific call
func (r *Retryer) ExecuteWithCustomIsRetryable(ctx context.Context, fn RetryableFunc, isRetryable IsRetryable) error {
	originalIsRetryable := r.isRetryable
	r.isRetryable = isRetryable
	defer func() {
		r.isRetryable = originalIsRetryable
	}()
	
	return r.Execute(ctx, fn)
}