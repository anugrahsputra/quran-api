# AGENTS.md - Quran API Development Guide

> **READ THIS FIRST**
> This document is the **single source of truth** for any AI agent working on the **Shared** package.
> It defines architectural intent, development constraints, and non-negotiable behavioral rules.

## Philosophy & Guidelines

### Core Philosophy

- **Safety First**
  Never risk user data, stability, or backward compatibility.
  When uncertain, stop and ask for clarification.

- **Incremental Progress**
  Break complex tasks into small, verifiable steps.
  Large, speculative changes are forbidden.

- **Clear Intent Over Cleverness**
  Prefer readable, boring, maintainable solutions.
  Clever hacks are a liability.

- **Native Performance Mindset**
  Optimize only when necessary and with evidence.
  Avoid premature optimization.

---

### Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:

- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

### Simplicity first

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

### Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:

- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:

- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

### Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:

- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:

```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

---

## Project Overview

- **Language**: Go 1.25.1
- **HTTP Framework**: Gin
- **Search Engine**: Bleve
- **Testing**: stretchr/testify
- **Logging**: op/go-logging, uber/zap
- **Module**: `github.com/anugrahsputra/go-quran-api`

---

## Build, Lint, and Test Commands

### Build

```bash
make build          # Build the application binary to tmp/quran-api
make run            # Run the application (go run cmd/main.go)
```

### Dependencies

```bash
make deps           # Download and tidy Go dependencies
```

### Code Quality

```bash
make format         # Format code with gofmt
make vet            # Run go vet
make lint           # Run golangci-lint (requires installation)
make install-linter # Install golangci-lint
```

### Testing

```bash
make test           # Run all unit tests (go test -v ./...)
make test-coverage # Run tests with coverage report (generates coverage.html)
```

#### Running a Single Test

```bash
# Run a specific test function
go test -v -run TestFunctionName ./path/to/package

# Run tests in a specific file
go test -v ./handler/... -run TestPing

# Run tests with verbose output and cover
go test -v -cover ./...
```

### Database/Search Index

```bash
make reindex        # Re-index Quran data for search
```

### Docker

```bash
make docker-build           # Build Docker image
make docker-build-with-index # Build Docker image with search index
make docker-run             # Start services with Docker Compose
make docker-down            # Stop Docker Compose services
```

---

## Code Style Guidelines

### Project Structure

```
.
├── cmd/              # Application entry points
├── common/           # Shared constants and utilities
├── config/           # Configuration loading
├── domain/           # Domain layer
│   ├── dto/          # Data Transfer Objects
│   ├── mapper/       # DTO <-> Model mappers
│   └── model/        # Domain models
├── handler/          # HTTP handlers (Gin)
├── repository/       # Data access layer
├── router/           # Route definitions
├── service/          # Business logic layer
└── utils/            # Utilities and helpers
    ├── helper/       # Helper functions
    └── middleware/   # Gin middleware
```

### Naming Conventions

- **Files**: `snake_case.go` (e.g., `health_handler.go`, `quran_service.go`)
- **Interfaces**: `PascalCase` with `I` prefix for interfaces (e.g., `IQuranService`, `IQuranRepository`)
- **Structs**: `PascalCase` (e.g., `HealthHandler`, `SurahHandler`)
- **Variables/Functions**: `PascalCase` for exported, `camelCase` for unexported
- **Constants**: `PascalCase` or `UPPER_SNAKE_CASE` for exported constants
- **Packages**: `snake_case` (e.g., `handler`, `repository`, `domain/dto`)

### Import Organization

Imports should be organized in three groups with blank lines between them:

1. **Standard library** - `net/http`, `fmt`, `context`, etc.
2. **External packages** - `github.com/gin-gonic/gin`, `github.com/stretchr/testify`, etc.
3. **Internal packages** - `github.com/anugrahsputra/go-quran-api/...`

```go
import (
    "context"
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/op/go-logging"
    "go.uber.org/zap"

    "github.com/anugrahsputra/go-quran-api/domain/dto"
    "github.com/anugrahsputra/go-quran-api/repository"
    "github.com/anugrahsputra/go-quran-api/service"
)
```

### Types and Go Patterns

- Use **interfaces** for dependencies (repository interfaces, service interfaces)
- Use **struct tags** for JSON serialization (e.g., `json:"status"`)
- Prefer **context.Context** as first parameter for service/repository methods
- Use **pointers** for struct receivers when modifying state, value receivers for read-only

```go
// Interface definition
type IQuranService interface {
    GetListSurah(ctx context.Context) ([]dto.SurahResp, error)
    GetSurahDetail(ctx context.Context, id int, page int, limit int) (dto.SurahDetailData, int, int, error)
}

