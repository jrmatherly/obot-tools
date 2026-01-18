# Obot Tools - Developer Guide

## Table of Contents

- [Getting Started](#getting-started)
- [Development Environment Setup](#development-environment-setup)
- [Building the Project](#building-the-project)
- [Testing](#testing)
- [Creating New Components](#creating-new-components)
- [Code Style and Best Practices](#code-style-and-best-practices)
- [Debugging](#debugging)
- [Contributing](#contributing)

---

## Getting Started

### Prerequisites

**Required**:

- Go 1.23 or later
- Git
- Make

**Optional** (depending on components):

- Python 3.13 or later
- Node.js 18+ and pnpm (for image tools)
- Docker (for containerized builds)
- PostgreSQL or SQLite (for knowledge base and credential stores)

### Clone the Repository

```bash
git clone https://github.com/obot-platform/tools.git
cd tools
```

### Quick Start

```bash
# Build all Go tools
make build

# Run tests
make test
```

---

## Development Environment Setup

### Go Development

1. **Install Go 1.23+**:

   ```bash
   # macOS
   brew install go@1.23
   
   # Linux (download from golang.org)
   wget https://go.dev/dl/go1.23.x.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.23.x.linux-amd64.tar.gz
   ```

2. **Set up Go environment**:

   ```bash
   export GOPATH=$HOME/go
   export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
   ```

3. **Verify installation**:

   ```bash
   go version
   ```

### Python Development

1. **Install Python 3.13+**:

   ```bash
   # macOS
   brew install python@3.13
   
   # Linux
   sudo apt-get install python3.13 python3.13-venv
   ```

2. **Create virtual environment** (for Python tools):

   ```bash
   cd voyage-model-provider
   python3 -m venv venv
   source venv/bin/activate
   pip install -r requirements.txt
   ```

### Node.js Development (for image tools)

1. **Install Node.js and pnpm**:

   ```bash
   # macOS
   brew install node pnpm
   
   # Linux
   curl -fsSL https://get.pnpm.io/install.sh | sh -
   ```

2. **Install dependencies**:

   ```bash
   cd images
   pnpm install
   ```

### Database Setup (for knowledge base)

**PostgreSQL** (recommended for production):

```bash
# macOS
brew install postgresql@15
brew services start postgresql@15

# Create database
createdb knowledge

# Install pgvector extension
psql knowledge -c "CREATE EXTENSION IF NOT EXISTS vector;"
```

**SQLite** (for development):

```bash
# macOS
brew install sqlite

# sqlite-vec extension will be embedded
```

### IDE Setup

**Recommended**: VS Code or GoLand

**VS Code Extensions**:

- Go (golang.go)
- Python (ms-python.python)
- ESLint (dbaeumer.vscode-eslint)
- YAML (redhat.vscode-yaml)

**VS Code Settings** (.vscode/settings.json):

```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "python.linting.enabled": true,
  "python.linting.pylintEnabled": true,
  "editor.formatOnSave": true
}
```

---

## Building the Project

### Build All Tools

```bash
make build
```

This runs `scripts/build.sh`, which:

1. Finds all directories with `main.go`
2. Skips `common` directories
3. Builds each to `bin/gptscript-go-tool` with optimizations (`-ldflags="-s -w"`)

### Build Specific Component

**Go Tool**:

```bash
cd <tool-directory>
go build -ldflags="-s -w" -o bin/gptscript-go-tool .
```

**Python Tool**:

```bash
cd <python-tool>
pip install -r requirements.txt
# No build step; run directly with Python
```

**TypeScript/Node.js Tool**:

```bash
cd images
pnpm install
pnpm build  # If build step exists
```

### Package Tools for Distribution

```bash
# Package all user-facing tools
make package-tools

# Package model and auth providers
make package-providers
```

### Docker Build

```bash
# Build multi-platform Docker image
make docker-build
```

**Multi-stage Dockerfile**:

- **base**: Wolfi base with Go, Python, Node.js, pnpm, uv
- **tools**: Builds all tools
- **providers**: Builds all providers

---

## Testing

### Run All Tests

```bash
make test
```

This runs:

1. Integration tests in `tests/`
2. Unit tests in `github-auth-provider/`
3. Unit tests in `google-auth-provider/`
4. Unit tests in `keycloak-auth-provider/`
5. Unit tests in `entra-auth-provider/`
6. Unit tests in `credential-stores/`

### Run Specific Test Suites

**Integration Tests**:

```bash
cd tests
GPTSCRIPT_TOOL_REMAP="github.com/obot-platform/tools=.." go test -v tool_test.go
```

**Component Tests**:

```bash
cd github-auth-provider
go test ./...
```

**Verbose Output**:

```bash
go test -v ./...
```

**With Coverage**:

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Manual Testing

**Test Model Provider**:

```bash
cd openai-model-provider
export OBOT_OPENAI_MODEL_PROVIDER_API_KEY=sk-...
export PORT=8000
export GPTSCRIPT_DEBUG=true
go run .

# In another terminal
curl http://localhost:8000/v1/models
```

**Test Auth Provider**:

```bash
cd github-auth-provider
export OBOT_GITHUB_AUTH_PROVIDER_CLIENT_ID=...
export OBOT_GITHUB_AUTH_PROVIDER_CLIENT_SECRET=...
export OBOT_AUTH_PROVIDER_COOKIE_SECRET=$(openssl rand -base64 32)
export OBOT_AUTH_PROVIDER_EMAIL_DOMAINS=example.com
export PORT=8000
go run .

# Visit http://localhost:8000/oauth2/start?rd=http://localhost:8000/
```

**Validate Credentials**:

```bash
cd openai-model-provider
export OBOT_OPENAI_MODEL_PROVIDER_API_KEY=sk-...
go run . validate
echo $?  # Should be 0 for success
```

---

## Creating New Components

### Creating a New Model Provider

1. **Create Directory Structure**:

   ```bash
   mkdir -p my-provider-model-provider/proxy
   cd my-provider-model-provider
   ```

2. **Initialize Go Module**:

   ```bash
   go mod init github.com/obot-platform/tools/my-provider-model-provider
   ```

3. **Create tool.gpt**:

   ```gptscript
   Name: MyProvider Model Provider
   Description: Integration with MyProvider API
   Metadata: envVars: OBOT_MYPROVIDER_MODEL_PROVIDER_API_KEY
   Metadata: icon: https://cdn.example.com/icon.svg
   
   #!${GPTSCRIPT_TOOL_DIR}/bin/gptscript-go-tool
   ```

4. **Create main.go**:

   ```go
   package main
   
   import (
       "fmt"
       "net/http"
       "os"
       
       "github.com/obot-platform/tools/openai-model-provider/proxy"
   )
   
   func main() {
       apiKey := os.Getenv("OBOT_MYPROVIDER_MODEL_PROVIDER_API_KEY")
       if apiKey == "" {
           fmt.Println("API key not set")
       }
       
       port := os.Getenv("PORT")
       if port == "" {
           port = "8000"
       }
       
       cfg := &proxy.Config{
           APIKey:               apiKey,
           PersonalAPIKeyHeader: "X-Obot-OBOT_MYPROVIDER_MODEL_PROVIDER_API_KEY",
           ListenPort:           port,
           BaseURL:              "https://api.myprovider.com/v1",
           RewriteModelsFn:      proxy.DefaultRewriteModelsResponse,
           Name:                 "MyProvider",
       }
       
       if len(os.Args) > 1 && os.Args[1] == "validate" {
           if err := cfg.Validate("/validate"); err != nil {
               os.Exit(1)
           }
           return
       }
       
       if err := proxy.Run(cfg); err != nil {
           panic(err)
       }
   }
   ```

5. **Add Dependencies**:

   ```bash
   go get github.com/obot-platform/tools/openai-model-provider/proxy
   go mod tidy
   ```

6. **Register in index.yaml**:

   ```yaml
   modelProviders:
     my-provider-model-provider:
       reference: ./my-provider-model-provider
   ```

7. **Test**:

   ```bash
   export OBOT_MYPROVIDER_MODEL_PROVIDER_API_KEY=your-key
   go run . validate
   go run .
   ```

### Creating a New Tool

1. **Create Directory**:

   ```bash
   mkdir my-tool
   cd my-tool
   ```

2. **Initialize Go Module**:

   ```bash
   go mod init github.com/obot-platform/tools/my-tool
   ```

3. **Create tool.gpt**:

   ```gptscript
   Name: My Tool
   Description: My tool description
   Metadata: icon: https://cdn.example.com/icon.svg
   Metadata: category: Capability
   
   ---
   Name: Do Something
   Description: Perform an action
   Param: input: Input parameter
   
   #!${GPTSCRIPT_TOOL_DIR}/bin/gptscript-go-tool action
   ```

4. **Create main.go**:

   ```go
   package main
   
   import (
       "fmt"
       "os"
   )
   
   func main() {
       if len(os.Args) < 2 {
           fmt.Println("Usage: my-tool <command>")
           os.Exit(1)
       }
       
       command := os.Args[1]
       
       switch command {
       case "action":
           doAction()
       default:
           fmt.Printf("Unknown command: %s\n", command)
           os.Exit(1)
       }
   }
   
   func doAction() {
       input := os.Getenv("INPUT")
       fmt.Printf("Performing action with input: %s\n", input)
   }
   ```

5. **Register in index.yaml**:

   ```yaml
   tools:
     my-tool:
       reference: ./my-tool
   ```

6. **Build and Test**:

   ```bash
   go build -o bin/gptscript-go-tool .
   export INPUT="test"
   ./bin/gptscript-go-tool action
   ```

### Creating a New Auth Provider

Follow the reference implementation in `google-auth-provider/`.

**Key Steps**:

1. Use `auth-providers-common` package
2. Implement required OAuth2 endpoints
3. Reference `placeholder-credential`
4. Add `envVars` metadata
5. Implement user validation logic
6. Test full OAuth2 flow

**Reference**: See `docs/auth-providers.md` for detailed requirements.

---

## Code Style and Best Practices

### Go Code Style

**Formatting**:

```bash
# Format code
go fmt ./...

# Use goimports for import organization
go install golang.org/x/tools/cmd/goimports@latest
goimports -w .
```

**Linting**:

```bash
# Install golangci-lint
brew install golangci-lint

# Run linter
golangci-lint run
```

**Best Practices**:

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use meaningful variable names
- Return errors, don't panic (except in main for fatal errors)
- Use context.Context for cancellation
- Document exported functions and types
- Write table-driven tests

**Example Test**:

```go
func TestDoAction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:  "valid input",
            input: "test",
            want:  "result",
        },
        {
            name:    "invalid input",
            input:   "",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := doAction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("doAction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("doAction() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Python Code Style

**Formatting**:

```bash
# Install formatters
pip install black isort

# Format code
black .
isort .
```

**Linting**:

```bash
# Install linters
pip install pylint flake8 mypy

# Run linters
pylint *.py
flake8 .
mypy .
```

**Best Practices**:

- Follow PEP 8
- Use type hints
- Use async/await for async operations
- Document functions with docstrings
- Use f-strings for formatting
- Handle exceptions properly

### GPTScript Tool Definition

**Best Practices**:

- Clear, concise tool names
- Descriptive descriptions
- Document all parameters
- Include metadata (icon, category, envVars)
- Use context sharing appropriately
- Separate sub-tools logically

**Example**:

```gptscript
Name: My Tool
Description: A tool that does something useful
Share Context: my_context
Metadata: icon: https://cdn.jsdelivr.net/npm/@phosphor-icons/core@2/assets/duotone/icon.svg
Metadata: category: Capability
Metadata: envVars: REQUIRED_ENV_VAR
Metadata: optionalEnvVars: OPTIONAL_ENV_VAR

---
Name: Sub Tool
Description: A specific action
Param: input: The input parameter
Param: optional: An optional parameter

#!${GPTSCRIPT_TOOL_DIR}/bin/gptscript-go-tool subtool
```

### Commit Message Guidelines

Follow [Conventional Commits](https://www.conventionalcommits.org/):

**Format**:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types**:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, no logic change)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples**:

```
feat(openai): add support for GPT-4 Turbo
fix(auth): resolve cookie encryption issue
docs: update API reference for knowledge tool
chore(deps): bump go version to 1.23
```

---

## Debugging

### Debug Mode

Enable debug logging:

```bash
export GPTSCRIPT_DEBUG=true
```

This enables:

- HTTP request/response logging
- Verbose error messages
- Internal state logging

### Go Debugging with Delve

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Start debugger
cd openai-model-provider
dlv debug

# Set breakpoints
(dlv) break main.main
(dlv) break proxy.Run

# Run
(dlv) continue
```

### Python Debugging with pdb

```python
import pdb

# Set breakpoint
pdb.set_trace()

# Or use breakpoint() in Python 3.7+
breakpoint()
```

### HTTP Debugging

**Test with curl**:

```bash
# List models
curl -v http://localhost:8000/v1/models

# Chat completion
curl -v http://localhost:8000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "X-Obot-OBOT_OPENAI_MODEL_PROVIDER_API_KEY: sk-..." \
  -d '{
    "model": "gpt-4",
    "messages": [{"role": "user", "content": "Hello"}]
  }'
```

**Use request/response logging**:

```go
// Add middleware in Go
app.Use(func(c *fiber.Ctx) error {
    log.Printf("Request: %s %s", c.Method(), c.Path())
    log.Printf("Body: %s", c.Body())
    
    err := c.Next()
    
    log.Printf("Response Status: %d", c.Response().StatusCode())
    log.Printf("Response Body: %s", c.Response().Body())
    
    return err
})
```

### Common Issues

**Issue**: `go build` fails with dependency errors  
**Solution**:

```bash
go mod tidy
go mod download
```

**Issue**: Python module not found  
**Solution**:

```bash
source venv/bin/activate
pip install -r requirements.txt
```

**Issue**: Port already in use  
**Solution**:

```bash
# Find process using port
lsof -i :8000

# Kill process
kill -9 <PID>

# Or use different port
export PORT=8001
```

**Issue**: OAuth2 callback fails  
**Solution**:

- Check `OBOT_SERVER_PUBLIC_URL` matches callback URL
- Verify OAuth app callback URL configuration
- Check `OBOT_AUTH_PROVIDER_COOKIE_SECRET` is set
- Ensure redirect URL (`rd` parameter) is properly encoded

---

## Contributing

### Development Workflow

1. **Fork and clone**:

   ```bash
   git clone https://github.com/your-username/tools.git
   cd tools
   git remote add upstream https://github.com/obot-platform/tools.git
   ```

2. **Create feature branch**:

   ```bash
   git checkout -b feature/my-feature
   ```

3. **Make changes**:
   - Write code following style guidelines
   - Add tests
   - Update documentation

4. **Test locally**:

   ```bash
   make test
   make build
   ```

5. **Commit changes**:

   ```bash
   git add .
   git commit -m "feat: add my feature"
   ```

6. **Push and create PR**:

   ```bash
   git push origin feature/my-feature
   # Open PR on GitHub
   ```

### Pull Request Guidelines

**Before submitting**:

- [ ] All tests pass (`make test`)
- [ ] Code builds successfully (`make build`)
- [ ] Code follows style guidelines
- [ ] Documentation updated (if needed)
- [ ] Commit messages follow conventional commits
- [ ] No debugging code or print statements left in

**PR Description**:

- Clear title and description
- Link to related issues
- Screenshots/examples (if applicable)
- Breaking changes noted (if any)

### Code Review Process

1. Automated checks run (tests, linting)
2. Maintainer review
3. Address feedback
4. Approval and merge

### Release Process

Releases are managed by maintainers:

1. Update version numbers
2. Update CHANGELOG
3. Create Git tag
4. Build and publish artifacts
5. Create GitHub release

---

## Additional Resources

- **Main Obot Repo**: https://github.com/obot-platform/obot
- **Issue Tracker**: Use `tools` label in main repo
- **Auth Provider Docs**: `docs/auth-providers.md`
- **Architecture Docs**: `docs/ARCHITECTURE.md`
- **API Reference**: `docs/API_REFERENCE.md`
- **Go Documentation**: https://go.dev/doc/
- **GPTScript Docs**: https://docs.gptscript.ai/
