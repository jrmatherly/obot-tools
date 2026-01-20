package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/sony/gobreaker"
)

// mockRoundTripper is a mock implementation of http.RoundTripper for testing
type mockRoundTripper struct {
	responses      []*http.Response
	errors         []error
	callCount      int
	statusCode     int
	returnError    error
	returnResponse *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.callCount++

	// Return pre-configured response/error if set
	if m.returnError != nil {
		return nil, m.returnError
	}
	if m.returnResponse != nil {
		return m.returnResponse, nil
	}

	// Return from response/error arrays
	if len(m.responses) > 0 && m.callCount <= len(m.responses) {
		return m.responses[m.callCount-1], nil
	}
	if len(m.errors) > 0 && m.callCount <= len(m.errors) {
		return nil, m.errors[m.callCount-1]
	}

	// Default: return OK response
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       http.NoBody,
		Request:    req,
	}, nil
}

func TestNewCircuitBreakerTransport(t *testing.T) {
	tests := []struct {
		name      string
		transport http.RoundTripper
		config    *CircuitBreakerConfig
		wantNil   bool
	}{
		{
			name:      "with valid transport and config",
			transport: http.DefaultTransport,
			config:    DefaultCircuitBreakerConfig("test-provider"),
			wantNil:   false,
		},
		{
			name:      "with nil transport uses default",
			transport: nil,
			config:    DefaultCircuitBreakerConfig("test-provider"),
			wantNil:   false,
		},
		{
			name:      "with nil config uses default",
			transport: http.DefaultTransport,
			config:    nil,
			wantNil:   false,
		},
		{
			name:      "with both nil",
			transport: nil,
			config:    nil,
			wantNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cbt := NewCircuitBreakerTransport(tt.transport, tt.config)

			if tt.wantNil && cbt != nil {
				t.Errorf("NewCircuitBreakerTransport() = %v, want nil", cbt)
			}
			if !tt.wantNil && cbt == nil {
				t.Error("NewCircuitBreakerTransport() = nil, want non-nil")
			}

			if cbt != nil {
				// Verify circuit breaker is initialized
				if cbt.cb == nil {
					t.Error("circuit breaker is nil")
				}
				if cbt.config == nil {
					t.Error("config is nil")
				}
				// Verify initial state is closed
				if cbt.State() != gobreaker.StateClosed {
					t.Errorf("initial state = %v, want %v", cbt.State(), gobreaker.StateClosed)
				}
			}
		})
	}
}

func TestDefaultCircuitBreakerConfig(t *testing.T) {
	provider := "test-provider"
	config := DefaultCircuitBreakerConfig(provider)

	if config.Provider != provider {
		t.Errorf("Provider = %v, want %v", config.Provider, provider)
	}
	if config.MaxRequests != 1 {
		t.Errorf("MaxRequests = %v, want 1", config.MaxRequests)
	}
	if config.Interval != 60*time.Second {
		t.Errorf("Interval = %v, want %v", config.Interval, 60*time.Second)
	}
	if config.Timeout != 60*time.Second {
		t.Errorf("Timeout = %v, want %v", config.Timeout, 60*time.Second)
	}
	if config.Threshold != 5 {
		t.Errorf("Threshold = %v, want 5", config.Threshold)
	}
}

