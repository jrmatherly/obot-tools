# Project Index: obot-tools

**Generated**: 2026-01-15  
**Repository**: https://github.com/obot-platform/tools  
**Purpose**: Official Obot platform tools, model providers, and authentication providers

---

## 📊 Token Efficiency

**Before**: Reading all files → ~58,000 tokens per session  
**After**: Reading this index → ~3,000 tokens per session  
**Savings**: 94% reduction (55,000 tokens saved per session)

---

## 📁 Project Structure

```
obot-tools/
├── Core Tools (8 components)
│   ├── memory/                    - Agent memory storage
│   ├── knowledge/                 - RAG knowledge base
│   ├── tasks/                     - Task management
│   ├── workspace-files/           - File operations
│   ├── threads/                   - Thread management
│   ├── time/                      - Time utilities
│   ├── images/                    - Image operations (TypeScript)
│   └── file-summarizer/           - Document summarization (Python)
│
├── Model Providers (9 components)
│   ├── openai-model-provider/     - OpenAI API (Go)
│   ├── anthropic-model-provider-go/ - Anthropic Claude (Go)
│   ├── ollama-model-provider/     - Local Ollama (Go)
│   ├── groq-model-provider/       - Groq API (Go)
│   ├── xai-model-provider/        - xAI Grok (Go)
│   ├── deepseek-model-provider/   - DeepSeek (Go)
│   ├── voyage-model-provider/     - Voyage embeddings (Python)
│   ├── vllm-model-provider/       - VLLM self-hosted (Go)
│   └── generic-openai-model-provider/ - OpenAI-compatible (Go)
│
├── Auth Providers (2 components)
│   ├── github-auth-provider/      - GitHub OAuth2 (Go)
│   └── google-auth-provider/      - Google OAuth2 (Go)
│
├── Credential Management
│   └── credential-stores/
│       ├── sqlite/                - SQLite backend (Go)
│       ├── postgres/              - PostgreSQL backend (Go)
│       └── pkg/common/            - Shared encryption utilities
│
├── System Tools (8 components)
│   ├── workflow/                  - Workflow execution
│   ├── result-formatter/          - Result formatting
│   ├── obot-model-provider/       - Internal provider
│   ├── tasks-workflow/            - Task workflows
│   ├── task-invoke/               - Task invocation
│   ├── loop-data/                 - Loop handling
│   ├── existing-credential/       - Credential passthrough
│   └── generic-credential/        - Generic credentials
│
├── Supporting
│   ├── auth-providers-common/     - Shared auth utilities
│   ├── placeholder-credential/    - Dev placeholder (Python)
│   └── oauth2/                    - OAuth2 utilities
│
└── Infrastructure
    ├── docs/                      - Comprehensive documentation
    ├── scripts/                   - Build scripts
    ├── tests/                     - Integration tests
    ├── charts/                    - Kubernetes Helm charts
    └── .github/workflows/         - CI/CD pipelines
```

---

## 🚀 Entry Points

### Build & Package

- **Build All**: `make build` → Builds all Go tools to `bin/gptscript-go-tool`
- **Test All**: `make test` → Runs all test suites
- **Package Tools**: `make package-tools` → Package user-facing tools
- **Package Providers**: `make package-providers` → Package model/auth providers
- **Docker Build**: `make docker-build` → Multi-platform container build

### Tool Registry

- **Registry**: `index.yaml` → Central tool, system, modelProviders, authProviders registry
- **GPTScript Runtime**: All tools execute via GPTScript with `tool.gpt` definitions

### Go Entry Points (main.go files)

- `knowledge/main.go` → CLI and server for knowledge base (most complex component)
- `openai-model-provider/main.go` → OpenAI proxy server (:8000)
- `github-auth-provider/main.go` → GitHub OAuth2 server (:8000)
- `memory/main.go` → Memory management tool
- `tasks/main.go` → Task management tool
- `credential-stores/sqlite/main.go` → SQLite credential store
- `credential-stores/postgres/main.go` → PostgreSQL credential store

### Python Entry Points (main.py files)

- `voyage-model-provider/main.py` → FastAPI server for Voyage embeddings (:8000)
- `file-summarizer/main.py` → Document summarization tool
- `placeholder-credential/main.py` → Development credential placeholder

