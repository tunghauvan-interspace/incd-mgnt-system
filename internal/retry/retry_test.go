package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryer_Execute(t *testing.T) {
	tests := []struct {
		name        string
		policy      *RetryPolicy
		fn          RetryableFunc
		wantErr     bool
		expectCalls int
	}{
		{
			name: "success on first try",
			policy: &RetryPolicy{
				MaxAttempts: 3,
				BaseDelay:   10 * time.Millisecond,
				MaxDelay:    100 * time.Millisecond,
				Multiplier:  2.0,
			},
			fn: func() error {
				return nil
			},
			wantErr:     false,
			expectCalls: 1,
		},
		{
			name: "success on second try",
			policy: &RetryPolicy{
				MaxAttempts: 3,
				BaseDelay:   10 * time.Millisecond,
				MaxDelay:    100 * time.Millisecond,
				Multiplier:  2.0,
			},
			fn: func() func() error {
				calls := 0
				return func() error {
					calls++
					if calls == 1 {
						return errors.New("temporary error")
					}
					return nil
				}
			}(),
			wantErr:     false,
			expectCalls: 2,
		},
		{
			name: "failure after max attempts",
			policy: &RetryPolicy{
				MaxAttempts: 2,
				BaseDelay:   10 * time.Millisecond,
				MaxDelay:    100 * time.Millisecond,
				Multiplier:  2.0,
			},
			fn: func() error {
				return errors.New("persistent error")
			},
			wantErr:     true,
			expectCalls: 2,
		},
		{
			name: "non-retryable error",
			policy: &RetryPolicy{
				MaxAttempts: 3,
				BaseDelay:   10 * time.Millisecond,
				MaxDelay:    100 * time.Millisecond,
				Multiplier:  2.0,
			},
			fn: func() error {
				return errors.New("validation error")
			},
			wantErr:     true,
			expectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			wrappedFn := func() error {
				calls++
				return tt.fn()
			}

			retryer := NewRetryer(tt.policy, DefaultIsRetryable)
			ctx := context.Background()

			err := retryer.Execute(ctx, wrappedFn)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			if calls != tt.expectCalls {
				t.Errorf("Expected %d calls, got %d", tt.expectCalls, calls)
			}
		})
	}
}

func TestRetryer_ExecuteWithContext(t *testing.T) {
	policy := &RetryPolicy{
		MaxAttempts: 5,
		BaseDelay:   50 * time.Millisecond,
		MaxDelay:    200 * time.Millisecond,
		Multiplier:  2.0,
	}

	retryer := NewRetryer(policy, DefaultIsRetryable)

	// Test context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	calls := 0
	fn := func() error {
		calls++
		return errors.New("temporary error")
	}

	err := retryer.Execute(ctx, fn)
	
	if err == nil {
		t.Error("Expected error due to context cancellation")
	}

	if calls == 0 {
		t.Error("Expected at least one call before cancellation")
	}

	if calls >= 5 {
		t.Error("Too many calls, context cancellation should have prevented this")
	}
}

func TestDefaultIsRetryable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "validation error",
			err:  errors.New("validation failed"),
			want: false,
		},
		{
			name: "invalid error",
			err:  errors.New("invalid input"),
			want: false,
		},
		{
			name: "network error",
			err:  errors.New("network timeout"),
			want: true,
		},
		{
			name: "database error",
			err:  errors.New("database connection failed"),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultIsRetryable(tt.err)
			if got != tt.want {
				t.Errorf("DefaultIsRetryable() = %v, want %v", got, tt.want)
			}
		})
	}
}