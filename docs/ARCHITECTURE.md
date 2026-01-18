# Obot Tools - Architecture Documentation

## Overview

The Obot Tools repository provides a comprehensive suite of tools, model providers, and authentication providers that extend the Obot platform. The architecture follows a modular design where each component is independently deployable while maintaining consistent interfaces.

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Obot Platform                            │
└─────────────────────┬───────────────────────────────────────┘
                      │
         ┌────────────┴──────────────────────────────────┐
         │   Tool Registry (index.yaml)                  │
         └────────────┬──────────────────────────────────┘
                      │
        ┌─────────────┼───────────┐
        │             │           │
   ┌────▼────┐  ┌────▼────┐  ┌────▼────┐
   │  Tools  │  │ Model   │  │  Auth   │
   │         │  │Providers│  │Providers│
   └─────────┘  └─────────┘  └─────────┘
```

## Component Categories

### 1. Core Tools (User-Facing Capabilities)

**Purpose**: Provide essential agent capabilities for memory, knowledge, tasks, and file operations.

**Components**:

- `memory/` - Long-term agent memory storage
- `knowledge/` - Knowledge base with RAG capabilities
- `tasks/` - Task management and execution
- `workspace-files/` - Workspace file operations
- `threads/` - Conversation thread management
- `time/` - Time and scheduling utilities
- `images/` - Image generation and analysis (TypeScript/Node.js)
- `file-summarizer/` - Document summarization (Python)

**Architecture Pattern**:

```
Tool Component
├── tool.gpt           # GPTScript tool definition
├── main.go/main.py    # Entry point
├── pkg/               # Package implementation (Go)
│   └── [domain]/      # Domain-specific logic
└── bin/               # Build output (gitignored)
```

### 2. System Tools (Internal Infrastructure)

**Purpose**: Internal system utilities used by the Obot platform.

**Components**:

- `workflow/` - Workflow execution engine
- `result-formatter/` - Result formatting and presentation
- `obot-model-provider/` - Internal model provider orchestration
- `tasks-workflow/` - Task workflow coordination
- `task-invoke/` - Task invocation system
- `loop-data/` - Loop data handling in workflows

**Key Characteristics**:

- Not directly exposed to users
- Referenced in index.yaml under `system:` section
- Tightly integrated with Obot core platform

### 3. Model Providers

**Purpose**: Standardized interfaces to various AI model providers.

**Architecture**:

```
Model Provider Component
├── tool.gpt              # Provider metadata and configuration
├── main.go/main.py       # HTTP server entry point
├── proxy/ (Go)           # Reverse proxy implementation
│   ├── proxy.go         # Core proxy logic
│   ├── validate.go      # Credential validation
│   └── rewrite.go       # Response rewriting
├── api/                  # API type definitions
└── server/ (Go)          # Server implementation
```

**Providers**:

- **OpenAI** (`openai-model-provider/`) - Go-based
- **Anthropic** (`anthropic-model-provider-go/`) - Go-based
- **Ollama** (`ollama-model-provider/`) - Go-based, local models
- **Groq** (`groq-model-provider/`) - Go-based
- **xAI** (`xai-model-provider/`) - Go-based
- **DeepSeek** (`deepseek-model-provider/`) - Go-based
- **Generic OpenAI** (`generic-openai-model-provider/`) - Compatible providers
- **Voyage** (`voyage-model-provider/`) - Python-based, embeddings
- **VLLM** (`vllm-model-provider/`) - Python-based, self-hosted

**Common Patterns**:

- HTTP server listening on configurable port (default: 8000)
- Support for global API keys via environment variables
- Per-request API keys via `X-Obot-<PROVIDER>_API_KEY` header
- `/v1/models` endpoint for model listing
- Model-specific endpoints (chat completions, embeddings, etc.)
- Validation mode for credential testing

### 4. Authentication Providers

**Purpose**: OAuth2-based user authentication for the Obot platform.

**Architecture**:

```
Auth Provider Component
├── tool.gpt              # Tool definition with envVars metadata
├── main.go               # HTTP server with OAuth2 flow
├── pkg/
│   └── profile/         # User profile management
│       ├── profile.go   # Profile retrieval
│       └── profile_test.go
└── Credential reference to placeholder-credential
```

**Providers**:

- **GitHub** (`github-auth-provider/`) - GitHub OAuth2
- **Google** (`google-auth-provider/`) - Google OAuth2 (reference implementation)

**Common OAuth2 Endpoints**:

- `/oauth2/start` - Initiate OAuth2 flow
- `/oauth2/callback` - Handle OAuth2 callback
- `/oauth2/sign_out` - Sign out and clear cookies
- `/obot-get-icon-url` - Retrieve user profile picture
- `/obot-get-state` - Get authenticated user state

**Security Features**:

- Encrypted `obot_access_token` cookie
- Cookie secret from `OBOT_AUTH_PROVIDER_COOKIE_SECRET`
- Secure cookies for HTTPS (`OBOT_SERVER_PUBLIC_URL` starts with https://)
- State parameter for CSRF protection
- User validation via email domains, teams, orgs, repos

**Common Utilities** (`auth-providers-common/`):

- State management (`pkg/state/`)
- Environment variable handling (`pkg/env/`)
- Icon URL retrieval (`pkg/icon/`)
- Error templates (`templates/`)

### 5. Credential Management

**Purpose**: Secure storage and retrieval of credentials.

**Components**:

- `credential-stores/` - Database-backed credential storage
  - `sqlite/` - SQLite implementation
  - `postgres/` - PostgreSQL implementation
  - `pkg/common/` - Shared encryption and DB logic
- `existing-credential/` - Use existing credentials from environment
- `generic-credential/` - Generic credential handling
- `placeholder-credential/` - Development placeholder (Python)

**Security Architecture**:

- AES-256 encryption for stored credentials
- Encryption keys from environment variables
- Database-backed persistence with encryption at rest
- Per-credential metadata and namespacing

## Knowledge Base (Deep Dive)

The knowledge base is one of the most complex components, providing RAG (Retrieval-Augmented Generation) capabilities.

### Architecture Layers

```
┌─────────────────────────────────────────────────┐
│         Knowledge Client (pkg/client/)          │
│  - Standalone mode                              │
│  - Metadata management                          │
│  - .knowignore support                          │
└─────────────────┬───────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────┐
│         Datastore (pkg/datastore/)              │
│  - Ingestion pipeline                           │
│  - Retrieval pipeline                           │
│  - Document processing                          │
└─────────┬───────────────────────────────────────┘
          │
    ┌─────┼─────┬────────┬─────────┐
    │     │     │        │         │