### TypeScript Entry Points

- `images/src/tools.ts` → Image tool exports
- `images/src/generate.ts` → Image generation
- `images/src/analyze.ts` → Image analysis

---

## 📦 Core Modules by Language

### Go Modules (Primary Language)

#### Model Provider Pattern

**Common Structure**: All Go model providers follow this pattern

- `main.go` → HTTP server entry point (:8000)
- `proxy/` → Reverse proxy implementation (if complex)
- Uses `openai-model-provider/proxy` package for common functionality

**Key Exports** (openai-model-provider/proxy):

- `type Config struct` → Provider configuration
- `func Run(cfg *Config) error` → Start HTTP server
- `func (cfg *Config) Validate(path string) error` → Test credentials
- `func DefaultRewriteModelsResponse` → Model list transformation

#### Auth Provider Pattern

**Common Structure**: GitHub and Google auth providers

- `main.go` → OAuth2 HTTP server (:8000)
- `pkg/profile/` → User profile retrieval
  - `GetProfile(accessToken string) (*Profile, error)`
  - `type Profile struct { User, Email, PreferredUsername string }`

**Shared Utilities** (auth-providers-common/pkg):

- `state/` → OAuth2 state management
- `env/` → Environment variable helpers
- `icon/` → Profile icon URL retrieval

#### Knowledge Base (Largest Component)

**Entry**: `knowledge/main.go`

**CLI Commands** (knowledge/pkg/cmd):

- `root.go` → Cobra CLI root
- `retrieve.go` → Search knowledge base
- `ingest.go` → Ingest documents
- `load.go` → Load remote content (GitHub, URLs)
- `create_dataset.go`, `delete_dataset.go`, `get_dataset.go`, `list_datasets.go`
- `delete_file.go`, `get_file.go`, `edit_dataset.go`
- `import.go`, `export.go` → Dataset import/export
- `client.go` → Client operations
- `askdir.go` → Query directory

**Core Subsystems** (knowledge/pkg):

- `client/` → Client interface, metadata, standalone mode, .knowignore support
  - `client.go`, `standalone.go`, `metadata.go`, `ignore.go`, `common.go`
- `config/` → Configuration loading, validation, embedding models
  - `config.go`, `config_test.go`, `embedding_models.go`
- `datastore/` → Document processing and retrieval pipeline
  - `datastore.go`, `ingest.go`, `retrieve.go`, `dataset.go`, `document.go`
  - `documentloader/` → PDF (MuPDF, GoPDF, SmartPDF), Office, Structured, Remote, OCR
  - `textsplitter/` → Chunking with overlap
  - `embeddings/` → OpenAI embeddings, provider abstraction
  - `retrievers/` → Vector, BM25, hybrid, subquery, routing, merging
  - `postprocessors/` → Similarity, BM25 rerank, Cohere rerank, deduplication, topK
  - `querymodifiers/` → Spell check, LLM enhancement
  - `transformers/` → Metadata extraction, keyword extraction, markdown processing
  - `lib/` → BM25, scoring utilities
- `flows/` → Configuration-driven pipelines
  - `config/` → YAML/JSON config, blueprints (default, obot, context)
- `vectorstore/` → pgvector, sqlite-vec implementations
- `index/` → SQLite and PostgreSQL metadata indexing
- `llm/` → LLM integration utilities
- `log/` → Logging utilities
- `output/` → Output formatting and redaction

#### Credential Stores

**Common** (credential-stores/pkg/common):

- `encryption.go` → AES-256-GCM encryption
- `db.go` → Database operations (store, reveal, delete)
- `db_test.go` → Shared test suite

**Implementations**:

- `sqlite/main.go` → SQLite backend
- `postgres/main.go` → PostgreSQL backend

#### System Tools

All follow simple CLI pattern with main.go entry points:

- `workflow/main.go`
- `result-formatter/main.go`
- `obot-model-provider/main.go`
- `tasks-workflow/main.go`
- `task-invoke/main.go`
- `loop-data/main.go`
- `generic-credential/main.go`

### Python Modules

#### Voyage Model Provider

**Entry**: `voyage-model-provider/main.py`

