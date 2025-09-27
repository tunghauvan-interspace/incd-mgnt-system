package circuitbreaker

import (
	"fmt"
	"sync"
	"time"
)

// State represents the circuit breaker state
type State int

const (
	// StateClosed means the circuit breaker is closed (normal operation)
	StateClosed State = iota
	// StateOpen means the circuit breaker is open (calls fail fast)
	StateOpen
	// StateHalfOpen means the circuit breaker is allowing limited requests
	StateHalfOpen
)

// String returns string representation of the state
func (s State) String() string {
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

// Config holds circuit breaker configuration
type Config struct {
	MaxRequests         uint32        // Maximum requests allowed in half-open state
	Interval           time.Duration // Window interval for failure counting
	Timeout            time.Duration // Time to wait before transitioning from open to half-open
	ReadyToTrip        func(counts Counts) bool // Function to determine if circuit should trip
	OnStateChange      func(name string, from State, to State) // Callback for state changes
	IsSuccessful       func(err error) bool // Function to determine if request is successful
}

// Counts holds the statistics for circuit breaker
type Counts struct {
	Requests             uint32
	TotalSuccesses       uint32
	TotalFailures        uint32
	ConsecutiveSuccesses uint32
	ConsecutiveFailures  uint32
}

// DefaultConfig returns a default circuit breaker configuration
func DefaultConfig() *Config {
	return &Config{
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     60 * time.Second,
		ReadyToTrip: func(counts Counts) bool {
			// Trip if failure rate is >= 50% and we have at least 3 requests
			return counts.Requests >= 3 && float64(counts.TotalFailures)/float64(counts.Requests) >= 0.5
		},
		OnStateChange: nil,
		IsSuccessful: func(err error) bool {
			return err == nil
		},
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	name        string
	config      *Config
	mutex       sync.RWMutex
	state       State
	counts      Counts
	expiry      time.Time
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(name string, config *Config) *CircuitBreaker {
	if config == nil {
		config = DefaultConfig()
	}
	
	cb := &CircuitBreaker{
		name:   name,
		config: config,
		state:  StateClosed,
		expiry: time.Now().Add(config.Interval),
	}
	
	return cb
}

// Call executes the given function if the circuit breaker allows it
func (cb *CircuitBreaker) Call(fn func() error) error {
	generation, err := cb.beforeCall()
	if err != nil {
		return err
	}
	
	defer func() {
		cb.afterCall(generation, err)
	}()
	
	err = fn()
	return err
}

// beforeCall handles pre-execution logic
func (cb *CircuitBreaker) beforeCall() (uint64, error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	now := time.Now()
	state, generation := cb.currentState(now)
	
	if state == StateOpen {
		return generation, fmt.Errorf("circuit breaker '%s' is OPEN", cb.name)
	} else if state == StateHalfOpen && cb.counts.Requests >= cb.config.MaxRequests {
		return generation, fmt.Errorf("circuit breaker '%s' is HALF_OPEN and max requests exceeded", cb.name)
	}
	
	cb.counts.Requests++
	return generation, nil
}

// afterCall handles post-execution logic
func (cb *CircuitBreaker) afterCall(generation uint64, err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	now := time.Now()
	state, currentGeneration := cb.currentState(now)
	
	if generation != currentGeneration {
		return // Ignore stale calls
	}
	
	if cb.config.IsSuccessful(err) {
		cb.onSuccess(state)
	} else {
		cb.onFailure(state)
	}
}

// onSuccess handles successful request
func (cb *CircuitBreaker) onSuccess(state State) {
	cb.counts.TotalSuccesses++
	cb.counts.ConsecutiveSuccesses++
	cb.counts.ConsecutiveFailures = 0
	
	if state == StateHalfOpen {
		cb.setState(StateClosed)
	}
}

// onFailure handles failed request
func (cb *CircuitBreaker) onFailure(state State) {
	cb.counts.TotalFailures++
	cb.counts.ConsecutiveFailures++
	cb.counts.ConsecutiveSuccesses = 0
	
	if cb.config.ReadyToTrip(cb.counts) {
		cb.setState(StateOpen)
	}
}

// currentState returns the current state and generation
func (cb *CircuitBreaker) currentState(now time.Time) (State, uint64) {
	switch cb.state {
	case StateClosed:
		if cb.expiry.Before(now) {
			cb.toNewGeneration(now)
		}
	case StateOpen:
		if cb.expiry.Before(now) {
			cb.setState(StateHalfOpen)
		}
	}
	
	return cb.state, cb.generation()
}

// setState changes the circuit breaker state
func (cb *CircuitBreaker) setState(state State) {
	if cb.state == state {
		return
	}
	
	prev := cb.state
	cb.state = state
	
	now := time.Now()
	switch state {
	case StateClosed:
		cb.toNewGeneration(now)
	case StateOpen:
		cb.expiry = now.Add(cb.config.Timeout)
	case StateHalfOpen:
		cb.expiry = time.Time{} // No expiry for half-open
	}
	
	if cb.config.OnStateChange != nil {
		cb.config.OnStateChange(cb.name, prev, state)
	}
}

// toNewGeneration resets counts for a new generation
func (cb *CircuitBreaker) toNewGeneration(expiry time.Time) {
	cb.counts = Counts{}
	cb.expiry = expiry.Add(cb.config.Interval)
}

// generation returns a unique identifier for the current generation
func (cb *CircuitBreaker) generation() uint64 {
	return uint64(cb.expiry.UnixNano())
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() State {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	
	now := time.Now()
	state, _ := cb.currentState(now)
	return state
}

// Counts returns the current counts
func (cb *CircuitBreaker) Counts() Counts {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	
	return cb.counts
}

// Name returns the circuit breaker name
func (cb *CircuitBreaker) Name() string {
	return cb.name
}