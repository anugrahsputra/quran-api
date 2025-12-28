# Shared Package Agent Documentation

> **READ THIS FIRST**
> This document is the **single source of truth** for any AI agent working on the **Shared** package.
> It defines architectural intent, development constraints, and non-negotiable behavioral rules.

---

## 1. Philosophy & Guidelines

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

### Eight Honors and Eight Shames

- **Shame** in guessing APIs, **Honor** in careful research
- **Shame** in vague execution, **Honor** in seeking confirmation
- **Shame** in assuming business logic, **Honor** in human verification
- **Shame** in creating new abstractions, **Honor** in reusing existing ones
- **Shame** in skipping validation, **Honor** in proactive testing
- **Shame** in breaking architecture, **Honor** in following specifications
- **Shame** in pretending to understand, **Honor** in honest ignorance
- **Shame** in blind modification, **Honor** in careful refactoring

---

## 2. Quality Standards

### Language & Style

- **English Only**
  All code, comments, documentation, identifiers, and error messages must be written in English.

- **No Unnecessary Comments**
  Do not comment _what_ the code does.
  Comments are allowed **only** to explain _why_ a decision exists.

---

## 3. Project Identity

- **Name**: Quran API
- **Description**: A production-ready RESTful API service for accessing Quranic data (Surahs, Verses, Tafsir), performing advanced full-text search, and retrieving prayer times.
- **Core Capabilities**:
  - **Surah Management**: List and detail views with pagination.
  - **Advanced Search**: Full-text search across Translations, Tafsir, and Topics using Bleve.
  - **Prayer Times**: Accurate prayer times based on location and date.
  - **Production Ready**: Includes health checks, rate limiting, security headers, and graceful shutdown.

---

## 4. Technology Stack

- **Language**: Go 1.25+
- **Web Framework**: Gin (v1)
- **Search Engine**: Bleve (v2)
- **Logging**: Zap
- **Configuration**: godotenv
- **Testing**: Testify
- **Containerization**: Docker

---

## 5. Repository Architecture

### Architectural Pattern

The project follows **Clean Architecture** principles, enforcing a strict separation of concerns and dependency flow:

`Router` -> `Handler` -> `Service` -> `Repository` -> `External API / Search Index`

### Directory Structure

- **`cmd/`**: Application entry point (`main.go`).
- **`config/`**: Configuration loading and logger setup.
- **`domain/`**:
  - **`model/`**: Internal domain entities and business logic structures.
  - **`dto/`**: Data Transfer Objects for API request/response definitions.
  - **`mapper/`**: Logic for transforming models to DTOs.
- **`handler/`**: HTTP controllers responsible for parsing requests and validating input.
- **`service/`**: Business logic implementation. Orchestrates data flow between handlers and repositories.
- **`repository/`**: Data access layer. Handles interactions with external APIs (Kemenag, Aladhan) and the Bleve search index.
- **`router/`**: Route definitions, middleware registration, and handler wiring.
- **`utils/`**:
  - **`helper/`**: Shared utility functions (Error handling, Caching).
  - **`middleware/`**: Cross-cutting concerns (CORS, Rate Limiting, Recovery, Security).

---

## 6. Key Workflows

### Development

1.  **Setup Environment**:
    - Copy `.env.dev.example` to `.env` for local development.
    - Run `make deps` to install Go dependencies.

2.  **Code Quality**:
    - Run `make format` to enforce standard Go formatting.
    - Run `make vet` to catch common errors.
    - Run `make lint` to check for stylistic and logical issues (requires `golangci-lint`).

3.  **Testing**:
    - Run `make test` to execute all unit tests.
    - Run `make test-coverage` to generate an HTML coverage report.
    - **Rule**: No feature is complete without passing tests.

4.  **Running Locally**:
    - Run `make run` to start the server on the configured port.
    - **Requirement**: A valid `quran.bleve` index must exist. If missing, run `make reindex` (Note: This takes 10-15 mins).

### Building

1.  **Binary Build**:
    - Run `make build` to compile the application to `tmp/quran-api`.
    - Verify the binary executes correctly by running it with `--help` or checking version output if available.

2.  **Docker Build**:
    - Standard: `make docker-build` (builds image `quran-api:latest`).
    - With Index: `make docker-build-with-index` (builds `quran-api:with-index` including the pre-built search index).

### Release

1.  **Pre-Release Verification**:
    - Ensure `make clean-all` followed by `make test` passes.
    - Verify `README.md` and API documentation matches implementation.

2.  **Artifact Generation**:
    - For containerized environments, prefer `docker-build-with-index` to bundle the search data, ensuring startup speed and consistency.

3.  **Deployment Configuration**:
    - Ensure `ENV=production` and `GIN_MODE=release` are set in the target environment.
    - Verify `SEARCH_INDEX_PATH` points to a valid, writable location if not using the bundled index.

---

## 7. Implementation Details

### Error Handling

- **Centralized Helper**: `utils/helper/error_helper.go`
- **Sanitization**:
  - **Production**: Masks internal errors (e.g., "An internal error occurred"). Returns generic messages for timeouts or upstream failures.
  - **Development**: Returns raw error strings for debugging.

### External APIs

- **Quran Data**: `quran.kemenag.go.id` (Text, Translation, Tafsir).
- **Prayer Times**: `api.aladhan.com` (Timings).
- **Resilience**: HTTP Clients configured with 10s timeouts.

### Search Configuration

- **Engine**: Bleve v2.
- **Index Storage**: File-system based (path defined in config).
- **Indexed Fields**: `SurahNumber`, `AyahNumber`, `Text` (Arabic), `Latin`, `Translation`, `Tafsir`, `Topic`.