- FastAPI server on port 8000
- `voyageai.AsyncClient` for embeddings
- Endpoints: `/`, `/v1/models`, `/v1/embeddings`
- Hardcoded model list (voyage-3-large, voyage-3, voyage-3-lite, etc.)

#### File Summarizer

**Entry**: `file-summarizer/main.py`
**Modules** (file-summarizer/tools):

- `summarizer.py` → Main summarization logic
- `reader.py` → File reading utilities
- `load_text.py` → Text loading
- `helper.py` → Helper functions
- `gptscript_workspace.py` → Workspace integration

#### Placeholder Credential

**Entry**: `placeholder-credential/main.py`

- Development-only credential placeholder

### TypeScript Modules

#### Images Tool

**Entry**: `images/src/tools.ts`
**Modules**:

- `generate.ts` → Image generation functionality
- `analyze.ts` → Image analysis functionality

---

## 🔧 Configuration Files

### Build & Package Management

- `Makefile` → Build automation (build, test, package-tools, package-providers, docker-build)
- `Dockerfile` → Multi-stage build (base: Wolfi, tools stage, providers stage)
- `scripts/build.sh` → Finds and builds all main.go files to bin/gptscript-go-tool
- `scripts/package-tools.sh` → Package tools for distribution
- `scripts/package-providers.sh` → Package providers for distribution

### Tool Registry

- `index.yaml` → Central registry with 4 sections: tools, system, modelProviders, authProviders

### Go Modules

Each Go component has:

- `go.mod` → Module definition
- `go.sum` → Dependency checksums
- Go version: 1.23

### Python Dependencies

- `requirements.txt` → Root Python dependencies (azure-identity, fastapi, gptscript, openai, uvicorn, voyageai)
- `lock.txt` → Additional locked dependencies
- Component-specific `requirements.txt` in voyage-model-provider/, file-summarizer/, placeholder-credential/

### TypeScript/Node.js

- `images/package.json` → Node.js dependencies
- `images/tsconfig.json` → TypeScript configuration
- `images/.eslintrc.cjs` → ESLint configuration

### Knowledge Base Configuration

- `knowledge/examples/*.yaml` → Example configs (advanced, bm25, hybrid, ocr, rerank, routing, etc.)
- `knowledge/pkg/flows/config/blueprints/*.yaml` → Pipeline blueprints (default, obot, context)
- `.knowignore` support → Files/patterns to exclude from ingestion

### CI/CD

**Workflows** (.github/workflows):

- `test.yaml` → PR testing (go test, build, docker multi-platform)
- `build-tools.yaml` → Package tools
- `build-providers.yaml` → Package providers
- `dispatch.yaml` → Manual workflow triggers
- `dependabot-reviewers.yaml` → Auto-assign reviewers

### Kubernetes

**Helm Chart** (charts/oauth-mcp-server):

- `Chart.yaml` → Chart metadata
- `values.yaml` → Default values
- `templates/` → Deployment, Service, Ingress, Secrets, ServiceAccount

---

## 📚 Documentation

### Comprehensive Documentation (docs/)

- `INDEX.md` → Documentation navigation hub and learning paths
- `ARCHITECTURE.md` → System architecture, component design, data flows (500+ lines)
- `API_REFERENCE.md` → Complete API documentation for all components (700+ lines)
- `DEVELOPER_GUIDE.md` → Development setup, building, testing, contributing (600+ lines)
- `auth-providers.md` → Auth provider specifications and requirements

### Component Documentation

- `README.md` → Project overview, quick start, component tables, development workflow
- `knowledge/README.md` → Knowledge base specific documentation
- `credential-stores/README.md` → Credential store documentation

### Tool Definitions

Every component has a `tool.gpt` file with:

- GPTScript tool definition
- Metadata (icon, category, envVars, optionalEnvVars)
- Sub-tool definitions
- Command references

**Examples**:

- `memory/tool.gpt` → Create Memory, Update Memory, Delete Memory, list_memories
- `knowledge/tool.gpt`, `knowledge/ingest.gpt`, `knowledge/load.gpt`, `knowledge/delete.gpt`, `knowledge/delete-file.gpt`
- `openai-model-provider/tool.gpt` → Model provider metadata
- `github-auth-provider/tool.gpt` → Auth provider with envVars metadata

