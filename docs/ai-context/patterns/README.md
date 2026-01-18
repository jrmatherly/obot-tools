# Obot-Tools Implementation Patterns

**Purpose:** Reusable implementation patterns for creating tools, model providers, and auth providers
**Last Updated:** 2026-01-16
**Versioning:** Patterns follow [Semantic Versioning](https://semver.org/) - see [CHANGELOG.md](./CHANGELOG.md)

---

## Available Patterns

### 1. [Model Provider Creation](./model-provider-creation.md) `v1.0.0`

**Pattern:** OpenAI-Compatible Model Provider Implementation
**Use Case:** Create new model provider for Obot platform (OpenAI, Anthropic, Ollama, etc.)
**Components:** Go HTTP server, GPTScript tool.gpt, OpenAI-compatible API

**When to Use:**

- Adding new LLM provider to Obot
- Implementing OpenAI-compatible proxy
- Creating custom model routing or caching
- Troubleshooting model provider connectivity or API issues

**Key Features:**

- ✅ OpenAI-compatible API format
- ✅ GPTScript tool.gpt integration with `Model Provider: true`
- ✅ `#!sys.daemon` for persistent server process
- ✅ `providerMeta` for model listing and capabilities

---

### 2. [Auth Provider Creation](./auth-provider-creation.md) `v1.0.0`

**Pattern:** OAuth2 Authentication Provider for Obot
**Use Case:** Create OAuth2 auth provider (GitHub, Google, custom OIDC)
**Components:** Go HTTP server, GPTScript tool.gpt, encrypted cookie session

**When to Use:**

- Adding new OAuth2 provider to Obot
- Implementing OIDC authentication flow
- Creating custom authorization logic
- Troubleshooting OAuth callback or token issues

**Key Features:**

- ✅ OAuth2/OIDC flow implementation
- ✅ GPTScript tool.gpt with `envVars` and `optionalEnvVars`
- ✅ Encrypted cookie-based session management (AES-256-GCM)
- ✅ Authorization endpoint, token endpoint, user info endpoint

---

### 3. [GPTScript Tool Development](./gptscript-tool-development.md) `v1.0.0`

**Pattern:** GPTScript .gpt File Format and Tool Creation
**Use Case:** Create GPTScript tools for Obot platform
**Components:** tool.gpt file, Go/Python/TypeScript implementation

**When to Use:**

- Creating new tool for Obot (knowledge, memory, file operations, API integrations)
- Understanding GPTScript .gpt file format
- Implementing context tools with `Type: context` and `Share Context`
- Debugging tool.gpt parsing or parameter handling

**Key Features:**

- ✅ Complete .gpt file format reference
- ✅ Tool types (standard, context, model provider, auth provider)
- ✅ Parameter handling (uppercase env var conversion)
- ✅ Tool directives (#!, Args, Share Context, etc.)

---

## Pattern Selection Guide

### I need to

**...create model provider:**
→ [Model Provider Creation](./model-provider-creation.md)

**...create auth provider:**
→ [Auth Provider Creation](./auth-provider-creation.md)

**...create GPTScript tool:**
→ [GPTScript Tool Development](./gptscript-tool-development.md)

**...understand .gpt file format:**
→ [GPTScript Tool Development](./gptscript-tool-development.md)

**...troubleshoot provider issues:**

- Model provider API → [Model Provider Creation](./model-provider-creation.md) (Troubleshooting section)
- Auth provider OAuth → [Auth Provider Creation](./auth-provider-creation.md) (OAuth Flow section)
- Tool.gpt parsing → [GPTScript Tool Development](./gptscript-tool-development.md) (Format section)

---

## Using Patterns

### Pattern Structure

Each pattern document includes:

1. **Overview**: What the pattern is and key features
2. **Architecture**: Component diagram showing tool/provider structure
3. **Prerequisites**: Required tools, dependencies, environment setup
4. **Procedure**: Step-by-step implementation guide
5. **Code Examples**: Complete tool.gpt examples and Go/Python code
6. **Testing**: How to test locally and validate functionality
7. **Troubleshooting**: Common issues and debugging techniques
8. **Related Documentation**: Links to canonical examples and references

### Integration with Component Docs

Obot-Tools documentation **references** these patterns instead of duplicating. This provides:

- ✅ **Single Source of Truth**: Updates apply to all components
- ✅ **Consistent Implementation**: Same patterns across providers
- ✅ **Easier Maintenance**: Edit once, benefits all documentation
- ✅ **Faster Learning**: Focus on what's unique, reference pattern for common parts

---

## Canonical Examples

Refer to these existing tools/providers as reference implementations:

### Model Providers

- `openai-model-provider/tool.gpt` - Standard OpenAI proxy pattern
- `anthropic-model-provider/tool.gpt` - Anthropic proxy with streaming
- `ollama-model-provider/tool.gpt` - Local model provider pattern

### Auth Providers

- `github-auth-provider/tool.gpt` - GitHub OAuth2 implementation
- `google-auth-provider/tool.gpt` - Google OIDC implementation

### Context Tools

- `memory/tool.gpt` - Context tool with `Type: context` and `Share Context`
- `knowledge/tool.gpt` - Complex tool with multiple features and RAG

### Standard Tools

- `echo/tool.gpt` - Simple tool pattern
- `file-operations/tool.gpt` - File system interaction pattern

---

## Contributing Patterns

### When to Create a New Pattern

Create a new pattern when:

1. **Common implementation across 3+ providers/tools**
2. **New tool type or provider category** not covered by existing patterns
3. **Complex procedure** requiring detailed step-by-step guide
4. **Best practice** that should be standardized

### Pattern Template

```markdown
# Pattern Name

**Pattern:** One-line description
**Use Case:** When to use this pattern
**Components:** What you'll create
**Version:** 1.0.0
**Last Updated:** January 2026

---

## Overview
What this pattern covers and benefits

## Architecture
Diagram showing component structure

## Prerequisites
- Go 1.23+ / Python 3.13+ / Node.js 18+
- GPTScript knowledge
- Required environment variables

## Procedure

### Step 1: Create tool.gpt
tool.gpt file structure and directives

### Step 2: Implement Server (if applicable)
Go/Python/TypeScript implementation

### Step 3: Test Locally
Local testing procedure

### Step 4: Integration
How to integrate with Obot

## Code Examples
Complete working examples

## Testing
Validation steps and test cases

## Troubleshooting
Common issues and solutions

## Canonical Examples
References to existing implementations

---

**Last Updated:** Date
**Pattern Version:** X.Y.Z
**Tested With:** GPTScript version, Go version, etc.
```

---

## Related Documentation

### Obot-Tools Documentation

- [CLAUDE.md](../../CLAUDE.md) - Complete tool/provider architecture
- [README.md](../../README.md) - Repository overview
- [docs/auth-providers.md](../../docs/auth-providers.md) - Auth provider requirements

### Serena Memories

- `.serena/memories/gptscript_tool_format.md` - Enhanced .gpt file format reference

### Component Directories

- Model Providers: `openai-model-provider/`, `anthropic-model-provider/`, `ollama-model-provider/`, etc.
- Auth Providers: `github-auth-provider/`, `google-auth-provider/`
- Tools: `memory/`, `knowledge/`, `echo/`, `file-operations/`
- Registry: `index.yaml` - Central tool/provider registry

---

**Total Patterns:** 3
**Coverage:** Model providers, auth providers, GPTScript tools
**Status:** Active development - patterns versioned and maintained