┌───▼─┐ ┌─▼──┐ ┌▼────┐ ┌─▼──────┐ ┌▼────────┐
│Index│ │Vec-│ │Text │ │Embed-  │ │Document │
│     │ │tor │ │Split│ │dings   │ │Loaders  │
│     │ │Store│ │ter │ │        │ │         │
└─────┘ └────┘ └─────┘ └────────┘ └─────────┘
```

### Key Subsystems

**Document Loading** (`pkg/datastore/documentloader/`):

- PDF processing (MuPDF, GoPDF, SmartPDF with OCR)
- Office document conversion (LibreOffice)
- Structured data (JSON, CSV)
- Remote loaders (GitHub, URLs)
- OCR support via OpenAI Vision API

**Text Splitting** (`pkg/datastore/textsplitter/`):

- Configurable chunk sizes and overlaps
- Language-aware splitting
- Metadata preservation

**Embeddings** (`pkg/datastore/embeddings/`):

- OpenAI embeddings
- Model provider abstraction
- Configurable dimensions

**Vector Stores** (`pkg/vectorstore/`):

- PostgreSQL + pgvector extension
- SQLite + sqlite-vec extension
- Cosine similarity search
- Metadata filtering

**Index Storage** (`pkg/index/`):

- SQLite and PostgreSQL backends
- Dataset management
- File tracking and metadata
- Incremental updates

**Retrievers** (`pkg/datastore/retrievers/`):

- Vector similarity retrieval
- BM25 keyword search
- Hybrid retrieval (vector + keyword)
- Subquery retrieval
- Routing retrieval
- Merging strategies

**Post-processors** (`pkg/datastore/postprocessors/`):

- Similarity score filtering
- BM25 reranking
- Cohere reranking
- Content substring matching
- Result deduplication
- Top-K reduction

**Query Modifiers** (`pkg/datastore/querymodifiers/`):

- Query enhancement with LLMs
- Spell checking
- Generic transformations

**Flows** (`pkg/flows/`):

- Configuration-driven pipelines
- Blueprint system (default, obot, context)
- YAML/JSON configuration
- Retriever and post-processor composition

## Build and Deployment Architecture

### Build System

**Multi-stage Docker Build**:

```dockerfile
Stage 1: base - Wolfi base with Go, Node.js, Python, pnpm, uv
Stage 2: tools - Build all tools (make package-tools)
Stage 3: providers - Build all providers (make package-providers)
```

**Makefile Targets**:

- `build` - Build all Go tools (finds main.go files)
- `test` - Run all tests (tools, auth providers, credential stores)
- `package-tools` - Package tools for distribution
- `package-providers` - Package providers for distribution
- `docker-build` - Build Docker image

**Build Script** (`scripts/build.sh`):

- Recursively finds all `main.go` files
- Skips `common` directories
- Builds to `bin/gptscript-go-tool` with optimized flags (`-ldflags="-s -w"`)

### Deployment Patterns

**Model Providers**:

- Run as HTTP daemons on configurable ports
- Environment-based configuration
- Stateless design for horizontal scaling
- Reverse proxy to upstream APIs

**Auth Providers**:

- Run as HTTP daemons
- Stateful OAuth2 flow with encrypted cookies
- Integration with OAuth2 Proxy library

**Tools**:

- Invoked via GPTScript runtime
- Binary execution model
- Context sharing via GPTScript

### CI/CD Pipeline

**GitHub Actions Workflows**:

1. **test.yaml** (Pull Requests):
   - Go tests with Go 1.23
   - Build verification (`make build`, `make test`)
   - Multi-platform Docker build (amd64, arm64)

2. **build-tools.yaml**:
   - Package tools for distribution
   - Artifact generation

3. **build-providers.yaml**:
   - Package providers for distribution
   - Artifact generation

4. **dispatch.yaml**:
   - Manual workflow triggers
   - Custom build orchestration

5. **dependabot-reviewers.yaml**:
   - Auto-assign reviewers to Dependabot PRs

## Integration Points

### GPTScript Integration

**Tool Definition Format** (`tool.gpt`):

```gptscript
Name: Tool Name
Description: Tool description
Share Context: context_name
Metadata: icon: https://...
Metadata: category: Capability
Metadata: envVars: REQUIRED_VAR1,REQUIRED_VAR2
Metadata: optionalEnvVars: OPTIONAL_VAR1