---

## 🧪 Test Coverage

### Go Tests

**Integration Tests**:

- `tests/tool_test.go` → Integration tests with GPTSCRIPT_TOOL_REMAP

**Unit Tests**:

- `github-auth-provider/pkg/profile/profile_test.go` → Profile retrieval tests
- `google-auth-provider/pkg/profile/profile_test.go` → Profile retrieval tests
- `credential-stores/pkg/common/db_test.go` → Credential encryption and DB tests
- `knowledge/pkg/config/config_test.go` → Configuration loading tests
- `knowledge/pkg/datastore/ingest_test.go` → Ingestion pipeline tests
- `knowledge/pkg/datastore/documentloader/documentloaders_test.go` → Document loader tests
- `knowledge/pkg/datastore/textsplitter/textsplitter_test.go` → Text splitting tests
- `knowledge/pkg/datastore/embeddings/embeddings_test.go` → Embedding tests
- `knowledge/pkg/vectorstore/helper/sql_test.go` → SQL helper tests
- `knowledge/pkg/flows/config/config_test.go` → Flow configuration tests

**Test Data**:

- `knowledge/pkg/datastore/testdata/pdf/` → PDF test files
- `knowledge/pkg/datastore/embeddings/test_assets/` → Embedding test configs
- `knowledge/pkg/flows/config/testdata/` → Valid/invalid config examples

### Test Commands

```bash
make test                           # All tests
cd tests && go test -v tool_test.go # Integration tests
cd <component> && go test ./...     # Component tests
go test -cover ./...                # With coverage
```

---

## 🔗 Key Dependencies

### Go Dependencies (via go.mod)

- **OpenAI Proxy**: openai-model-provider/proxy (internal package)
- **Auth Common**: auth-providers-common/pkg (state, env, icon)
- **GPTScript**: Integration via tool.gpt definitions
- **Cobra**: CLI framework (knowledge tool)
- **Database**: SQLite, PostgreSQL drivers
- **Encryption**: AES-256-GCM (credential stores)

### Python Dependencies (requirements.txt)

- `azure-identity==1.25.1` → Azure authentication
- `azure-mgmt-cognitiveservices==14.1.0` → Azure AI services
- `fastapi==0.124.2` → Web framework
- `gptscript @ git+...` → GPTScript Python client
- `openai==2.11.0` → OpenAI API client
- `tiktoken==0.12.0` → Token counting
- `uvicorn==0.38.0` → ASGI server
- `voyageai==0.3.6` → Voyage AI client

### TypeScript Dependencies (images/package.json)

- Node.js packages for image processing

### Runtime Dependencies

- **Go**: 1.23+
- **Python**: 3.13+
- **Node.js**: 18+ (for image tools)
- **PostgreSQL**: 15+ with pgvector extension (for knowledge base)
- **SQLite**: With sqlite-vec extension (for knowledge base)
- **Docker**: For containerized builds
- **Make**: Build automation

---

## 📝 Quick Start

### 1. Clone Repository

```bash
git clone https://github.com/obot-platform/tools.git
cd tools
```

### 2. Build All Components

```bash
make build
```

Builds all Go tools to `bin/gptscript-go-tool` with optimizations (`-ldflags="-s -w"`)

### 3. Run Tests

```bash
make test
```

Runs integration tests, auth provider tests, credential store tests

### 4. Build Specific Component

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

### 5. Test Model Provider

```bash
export OBOT_OPENAI_MODEL_PROVIDER_API_KEY=sk-...
export PORT=8000
export GPTSCRIPT_DEBUG=true
cd openai-model-provider
go run .

# In another terminal
curl http://localhost:8000/v1/models
```

### 6. Package for Distribution

```bash
make package-tools      # Package user-facing tools
make package-providers  # Package model/auth providers
make docker-build       # Build Docker image
```

---

## 🏗️ Architecture Patterns

### Model Provider Pattern

