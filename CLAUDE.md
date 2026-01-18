# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the **obot-tools** repository - the official home of Obot platform tools, model providers, and authentication providers. The project uses GPTScript as its tool execution runtime and follows a modular plugin architecture where each component is independently deployable.

## Build & Test Commands

### Building

```bash
# Build all Go components
make build

# Build specific component
cd <component-name>
go build -ldflags="-s -w" -o bin/gptscript-go-tool .

# Package for distribution
make package-tools      # User-facing tools
make package-providers  # Model and auth providers

# Docker build
make docker-build
```

### Testing

```bash
# Run all tests
make test

# Test specific component
cd <component-name>
go test ./...

# Integration tests with tool remapping
cd tests
GPTSCRIPT_TOOL_REMAP="github.com/obot-platform/tools=.." go test -v tool_test.go

# Run with coverage
go test -cover ./...
```

### Running Components Locally

**Model Provider** (default port 8000):

```bash
export OBOT_<PROVIDER>_MODEL_PROVIDER_API_KEY=...
export PORT=8000
export GPTSCRIPT_DEBUG=true  # Enable debug logging
cd <provider-name>
go run .                      # Start server
go run . validate             # Test credentials
```

**Auth Provider** (default port 8000):

```bash
export OBOT_<PROVIDER>_AUTH_PROVIDER_CLIENT_ID=...
export OBOT_<PROVIDER>_AUTH_PROVIDER_CLIENT_SECRET=...
export OBOT_AUTH_PROVIDER_COOKIE_SECRET=$(openssl rand -base64 32)
export OBOT_AUTH_PROVIDER_EMAIL_DOMAINS=example.com
cd <auth-provider>
go run .
```

**Python Components**:

```bash
cd voyage-model-provider
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python main.py
```

## Architecture & Design Patterns

### Component Organization

The repository is organized by component type:

- **Core Tools** (8): `memory/`, `knowledge/`, `tasks/`, `workspace-files/`, `threads/`, `time/`, `images/`, `file-summarizer/`
- **Model Providers** (9): `*-model-provider/` directories
- **Auth Providers** (2): `github-auth-provider/`, `google-auth-provider/`
- **Credential Stores** (2): `credential-stores/sqlite/`, `credential-stores/postgres/`
- **System Tools** (8): `workflow/`, `result-formatter/`, `obot-model-provider/`, etc.
- **Supporting**: `auth-providers-common/`, `placeholder-credential/`, `oauth2/`

**Central Registry**: All components register in `index.yaml` under sections: `tools`, `system`, `modelProviders`, `authProviders`.

### Model Provider Pattern

All model providers follow a consistent HTTP server pattern:

**Structure**:

```
<provider>-model-provider/
├── main.go              # HTTP server entry point
├── proxy/               # Reverse proxy implementation (if needed)
├── tool.gpt             # GPTScript tool definition
└── go.mod
```

**Key Characteristics**:

- Listen on port 8000 (configurable via `PORT` env var)
- Global API key: `OBOT_<PROVIDER>_MODEL_PROVIDER_API_KEY`
- Per-request API key: `X-Obot-OBOT_<PROVIDER>_MODEL_PROVIDER_API_KEY` header
- Validation mode: `go run . validate` (exit 0 on success)
- Standard endpoints: `/`, `/v1/models`, provider-specific endpoints
- Reverse proxy to upstream API with optional response rewriting

**Common Proxy Package**: Most Go model providers use `openai-model-provider/proxy` for shared functionality (`Config`, `Run()`, `Validate()`, `DefaultRewriteModelsResponse`).

### Auth Provider Pattern (OAuth2)

**Structure**:

```
<provider>-auth-provider/
├── main.go              # OAuth2 HTTP server
├── pkg/profile/         # User profile retrieval
├── tool.gpt             # With envVars metadata
└── go.mod
```

**Key Characteristics**:

- OAuth2 authorization code flow
- Encrypted cookie storage (`obot_access_token`) using `OBOT_AUTH_PROVIDER_COOKIE_SECRET`
- Required endpoints: `/oauth2/start`, `/oauth2/callback`, `/oauth2/sign_out`, `/obot-get-icon-url`, `/obot-get-state`
- User validation (email domains, teams, orgs, repos)
- Reference implementation: `google-auth-provider/` (follow this pattern for new providers)
- Shared utilities in `auth-providers-common/pkg`: state management, env helpers, icon retrieval

### GPTScript Tool Definition (.gpt files)

**Every component has a `tool.gpt` file** that defines the tool for GPTScript runtime.

**Structure**:

```gptscript
Name: Tool Name
Description: Tool description
Metadata: icon: https://...
Metadata: category: Capability
Metadata: envVars: REQUIRED_VAR1,REQUIRED_VAR2
Metadata: optionalEnvVars: OPTIONAL_VAR1

---
Name: Sub Tool
Description: Sub-tool action
Param: input: Parameter description

#!${GPTSCRIPT_TOOL_DIR}/bin/gptscript-go-tool command
```