func TestLoadCircuitBreakerConfigFromEnv(t *testing.T) {
	tests := []struct {
		name          string
		envVars       map[string]string
		provider      string
		envPrefix     string
		wantThreshold uint32
		wantTimeout   time.Duration
		wantInterval  time.Duration
	}{
		{
			name:          "default values when no env vars",
			envVars:       map[string]string{},
			provider:      "openai",
			envPrefix:     "OBOT_OPENAI_MODEL_PROVIDER_",
			wantThreshold: 5,
			wantTimeout:   60 * time.Second,
			wantInterval:  60 * time.Second,
		},
		{
			name: "custom threshold from env",
			envVars: map[string]string{
				"OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_THRESHOLD": "10",
			},
			provider:      "openai",
			envPrefix:     "OBOT_OPENAI_MODEL_PROVIDER_",
			wantThreshold: 10,
			wantTimeout:   60 * time.Second,
			wantInterval:  60 * time.Second,
		},
		{
			name: "custom timeout from env",
			envVars: map[string]string{
				"OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_TIMEOUT": "30s",
			},
			provider:      "openai",
			envPrefix:     "OBOT_OPENAI_MODEL_PROVIDER_",
			wantThreshold: 5,
			wantTimeout:   30 * time.Second,
			wantInterval:  60 * time.Second,
		},
		{
			name: "custom interval from env",
			envVars: map[string]string{
				"OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_INTERVAL": "2m",
			},
			provider:      "openai",
			envPrefix:     "OBOT_OPENAI_MODEL_PROVIDER_",
			wantThreshold: 5,
			wantTimeout:   60 * time.Second,
			wantInterval:  2 * time.Minute,
		},
		{
			name: "all custom values",
			envVars: map[string]string{
				"OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_THRESHOLD": "3",
				"OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_TIMEOUT":   "45s",
				"OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_INTERVAL":  "90s",
			},
			provider:      "openai",
			envPrefix:     "OBOT_OPENAI_MODEL_PROVIDER_",
			wantThreshold: 3,
			wantTimeout:   45 * time.Second,
			wantInterval:  90 * time.Second,
		},
		{
			name: "invalid values ignored",
			envVars: map[string]string{
				"OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_THRESHOLD": "invalid",
				"OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_TIMEOUT":   "invalid",
			},
			provider:      "openai",
			envPrefix:     "OBOT_OPENAI_MODEL_PROVIDER_",
			wantThreshold: 5,
			wantTimeout:   60 * time.Second,
			wantInterval:  60 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Clearenv()

			// Set test environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer os.Clearenv()

			config := LoadCircuitBreakerConfigFromEnv(tt.provider, tt.envPrefix)

			if config.Threshold != tt.wantThreshold {
				t.Errorf("Threshold = %v, want %v", config.Threshold, tt.wantThreshold)
			}
			if config.Timeout != tt.wantTimeout {
				t.Errorf("Timeout = %v, want %v", config.Timeout, tt.wantTimeout)
			}
			if config.Interval != tt.wantInterval {
				t.Errorf("Interval = %v, want %v", config.Interval, tt.wantInterval)
			}
			if config.Provider != tt.provider {
				t.Errorf("Provider = %v, want %v", config.Provider, tt.provider)
			}
		})
	}
}

func TestCircuitBreakerTransport_SuccessfulRequests(t *testing.T) {
	mock := &mockRoundTripper{
		statusCode: http.StatusOK,
	}

	config := &CircuitBreakerConfig{
		Provider:    "test",
		MaxRequests: 1,
		Interval:    60 * time.Second,
		Timeout:     60 * time.Second,
		Threshold:   5,
	}

	cbt := NewCircuitBreakerTransport(mock, config)

	// Make multiple successful requests
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "http://example.com", nil)
		resp, err := cbt.RoundTrip(req)

		if err != nil {
			t.Errorf("request %d: unexpected error: %v", i+1, err)
		}
		if resp == nil {
			t.Errorf("request %d: response is nil", i+1)
		}
		if resp != nil && resp.StatusCode != http.StatusOK {
			t.Errorf("request %d: StatusCode = %v, want %v", i+1, resp.StatusCode, http.StatusOK)
		}
	}

	// Verify circuit breaker remains closed
	if cbt.State() != gobreaker.StateClosed {
		t.Errorf("State() = %v, want %v", cbt.State(), gobreaker.StateClosed)
	}
}

func TestCircuitBreakerTransport_CircuitOpensAfterThreshold(t *testing.T) {
	mock := &mockRoundTripper{
		returnResponse: &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Status:     "503 Service Unavailable",
			Body:       http.NoBody,
		},
	}

	config := &CircuitBreakerConfig{
		Provider:    "test",
		MaxRequests: 1,
		Interval:    60 * time.Second,
		Timeout:     1 * time.Second,
		Threshold:   3,
	}

	cbt := NewCircuitBreakerTransport(mock, config)

	// Verify initial state
	if cbt.State() != gobreaker.StateClosed {
		t.Errorf("initial State() = %v, want %v", cbt.State(), gobreaker.StateClosed)
	}

	// Make requests that trigger failures (threshold = 3)
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "http://example.com", nil)
		resp, err := cbt.RoundTrip(req)

		if err == nil {
			t.Errorf("request %d: expected error for 503 status, got nil", i+1)
		}
		if resp != nil {
			t.Errorf("request %d: expected nil response for circuit breaker failure, got %v", i+1, resp)
		}
	}

	// Verify circuit is now open
	if cbt.State() != gobreaker.StateOpen {
		t.Errorf("after threshold failures, State() = %v, want %v", cbt.State(), gobreaker.StateOpen)
	}

	// Next request should be rejected immediately
	req := httptest.NewRequest("GET", "http://example.com", nil)
	resp, err := cbt.RoundTrip(req)

	if err == nil {
		t.Error("expected error when circuit is open, got nil")
	}
	if resp != nil {
		t.Errorf("expected nil response when circuit is open, got %v", resp)
	}
	if mock.callCount != 3 {
		t.Errorf("mock was called %d times, want 3 (rejected request should not call transport)", mock.callCount)
	}
}

