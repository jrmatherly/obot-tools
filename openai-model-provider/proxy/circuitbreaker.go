package proxy

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sony/gobreaker"
)

var (
	// CircuitBreakerState tracks the current state of the circuit breaker by provider
	CircuitBreakerState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "obot_model_provider_circuit_breaker_state",
			Help: "Current state of the circuit breaker by provider (0=closed, 1=half-open, 2=open)",
		},
		[]string{"provider"},
	)

	// CircuitBreakerRequests tracks requests by provider and state
	CircuitBreakerRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "obot_model_provider_circuit_breaker_requests_total",
			Help: "Total number of requests by provider and outcome (success/failure/rejected)",
		},
		[]string{"provider", "outcome"},
	)

	// CircuitBreakerLatency tracks request latency distribution
	CircuitBreakerLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "obot_model_provider_circuit_breaker_latency_seconds",
			Help:    "Distribution of request latencies in seconds by provider",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60, 120},
		},
		[]string{"provider"},
	)
)

// CircuitBreakerConfig holds circuit breaker configuration for HTTP transport
type CircuitBreakerConfig struct {
	// Provider name for metrics
	Provider string

	// MaxRequests is the maximum number of requests allowed to pass through
	// when the CircuitBreaker is half-open. Default: 1
	MaxRequests uint32

	// Interval is the cyclic period of the closed state for the CircuitBreaker
	// to clear the internal Counts. Default: 60s
	Interval time.Duration

	// Timeout is the period of the open state, after which the state becomes half-open.
	// Default: 60s
	Timeout time.Duration

	// Threshold is the number of consecutive failures before the circuit opens.
	// Default: 5
	Threshold uint32
}

// DefaultCircuitBreakerConfig returns a CircuitBreakerConfig with sensible defaults
func DefaultCircuitBreakerConfig(provider string) *CircuitBreakerConfig {
	return &CircuitBreakerConfig{
		Provider:    provider,
		MaxRequests: 1,
		Interval:    60 * time.Second,
		Timeout:     60 * time.Second,
		Threshold:   5,
	}
}

// LoadCircuitBreakerConfigFromEnv loads circuit breaker configuration from environment variables
// Environment variables should be prefixed with the provider-specific prefix (e.g., OBOT_OPENAI_MODEL_PROVIDER_)
func LoadCircuitBreakerConfigFromEnv(provider, envPrefix string) *CircuitBreakerConfig {
	config := DefaultCircuitBreakerConfig(provider)

	if val := os.Getenv(envPrefix + "CIRCUIT_BREAKER_THRESHOLD"); val != "" {
		if threshold, err := strconv.ParseUint(val, 10, 32); err == nil {
			config.Threshold = uint32(threshold)
		}
	}

	if val := os.Getenv(envPrefix + "CIRCUIT_BREAKER_TIMEOUT"); val != "" {
		if timeout, err := time.ParseDuration(val); err == nil {
			config.Timeout = timeout
		}
	}

	if val := os.Getenv(envPrefix + "CIRCUIT_BREAKER_INTERVAL"); val != "" {
		if interval, err := time.ParseDuration(val); err == nil {
			config.Interval = interval
		}
	}

	return config
}

// CircuitBreakerTransport wraps an http.RoundTripper with circuit breaker protection
type CircuitBreakerTransport struct {
	transport http.RoundTripper
	cb        *gobreaker.CircuitBreaker
	config    *CircuitBreakerConfig
}

