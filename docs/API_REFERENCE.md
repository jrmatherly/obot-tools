# Obot Tools - API Reference

## Table of Contents

- [Model Providers](#model-providers)
  - [Common Model Provider API](#common-model-provider-api)
  - [OpenAI Model Provider](#openai-model-provider)
  - [Anthropic Model Provider](#anthropic-model-provider)
  - [Ollama Model Provider](#ollama-model-provider)
  - [Voyage Model Provider](#voyage-model-provider)
- [Authentication Providers](#authentication-providers)
  - [GitHub Auth Provider](#github-auth-provider)
  - [Google Auth Provider](#google-auth-provider)
  - [Keycloak Auth Provider](#keycloak-auth-provider)
  - [Entra ID Auth Provider](#entra-id-auth-provider)
- [Core Tools](#core-tools)
  - [Memory Tool](#memory-tool)
  - [Knowledge Tool](#knowledge-tool)
  - [Tasks Tool](#tasks-tool)
  - [Workspace Files Tool](#workspace-files-tool)
- [Credential Stores](#credential-stores)
  - [SQLite Credential Store](#sqlite-credential-store)
  - [PostgreSQL Credential Store](#postgresql-credential-store)

---

## Model Providers

### Common Model Provider API

All model providers implement a consistent HTTP API based on the OpenAI API specification.

#### Base Configuration

**Default Port**: `8000`  
**Protocol**: HTTP  
**Base Path**: `/v1`

#### Environment Variables

All providers support:

- `PORT` - Server port (default: `8000`)
- `GPTSCRIPT_DEBUG` - Enable debug logging (`true`/`false`)
- `OBOT_<PROVIDER>_MODEL_PROVIDER_API_KEY` - Global API key for the provider

#### Authentication

**Global API Key** (environment variable):

```bash
export OBOT_OPENAI_MODEL_PROVIDER_API_KEY=sk-...
```

**Per-Request API Key** (HTTP header):

```
X-Obot-OBOT_OPENAI_MODEL_PROVIDER_API_KEY: sk-...
```

Per-request API keys override global API keys.

#### Common Endpoints

##### GET /

**Description**: Returns the server URI.

**Response**:

```
http://127.0.0.1:8000
```

##### GET /v1/models

**Description**: List available models from the provider.

**Response**:

```json
{
  "object": "list",
  "data": [
    {
      "id": "model-id",
      "object": "model",
      "created": 1234567890,
      "owned_by": "provider-name"
    }
  ]
}
```

##### POST /v1/chat/completions

**Description**: Create a chat completion (for chat-capable models).

**Request**:

```json
{
  "model": "gpt-4",
  "messages": [
    {
      "role": "user",
      "content": "Hello!"
    }
  ],
  "temperature": 0.7,
  "max_tokens": 100,
  "stream": false
}
```

**Response**:

```json
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1234567890,
  "model": "gpt-4",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Hello! How can I help you?"
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 20,
    "total_tokens": 30
  }
}
```

##### POST /v1/embeddings

**Description**: Create embeddings for text (for embedding-capable models).

**Request**:

```json
{
  "model": "text-embedding-ada-002",
  "input": "The quick brown fox"
}
```

**Response**:

```json
{
  "object": "list",
  "data": [
    {
      "object": "embedding",
      "embedding": [0.1, 0.2, 0.3, ...],
      "index": 0
    }
  ],
  "model": "text-embedding-ada-002",
  "usage": {
    "prompt_tokens": 5,
    "total_tokens": 5
  }
}
```

#### Validation Mode

All providers support validation to test credentials.

**Command**:

```bash
cd <provider-directory>
go run . validate
```

**Behavior**: Makes a test API call and exits with code 0 (success) or 1 (failure).

---

### OpenAI Model Provider

**Directory**: `openai-model-provider/`  
**Language**: Go  
**Upstream API**: `https://api.openai.com/v1`

#### Environment Variables

- `OBOT_OPENAI_MODEL_PROVIDER_API_KEY` - OpenAI API key

#### Supported Endpoints

- `GET /v1/models`
- `POST /v1/chat/completions`
- `POST /v1/completions`
- `POST /v1/embeddings`
- `POST /v1/images/generations`
- All other OpenAI API endpoints (proxied)

#### Model Rewriting

The provider rewrites the `/v1/models` response to include only relevant models and adds metadata.

#### Custom Paths

- `/v1/*` - Reverse proxy to OpenAI API

---

### Anthropic Model Provider

**Directory**: `anthropic-model-provider-go/`  
**Language**: Go  
**Upstream API**: `https://api.anthropic.com`

#### Environment Variables

- `OBOT_ANTHROPIC_MODEL_PROVIDER_API_KEY` - Anthropic API key

#### Supported Endpoints

- `GET /v1/models` - Lists Anthropic models in OpenAI-compatible format
- `POST /v1/messages` - Anthropic messages API (native format)
- `POST /v1/chat/completions` - OpenAI-compatible chat completions (converted to Anthropic format)

#### Model Translation

Translates between OpenAI chat completions format and Anthropic messages format:

- Request: OpenAI → Anthropic
- Response: Anthropic → OpenAI

---

### Ollama Model Provider

**Directory**: `ollama-model-provider/`  
**Language**: Go  
**Upstream API**: Ollama local server

#### Environment Variables

- `OLLAMA_HOST` - Ollama server URL (default: `http://localhost:11434`)

#### Supported Endpoints

- `GET /v1/models` - Lists locally available Ollama models
- `POST /v1/chat/completions` - Chat completions via Ollama
- `POST /v1/embeddings` - Embeddings via Ollama

#### Notes

- No API key required (local server)
- Requires Ollama server running locally or accessible via network

---

### Voyage Model Provider

**Directory**: `voyage-model-provider/`  
**Language**: Python (FastAPI)  
**Upstream API**: Voyage AI API

#### Environment Variables

- `OBOT_VOYAGE_MODEL_PROVIDER_API_KEY` - Voyage AI API key

#### Supported Endpoints

- `GET /v1/models` - Lists Voyage embedding models
- `POST /v1/embeddings` - Create embeddings

#### Available Models

As of the latest update:

- `voyage-3-large`
- `voyage-3`
- `voyage-3-lite`
- `voyage-code-3`
- `voyage-finance-2`
- `voyage-multilingual-2`
- `voyage-law-2`
- `voyage-code-2`

**Note**: Model list is hardcoded and should be updated when Voyage adds new models.

#### Embeddings Request

**Request**:

```json
{
  "model": "voyage-3",
  "input": "Text to embed"
}
```

**Response**:

```json
{
  "data": [
    {
      "embedding": [0.1, 0.2, ...]
    }
  ]
}
```

---

## Authentication Providers

### Common Auth Provider API

All authentication providers implement the Obot auth provider protocol.

#### Required Endpoints

##### GET /oauth2/start

**Description**: Initiate OAuth2 authorization flow.

**Query Parameters**:

- `rd` (required) - Redirect URL after successful authentication

**Behavior**:

1. Generates OAuth2 state parameter
2. Stores redirect URL with state
3. Redirects user to OAuth provider's authorization URL

**Example**:

```
GET /oauth2/start?rd=https://obot.example.com/dashboard
```

##### GET /oauth2/callback

**Description**: Handle OAuth2 callback from provider.

**Query Parameters**:

- `code` (required) - Authorization code from OAuth provider
- `state` (required) - State parameter for CSRF protection

**Behavior**:

1. Validates state parameter
2. Exchanges authorization code for access token
3. Retrieves user profile
4. Validates user (email domain, team, org, etc.)
5. Encrypts and sets `obot_access_token` cookie
6. Redirects to URL from `rd` parameter

**Example**:

```
GET /oauth2/callback?code=abc123&state=xyz789
```

##### GET /oauth2/sign_out

**Description**: Sign out user and clear session.

**Query Parameters**:

- `rd` (required) - Redirect URL after sign out

**Behavior**:

1. Clears `obot_access_token` cookie
2. Redirects to URL from `rd` parameter

**Example**:

```
GET /oauth2/sign_out?rd=https://obot.example.com/
```

##### GET /obot-get-icon-url

**Description**: Get user's profile picture URL.

**Headers**:

- `Authorization: Bearer <access_token>` (required)

**Response**:

```json
{
  "iconURL": "https://avatars.githubusercontent.com/u/123456"
}
```

##### POST /obot-get-state

**Description**: Get authenticated user state from HTTP request.

**Request Body**:

```json
{
  "method": "GET",
  "url": "https://obot.example.com/api/...",
  "header": {
    "Cookie": ["obot_access_token=encrypted_value"],
    "User-Agent": ["Mozilla/5.0 ..."]
  }
}
```

**Response** (Success):

```json
{
  "accessToken": "github_token_xyz",
  "preferredUsername": "johndoe",
  "user": "johndoe",
  "email": "johndoe@example.com"
}
```

**Response** (Unauthenticated):

```
HTTP 400 Bad Request
```

#### Cookie Management

**Cookie Name**: `obot_access_token`

**Cookie Attributes**:

- `HttpOnly`: Always true
- `Secure`: True if `OBOT_SERVER_PUBLIC_URL` starts with `https://`
- `SameSite`: Lax
- `Path`: /
- `Encryption`: AES-256 using `OBOT_AUTH_PROVIDER_COOKIE_SECRET`

---

### GitHub Auth Provider

**Directory**: `github-auth-provider/`  
**Language**: Go  
**OAuth Provider**: GitHub

#### Environment Variables

**Required**:

- `OBOT_GITHUB_AUTH_PROVIDER_CLIENT_ID` - GitHub OAuth app client ID
- `OBOT_GITHUB_AUTH_PROVIDER_CLIENT_SECRET` - GitHub OAuth app client secret
- `OBOT_AUTH_PROVIDER_COOKIE_SECRET` - Cookie encryption secret (32+ bytes)
- `OBOT_AUTH_PROVIDER_EMAIL_DOMAINS` - Comma-separated allowed email domains

**Optional**:

- `OBOT_GITHUB_AUTH_PROVIDER_TEAMS` - Comma-separated allowed teams (format: `org/team`)
- `OBOT_GITHUB_AUTH_PROVIDER_ORG` - Required organization
- `OBOT_GITHUB_AUTH_PROVIDER_REPO` - Required repository access (format: `org/repo`)
- `OBOT_GITHUB_AUTH_PROVIDER_TOKEN` - GitHub token for additional API calls
- `OBOT_GITHUB_AUTH_PROVIDER_ALLOW_USERS` - Comma-separated allowed usernames

#### User Validation

Users must satisfy **all** configured validations:

1. Email domain matches one of `OBOT_AUTH_PROVIDER_EMAIL_DOMAINS`
2. If `OBOT_GITHUB_AUTH_PROVIDER_ORG` is set, user must be member of org
3. If `OBOT_GITHUB_AUTH_PROVIDER_TEAMS` is set, user must be member of at least one team
4. If `OBOT_GITHUB_AUTH_PROVIDER_REPO` is set, user must have access to repo
5. If `OBOT_GITHUB_AUTH_PROVIDER_ALLOW_USERS` is set, username must be in list

#### Profile API

**Package**: `pkg/profile`

**Function**: `GetProfile(accessToken string) (*Profile, error)`

**Profile Structure**:

```go
type Profile struct {
    User              string
    Email             string
    PreferredUsername string
}
```

---

### Google Auth Provider

**Directory**: `google-auth-provider/`  
**Language**: Go  
**OAuth Provider**: Google  
**Status**: Reference implementation

#### Environment Variables

**Required**:

- `OBOT_GOOGLE_AUTH_PROVIDER_CLIENT_ID` - Google OAuth client ID
- `OBOT_GOOGLE_AUTH_PROVIDER_CLIENT_SECRET` - Google OAuth client secret
- `OBOT_AUTH_PROVIDER_COOKIE_SECRET` - Cookie encryption secret
- `OBOT_AUTH_PROVIDER_EMAIL_DOMAINS` - Comma-separated allowed email domains

**Optional**:

- `OBOT_GOOGLE_AUTH_PROVIDER_ALLOW_USERS` - Comma-separated allowed email addresses

#### User Validation

Users must satisfy:

1. Email domain matches one of `OBOT_AUTH_PROVIDER_EMAIL_DOMAINS`
2. If `OBOT_GOOGLE_AUTH_PROVIDER_ALLOW_USERS` is set, email must be in list

#### Profile API

**Package**: `pkg/profile`

**Function**: `GetProfile(accessToken string) (*Profile, error)`

---

### Keycloak Auth Provider

**Directory**: `keycloak-auth-provider/`
**Language**: Go
**OAuth Provider**: Keycloak / Red Hat SSO (OIDC)

#### Environment Variables

**Required**:

- `OBOT_KEYCLOAK_AUTH_PROVIDER_CLIENT_ID` - Keycloak client ID
- `OBOT_KEYCLOAK_AUTH_PROVIDER_CLIENT_SECRET` - Keycloak client secret
- `OBOT_KEYCLOAK_AUTH_PROVIDER_URL` - Keycloak base URL (e.g., `https://keycloak.example.com`)
- `OBOT_KEYCLOAK_AUTH_PROVIDER_REALM` - Keycloak realm name
- `OBOT_AUTH_PROVIDER_COOKIE_SECRET` - Cookie encryption secret (32+ bytes)
- `OBOT_AUTH_PROVIDER_EMAIL_DOMAINS` - Comma-separated allowed email domains (use `*` for all)

**Optional**:

- `OBOT_KEYCLOAK_AUTH_PROVIDER_ALLOWED_GROUPS` - Comma-separated allowed Keycloak groups
- `OBOT_KEYCLOAK_AUTH_PROVIDER_ALLOWED_ROLES` - Comma-separated allowed roles (format: `role` or `client-id:role`)
- `OBOT_KEYCLOAK_AUTH_PROVIDER_GROUP_CACHE_TTL` - Group cache duration (default: `1h`)
- `OBOT_AUTH_PROVIDER_POSTGRES_CONNECTION_DSN` - PostgreSQL for session storage
- `OBOT_AUTH_PROVIDER_TOKEN_REFRESH_DURATION` - Token refresh interval (default: `1h`)

#### User Validation

Users must satisfy **all** configured validations:

1. Email domain matches one of `OBOT_AUTH_PROVIDER_EMAIL_DOMAINS`
2. If `OBOT_KEYCLOAK_AUTH_PROVIDER_ALLOWED_GROUPS` is set, user must be member of at least one group
3. If `OBOT_KEYCLOAK_AUTH_PROVIDER_ALLOWED_ROLES` is set, user must have at least one role

#### Profile API

**Package**: `pkg/profile`

**Function**: `GetProfile(accessToken string) (*Profile, error)`

---

### Entra ID Auth Provider

**Directory**: `entra-auth-provider/`
**Language**: Go
**OAuth Provider**: Microsoft Entra ID (Azure AD)

#### Environment Variables

**Required**:

- `OBOT_ENTRA_AUTH_PROVIDER_CLIENT_ID` - Azure App Registration client ID
- `OBOT_ENTRA_AUTH_PROVIDER_CLIENT_SECRET` - Azure App Registration client secret
- `OBOT_ENTRA_AUTH_PROVIDER_TENANT_ID` - Azure tenant ID (or `common`/`organizations` for multi-tenant)
- `OBOT_AUTH_PROVIDER_COOKIE_SECRET` - Cookie encryption secret (32+ bytes)
- `OBOT_AUTH_PROVIDER_EMAIL_DOMAINS` - Comma-separated allowed email domains (use `*` for all)

**Optional**:

- `OBOT_ENTRA_AUTH_PROVIDER_ALLOWED_GROUPS` - Comma-separated Azure AD group IDs
- `OBOT_ENTRA_AUTH_PROVIDER_ALLOWED_TENANTS` - Comma-separated tenant IDs (required for multi-tenant)
- `OBOT_ENTRA_AUTH_PROVIDER_USE_WORKLOAD_IDENTITY` - Enable Azure Workload Identity (`true`/`false`)
- `OBOT_ENTRA_AUTH_PROVIDER_GROUP_CACHE_TTL` - Group cache duration (default: `1h`)
- `OBOT_ENTRA_AUTH_PROVIDER_ICON_CACHE_TTL` - Profile picture cache duration (default: `24h`)
- `OBOT_AUTH_PROVIDER_POSTGRES_CONNECTION_DSN` - PostgreSQL for session storage
- `OBOT_AUTH_PROVIDER_TOKEN_REFRESH_DURATION` - Token refresh interval (default: `1h`)

#### User Validation

Users must satisfy **all** configured validations:

1. Email domain matches one of `OBOT_AUTH_PROVIDER_EMAIL_DOMAINS`
2. If `OBOT_ENTRA_AUTH_PROVIDER_ALLOWED_GROUPS` is set, user must be member of at least one group
3. If `OBOT_ENTRA_AUTH_PROVIDER_ALLOWED_TENANTS` is set (multi-tenant), user's tenant must be in list

#### Profile API

**Package**: `pkg/profile`

**Function**: `GetProfile(accessToken string) (*Profile, error)`

---

## Core Tools

### Memory Tool

**Directory**: `memory/`  
**Language**: Go  
**Purpose**: Long-term agent memory storage

#### Tool Definition

**Name**: Memory  
**Description**: Interact with agent memory to store and retrieve information  
**Context**: `memory_context`  
**Category**: Capability

#### Sub-Tools

##### Create Memory

**Name**: Create Memory  
**Description**: Store information in agent memory

**Parameters**:

- `content` (string, required) - The content to remember

**Usage**:

```gptscript
Create Memory with content="User prefers Python over JavaScript"
```

##### Update Memory

**Name**: Update Memory  
**Description**: Update information in agent memory

**Parameters**:

- `memory_id` (string, required) - The ID of the memory to update
- `content` (string, required) - The updated content

**Usage**:

```gptscript
Update Memory with memory_id="123" and content="User prefers Python and Go"
```

##### Delete Memory

**Name**: Delete Memory  
**Description**: Delete information from agent memory

**Parameters**:

- `memory_id` (string, required) - The ID of the memory to delete

**Usage**:

```gptscript
Delete Memory with memory_id="123"
```

##### list_memories

**Name**: list_memories  
**Type**: context  
**Description**: List all memories

**Returns**: JSON array of all stored memories

---

### Knowledge Tool

**Directory**: `knowledge/`  
**Language**: Go  
**Purpose**: Knowledge base with RAG capabilities

#### Tool Definition

**Name**: Knowledge  
**Description**: Knowledge retrieval system  
**Category**: Capability

#### GPTScript Tools

##### Main Knowledge Tool

**File**: `tool.gpt`  
**Description**: Retrieve information from knowledge base

**Parameters**:

- `query` (string, required) - The question or search query
- `dataset` (string, optional) - Dataset name to search
- `topK` (int, optional) - Number of results to return

##### Knowledge Ingest

**File**: `ingest.gpt`  
**Description**: Ingest documents into knowledge base

**Parameters**:

- `directory` (string, required) - Directory path to ingest
- `dataset` (string, optional) - Dataset name to create/update
- `config` (string, optional) - Path to configuration file

##### Knowledge Load

**File**: `load.gpt`  
**Description**: Load remote content into knowledge base

**Parameters**:

- `url` (string, required) - URL to load (supports GitHub repos)
- `dataset` (string, optional) - Dataset name

##### Knowledge Delete

**File**: `delete.gpt`  
**Description**: Delete a dataset

**Parameters**:

- `dataset` (string, required) - Dataset name to delete

##### Knowledge Delete File

**File**: `delete-file.gpt`  
**Description**: Delete a file from a dataset

**Parameters**:

- `dataset` (string, required) - Dataset name
- `file` (string, required) - File path to delete

#### Command-Line Interface

The knowledge tool also provides a comprehensive CLI.

##### Commands

###### Retrieve

**Command**: `knowledge retrieve`

**Flags**:

- `--query` - Search query
- `--dataset` - Dataset name
- `--config` - Config file path
- `--top-k` - Number of results

**Example**:

```bash
knowledge retrieve --query "How to configure embeddings?" --dataset docs --top-k 5
```

###### Ingest

**Command**: `knowledge ingest`

**Flags**:

- `--directory` - Directory to ingest
- `--dataset` - Dataset name
- `--config` - Config file path
- `--prune` - Remove deleted files

**Example**:

```bash
knowledge ingest --directory ./docs --dataset documentation --prune
```

###### Load

**Command**: `knowledge load`

**Flags**:

- `--url` - URL to load
- `--dataset` - Dataset name
- `--config` - Config file path

**Example**:

```bash
knowledge load --url https://github.com/org/repo --dataset repo-docs
```

###### List Datasets

**Command**: `knowledge list-datasets`

**Example**:

```bash
knowledge list-datasets
```

###### Get Dataset

**Command**: `knowledge get-dataset`

**Flags**:

- `--dataset` - Dataset name

**Example**:

```bash
knowledge get-dataset --dataset documentation
```

###### Delete Dataset

**Command**: `knowledge delete-dataset`

**Flags**:

- `--dataset` - Dataset name

**Example**:

```bash
knowledge delete-dataset --dataset old-docs
```

#### Configuration

**Format**: YAML or JSON

**Example Configuration**:

```yaml
# Embedding provider
embeddingProvider:
  type: openai
  model: text-embedding-3-small
  apiKey: ${OPENAI_API_KEY}

# Vector store
vectorStore:
  type: pgvector
  connectionString: postgresql://localhost/knowledge

# Index store
indexStore:
  type: postgres
  connectionString: postgresql://localhost/knowledge

# Text splitter
textSplitter:
  chunkSize: 1000
  chunkOverlap: 200

# Retrievers
retrievers:
  - type: hybrid
    vectorWeight: 0.7
    keywordWeight: 0.3

# Post-processors
postProcessors:
  - type: similarityScore
    threshold: 0.7
  - type: topK
    k: 5
```

#### Environment Variables

- `OBOT_KNOWLEDGE_CONFIG` - Path to config file
- `OBOT_KNOWLEDGE_DATASET` - Default dataset name
- `OBOT_KNOWLEDGE_POSTGRES_URL` - PostgreSQL connection string
- `OBOT_KNOWLEDGE_SQLITE_PATH` - SQLite database path
- `OPENAI_API_KEY` - OpenAI API key for embeddings

---

### Tasks Tool

**Directory**: `tasks/`  
**Language**: Go  
**Purpose**: Task management and execution

#### Tool Definition

**Name**: Tasks  
**Description**: Task management capabilities  
**Category**: Capability

#### Features

- Create and manage tasks
- Track task status
- Execute task workflows
- Task dependencies

---

### Workspace Files Tool

**Directory**: `workspace-files/`  
**Language**: Go  
**Purpose**: File operations in agent workspace

#### Tool Definition

**Name**: Workspace Files  
**Description**: File operations in workspace  
**Category**: Capability

#### Features

- Read files from workspace
- Write files to workspace
- List workspace files
- Delete workspace files

---

## Credential Stores

### Common Credential Store API

Credential stores implement the following tool commands:

#### store

**Description**: Store a credential

**Parameters**:

- `name` (string, required) - Credential name
- `env` (map, required) - Environment variables as key-value pairs

**Example**:

```gptscript
store with name="my-api-key" and env={"API_KEY": "secret123"}
```

#### reveal

**Description**: Retrieve and decrypt a credential

**Parameters**:

- `name` (string, required) - Credential name

**Returns**: Environment variables as key-value pairs

**Example**:

```gptscript
reveal with name="my-api-key"
```

#### delete

**Description**: Delete a credential

**Parameters**:

- `name` (string, required) - Credential name

**Example**:

```gptscript
delete with name="my-api-key"
```

---

### SQLite Credential Store

**Directory**: `credential-stores/sqlite/`  
**Language**: Go  
**Backend**: SQLite database

#### Environment Variables

- `OBOT_CREDENTIAL_STORE_SQLITE_FILE` - SQLite database file path (default: `credentials.db`)
- `OBOT_CREDENTIAL_STORE_ENCRYPTION_KEY` - Encryption key (required, 32+ bytes)

#### Database Schema

```sql
CREATE TABLE credentials (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    encrypted_data BLOB NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### Encryption

- Algorithm: AES-256-GCM
- Key derivation: PBKDF2 with salt
- Unique IV per credential

---

### PostgreSQL Credential Store

**Directory**: `credential-stores/postgres/`  
**Language**: Go  
**Backend**: PostgreSQL database

#### Environment Variables

- `OBOT_CREDENTIAL_STORE_POSTGRES_URL` - PostgreSQL connection string (required)
- `OBOT_CREDENTIAL_STORE_ENCRYPTION_KEY` - Encryption key (required, 32+ bytes)

#### Connection String Format

```
postgresql://user:password@host:port/database?sslmode=require
```

#### Database Schema

```sql
CREATE TABLE credentials (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    encrypted_data BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Encryption

- Algorithm: AES-256-GCM
- Key derivation: PBKDF2 with salt
- Unique IV per credential

---

## Error Codes and Responses

### HTTP Status Codes

#### Model Providers

- `200 OK` - Successful request
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Missing or invalid API key
- `404 Not Found` - Model or endpoint not found
- `429 Too Many Requests` - Rate limit exceeded (from upstream)
- `500 Internal Server Error` - Provider or proxy error
- `502 Bad Gateway` - Upstream API error
- `503 Service Unavailable` - Provider temporarily unavailable

#### Authentication Providers

- `200 OK` - Successful authentication state retrieval
- `302 Found` - Redirect (OAuth flow)
- `400 Bad Request` - Invalid request or unauthenticated
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - User validation failed
- `500 Internal Server Error` - Server error

### Error Response Format

**Model Providers** (OpenAI-compatible):

```json
{
  "error": {
    "message": "Error description",
    "type": "invalid_request_error",
    "code": "invalid_api_key"
  }
}
```

**Authentication Providers**:

```json
{
  "error": "Error description",
  "details": "Additional context"
}
```

---

## Rate Limits and Quotas

Rate limits are determined by upstream providers:

- **OpenAI**: Per-tier limits (see OpenAI documentation)
- **Anthropic**: Per-tier limits (see Anthropic documentation)
- **Ollama**: No rate limits (local)
- **Voyage**: Per-account limits (see Voyage documentation)

Obot tool providers do not impose additional rate limits.

---

## Versioning

### API Versioning

Model provider APIs follow OpenAI API versioning (`/v1`).

### Tool Versioning

Tools are versioned implicitly through Git tags and releases.

### Compatibility

- **Backward compatibility**: Minor version changes maintain compatibility
- **Breaking changes**: Major version changes may break compatibility
- **Deprecation policy**: Deprecated endpoints supported for at least 6 months

---

## Support and Documentation

- **Main repository**: https://github.com/obot-platform/obot
- **Issues**: Use the `tools` label in the main Obot repository
- **Auth provider reference**: See `docs/auth-providers.md`