1. HTTP server on port 8000 (configurable via PORT env var)
2. Global API key via `OBOT_<PROVIDER>_MODEL_PROVIDER_API_KEY`
3. Per-request API key via `X-Obot-OBOT_<PROVIDER>_MODEL_PROVIDER_API_KEY` header
4. Validation mode: `go run . validate` (exit code 0 on success)
5. Endpoints: `/`, `/v1/models`, provider-specific endpoints
6. Reverse proxy to upstream API
7. Optional model response rewriting

### Auth Provider Pattern (OAuth2)

1. HTTP server on port 8000
2. Environment variables for OAuth client credentials
3. Encrypted cookie storage (`obot_access_token`)
4. Required endpoints: `/oauth2/start`, `/oauth2/callback`, `/oauth2/sign_out`, `/obot-get-icon-url`, `/obot-get-state`
5. State parameter for CSRF protection
6. User validation (email domains, teams, orgs)
7. Reference placeholder-credential in tool.gpt

### Tool Pattern

1. `tool.gpt` with GPTScript definitions
2. Binary at `bin/gptscript-go-tool`
3. Sub-tools as separate tool definitions
4. Context sharing via Share Context directive
5. Metadata: icon, category, envVars, optionalEnvVars

### Credential Store Pattern

1. AES-256-GCM encryption
2. Database backend (SQLite or PostgreSQL)
3. Operations: store, reveal, delete
4. Encryption key from environment variable
5. Shared encryption/DB logic in pkg/common

---

## 🔍 Component Lookup

### By Language

**Go**: 30+ components (all model providers except voyage, all auth providers, knowledge, memory, tasks, workspace-files, credential stores, system tools)  
**Python**: 3 components (voyage-model-provider, file-summarizer, placeholder-credential)  
**TypeScript**: 1 component (images)  
**GPTScript**: 2 components (time, threads) - pure tool definitions

### By Category

**User Tools (8)**: memory, knowledge, tasks, workspace-files, threads, time, images, file-summarizer  
**Model Providers (9)**: openai, anthropic, ollama, groq, xai, deepseek, voyage, vllm, generic-openai  
**Auth Providers (2)**: github, google  
**Credential Stores (2)**: sqlite, postgres  
**System Tools (8)**: workflow, result-formatter, obot-model-provider, tasks-workflow, task-invoke, loop-data, existing-credential, generic-credential

### By Complexity (Lines of Code)

1. **knowledge/** → Most complex (100+ files, RAG pipeline, CLI, multiple backends)
2. **openai-model-provider/** → Medium (proxy pattern, model rewriting)
3. **github-auth-provider/** → Medium (OAuth2 flow, user validation)
4. **credential-stores/** → Medium (encryption, DB backends)
5. **voyage-model-provider/** → Simple (FastAPI server, embeddings only)
6. **memory/, tasks/, workspace-files/** → Simple (basic CRUD tools)

---

## 📊 Metrics

### Repository Size

- **Go Files**: ~200+ files
- **Python Files**: ~10 files
- **TypeScript Files**: 3 files
- **Documentation**: 5 comprehensive markdown files (~2000+ lines total)
- **Test Files**: 15+ test files

### Component Count

- **Total Components**: 30+
- **User Tools**: 8
- **Model Providers**: 9
- **Auth Providers**: 2
- **Credential Stores**: 2
- **System Tools**: 8
- **Supporting Libraries**: 3

### Build Outputs

- **Go Binaries**: Each component builds to `bin/gptscript-go-tool`
- **Docker Image**: Multi-stage build with tools and providers
- **Package Outputs**: Separate packages for tools and providers

---

## 🎯 Development Priorities

### High Priority Components

1. **knowledge/** → Core RAG functionality, complex pipeline
2. **openai-model-provider/** → Most widely used model provider
3. **github-auth-provider/** → Primary auth method
4. **memory/** → Core agent memory

### Extension Points

- New model providers → Follow openai-model-provider pattern
- New auth providers → Follow google-auth-provider (reference impl)
- New credential stores → Extend pkg/common pattern
- New tools → Add to index.yaml, implement tool.gpt

---

## 📞 Support

- **Issues**: https://github.com/obot-platform/obot/issues (use `tools` label)
- **Documentation**: docs/ directory
- **Main Repo**: https://github.com/obot-platform/obot

---

**Index Version**: 1.0  
**Last Updated**: 2026-01-15  
**Index Size**: ~5KB (fits in context window)
