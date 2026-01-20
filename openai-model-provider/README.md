# OpenAI Model Provider

A production-ready OpenAI API proxy with circuit breaker protection for resilient LLM integrations.

## Features

- ✅ **OpenAI API Proxy** - Full compatibility with OpenAI API endpoints
- ✅ **Circuit Breaker Protection** - Automatic failure detection and recovery
- ✅ **Prometheus Metrics** - Complete observability of API health and circuit state
- ✅ **Environment Configuration** - Zero-code configuration via environment variables

## Quick Start

### Basic Usage

```bash
# Set your OpenAI API key
export OBOT_OPENAI_MODEL_PROVIDER_API_KEY="sk-..."

# Start the provider
./openai-model-provider
```

The provider will start on port 8080 by default.

### With Circuit Breaker Protection

```bash
# Enable circuit breaker for production resilience
export OBOT_OPENAI_MODEL_PROVIDER_API_KEY="sk-..."
export OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_ENABLED=true

# Optional: Customize circuit breaker settings
export OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_THRESHOLD=5
export OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_TIMEOUT=60s
export OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_INTERVAL=60s

# Start the provider
./openai-model-provider
```

**Expected startup log when circuit breaker is enabled:**
```
[model-provider: OpenAI] Circuit breaker enabled (threshold=5, timeout=60s, interval=60s)
```

## Circuit Breaker Configuration

The circuit breaker protects against cascading failures when the OpenAI API becomes unavailable, rate-limited, or experiences high latency.

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_ENABLED` | `false` | Enable circuit breaker protection (opt-in) |
| `OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_THRESHOLD` | `5` | Number of consecutive failures before opening circuit |
| `OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_TIMEOUT` | `60s` | How long circuit stays open before testing recovery |
| `OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_INTERVAL` | `60s` | Time window for counting failures |

### How It Works

```
CLOSED (normal) ──[5 failures]──> OPEN (protecting)
    ↑                                    │
    │                            [wait 60s]
    │                                    ↓
    └────[success]──── HALF-OPEN (testing)
```

1. **CLOSED**: Normal operation, all requests go through
2. **OPEN**: Circuit breaker protecting upstream, requests fail fast (<100ms)
3. **HALF-OPEN**: Testing recovery, single request allowed

### What Counts as a Failure

- **HTTP Status Codes**: 429 (Rate Limited), 500, 502, 503, 504 (Server Errors)
- **Network Errors**: Connection refused, timeout, DNS resolution failures
- **Client Errors (4xx)**: NOT counted (these are application issues, not service issues)

### Prometheus Metrics

When circuit breaker is enabled, the following metrics are exposed:

```prometheus
# Circuit state: 0=closed, 1=half-open, 2=open
obot_model_provider_circuit_breaker_state{provider="OpenAI"}

# Request counts by outcome
obot_model_provider_circuit_breaker_requests_total{provider="OpenAI",outcome="success"}
obot_model_provider_circuit_breaker_requests_total{provider="OpenAI",outcome="failure"}
obot_model_provider_circuit_breaker_requests_total{provider="OpenAI",outcome="rejected"}

# Request latency histogram
obot_model_provider_circuit_breaker_latency_seconds{provider="OpenAI"}
```

### Configuration Examples

#### Production (Conservative)

```bash
# Tolerate some failures before protecting
export OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_ENABLED=true
export OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_THRESHOLD=5
export OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_TIMEOUT=60s
export OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_INTERVAL=60s
```

#### Aggressive (Unstable Upstreams)

```bash
# Open quickly on failures, recover faster
export OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_ENABLED=true
export OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_THRESHOLD=3
export OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_TIMEOUT=30s
export OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_INTERVAL=30s
```

#### Kubernetes Deployment

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: openai-provider-config
data:
  OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_ENABLED: "true"
  OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_THRESHOLD: "5"
  OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_TIMEOUT: "60s"
  OBOT_OPENAI_MODEL_PROVIDER_CIRCUIT_BREAKER_INTERVAL: "60s"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: openai-model-provider
spec:
  template:
    spec:
      containers:
      - name: provider
        image: openai-model-provider:latest
        envFrom:
        - configMapRef:
            name: openai-provider-config
        - secretRef:
            name: openai-api-key  # Contains OBOT_OPENAI_MODEL_PROVIDER_API_KEY
```

### Monitoring and Alerting

**Grafana Dashboard Query - Circuit State:**
```promql
obot_model_provider_circuit_breaker_state{provider="OpenAI"}
```

**Alert Example - Circuit Breaker Open:**
```yaml
- alert: OpenAICircuitBreakerOpen
  expr: obot_model_provider_circuit_breaker_state{provider="OpenAI"} == 2
  for: 1m
  labels:
    severity: warning
  annotations:
    summary: "OpenAI circuit breaker is open"
    description: "Circuit breaker has been open for 1 minute, OpenAI API may be unhealthy"
```

### Troubleshooting

**Circuit opens unexpectedly:**
```bash
# Check metrics
curl http://localhost:9090/metrics | grep circuit_breaker

# Review logs for error details
kubectl logs -f <provider-pod>

# Verify OpenAI API status
curl https://status.openai.com/api/v2/status.json
```

**Circuit doesn't open on failures:**
```bash
# Verify circuit breaker is enabled
kubectl logs <provider-pod> | grep "Circuit breaker"

# Check if it's HTTP 4xx errors (not counted as failures)
curl http://localhost:9090/metrics | grep outcome=\"failure\"
```

## Additional Configuration

### API Key

```bash
export OBOT_OPENAI_MODEL_PROVIDER_API_KEY="sk-..."
```

### Base URL (Custom OpenAI Endpoint)

```bash
export OBOT_OPENAI_MODEL_PROVIDER_BASE_URL="https://api.openai.com"
```

### Port

```bash
export OBOT_OPENAI_MODEL_PROVIDER_PORT=8080
```

## Testing

### Unit Tests

```bash
cd /path/to/openai-model-provider
go test ./proxy -v
```

### Integration Test - Circuit Breaker

Test that circuit breaker opens on failures and recovers:

```bash
# See docs/CIRCUIT_BREAKER.md for comprehensive testing guide
cd /path/to/AI
./run-circuit-breaker-test.sh
./run-circuit-breaker-recovery-test.sh
```

## Building

```bash
# Build binary
go build -o openai-model-provider .

# Build with Docker
docker build -t openai-model-provider:latest .
```

## Documentation

For comprehensive circuit breaker documentation, including:
- Complete configuration reference for all providers
- Prometheus metrics and Grafana dashboards
- Troubleshooting guide
- Testing procedures
- Best practices and production recommendations

See: **[docs/CIRCUIT_BREAKER.md](../../docs/CIRCUIT_BREAKER.md)**

## Support

- **Circuit Breaker Issues**: See [docs/CIRCUIT_BREAKER.md](../../docs/CIRCUIT_BREAKER.md)
- **Provider Issues**: Check provider logs and OpenAI API status
- **Metrics**: Verify Prometheus endpoint at `http://localhost:9090/metrics`

## License

Apache License 2.0
