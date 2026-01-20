package metrics

import (
	"errors"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestRecordRequestSuccess(t *testing.T) {
	// Reset the counter before testing
	ModelProviderRequestsTotal.Reset()

	provider := "openai"
	model := "gpt-4"
	duration := 100 * time.Millisecond

	initialCount := testutil.ToFloat64(ModelProviderRequestsTotal.WithLabelValues(provider, model, "success"))

	RecordRequestSuccess(provider, model, duration)

	finalCount := testutil.ToFloat64(ModelProviderRequestsTotal.WithLabelValues(provider, model, "success"))

	if finalCount != initialCount+1 {
		t.Errorf("Expected counter to increment by 1, got %f -> %f", initialCount, finalCount)
	}
}

func TestRecordRequestFailure(t *testing.T) {
	// Reset the counter before testing
	ModelProviderRequestsTotal.Reset()

	provider := "anthropic"
	model := "claude-3-5-sonnet"
	duration := 200 * time.Millisecond

	initialCount := testutil.ToFloat64(ModelProviderRequestsTotal.WithLabelValues(provider, model, "failure"))

	RecordRequestFailure(provider, model, duration)

	finalCount := testutil.ToFloat64(ModelProviderRequestsTotal.WithLabelValues(provider, model, "failure"))

	if finalCount != initialCount+1 {
		t.Errorf("Expected counter to increment by 1, got %f -> %f", initialCount, finalCount)
	}
}

func TestRecordRequestError(t *testing.T) {
	// Reset the counter before testing
	ModelProviderRequestsTotal.Reset()

	provider := "vllm"
	model := "llama-2"
	duration := 50 * time.Millisecond

	initialCount := testutil.ToFloat64(ModelProviderRequestsTotal.WithLabelValues(provider, model, "error"))

	RecordRequestError(provider, model, duration)

	finalCount := testutil.ToFloat64(ModelProviderRequestsTotal.WithLabelValues(provider, model, "error"))

	if finalCount != initialCount+1 {
		t.Errorf("Expected counter to increment by 1, got %f -> %f", initialCount, finalCount)
	}
}

func TestRecordError(t *testing.T) {
	// Reset the counter before testing
	ModelProviderErrorsTotal.Reset()

	provider := "openai"
	model := "gpt-4"
	errorType := "rate_limit"

	initialCount := testutil.ToFloat64(ModelProviderErrorsTotal.WithLabelValues(provider, model, errorType))

	RecordError(provider, model, errorType)

	finalCount := testutil.ToFloat64(ModelProviderErrorsTotal.WithLabelValues(provider, model, errorType))

	if finalCount != initialCount+1 {
		t.Errorf("Expected error counter to increment by 1, got %f -> %f", initialCount, finalCount)
	}
}

func TestRecordInputTokens(t *testing.T) {
	// Reset the counter before testing
	ModelProviderTokensTotal.Reset()

	provider := "anthropic"
	model := "claude-3-5-sonnet"
	tokens := 150

	initialCount := testutil.ToFloat64(ModelProviderTokensTotal.WithLabelValues(provider, model, "input"))

	RecordInputTokens(provider, model, tokens)

	finalCount := testutil.ToFloat64(ModelProviderTokensTotal.WithLabelValues(provider, model, "input"))

	expected := initialCount + float64(tokens)
	if finalCount != expected {
		t.Errorf("Expected token counter to be %f, got %f", expected, finalCount)
	}
}

func TestRecordOutputTokens(t *testing.T) {
	// Reset the counter before testing
	ModelProviderTokensTotal.Reset()

	provider := "vllm"
	model := "llama-2"
	tokens := 250

	initialCount := testutil.ToFloat64(ModelProviderTokensTotal.WithLabelValues(provider, model, "output"))

	RecordOutputTokens(provider, model, tokens)

	finalCount := testutil.ToFloat64(ModelProviderTokensTotal.WithLabelValues(provider, model, "output"))

	expected := initialCount + float64(tokens)
	if finalCount != expected {
		t.Errorf("Expected token counter to be %f, got %f", expected, finalCount)
	}
}

func TestRecordTokens(t *testing.T) {
	// Reset the counter before testing
	ModelProviderTokensTotal.Reset()

	provider := "openai"
	model := "gpt-4"
	inputTokens := 100
	outputTokens := 200

	initialInputCount := testutil.ToFloat64(ModelProviderTokensTotal.WithLabelValues(provider, model, "input"))
	initialOutputCount := testutil.ToFloat64(ModelProviderTokensTotal.WithLabelValues(provider, model, "output"))

	RecordTokens(provider, model, inputTokens, outputTokens)

	finalInputCount := testutil.ToFloat64(ModelProviderTokensTotal.WithLabelValues(provider, model, "input"))
	finalOutputCount := testutil.ToFloat64(ModelProviderTokensTotal.WithLabelValues(provider, model, "output"))

	expectedInput := initialInputCount + float64(inputTokens)
	expectedOutput := initialOutputCount + float64(outputTokens)

	if finalInputCount != expectedInput {
		t.Errorf("Expected input token counter to be %f, got %f", expectedInput, finalInputCount)
	}
	if finalOutputCount != expectedOutput {
		t.Errorf("Expected output token counter to be %f, got %f", expectedOutput, finalOutputCount)
	}
}

func TestTrackRequestDuration_Success(t *testing.T) {
	// Reset metrics before testing
	ModelProviderActiveRequests.Reset()
	ModelProviderRequestsTotal.Reset()

	provider := "anthropic"
	model := "claude-3-5-sonnet"
	startTime := time.Now()

	initialActive := testutil.ToFloat64(ModelProviderActiveRequests.WithLabelValues(provider, model))
	initialSuccess := testutil.ToFloat64(ModelProviderRequestsTotal.WithLabelValues(provider, model, "success"))

	// Track request duration
	done := TrackRequestDuration(provider, model, startTime)

	// Check that active requests were incremented
	activeCount := testutil.ToFloat64(ModelProviderActiveRequests.WithLabelValues(provider, model))
	if activeCount != initialActive+1 {
		t.Errorf("Expected active requests to increment by 1, got %f -> %f", initialActive, activeCount)
	}

	// Simulate some work
	time.Sleep(10 * time.Millisecond)

	// Complete the request (no error)
	done(nil)

	// Check that active requests were decremented
	finalActive := testutil.ToFloat64(ModelProviderActiveRequests.WithLabelValues(provider, model))
	if finalActive != initialActive {
		t.Errorf("Expected active requests to return to initial value %f, got %f", initialActive, finalActive)
	}

	// Check that success counter was incremented
	finalSuccess := testutil.ToFloat64(ModelProviderRequestsTotal.WithLabelValues(provider, model, "success"))
	if finalSuccess != initialSuccess+1 {
		t.Errorf("Expected success counter to increment by 1, got %f -> %f", initialSuccess, finalSuccess)
	}
}

func TestTrackRequestDuration_Error(t *testing.T) {
	// Reset metrics before testing
	ModelProviderActiveRequests.Reset()
	ModelProviderRequestsTotal.Reset()

	provider := "openai"
	model := "gpt-4"
	startTime := time.Now()

	initialActive := testutil.ToFloat64(ModelProviderActiveRequests.WithLabelValues(provider, model))
	initialError := testutil.ToFloat64(ModelProviderRequestsTotal.WithLabelValues(provider, model, "error"))

	// Track request duration
	done := TrackRequestDuration(provider, model, startTime)

	// Check that active requests were incremented
	activeCount := testutil.ToFloat64(ModelProviderActiveRequests.WithLabelValues(provider, model))
	if activeCount != initialActive+1 {
		t.Errorf("Expected active requests to increment by 1, got %f -> %f", initialActive, activeCount)
	}

	// Simulate some work
	time.Sleep(10 * time.Millisecond)

	// Complete the request with an error
	done(errors.New("test error"))

	// Check that active requests were decremented
	finalActive := testutil.ToFloat64(ModelProviderActiveRequests.WithLabelValues(provider, model))
	if finalActive != initialActive {
		t.Errorf("Expected active requests to return to initial value %f, got %f", initialActive, finalActive)
	}

	// Check that error counter was incremented
	finalError := testutil.ToFloat64(ModelProviderRequestsTotal.WithLabelValues(provider, model, "error"))
	if finalError != initialError+1 {
		t.Errorf("Expected error counter to increment by 1, got %f -> %f", initialError, finalError)
	}
}

func TestHistogramRecordsDuration(t *testing.T) {
	// Reset histogram before testing
	ModelProviderRequestDuration.Reset()

	provider := "vllm"
	model := "llama-2"
	duration := 500 * time.Millisecond

	// Record a request
	RecordRequestSuccess(provider, model, duration)

	// Verify that the histogram has observations by collecting metrics
	// We cannot use ToFloat64 with histograms, so we check by collecting and counting
	count := testutil.CollectAndCount(ModelProviderRequestDuration)
	if count < 1 {
		t.Errorf("Expected at least 1 metric collected from histogram, got %d", count)
	}
}

func TestMetricsLabelsAreCorrect(t *testing.T) {
	// Reset all metrics
	ModelProviderRequestsTotal.Reset()
	ModelProviderErrorsTotal.Reset()
	ModelProviderTokensTotal.Reset()

	tests := []struct {
		name     string
		provider string
		model    string
	}{
		{
			name:     "openai provider",
			provider: "openai",
			model:    "gpt-4",
		},
		{
			name:     "anthropic provider",
			provider: "anthropic",
			model:    "claude-3-5-sonnet",
		},
		{
			name:     "vllm provider",
			provider: "vllm",
			model:    "llama-2",
		},
		{
			name:     "obot provider",
			provider: "obot",
			model:    "custom-model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Record metrics with labels
			RecordRequestSuccess(tt.provider, tt.model, 100*time.Millisecond)
			RecordError(tt.provider, tt.model, "timeout")
			RecordTokens(tt.provider, tt.model, 50, 100)

			// Verify metrics can be retrieved with correct labels
			successCount := testutil.ToFloat64(ModelProviderRequestsTotal.WithLabelValues(tt.provider, tt.model, "success"))
			if successCount < 1 {
				t.Errorf("Expected at least 1 success count for provider=%s, model=%s", tt.provider, tt.model)
			}

			errorCount := testutil.ToFloat64(ModelProviderErrorsTotal.WithLabelValues(tt.provider, tt.model, "timeout"))
			if errorCount < 1 {
				t.Errorf("Expected at least 1 error count for provider=%s, model=%s", tt.provider, tt.model)
			}

			inputTokens := testutil.ToFloat64(ModelProviderTokensTotal.WithLabelValues(tt.provider, tt.model, "input"))
			if inputTokens < 50 {
				t.Errorf("Expected at least 50 input tokens for provider=%s, model=%s", tt.provider, tt.model)
			}

			outputTokens := testutil.ToFloat64(ModelProviderTokensTotal.WithLabelValues(tt.provider, tt.model, "output"))
			if outputTokens < 100 {
				t.Errorf("Expected at least 100 output tokens for provider=%s, model=%s", tt.provider, tt.model)
			}
		})
	}
}