func TestCircuitBreakerTransport_CircuitRecovery(t *testing.T) {
	mock := &mockRoundTripper{}

	config := &CircuitBreakerConfig{
		Provider:    "test",
		MaxRequests: 1,
		Interval:    60 * time.Second,
		Timeout:     500 * time.Millisecond, // Short timeout for testing
		Threshold:   3,
	}

	cbt := NewCircuitBreakerTransport(mock, config)

	// Trigger circuit to open
	mock.returnResponse = &http.Response{
		StatusCode: http.StatusInternalServerError,
		Status:     "500 Internal Server Error",
		Body:       http.NoBody,
	}

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "http://example.com", nil)
		cbt.RoundTrip(req)
	}

	if cbt.State() != gobreaker.StateOpen {
		t.Fatalf("circuit should be open, got state: %v", cbt.State())
	}

	// Wait for timeout to transition to half-open
	time.Sleep(600 * time.Millisecond)

	// Now return successful responses
	mock.returnResponse = &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       http.NoBody,
	}

	// Make a successful request in half-open state
	req := httptest.NewRequest("GET", "http://example.com", nil)
	resp, err := cbt.RoundTrip(req)

	if err != nil {
		t.Errorf("expected success in half-open state, got error: %v", err)
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK response, got: %v", resp)
	}

	// Circuit should now be closed
	if cbt.State() != gobreaker.StateClosed {
		t.Errorf("after successful request in half-open, State() = %v, want %v", cbt.State(), gobreaker.StateClosed)
	}
}

func TestCircuitBreakerTransport_FailureStatusCodes(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		shouldTrigger  bool
		description    string
	}{
		{
			name:          "200 OK does not trigger",
			statusCode:    http.StatusOK,
			shouldTrigger: false,
			description:   "2xx success codes should not trigger circuit breaker",
		},
		{
			name:          "201 Created does not trigger",
			statusCode:    http.StatusCreated,
			shouldTrigger: false,
			description:   "2xx success codes should not trigger circuit breaker",
		},
		{
			name:          "400 Bad Request does not trigger",
			statusCode:    http.StatusBadRequest,
			shouldTrigger: false,
			description:   "4xx client errors (except 429) should not trigger circuit breaker",
		},
		{
			name:          "404 Not Found does not trigger",
			statusCode:    http.StatusNotFound,
			shouldTrigger: false,
			description:   "4xx client errors (except 429) should not trigger circuit breaker",
		},
		{
			name:          "429 Rate Limited triggers",
			statusCode:    http.StatusTooManyRequests,
			shouldTrigger: true,
			description:   "429 rate limit should trigger circuit breaker",
		},
		{
			name:          "500 Internal Server Error triggers",
			statusCode:    http.StatusInternalServerError,
			shouldTrigger: true,
			description:   "5xx server errors should trigger circuit breaker",
		},
		{
			name:          "502 Bad Gateway triggers",
			statusCode:    http.StatusBadGateway,
			shouldTrigger: true,
			description:   "5xx server errors should trigger circuit breaker",
		},
		{
			name:          "503 Service Unavailable triggers",
			statusCode:    http.StatusServiceUnavailable,
			shouldTrigger: true,
			description:   "5xx server errors should trigger circuit breaker",
		},
		{
			name:          "504 Gateway Timeout triggers",
			statusCode:    http.StatusGatewayTimeout,
			shouldTrigger: true,
			description:   "5xx server errors should trigger circuit breaker",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRoundTripper{
				returnResponse: &http.Response{
					StatusCode: tt.statusCode,
					Status:     fmt.Sprintf("%d", tt.statusCode),
					Body:       http.NoBody,
				},
			}

			config := &CircuitBreakerConfig{
				Provider:    "test",
				MaxRequests: 1,
				Interval:    60 * time.Second,
				Timeout:     60 * time.Second,
				Threshold:   5,
			}

			cbt := NewCircuitBreakerTransport(mock, config)

			req := httptest.NewRequest("GET", "http://example.com", nil)
			resp, err := cbt.RoundTrip(req)

			if tt.shouldTrigger {
				if err == nil {
					t.Errorf("%s: expected error for status %d, got nil", tt.description, tt.statusCode)
				}
				if resp != nil {
					t.Errorf("%s: expected nil response for circuit breaker failure, got %v", tt.description, resp)
				}
			} else {
				if err != nil {
					t.Errorf("%s: unexpected error for status %d: %v", tt.description, tt.statusCode, err)
				}
				if resp == nil {
					t.Errorf("%s: expected response for status %d", tt.description, tt.statusCode)
				}
			}
		})
	}
}

