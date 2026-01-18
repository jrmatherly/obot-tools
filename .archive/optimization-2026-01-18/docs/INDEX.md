# Obot Tools - Documentation Index

Welcome to the Obot Tools documentation. This index provides quick access to all documentation resources.

## 📖 Core Documentation

### [Architecture Guide](ARCHITECTURE.md)

**Comprehensive system architecture and design documentation**

Covers:

- System architecture overview and component organization
- Detailed component categories (tools, model providers, auth providers, credentials)
- Knowledge base deep dive with RAG architecture
- Build and deployment architecture
- Integration points with Obot platform and GPTScript
- Data flow examples and sequence diagrams
- Technology choices and rationale
- Security considerations and observability patterns
- Performance characteristics and scalability

**Audience**: Architects, senior developers, DevOps engineers

---

### [API Reference](API_REFERENCE.md)

**Complete API documentation for all components**

Covers:

- Model provider APIs (OpenAI, Anthropic, Ollama, Voyage, etc.)
  - Common API patterns and authentication
  - Endpoint specifications and request/response formats
  - Model listings and capabilities
- Authentication provider APIs (GitHub, Google)
  - OAuth2 flow endpoints
  - Cookie management and security
  - User validation and state management
- Core tools (Memory, Knowledge, Tasks, Workspace Files)
  - Tool definitions and sub-tools
  - Command-line interfaces
  - Configuration formats
- Credential stores (SQLite, PostgreSQL)
  - Storage and retrieval APIs
  - Encryption and security
- Error codes, rate limits, and versioning

**Audience**: Developers integrating with Obot tools, API consumers

---

### [Developer Guide](DEVELOPER_GUIDE.md)

**Hands-on guide for developing with and contributing to Obot tools**

Covers:

- Getting started and prerequisites
- Development environment setup (Go, Python, Node.js, databases)
- Building the project (all tools, specific components, Docker)
- Testing strategies and commands
- Creating new components
  - Model providers
  - Tools
  - Auth providers
- Code style and best practices
  - Go conventions
  - Python conventions
  - GPTScript tool definitions
  - Commit message guidelines
- Debugging techniques
- Contributing workflow and PR guidelines

**Audience**: Contributors, developers building new tools and providers

---

### [Auth Provider Requirements](auth-providers.md)

**Detailed specifications for building authentication providers**

Covers:

- Daemon tool requirements
- Metadata specifications (envVars, optionalEnvVars)
- Placeholder credential references
- OAuth2 implementation requirements
- Required URL paths and endpoints
- Token cookie management and security
- User state retrieval via `/obot-get-state`
- Request/response schemas
- Reference implementation (Google auth provider)

**Audience**: Developers creating new authentication providers

---

## 🎯 Quick Reference Guides

### Component Overview

#### Core Tools (User-Facing)

- **[memory/](../memory/)** - Agent memory storage
- **[knowledge/](../knowledge/)** - RAG knowledge base
- **[tasks/](../tasks/)** - Task management
- **[workspace-files/](../workspace-files/)** - File operations
- **[threads/](../threads/)** - Thread management
- **[time/](../time/)** - Time utilities
- **[images/](../images/)** - Image operations
- **[file-summarizer/](../file-summarizer/)** - Document summarization

#### Model Providers

- **[openai-model-provider/](../openai-model-provider/)** - OpenAI API
- **[anthropic-model-provider-go/](../anthropic-model-provider-go/)** - Anthropic Claude
- **[ollama-model-provider/](../ollama-model-provider/)** - Local Ollama models
- **[groq-model-provider/](../groq-model-provider/)** - Groq API
- **[xai-model-provider/](../xai-model-provider/)** - xAI Grok
- **[deepseek-model-provider/](../deepseek-model-provider/)** - DeepSeek
- **[voyage-model-provider/](../voyage-model-provider/)** - Voyage embeddings
- **[vllm-model-provider/](../vllm-model-provider/)** - VLLM self-hosted
- **[generic-openai-model-provider/](../generic-openai-model-provider/)** - OpenAI-compatible

#### Authentication Providers

- **[github-auth-provider/](../github-auth-provider/)** - GitHub OAuth2
- **[google-auth-provider/](../google-auth-provider/)** - Google OAuth2

#### Credential Stores