// Struct with interface field
type SurahHandler struct {
    quranService service.IQuranService
}

// Constructor returns interface
func NewSurahHandler(surahService service.IQuranService) *SurahHandler {
    return &SurahHandler{
        quranService: surahService,
    }
}
```

### Error Handling

- Return errors from service/repository layers
- Handle errors in handlers with appropriate HTTP status codes
- Use `helper.SanitizeError()` for production-safe error messages
- Use structured logging with context (`zap` for structured logging, `op/go-logging` for simple logging)

```go
func (h *SurahHandler) GetListSurah(c *gin.Context) {
    response, err := h.quranService.GetListSurah(c.Request.Context())
    if err != nil {
        logger.Errorf("HTTP request failed - Error: %s", err.Error())
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
            Status:  http.StatusInternalServerError,
            Message: helper.SanitizeError(err),
        })
        return
    }

    c.JSON(http.StatusOK, dto.SurahListResp{
        Status:  http.StatusOK,
        Message: "success",
        Data:    response,
    })
}
```

### Logging

- Use `op/go-logging` for simple logging in handlers and repositories
- Use `uber/zap` for structured logging in services when additional context is needed
- Log HTTP requests in handlers with method, path, client IP, and user agent
- Log errors with sufficient context for debugging

```go
var logger = logging.MustGetLogger("handler")

logger.Infof("HTTP request received - Method: %s, Path: %s", c.Request.Method, c.Request.URL.Path)
logger.Errorf("HTTP request failed - Error: %s", err.Error())
```

### Response DTOs

Use consistent response structure:

```go
// Success response
type SurahListResp struct {
    Status  int         `json:"status"`
    Message string      `json:"message"`
    Data    []SurahResp `json:"data"`
}

// Error response
type ErrorResponse struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
}
```

### Middleware Patterns

Middleware functions should:

- Accept `gin.HandlerFunc`
- Process request/response
- Call `c.Next()` to continue chain

```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Next()
    }
}
```

### Testing

- Use `stretchr/testify` with `assert` package
- Use `httptest` for HTTP handler testing
- Create mock structs implementing interfaces for testing

```go
import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestPing(t *testing.T) {
    gin.SetMode(gin.TestMode)

    mockRepo := new(MockQuranSearchRepository)
    h := NewHealthHandler(mockRepo)

    r := gin.Default()
    r.GET("/ping", h.Ping)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/ping", nil)
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.JSONEq(t, `{"message":"pong"}`, w.Body.String())
}
```

### Configuration

- Use environment variables with `.env` files for local development
- Use `config.LoadConfig()` to load configuration
- Never commit secrets - use `.env.example` as template

---

## Linting Configuration

The project uses `.golangci.yml` with these linters enabled:

- govet
- errcheck
- staticcheck
- unused
- gofmt
- goimports

Run linting with: `make lint`

---

## Additional Notes

- The project uses Bleve for full-text search indexing
- Quran data is fetched from external API (Kemenag)
- Health checks are available at `/api/v1/health`, `/api/v1/readiness`, `/api/v1/liveness`
- Search indexing can be triggered with `-reindex` flag or `AUTO_INDEX=true` env var