#!${GPTSCRIPT_TOOL_DIR}/bin/gptscript-go-tool command
```

**Context Sharing**:

- Tools can share context with other tools
- Context tools provide shared state
- Type: context annotation for context-only tools

**Credentials**:

- Referenced via `Credential:` directive
- Supports relative paths
- Integration with credential stores

### Obot Platform Integration

**Tool Registry** (`index.yaml`):

```yaml
tools:
  tool-name:
    reference: ./tool-directory

system:
  system-tool-name:
    reference: ./system-tool-directory

modelProviders:
  provider-name:
    reference: ./provider-directory

authProviders:
  auth-provider-name:
    reference: ./auth-provider-directory
```

**Platform Expectations**:

- Tools implement GPTScript tool protocol
- Model providers implement OpenAI-compatible APIs
- Auth providers implement Obot auth protocol
- Consistent metadata in tool.gpt files

## Data Flow Examples

### Model Provider Request Flow

```
User Request
    │
    ▼
Obot Platform
    │
    ▼
Model Provider HTTP Server (:8000)
    │
    ├─ Check X-Obot-*_API_KEY header (per-request)
    ├─ Fallback to OBOT_*_API_KEY env (global)
    │
    ▼
Proxy Layer
    │
    ├─ Rewrite request (if needed)
    ├─ Add authentication
    │
    ▼
Upstream API (OpenAI, Anthropic, etc.)
    │
    ▼
Response Rewriting (if needed)
    │
    ▼
Return to Obot Platform
```

### Authentication Flow

```
User (Browser)
    │
    ▼
/oauth2/start?rd=<redirect-url>
    │
    ├─ Generate state parameter
    ├─ Store redirect URL
    │
    ▼
Redirect to OAuth Provider (GitHub/Google)
    │
    ▼
User Authorizes
    │
    ▼
/oauth2/callback?code=...&state=...
    │
    ├─ Validate state
    ├─ Exchange code for token
    ├─ Encrypt and set obot_access_token cookie
    │
    ▼
Redirect to original URL (rd parameter)
    │
    ▼
Subsequent Requests
    │
    ▼
/obot-get-state
    │
    ├─ Decrypt obot_access_token cookie
    ├─ Return user profile (user, email, preferredUsername)
    │
    ▼
Authenticated Session
```

### Knowledge Ingestion Flow

```
File/Directory Input
    │
    ▼
Client (.knowignore filtering)
    │
    ▼
Document Loader
    │
    ├─ Detect file type
    ├─ Convert (Office → PDF → Text)
    ├─ OCR (if needed)
    │
    ▼
Transformers
    │
    ├─ Extract metadata
    ├─ Extract keywords
    ├─ Markdown processing
    │
    ▼
Text Splitter
    │
    ├─ Chunk documents
    ├─ Maintain metadata
    │
    ▼