**Key Directives**:

- `Name`: Tool name (user-facing)
- `Description`: Critical - used by LLM to decide when to call tool
- `Param`: Parameters (format: `param-name: description`)
- `Metadata: envVars`: Required environment variables (comma-separated)
- `Metadata: optionalEnvVars`: Optional environment variables
- `Credential`: Reference to credential file (e.g., `../placeholder-credential as ...`)
- `Share Context`: Context shared with other tools
- `Share Tools`: Tools shared by this tool
- `Type: context`: Marks tool as context-only (auto-called for LLM prompting)

**Tool Separation**: Use `---` to separate multiple tool definitions in one file.

**Parameter Access**: Parameters become uppercase environment variables (e.g., `param: color` → `$COLOR`).

## GPTScript Tool Format - Complete Reference

For comprehensive GPTScript .gpt file format documentation, including:

- Complete directive reference (standard + obot-tools-specific)
- Model Provider and Auth Provider patterns
- Context tool patterns
- Provider metadata structure (providerMeta)
- Credential syntax and mapping
- Command types (#!sys.daemon, #!sys.echo, shell scripts)
- Parameter handling (become UPPERCASE env vars)
- Tool bundles and sharing
- Registration in index.yaml

**Load the Serena memory**: `gptscript_tool_format`

**Study canonical examples in this repo**:

- `memory/tool.gpt` - Context tools, Share Tools/Context, multiple sub-tools
- `openai-model-provider/tool.gpt` - Model provider, #!sys.daemon, providerMeta
- `github-auth-provider/tool.gpt` - Auth provider, envVars/optionalEnvVars
- `knowledge/tool.gpt` - Complex tool with multiple features, output filter

**Auth Provider Requirements**: See `docs/auth-providers.md` for complete specification including:

- OAuth2 flow implementation
- Required endpoints and API contracts
- Cookie handling and encryption
- Token refresh mechanism
- JSON schemas for request/response
- Reference implementation: `google-auth-provider/`

**Official Documentation**: https://docs.gptscript.ai

### Knowledge Base Architecture (Most Complex Component)

The knowledge tool (`knowledge/`) is the largest component with a sophisticated RAG pipeline:

**Key Subsystems**:

1. **Client** (`pkg/client/`): Standalone mode, metadata, `.knowignore` support
2. **Datastore** (`pkg/datastore/`): Core ingestion and retrieval pipeline
   - **Document Loaders**: PDF (MuPDF, GoPDF, SmartPDF with OCR), Office, Structured data, Remote (GitHub, URLs)
   - **Text Splitter**: Configurable chunking with overlap
   - **Embeddings**: OpenAI embeddings, provider abstraction
   - **Retrievers**: Vector similarity, BM25 keyword, hybrid, subquery, routing, merging
   - **Post-processors**: Similarity filtering, BM25/Cohere reranking, deduplication, top-K
   - **Query Modifiers**: Spell check, LLM enhancement
   - **Transformers**: Metadata extraction, keyword extraction, markdown processing
3. **Vector Stores** (`pkg/vectorstore/`): pgvector (PostgreSQL), sqlite-vec
4. **Index Storage** (`pkg/index/`): SQLite and PostgreSQL metadata indexing
5. **Flows** (`pkg/flows/`): YAML/JSON configuration-driven pipelines with blueprints

**CLI Commands** (`pkg/cmd/`): `retrieve`, `ingest`, `load`, `create-dataset`, `delete-dataset`, `get-dataset`, `list-datasets`, `delete-file`, `get-file`, `import`, `export`, `client`, `askdir`

**Configuration**: Supports YAML/JSON with blueprints in `pkg/flows/config/blueprints/` (default, obot, context).

### Credential Store Pattern

**Structure**:

```
credential-stores/
├── pkg/common/          # Shared encryption and DB logic
│   ├── encryption.go   # AES-256-GCM
│   ├── db.go           # Store, reveal, delete operations
│   └── db_test.go
├── sqlite/main.go      # SQLite backend
└── postgres/main.go    # PostgreSQL backend
```

**Operations**: `store`, `reveal`, `delete` (defined in tool.gpt).

**Security**: AES-256-GCM encryption with key from `OBOT_CREDENTIAL_STORE_ENCRYPTION_KEY`.

## Code Conventions

### Go Code

**Formatting**:

- Use `go fmt ./...` before committing
- Tabs for indentation (Go standard)
- Build with optimizations: `go build -ldflags="-s -w"`

**Naming**:

- Exported: PascalCase (`Run`, `Config`)
- Unexported: camelCase (`apiKey`, `port`)
- Environment variables: `OBOT_` prefix in SCREAMING_SNAKE_CASE

**Patterns**:

- Return errors, don't panic (except in main for fatal errors)
- Use `fmt.Println` for informational messages to stdout
- Default values for optional config (e.g., `PORT` defaults to `"8000"`)

### Python Code

**Formatting**:

- Follow PEP 8
- Use `black` and `isort` for formatting
- Type hints where applicable

**FastAPI Patterns** (for model providers):

- Async/await for async operations
- Middleware for logging/debugging (check `GPTSCRIPT_DEBUG`)
- JSONResponse for API endpoints
- HTTPException for errors

### Environment Variables

**Naming Convention**:

- Model providers: `OBOT_<PROVIDER>_MODEL_PROVIDER_API_KEY`
- Auth providers: `OBOT_<PROVIDER>_AUTH_PROVIDER_*`
- Common: `PORT`, `GPTSCRIPT_DEBUG`, `OBOT_SERVER_PUBLIC_URL`, `OBOT_AUTH_PROVIDER_COOKIE_SECRET`

### Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `chore:` - Maintenance (deps, build)

## Creating New Components

### New Model Provider

1. Create directory: `<provider>-model-provider/`
2. Initialize Go module: `go mod init github.com/obot-platform/tools/<provider>-model-provider`
3. Create `main.go` using `openai-model-provider/proxy` package for common functionality
4. Create `tool.gpt` with metadata including `envVars`
5. Register in `index.yaml` under `modelProviders:`
6. Test with `go run . validate` and `go run .`

### New Tool

1. Create directory: `<tool-name>/`
2. Create `tool.gpt` with tool definitions (use `---` to separate sub-tools)
3. Implement business logic in Go/Python
4. Build to `bin/gptscript-go-tool`
5. Register in `index.yaml` under `tools:` or `system:`

### New Auth Provider

**Follow the reference implementation**: `google-auth-provider/`

1. Use `auth-providers-common/pkg` utilities
2. Implement required OAuth2 endpoints (see docs/auth-providers.md)
3. Reference `placeholder-credential` in tool.gpt
4. Add `envVars` metadata
5. Implement user validation logic
6. Register in `index.yaml` under `authProviders:`

## Important Files & Locations

- `index.yaml` - Central tool registry (4 sections: tools, system, modelProviders, authProviders)
- `Makefile` - Build automation (build, test, package-tools, package-providers, docker-build)
- `scripts/build.sh` - Finds all `main.go` files and builds to `bin/gptscript-go-tool`
- `Dockerfile` - Multi-stage build (base → tools → providers)
- `.github/workflows/` - CI/CD pipelines (test.yaml, build-tools.yaml, build-providers.yaml)
- `docs/` - Comprehensive documentation (ARCHITECTURE.md, API_REFERENCE.md, DEVELOPER_GUIDE.md)
- `PROJECT_INDEX.md` - Quick reference index of entire codebase structure

## Debug & Development

**Enable Debug Logging**:

```bash
export GPTSCRIPT_DEBUG=true
```

**Validation Mode** (for providers):

```bash
cd <provider>
go run . validate  # Exit code 0 on success, 1 on failure
```

**Test with curl**:

```bash
# List models
curl http://localhost:8000/v1/models

# Test endpoint with per-request API key
curl http://localhost:8000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "X-Obot-OBOT_OPENAI_MODEL_PROVIDER_API_KEY: sk-..." \
  -d '{"model":"gpt-4","messages":[{"role":"user","content":"Hello"}]}'
```

## Dependencies & Versions

- **Go**: 1.23+ (all `go.mod` files specify this)
- **Python**: 3.13+ (FastAPI, Uvicorn, voyageai, openai)
- **Node.js**: 18+ (for images tool only)
- **PostgreSQL**: 15+ with pgvector extension (for knowledge base)
- **SQLite**: With sqlite-vec extension (for knowledge base)

## Issues & Support

- File issues in [main Obot repo](https://github.com/obot-platform/obot/issues) with `tools` label
- See comprehensive docs in `docs/` directory for architecture, API reference, and development guide

## Workspace Integration

This project is part of the AI workspace. Additional resources:

- **Claude Code commands**: `AI/.claude/commands/` (expert-mode, etc.)
- **Shared agents**: `AI/.claude/agents/` (explore, security-audit, etc.)
- **SuperClaude skills**: `/sc:analyze`, `/sc:test`, `/sc:git`, etc.
- **Serena memories**: `AI/.serena/memories/` (task_completion_checklist, gptscript_tool_format, etc.)
- **GitHub Actions**: Workspace-level PR review and issue triage

For session initialization with full context, run `/expert-mode` from the workspace root.

**GPTScript Development**: Load the `gptscript_tool_format` Serena memory when working on `.gpt` files for comprehensive format reference.
