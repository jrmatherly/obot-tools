module github.com/obot-platform/tools/obot-model-provider

go 1.23.4

require (
	github.com/gptscript-ai/chat-completion-client v0.0.0-20241216203633-5c0178fb89ed
	github.com/obot-platform/tools/pkg/metrics v0.0.0-00010101000000-000000000000
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_golang v1.20.5 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	golang.org/x/sys v0.22.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/obot-platform/tools/pkg/metrics => ../pkg/metrics