Embeddings Generation
    │
    ├─ Call embedding provider
    ├─ Batch processing
    │
    ▼
Vector Store + Index
    │
    ├─ Store vectors (pgvector/sqlite-vec)
    ├─ Store metadata (postgres/sqlite)
    ├─ Update file tracking
```

### Knowledge Retrieval Flow

```
Query Input
    │
    ▼
Query Modifiers (optional)
    │
    ├─ Spell check
    ├─ LLM enhancement
    │
    ▼
Retriever(s)
    │
    ├─ Vector retrieval
    ├─ Keyword (BM25) retrieval
    ├─ Hybrid retrieval
    │
    ▼
Post-processors
    │
    ├─ Deduplication
    ├─ Similarity filtering
    ├─ Reranking (BM25/Cohere)
    ├─ Top-K reduction
    │
    ▼
Ranked Results
```

## Technology Choices and Rationale

### Go for Core Infrastructure

- **Rationale**: Performance, concurrency, static typing
- **Use cases**: Model providers, auth providers, system tools
- **Benefits**: Single binary deployment, cross-compilation

### Python for Specialized Tools

- **Rationale**: Rich ML/AI ecosystem
- **Use cases**: Embedding providers, document processing
- **Benefits**: Fast prototyping, library availability

### TypeScript for Image Tools

- **Rationale**: Node.js ecosystem for image processing
- **Use cases**: Image generation and analysis
- **Benefits**: npm package ecosystem

### SQLite and PostgreSQL

- **Rationale**: Embeddable (SQLite) and scalable (PostgreSQL)
- **Use cases**: Vector storage, metadata indexing, credential storage
- **Benefits**: SQL compatibility, pgvector/sqlite-vec extensions

### FastAPI for Python Services

- **Rationale**: Modern async Python web framework
- **Use cases**: Python-based model providers
- **Benefits**: Automatic OpenAPI docs, type validation, async support

### GPTScript as Tool Runtime

- **Rationale**: Obot platform's tool execution engine
- **Use cases**: All tool definitions and invocations
- **Benefits**: Standardized tool protocol, credential management

## Security Considerations

### Credential Security

- All credentials encrypted at rest (AES-256)
- Environment variable isolation
- Per-request credential override support
- Placeholder credentials for development only

### Authentication

- OAuth2 standard compliance
- Encrypted cookie storage
- CSRF protection via state parameter
- Configurable user validation (email domains, teams, orgs)

### Network Security

- HTTPS enforcement for auth cookies
- Reverse proxy isolation from upstream APIs
- No credential logging or exposure

### Dependency Management

- Regular Dependabot updates
- Go module verification
- Python requirements.txt with locked versions

## Observability

### Logging

- Debug mode via `GPTSCRIPT_DEBUG=true`
- Structured logging in Go components
- HTTP request/response logging (debug mode)

### Monitoring

- Health check endpoints in model providers
- Validation endpoints for credential testing

### Error Handling

- Consistent error formats
- HTTP status code conventions
- Detailed error messages in debug mode

## Extensibility

### Adding New Model Providers

1. Create provider directory with `tool.gpt`
2. Implement HTTP server (port 8000)
3. Support global and per-request API keys
4. Implement `/v1/models` endpoint
5. Add validation mode
6. Register in `index.yaml`

### Adding New Tools

1. Create tool directory with `tool.gpt`
2. Implement business logic in Go/Python
3. Build to `bin/gptscript-go-tool`
4. Add tests
5. Register in `index.yaml`

### Adding New Auth Providers

1. Follow Google auth provider reference implementation
2. Use `auth-providers-common` utilities
3. Implement required OAuth2 endpoints
4. Reference `placeholder-credential`
5. Add envVars metadata
6. Register in `index.yaml`

## Performance Characteristics

### Model Providers

- Proxy overhead: < 10ms typically
- Horizontal scaling: Stateless design
- Concurrent requests: Limited by upstream API

### Knowledge Base

- Ingestion: Parallel document processing
- Retrieval: Vector similarity O(n) with ANN optimizations
- Hybrid retrieval: Combines vector + BM25 efficiently

### Credential Stores

- SQLite: Embedded, single-file, microsecond latency
- PostgreSQL: Client-server, connection pooling, millisecond latency

## Future Architecture Considerations

### Potential Enhancements

- Distributed caching for model provider responses
- Advanced vector search (HNSW, IVF indices)
- Multi-tenant credential isolation
- Event-driven architecture for async tool execution
- gRPC support for high-performance tool communication

### Scalability Patterns

- Load balancing for model providers
- Sharding for knowledge base
- Read replicas for credential stores
- Edge deployment for auth providers
