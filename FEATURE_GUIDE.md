# Feature Development Guide

This guide explains how to add new features to the Quran API following the existing architecture.

## Architecture Overview

The application follows a **Clean Architecture** pattern with clear separation of concerns:

`Router` -> `Handler` -> `Service` -> `Repository` -> `External API / Search Index`

## Directory Structure

```
├── domain/
│   ├── dto/          # Data Transfer Objects (API responses)
│   └── model/        # Domain models (internal data structures)
├── handler/          # HTTP handlers (request/response handling)
├── service/          # Business logic layer
├── repository/       # Data access layer (external APIs, databases)
├── router/           # Route definitions
└── utils/            # Utilities and middleware
```

## Step-by-Step: Adding a New Feature

### Example: Adding a "Juz" (Chapter) Feature

Let's say you want to add an endpoint to get verses by Juz number.

### Step 1: Define Domain Models (`domain/model/`)

Create `domain/model/juz.go`:

```go
package model

type Juz struct {
    ID          int    `json:"id"`
    StartSurah int    `json:"start_surah"`
    StartAyah  int    `json:"start_ayah"`
    EndSurah   int    `json:"end_surah"`
    EndAyah    int    `json:"end_ayah"`
}
```

### Step 2: Define DTOs (`domain/dto/`)

Create `domain/dto/juz_resp.go`:

```go
package dto

import "github.com/anugrahsputra/go-quran-api/domain/model"

type JuzResponse struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
    Meta    Meta   `json:"meta"`
    Data    JuzData `json:"data"`
}

type JuzData struct {
    JuzID      int           `json:"juz_id"`
    Verses     []model.Verse `json:"verses"`
}
```

### Step 3: Create Repository (`repository/`)

Create `repository/juz_repository.go`:

```go
package repository

import (
    "context"
    "github.com/anugrahsputra/go-quran-api/domain/model"
)

type IJuzRepository interface {
    GetJuzVerses(ctx context.Context, juzID int, page, limit int) (model.JuzApi, error)
}

type juzRepository struct {
    // Add dependencies (HTTP client, cache, etc.)
}

func NewJuzRepository(cfg *config.Config) IJuzRepository {
    return &juzRepository{
        // Initialize dependencies
    }
}

func (r *juzRepository) GetJuzVerses(ctx context.Context, juzID int, page, limit int) (model.JuzApi, error) {
    // Implement data fetching logic
    // Call external API, query database, etc.
}
```

### Step 4: Create Service (`service/`)

Create `service/juz_service.go`:

```go
package service

import (
    "context"
    "github.com/anugrahsputra/go-quran-api/domain/dto"
    "github.com/anugrahsputra/go-quran-api/repository"
)

type IJuzService interface {
    GetJuzVerses(ctx context.Context, juzID int, page, limit int) (dto.JuzData, int, int, error)
}

type juzService struct {
    repository repository.IJuzRepository
}

func NewJuzService(r repository.IJuzRepository) IJuzService {
    return &juzService{
        repository: r,
    }
}

func (s *juzService) GetJuzVerses(ctx context.Context, juzID int, page, limit int) (dto.JuzData, int, int, error) {
    // Business logic:
    // 1. Validate inputs
    // 2. Call repository
    // 3. Transform data
    // 4. Calculate pagination
    // 5. Return result
}
```

### Step 5: Create Handler (`handler/`)

Create `handler/juz_handler.go`:

```go
package handler

import (
    "net/http"
    "strconv"
    "github.com/anugrahsputra/go-quran-api/domain/dto"
    "github.com/anugrahsputra/go-quran-api/service"
    "github.com/gin-gonic/gin"
)

type JuzHandler struct {
    juzService service.IJuzService
}

func NewJuzHandler(juzService service.IJuzService) *JuzHandler {
    return &JuzHandler{
        juzService: juzService,
    }
}

func (h *JuzHandler) GetJuzVerses(c *gin.Context) {
    // 1. Parse parameters
    juzIDStr := c.Param("juz_id")
    juzID, err := strconv.Atoi(juzIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Status:  http.StatusBadRequest,
            Message: "invalid juz_id",
        })
        return
    }

    // 2. Parse pagination
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    // 3. Validate
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }

    // 4. Call service
    data, total, totalPages, err := h.juzService.GetJuzVerses(c.Request.Context(), juzID, page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
            Status:  http.StatusInternalServerError,
            Message: err.Error(),
        })
        return
    }

    // 5. Return response
    c.JSON(http.StatusOK, dto.JuzResponse{
        Status:  http.StatusOK,
        Message: "success",
        Meta: dto.Meta{
            Total:      total,
            Page:       page,
            Limit:      limit,
            TotalPages: totalPages,
        },
        Data: data,
    })
}
```

### Step 6: Create Router (`router/`)

Create `router/juz_route.go`:

```go
package router

import (
    "github.com/anugrahsputra/go-quran-api/handler"
    "github.com/anugrahsputra/go-quran-api/utils/middleware"
    "github.com/gin-gonic/gin"
)

func JuzRoute(r *gin.RouterGroup, juzHandler *handler.JuzHandler, rateLimiter *middleware.RateLimiter) {
    juzGroup := r.Group("/juz", rateLimiter.Middleware())
    {
        juzGroup.GET("/:juz_id", juzHandler.GetJuzVerses)
    }
}
```

### Step 7: Register in Main Router (`router/main_route.go`)

Add to `SetupRoute()`:

```go
// After other routes...
juzRepo := repository.NewJuzRepository(cfg)
juzService := service.NewJuzService(juzRepo)
juzHandler := handler.NewJuzHandler(juzService)
JuzRoute(apiV1, juzHandler, rateLimiter)
```

## Quick Checklist

When adding a new feature:

- [ ] **Domain Model** - Define data structures in `domain/model/`
- [ ] **DTO** - Define API response structures in `domain/dto/`
- [ ] **Repository** - Create data access layer in `repository/`
- [ ] **Service** - Create business logic in `service/`
- [ ] **Handler** - Create HTTP handler in `handler/`
- [ ] **Router** - Create route definitions in `router/`
- [ ] **Register** - Add routes to `router/main_route.go`
- [ ] **Tests** - Add unit tests (optional but recommended)

## Common Patterns

### Error Handling

```go
// In Handler
if err != nil {
    c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
        Status:  http.StatusInternalServerError,
        Message: err.Error(),
    })
    return
}
```

### Pagination

```go
// Parse
page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

// Validate
if page < 1 {
    page = 1
}
if limit < 1 || limit > 100 {
    limit = 10
}

// Calculate
totalPages := (total + limit - 1) / limit
```

### Logging

```go
logger.Infof("Processing request: %s %s", c.Request.Method, c.Request.URL.Path)
```

### Context Usage

Always pass `context.Context` from the handler through service to repository:

```go
// Handler
ctx := c.Request.Context()

// Service
func (s *service) Method(ctx context.Context, ...) { }

// Repository
func (r *repo) Method(ctx context.Context, ...) { }
```

## Testing Your Feature

1. **Start the server:**
   ```bash
   go run main.go
   ```

2. **Test the endpoint:**
   ```bash
   curl http://localhost:8080/api/v1/juz/1?page=1&limit=10
   ```

3. **Check logs** for any errors

## Tips

1. **Follow existing patterns** - Look at `surah_handler.go`, `search_service.go`, etc. for examples
2. **Keep it simple** - Start with basic functionality, add complexity later
3. **Use interfaces** - Makes testing easier
4. **Handle errors gracefully** - Always return proper HTTP status codes
5. **Add logging** - Helps with debugging
6. **Validate inputs** - Check parameters before processing

## Need Help?

- Check existing implementations in `handler/`, `service/`, `repository/`
- Look at test files (e.g., `surah_handler_test.go`) for examples
- Review the router setup in `router/main_route.go`
