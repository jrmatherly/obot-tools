# Anthropic Model Provider

A production-ready Anthropic API proxy with circuit breaker protection for resilient Claude LLM integrations.

## Features

- ✅ **Anthropic API Proxy** - Full compatibility with Anthropic Claude API endpoints
- ✅ **Circuit Breaker Protection** - Automatic failure detection and recovery
- ✅ **Prometheus Metrics** - Complete observability of API health and circuit state
- ✅ **Environment Configuration** - Zero-code configuration via environment variables

## Quick Start

### Basic Usage

```bash
# Set your Anthropic API key
export OBOT_ANTHROPIC_MODEL_PROVIDER_API_KEY="sk-ant-..."

# Start the provider
./anthropic-model-provider-go
```

The provider will start on port 8080 by default.

### With Circuit Breaker Protection

```bash
# Enable circuit breaker for production resilience
export OBOT_ANTHROPIC_MODEL_PROVIDER_API_KEY="sk-ant-..."
export OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_ENABLED=true

# Optional: Customize circuit breaker settings
export OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_THRESHOLD=5
export OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_TIMEOUT=60s
export OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_INTERVAL=60s

# Start the provider
./anthropic-model-provider-go
```

**Expected startup log when circuit breaker is enabled:**
```
[model-provider: Anthropic] Circuit breaker enabled (threshold=5, timeout=60s, interval=60s)
```

## Circuit Breaker Configuration

The circuit breaker protects against cascading failures when the Anthropic API becomes unavailable, rate-limited, or experiences high latency.

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_ENABLED` | `false` | Enable circuit breaker protection (opt-in) |
| `OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_THRESHOLD` | `5` | Number of consecutive failures before opening circuit |
| `OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_TIMEOUT` | `60s` | How long circuit stays open before testing recovery |
| `OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_INTERVAL` | `60s` | Time window for counting failures |

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
obot_model_provider_circuit_breaker_state{provider="Anthropic"}

# Request counts by outcome
obot_model_provider_circuit_breaker_requests_total{provider="Anthropic",outcome="success"}
obot_model_provider_circuit_breaker_requests_total{provider="Anthropic",outcome="failure"}
obot_model_provider_circuit_breaker_requests_total{provider="Anthropic",outcome="rejected"}

# Request latency histogram
obot_model_provider_circuit_breaker_latency_seconds{provider="Anthropic"}
```

### Configuration Examples

#### Production (Conservative)

```bash
# Tolerate some failures before protecting
export OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_ENABLED=true
export OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_THRESHOLD=5
export OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_TIMEOUT=60s
export OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_INTERVAL=60s
```

#### Aggressive (Unstable Upstreams)

```bash
# Open quickly on failures, recover faster
export OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_ENABLED=true
export OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_THRESHOLD=3
export OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_TIMEOUT=30s
export OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_INTERVAL=30s
```

#### Kubernetes Deployment

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: anthropic-provider-config
data:
  OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_ENABLED: "true"
  OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_THRESHOLD: "5"
  OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_TIMEOUT: "60s"
  OBOT_ANTHROPIC_MODEL_PROVIDER_CIRCUIT_BREAKER_INTERVAL: "60s"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: anthropic-model-provider
spec:
  template:
    spec:
      containers:
      - name: provider
        image: anthropic-model-provider-go:latest
        envFrom:
        - configMapRef:
            name: anthropic-provider-config
        - secretRef:
            name: anthropic-api-key  # Contains OBOT_ANTHROPIC_MODEL_PROVIDER_API_KEY
```

### Monitoring and Alerting

**Grafana Dashboard Query - Circuit State:**
```promql
obot_model_provider_circuit_breaker_state{provider="Anthropic"}
```

**Alert Example - Circuit Breaker Open:**
```yaml
- alert: AnthropicCircuitBreakerOpen
  expr: obot_model_provider_circuit_breaker_state{provider="Anthropic"} == 2
  for: 1m
  labels:
    severity: warning
  annotations:
    summary: "Anthropic circuit breaker is open"
    description: "Circuit breaker has been open for 1 minute, Anthropic API may be unhealthy"
```

### Troubleshooting

**Circuit opens unexpectedly:**
```bash
# Check metrics
curl http://localhost:9090/metrics | grep circuit_breaker

# Review logs for error details
kubectl logs -f <provider-pod>

# Verify Anthropic API status
curl https://status.anthropic.com/api/v2/status.json
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
export OBOT_ANTHROPIC_MODEL_PROVIDER_API_KEY="sk-ant-..."
```

### Base URL (Custom Anthropic Endpoint)

```bash
export OBOT_ANTHROPIC_MODEL_PROVIDER_BASE_URL="https://api.anthropic.com"
```

### Port

```bash
export OBOT_ANTHROPIC_MODEL_PROVIDER_PORT=8080
```

## Testing

### Unit Tests

```bash
cd /path/to/anthropic-model-provider-go
go test ./... -v
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
go build -o anthropic-model-provider-go .

# Build with Docker
docker build -t anthropic-model-provider-go:latest .
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
- **Provider Issues**: Check provider logs and Anthropic API status
- **Metrics**: Verify Prometheus endpoint at `http://localhost:9090/metrics`

## License

Apache License 2.0
