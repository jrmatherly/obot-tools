package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ModelProviderRequestsTotal tracks total requests per provider and model
	ModelProviderRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "obot_model_provider_requests_total",
			Help: "Total number of model provider requests by provider, model, and status",
		},
		[]string{"provider", "model", "status"},
	)

	// ModelProviderRequestDuration tracks request latency distribution
	ModelProviderRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "obot_model_provider_request_duration_seconds",
			Help:    "Model provider request latency in seconds by provider and model",
			Buckets: prometheus.DefBuckets, // [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]
		},
		[]string{"provider", "model"},
	)

	// ModelProviderActiveRequests tracks concurrent requests
	ModelProviderActiveRequests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "obot_model_provider_active_requests",
			Help: "Current number of active model provider requests by provider and model",
		},
		[]string{"provider", "model"},
	)

	// ModelProviderErrorsTotal tracks errors by provider, model, and error type
	ModelProviderErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "obot_model_provider_errors_total",
			Help: "Total number of model provider errors by provider, model, and error type",
		},
		[]string{"provider", "model", "error_type"},
	)

	// ModelProviderTokensTotal tracks token usage (input/output)
	ModelProviderTokensTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "obot_model_provider_tokens_total",
			Help: "Total number of tokens processed by provider, model, and token type (input/output)",
		},
		[]string{"provider", "model", "token_type"},
	)
)

// RecordRequestSuccess records a successful request
func RecordRequestSuccess(provider, model string, duration time.Duration) {
	ModelProviderRequestsTotal.WithLabelValues(provider, model, "success").Inc()
	ModelProviderRequestDuration.WithLabelValues(provider, model).Observe(duration.Seconds())
}

// RecordRequestFailure records a failed request
func RecordRequestFailure(provider, model string, duration time.Duration) {
	ModelProviderRequestsTotal.WithLabelValues(provider, model, "failure").Inc()
	ModelProviderRequestDuration.WithLabelValues(provider, model).Observe(duration.Seconds())
}

// RecordRequestError records a request error
func RecordRequestError(provider, model string, duration time.Duration) {
	ModelProviderRequestsTotal.WithLabelValues(provider, model, "error").Inc()
	ModelProviderRequestDuration.WithLabelValues(provider, model).Observe(duration.Seconds())
}

// RecordError records an error by type
func RecordError(provider, model, errorType string) {
	ModelProviderErrorsTotal.WithLabelValues(provider, model, errorType).Inc()
}

// RecordInputTokens records input token usage
func RecordInputTokens(provider, model string, tokens int) {
	ModelProviderTokensTotal.WithLabelValues(provider, model, "input").Add(float64(tokens))
}

// RecordOutputTokens records output token usage
func RecordOutputTokens(provider, model string, tokens int) {
	ModelProviderTokensTotal.WithLabelValues(provider, model, "output").Add(float64(tokens))
}

// RecordTokens records both input and output token usage
func RecordTokens(provider, model string, inputTokens, outputTokens int) {
	RecordInputTokens(provider, model, inputTokens)
	RecordOutputTokens(provider, model, outputTokens)
}

// TrackRequestDuration returns a function to be deferred that tracks request duration
func TrackRequestDuration(provider, model string, startTime time.Time) func(error) {
	ModelProviderActiveRequests.WithLabelValues(provider, model).Inc()

	return func(err error) {
		defer ModelProviderActiveRequests.WithLabelValues(provider, model).Dec()

		duration := time.Since(startTime)
		if err != nil {
			RecordRequestError(provider, model, duration)
		} else {
			RecordRequestSuccess(provider, model, duration)
		}
	}
}
