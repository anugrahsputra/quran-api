# Feature Template

Use this as a starting point when creating a new feature. Copy and modify as needed.

## Feature Name: [Your Feature Name]

### 1. Domain Model (`domain/model/[feature].go`)

```go
package model

type [Feature] struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    // Add your fields
}
```

### 2. DTO (`domain/dto/[feature]_resp.go`)

```go
package dto

type [Feature]Response struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
    Meta    Meta   `json:"meta,omitempty"` // Include if paginated
    Data    interface{} `json:"data"`
}

type [Feature]Data struct {
    // Define your response data structure
}
```

### 3. Repository Interface (`repository/[feature]_repository.go`)

```go
package repository

import (
    "context"
    "github.com/anugrahsputra/go-quran-api/domain/model"
)

type I[Feature]Repository interface {
    Get[Feature](ctx context.Context, id int) (model.[Feature], error)
    // Add more methods as needed
}

type [feature]Repository struct {
    // Add dependencies
}

func New[Feature]Repository(cfg *config.Config) I[Feature]Repository {
    return &[feature]Repository{
        // Initialize
    }
}

func (r *[feature]Repository) Get[Feature](ctx context.Context, id int) (model.[Feature], error) {
    // Implement data fetching
    return model.[Feature]{}, nil
}
```

### 4. Service (`service/[feature]_service.go`)

```go
package service

import (
    "context"
    "github.com/anugrahsputra/go-quran-api/domain/dto"
    "github.com/anugrahsputra/go-quran-api/repository"
)

type I[Feature]Service interface {
    Get[Feature](ctx context.Context, id int) (dto.[Feature]Data, error)
}

type [feature]Service struct {
    repository repository.I[Feature]Repository
}

func New[Feature]Service(r repository.I[Feature]Repository) I[Feature]Service {
    return &[feature]Service{
        repository: r,
    }
}

func (s *[feature]Service) Get[Feature](ctx context.Context, id int) (dto.[Feature]Data, error) {
    // Business logic
    data, err := s.repository.Get[Feature](ctx, id)
    if err != nil {
        return dto.[Feature]Data{}, err
    }
    
    // Transform and return
    return dto.[Feature]Data{}, nil
}
```

### 5. Handler (`handler/[feature]_handler.go`)

```go
package handler

import (
    "net/http"
    "strconv"
    "github.com/anugrahsputra/go-quran-api/domain/dto"
    "github.com/anugrahsputra/go-quran-api/service"
    "github.com/gin-gonic/gin"
)

type [Feature]Handler struct {
    [feature]Service service.I[Feature]Service
}

func New[Feature]Handler([feature]Service service.I[Feature]Service) *[Feature]Handler {
    return &[Feature]Handler{
        [feature]Service: [feature]Service,
    }
}

func (h *[Feature]Handler) Get[Feature](c *gin.Context) {
    // Parse parameters
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Status:  http.StatusBadRequest,
            Message: "invalid id",
        })
        return
    }

    // Call service
    data, err := h.[feature]Service.Get[Feature](c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
            Status:  http.StatusInternalServerError,
            Message: err.Error(),
        })
        return
    }

    // Return response
    c.JSON(http.StatusOK, dto.[Feature]Response{
        Status:  http.StatusOK,
        Message: "success",
        Data:    data,
    })
}
```

### 6. Router (`router/[feature]_route.go`)

```go
package router

import (
    "github.com/anugrahsputra/go-quran-api/handler"
    "github.com/anugrahsputra/go-quran-api/utils/middleware"
    "github.com/gin-gonic/gin"
)

func [Feature]Route(r *gin.RouterGroup, [feature]Handler *handler.[Feature]Handler, rateLimiter *middleware.RateLimiter) {
    [feature]Group := r.Group("/[feature]", rateLimiter.Middleware())
    {
        [feature]Group.GET("/:id", [feature]Handler.Get[Feature])
    }
}
```

### 7. Register in Main Router

Add to `router/main_route.go` in `SetupRoute()`:

```go
// Initialize
[feature]Repo := repository.New[Feature]Repository(cfg)
[feature]Service := service.New[Feature]Service([feature]Repo)
[feature]Handler := handler.New[Feature]Handler([feature]Service)

// Register routes
[Feature]Route(apiV1, [feature]Handler, rateLimiter)
```

## Notes

- Replace `[Feature]` with your actual feature name (PascalCase)
- Replace `[feature]` with your actual feature name (camelCase)
- Add error handling, logging, and validation as needed
- Follow existing code style and patterns