func TestGaugeIncreasesAndDecreases(t *testing.T) {
	// Reset gauge before testing
	ModelProviderActiveRequests.Reset()

	provider := "openai"
	model := "gpt-4"

	initial := testutil.ToFloat64(ModelProviderActiveRequests.WithLabelValues(provider, model))

	// Increase gauge
	ModelProviderActiveRequests.WithLabelValues(provider, model).Inc()
	afterInc := testutil.ToFloat64(ModelProviderActiveRequests.WithLabelValues(provider, model))

	if afterInc != initial+1 {
		t.Errorf("Expected gauge to increase by 1, got %f -> %f", initial, afterInc)
	}

	// Decrease gauge
	ModelProviderActiveRequests.WithLabelValues(provider, model).Dec()
	afterDec := testutil.ToFloat64(ModelProviderActiveRequests.WithLabelValues(provider, model))

	if afterDec != initial {
		t.Errorf("Expected gauge to return to initial value %f, got %f", initial, afterDec)
	}
}

func TestMultipleConcurrentRequests(t *testing.T) {
	// Reset metrics before testing
	ModelProviderActiveRequests.Reset()

	provider := "anthropic"
	model := "claude-3-5-sonnet"

	initial := testutil.ToFloat64(ModelProviderActiveRequests.WithLabelValues(provider, model))

	// Simulate 3 concurrent requests
	done1 := TrackRequestDuration(provider, model, time.Now())
	done2 := TrackRequestDuration(provider, model, time.Now())
	done3 := TrackRequestDuration(provider, model, time.Now())

	// Check that active requests shows 3 concurrent requests
	concurrent := testutil.ToFloat64(ModelProviderActiveRequests.WithLabelValues(provider, model))
	if concurrent != initial+3 {
		t.Errorf("Expected 3 active requests, got %f", concurrent-initial)
	}

	// Complete all requests
	done1(nil)
	done2(nil)
	done3(nil)

	// Check that active requests is back to initial
	final := testutil.ToFloat64(ModelProviderActiveRequests.WithLabelValues(provider, model))
	if final != initial {
		t.Errorf("Expected active requests to return to initial value %f, got %f", initial, final)
	}
}