- **[credential-stores/sqlite/](../credential-stores/sqlite/)** - SQLite backend
- **[credential-stores/postgres/](../credential-stores/postgres/)** - PostgreSQL backend

---

## 🚀 Getting Started Paths

### For New Users

1. Read the project [README](../README.md)
2. Review [Architecture Guide](ARCHITECTURE.md) - Overview section
3. Check [API Reference](API_REFERENCE.md) - Common Model Provider API
4. Explore component directories

### For Contributors

1. Review [Developer Guide](DEVELOPER_GUIDE.md) - Getting Started
2. Set up development environment (Developer Guide - Development Environment Setup)
3. Build and test (Developer Guide - Building the Project, Testing)
4. Review code style (Developer Guide - Code Style and Best Practices)
5. Follow contributing workflow (Developer Guide - Contributing)

### For Integration Developers

1. Review [API Reference](API_REFERENCE.md) - Relevant component APIs
2. Check [Architecture Guide](ARCHITECTURE.md) - Integration Points
3. Study example implementations in component directories
4. Test with development environment

### For Platform Architects

1. Read [Architecture Guide](ARCHITECTURE.md) completely
2. Review [API Reference](API_REFERENCE.md) - All sections
3. Understand GPTScript integration (Architecture Guide - GPTScript Integration)
4. Study data flows (Architecture Guide - Data Flow Examples)

---

## 📋 Common Tasks

### Building and Testing