// NewCircuitBreakerTransport creates a new circuit breaker HTTP transport wrapper
func NewCircuitBreakerTransport(transport http.RoundTripper, config *CircuitBreakerConfig) *CircuitBreakerTransport {
	if transport == nil {
		transport = http.DefaultTransport
	}

	if config == nil {
		config = DefaultCircuitBreakerConfig("unknown")
	}

	settings := gobreaker.Settings{
		Name:        fmt.Sprintf("%s-circuit-breaker", config.Provider),
		MaxRequests: config.MaxRequests,
		Interval:    config.Interval,
		Timeout:     config.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= config.Threshold
		},
		OnStateChange: func(_ string, _ gobreaker.State, to gobreaker.State) {
			// Update metrics based on new state
			var stateValue float64
			switch to {
			case gobreaker.StateClosed:
				stateValue = 0
			case gobreaker.StateHalfOpen:
				stateValue = 1
			case gobreaker.StateOpen:
				stateValue = 2
			}
			CircuitBreakerState.WithLabelValues(config.Provider).Set(stateValue)
		},
	}

	cbt := &CircuitBreakerTransport{
		transport: transport,
		cb:        gobreaker.NewCircuitBreaker(settings),
		config:    config,
	}

	// Initialize state metric
	CircuitBreakerState.WithLabelValues(config.Provider).Set(0) // Closed state

	return cbt
}

// RoundTrip implements http.RoundTripper interface with circuit breaker protection
func (cbt *CircuitBreakerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	result, err := cbt.cb.Execute(func() (interface{}, error) {
		return cbt.transport.RoundTrip(req)
	})

	duration := time.Since(start).Seconds()
	CircuitBreakerLatency.WithLabelValues(cbt.config.Provider).Observe(duration)

	if err != nil {
		// Check if circuit breaker rejected the request
		if err == gobreaker.ErrOpenState {
			CircuitBreakerRequests.WithLabelValues(cbt.config.Provider, "rejected").Inc()
			return nil, fmt.Errorf("circuit breaker open for provider %s: too many failures", cbt.config.Provider)
		}

		// Check if it's a retryable HTTP error
		CircuitBreakerRequests.WithLabelValues(cbt.config.Provider, "failure").Inc()
		return nil, err
	}

	resp, ok := result.(*http.Response)
	if !ok {
		CircuitBreakerRequests.WithLabelValues(cbt.config.Provider, "failure").Inc()
		return nil, fmt.Errorf("unexpected response type from circuit breaker")
	}

	// Check HTTP status code to determine if this should count as a failure
	if isFailureStatusCode(resp.StatusCode) {
		CircuitBreakerRequests.WithLabelValues(cbt.config.Provider, "failure").Inc()
		// Return an error so circuit breaker counts this as a failure
		return resp, fmt.Errorf("HTTP %d from provider %s", resp.StatusCode, cbt.config.Provider)
	}

	CircuitBreakerRequests.WithLabelValues(cbt.config.Provider, "success").Inc()
	return resp, nil
}

// isFailureStatusCode determines if an HTTP status code should be counted as a circuit breaker failure
func isFailureStatusCode(statusCode int) bool {
	// 5xx errors and 429 (rate limit) should trigger circuit breaker
	return statusCode >= 500 || statusCode == 429
}

// State returns the current state of the circuit breaker
func (cbt *CircuitBreakerTransport) State() gobreaker.State {
	return cbt.cb.State()
}

// Counts returns the current counts of the circuit breaker
func (cbt *CircuitBreakerTransport) Counts() gobreaker.Counts {
	return cbt.cb.Counts()
}

// isRetryableHTTPError determines if an HTTP error should be retried
func isRetryableHTTPError(err error) bool {
	if err == nil {
		return false
	}

	// Circuit breaker open error is not retryable
	if err == gobreaker.ErrOpenState {
		return false
	}

	// Context errors are not retryable
	if err == context.Canceled || err == context.DeadlineExceeded {
		return false
	}

	// Network errors are generally retryable
	errMsg := err.Error()
	retryablePatterns := []string{
		"connection refused",
		"connection reset",
		"timeout",
		"temporary failure",
		"503",
		"502",
		"504",
		"429", // Rate limit
	}

	for _, pattern := range retryablePatterns {
		if strings.Contains(strings.ToLower(errMsg), strings.ToLower(pattern)) {
			return true
		}
	}

	return false
}
