# Obot Tools

The official repository for Obot platform tools, model providers, and authentication providers.

[![Build Status](https://github.com/obot-platform/tools/workflows/Test/badge.svg)](https://github.com/obot-platform/tools/actions)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://go.dev/)
[![Python Version](https://img.shields.io/badge/Python-3.13+-3776AB?logo=python)](https://python.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

## Overview

This repository contains a comprehensive suite of components that extend the [Obot platform](https://github.com/obot-platform/obot):

- **🛠️ Core Tools**: Memory, knowledge base (RAG), tasks, file operations, and more
- **🤖 Model Providers**: Integrations with OpenAI, Anthropic, Ollama, Groq, xAI, DeepSeek, Voyage, and more
- **🔐 Auth Providers**: OAuth2-based authentication for GitHub, Google, Keycloak, and Microsoft Entra ID
- **🔑 Credential Management**: Secure credential storage with SQLite and PostgreSQL backends

## Quick Start

### Prerequisites

- Go 1.23+
- Make
- Git

### Build All Tools

```bash
git clone https://github.com/obot-platform/tools.git
cd tools
make build
```

### Run Tests

```bash
make test
```

### Package for Distribution

```bash
make package-tools
make package-providers
```

## Documentation

### 📚 Core Documentation

- **[Architecture Guide](docs/ARCHITECTURE.md)** - System architecture, component design, and data flows
- **[API Reference](docs/API_REFERENCE.md)** - Complete API documentation for all components
- **[Developer Guide](docs/DEVELOPER_GUIDE.md)** - Development setup, building, testing, and contributing
- **[Auth Provider Requirements](docs/auth-providers.md)** - Specifications for building auth providers

### 🔧 Component Categories

#### Core Tools

User-facing capabilities that provide essential agent functionality:

| Tool | Description | Language |
| ------ | ------------- | ---------- |
| [memory](memory/) | Long-term agent memory storage | Go |
| [knowledge](knowledge/) | Knowledge base with RAG capabilities | Go |
| [tasks](tasks/) | Task management and execution | Go |
| [workspace-files](workspace-files/) | Workspace file operations | Go |
| [threads](threads/) | Conversation thread management | GPTScript |
| [time](time/) | Time and scheduling utilities | GPTScript |
| [images](images/) | Image generation and analysis | TypeScript |
| [file-summarizer](file-summarizer/) | Document summarization | Python |

#### Model Providers

AI model integrations with standardized APIs:

| Provider | Language | Models |
| ---------- | ---------- | -------- |
| [openai-model-provider](openai-model-provider/) | Go | GPT-4, GPT-3.5, embeddings, images |
| [anthropic-model-provider-go](anthropic-model-provider-go/) | Go | Claude 3.5 Sonnet, Opus, Haiku |
| [ollama-model-provider](ollama-model-provider/) | Go | Local models (Llama, Mistral, etc.) |
| [groq-model-provider](groq-model-provider/) | Go | Groq-hosted models |
| [xai-model-provider](xai-model-provider/) | Go | xAI Grok models |
| [deepseek-model-provider](deepseek-model-provider/) | Go | DeepSeek models |
| [voyage-model-provider](voyage-model-provider/) | Python | Voyage embeddings |
| [vllm-model-provider](vllm-model-provider/) | Go | Self-hosted VLLM models |
| [generic-openai-model-provider](generic-openai-model-provider/) | Go | OpenAI-compatible APIs |

#### Authentication Providers

OAuth2-based user authentication:

| Provider | Description |
| ---------- | ------------- |
| [github-auth-provider](github-auth-provider/) | GitHub OAuth2 with org/team/repo validation |
| [google-auth-provider](google-auth-provider/) | Google OAuth2 (reference implementation) |
| [keycloak-auth-provider](keycloak-auth-provider/) | Keycloak/RHSSO OAuth2 with realm/group validation |
| [entra-auth-provider](entra-auth-provider/) | Microsoft Entra ID (Azure AD) with group validation |

#### Credential Management

Secure credential storage solutions:

| Store | Backend |
| ---------- | --------- |
| [sqlite](credential-stores/sqlite/) | SQLite with AES-256 encryption |
| [postgres](credential-stores/postgres/) | PostgreSQL with AES-256 encryption |

## Development

### Building a Specific Component

**Go Component**:

```bash
cd openai-model-provider
go build -ldflags="-s -w" -o bin/gptscript-go-tool .
```

**Python Component**:

```bash
cd voyage-model-provider
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python main.py
```

### Running Tests

**All Tests**:

```bash
make test
```

**Specific Component**:

```bash
cd github-auth-provider
go test ./...
```

**With Coverage**:

```bash
go test -cover ./...
```

### Environment Configuration

**Model Provider Example**:

```bash
export OBOT_OPENAI_MODEL_PROVIDER_API_KEY=sk-...
export PORT=8000
export GPTSCRIPT_DEBUG=true
cd openai-model-provider
go run .
```

**Auth Provider Example**:

```bash
export OBOT_GITHUB_AUTH_PROVIDER_CLIENT_ID=...
export OBOT_GITHUB_AUTH_PROVIDER_CLIENT_SECRET=...
export OBOT_AUTH_PROVIDER_COOKIE_SECRET=$(openssl rand -base64 32)
export OBOT_AUTH_PROVIDER_EMAIL_DOMAINS=example.com
cd github-auth-provider
go run .
```

## Architecture Highlights

### Modular Design

Each component is independently deployable while maintaining consistent interfaces:

```
┌─────────────────────────────────────┐
│       Obot Platform                 │
└─────────────┬───────────────────────┘
              │
       ┌──────┴──────┐
       │ index.yaml  │ (Tool Registry)
       └──────┬──────┘
              │
    ┌─────────┼─────────┐
    │         │         │
┌───▼───┐ ┌──▼──┐ ┌───▼────┐
│ Tools │ │Model│ │  Auth  │
│       │ │Prvdrs│ │Prvdrs │
└───────┘ └─────┘ └────────┘
```

### Technology Stack

- **Go 1.23**: Core infrastructure (model providers, auth, system tools)
- **Python 3.13**: Specialized ML/AI tools (embeddings, document processing)
- **TypeScript/Node.js**: Image processing tools
- **FastAPI**: Python web services
- **PostgreSQL/SQLite**: Vector storage and metadata indexing
- **Docker**: Multi-stage containerized builds

### Key Features

- **Standardized APIs**: OpenAI-compatible model provider APIs
- **Security**: AES-256 credential encryption, OAuth2 authentication, encrypted cookies
- **Scalability**: Stateless design, horizontal scaling support
- **Extensibility**: Plugin architecture via GPTScript tool definitions

## Docker Deployment

```bash
# Build Docker image
make docker-build

# Or manually
docker build -t obot-platform/tools:latest --target tools .
```

Multi-stage build includes:

- **base**: Wolfi base with Go, Python, Node.js
- **tools**: All user-facing tools
- **providers**: Model and auth providers

## Contributing

We welcome contributions! Please see our [Developer Guide](docs/DEVELOPER_GUIDE.md) for details.

### Development Workflow

1. Fork and clone the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make changes and add tests
4. Run tests (`make test`)
5. Build (`make build`)
6. Commit with conventional commits (`feat: add feature`)
7. Push and create a pull request

### Code Style

- **Go**: Follow [Effective Go](https://go.dev/doc/effective_go), use `go fmt`
- **Python**: Follow PEP 8, use `black` and `isort`
- **Commits**: Use [Conventional Commits](https://www.conventionalcommits.org/)

## Issues and Support

- **Bug Reports & Feature Requests**: Use the [main Obot repo](https://github.com/obot-platform/obot/issues) with the `tools` label
- **Documentation**: See [docs/](docs/) directory
- **Questions**: Open a discussion in the main Obot repository

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.

## Related Projects

- **[Obot Platform](https://github.com/obot-platform/obot)** - Main Obot platform repository
- **[GPTScript](https://github.com/gptscript-ai/gptscript)** - Tool execution runtime

## Tool Registry

All tools are registered in [`index.yaml`](index.yaml) with four main sections:

- `tools:` - User-facing capabilities
- `system:` - Internal system utilities
- `modelProviders:` - AI model integrations
- `authProviders:` - Authentication providers

See the [Architecture Guide](docs/ARCHITECTURE.md) for detailed information on the registry and tool definitions.