func TestCircuitBreakerTransport_NetworkErrors(t *testing.T) {
	tests := []struct {
		name  string
		err   error
	}{
		{
			name: "connection refused",
			err:  errors.New("connection refused"),
		},
		{
			name: "timeout",
			err:  errors.New("timeout exceeded"),
		},
		{
			name: "generic network error",
			err:  errors.New("network error occurred"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRoundTripper{
				returnError: tt.err,
			}

			config := &CircuitBreakerConfig{
				Provider:    "test",
				MaxRequests: 1,
				Interval:    60 * time.Second,
				Timeout:     60 * time.Second,
				Threshold:   3,
			}

			cbt := NewCircuitBreakerTransport(mock, config)

			// Make requests to trigger failures
			for i := 0; i < 3; i++ {
				req := httptest.NewRequest("GET", "http://example.com", nil)
				resp, err := cbt.RoundTrip(req)

				if err == nil {
					t.Errorf("request %d: expected error, got nil", i+1)
				}
				if resp != nil {
					t.Errorf("request %d: expected nil response for network error, got %v", i+1, resp)
				}
			}

			// Verify circuit opened after threshold failures
			if cbt.State() != gobreaker.StateOpen {
				t.Errorf("after network errors, State() = %v, want %v", cbt.State(), gobreaker.StateOpen)
			}
		})
	}
}

func TestCircuitBreakerTransport_StateAndCounts(t *testing.T) {
	mock := &mockRoundTripper{
		statusCode: http.StatusOK,
	}

	config := &CircuitBreakerConfig{
		Provider:    "test",
		MaxRequests: 1,
		Interval:    60 * time.Second,
		Timeout:     60 * time.Second,
		Threshold:   5,
	}

	cbt := NewCircuitBreakerTransport(mock, config)

	// Test State() method
	state := cbt.State()
	if state != gobreaker.StateClosed {
		t.Errorf("State() = %v, want %v", state, gobreaker.StateClosed)
	}

	// Test Counts() method
	counts := cbt.Counts()
	if counts.Requests != 0 {
		t.Errorf("initial Counts().Requests = %v, want 0", counts.Requests)
	}
	if counts.TotalSuccesses != 0 {
		t.Errorf("initial Counts().TotalSuccesses = %v, want 0", counts.TotalSuccesses)
	}
	if counts.TotalFailures != 0 {
		t.Errorf("initial Counts().TotalFailures = %v, want 0", counts.TotalFailures)
	}

	// Make a successful request
	req := httptest.NewRequest("GET", "http://example.com", nil)
	cbt.RoundTrip(req)

	// Check counts after successful request
	counts = cbt.Counts()
	if counts.Requests != 1 {
		t.Errorf("after 1 request, Counts().Requests = %v, want 1", counts.Requests)
	}
	if counts.TotalSuccesses != 1 {
		t.Errorf("after 1 success, Counts().TotalSuccesses = %v, want 1", counts.TotalSuccesses)
	}
}

func TestIsFailureStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"200 OK", http.StatusOK, false},
		{"201 Created", http.StatusCreated, false},
		{"204 No Content", http.StatusNoContent, false},
		{"400 Bad Request", http.StatusBadRequest, false},
		{"401 Unauthorized", http.StatusUnauthorized, false},
		{"403 Forbidden", http.StatusForbidden, false},
		{"404 Not Found", http.StatusNotFound, false},
		{"429 Too Many Requests", http.StatusTooManyRequests, true},
		{"500 Internal Server Error", http.StatusInternalServerError, true},
		{"502 Bad Gateway", http.StatusBadGateway, true},
		{"503 Service Unavailable", http.StatusServiceUnavailable, true},
		{"504 Gateway Timeout", http.StatusGatewayTimeout, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isFailureStatusCode(tt.statusCode)
			if got != tt.want {
				t.Errorf("isFailureStatusCode(%d) = %v, want %v", tt.statusCode, got, tt.want)
			}
		})
	}
}

func TestIsRetryableHTTPError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"circuit breaker open", gobreaker.ErrOpenState, false},
		{"connection refused", errors.New("connection refused"), true},
		{"connection reset", errors.New("connection reset by peer"), true},
		{"timeout", errors.New("request timeout"), true},
		{"503 error", errors.New("HTTP 503 Service Unavailable"), true},
		{"502 error", errors.New("HTTP 502 Bad Gateway"), true},
		{"504 error", errors.New("HTTP 504 Gateway Timeout"), true},
		{"429 error", errors.New("HTTP 429 Too Many Requests"), true},
		{"generic error", errors.New("something went wrong"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRetryableHTTPError(tt.err)
			if got != tt.want {
				t.Errorf("isRetryableHTTPError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