- **Build all tools**: See [Developer Guide - Building the Project](DEVELOPER_GUIDE.md#building-the-project)
- **Run tests**: See [Developer Guide - Testing](DEVELOPER_GUIDE.md#testing)
- **Package for distribution**: See [Developer Guide - Package Tools](DEVELOPER_GUIDE.md#package-tools-for-distribution)

### Creating Components

- **New model provider**: See [Developer Guide - Creating a New Model Provider](DEVELOPER_GUIDE.md#creating-a-new-model-provider)
- **New tool**: See [Developer Guide - Creating a New Tool](DEVELOPER_GUIDE.md#creating-a-new-tool)
- **New auth provider**: See [Auth Provider Requirements](auth-providers.md)

### Configuration and Deployment

- **Environment variables**: See [API Reference - Environment Variables](API_REFERENCE.md#environment-variables)
- **Docker deployment**: See [Architecture Guide - Build and Deployment](ARCHITECTURE.md#build-and-deployment-architecture)
- **Database setup**: See [Developer Guide - Database Setup](DEVELOPER_GUIDE.md#database-setup-for-knowledge-base)

### Debugging

- **Debug mode**: See [Developer Guide - Debug Mode](DEVELOPER_GUIDE.md#debug-mode)
- **HTTP debugging**: See [Developer Guide - HTTP Debugging](DEVELOPER_GUIDE.md#http-debugging)
- **Common issues**: See [Developer Guide - Common Issues](DEVELOPER_GUIDE.md#common-issues)

---

## 🔍 Search and Navigation Tips

### Finding Information

**By Component Type**:

- Tools → [Architecture Guide - Core Tools](ARCHITECTURE.md#1-core-tools-user-facing-capabilities)
- Model Providers → [Architecture Guide - Model Providers](ARCHITECTURE.md#3-model-providers)
- Auth Providers → [Architecture Guide - Authentication Providers](ARCHITECTURE.md#4-authentication-providers)

**By Task**:

- API Integration → [API Reference](API_REFERENCE.md)
- Development Setup → [Developer Guide - Development Environment Setup](DEVELOPER_GUIDE.md#development-environment-setup)
- Architecture Understanding → [Architecture Guide](ARCHITECTURE.md)
- Creating New Components → [Developer Guide - Creating New Components](DEVELOPER_GUIDE.md#creating-new-components)

**By Technology**:

- Go → [Developer Guide - Go Development](DEVELOPER_GUIDE.md#go-development)
- Python → [Developer Guide - Python Development](DEVELOPER_GUIDE.md#python-development)
- GPTScript → [Architecture Guide - GPTScript Integration](ARCHITECTURE.md#gptscript-integration)
- Docker → [Architecture Guide - Docker Build](ARCHITECTURE.md#build-and-deployment-architecture)

### Cross-References

Documents are heavily cross-referenced:

- Architecture Guide references API specifications
- API Reference links to Architecture for context
- Developer Guide references both for implementation details
- Auth Provider Requirements links to reference implementations

---

## 📚 Additional Resources

### External Links

- **[Obot Platform](https://github.com/obot-platform/obot)** - Main platform repository
- **[GPTScript Documentation](https://docs.gptscript.ai/)** - Tool runtime documentation
- **[Go Documentation](https://go.dev/doc/)** - Go language resources
- **[FastAPI Documentation](https://fastapi.tiangolo.com/)** - Python web framework

### File References

- **[index.yaml](../index.yaml)** - Tool registry
- **[Makefile](../Makefile)** - Build automation
- **[Dockerfile](../Dockerfile)** - Container build configuration
- **[LICENSE](../LICENSE)** - Apache 2.0 license

### Tool Definitions

Each component directory contains a `tool.gpt` file with GPTScript tool definitions:

- Example: [memory/tool.gpt](../memory/tool.gpt)
- Example: [openai-model-provider/tool.gpt](../openai-model-provider/tool.gpt)
- Example: [github-auth-provider/tool.gpt](../github-auth-provider/tool.gpt)

---

## 🛠️ Maintenance and Updates

### Documentation Standards

All documentation follows these principles:

- **Clear structure** with table of contents
- **Cross-referencing** between related sections
- **Code examples** for practical guidance
- **Audience targeting** (beginners, developers, architects)
- **Up-to-date** with latest codebase changes

### Updating Documentation

When making changes:

1. Update relevant documentation files
2. Check cross-references remain valid
3. Update code examples if APIs change
4. Review README.md for high-level changes
5. Update this INDEX.md if structure changes

### Documentation Gaps

If you find missing or unclear documentation:

1. Open an issue in the [main Obot repo](https://github.com/obot-platform/obot/issues) with `tools` and `documentation` labels
2. Include specific sections that need improvement
3. Propose additions or clarifications

---

## 📞 Getting Help

### Channels

- **Issues**: [Obot repository issues](https://github.com/obot-platform/obot/issues) with `tools` label
- **Discussions**: GitHub Discussions in main repository
- **Documentation**: This documentation set

### Issue Templates

When opening an issue:

- **Bug Report**: Include component, steps to reproduce, expected vs actual behavior
- **Feature Request**: Describe use case, proposed solution, alternatives considered
- **Documentation**: Specify unclear sections, propose improvements

---

## 🎓 Learning Path

### Beginner Path

1. ✅ Read [README](../README.md)
2. ✅ Explore [Quick Start](../README.md#quick-start)
3. ✅ Review [Component Overview](../README.md#component-categories)
4. ✅ Study [Architecture Guide - Overview](ARCHITECTURE.md#overview)
5. ✅ Try building a component ([Developer Guide](DEVELOPER_GUIDE.md))

### Intermediate Path

1. ✅ Complete Beginner Path
2. ✅ Study [Architecture Guide](ARCHITECTURE.md) in depth
3. ✅ Review [API Reference](API_REFERENCE.md) for components of interest
4. ✅ Create a simple model provider ([Developer Guide](DEVELOPER_GUIDE.md#creating-a-new-model-provider))
5. ✅ Contribute a bug fix or small feature

### Advanced Path

1. ✅ Complete Intermediate Path
2. ✅ Study [Knowledge Base Deep Dive](ARCHITECTURE.md#knowledge-base-deep-dive)
3. ✅ Review [Security Considerations](ARCHITECTURE.md#security-considerations)
4. ✅ Create a complex component (auth provider, credential store)
5. ✅ Optimize performance or add advanced features

---

## 📊 Documentation Map

```
docs/
├── INDEX.md                    ← You are here
├── ARCHITECTURE.md             ← System design and architecture
├── API_REFERENCE.md            ← Complete API documentation
├── DEVELOPER_GUIDE.md          ← Development and contribution guide
└── auth-providers.md           ← Auth provider specifications

../
├── README.md                   ← Project overview and quick start
├── index.yaml                  ← Tool registry
├── Makefile                    ← Build automation
├── Dockerfile                  ← Container builds
├── [component]/                ← Individual component directories
│   ├── tool.gpt               ← GPTScript tool definition
│   ├── README.md (if present) ← Component-specific docs
│   └── ...
└── ...
```

---

**Last Updated**: January 2026  
**Documentation Version**: 1.0  
**Obot Tools Version**: Latest (main branch)
