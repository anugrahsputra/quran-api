# Quran API

A production-ready RESTful API service for accessing Quranic data including surahs (chapters), verses, full-text search functionality (including Tafsir and Topics), and prayer times. Built with Go and the Gin framework, featuring efficient full-text search capabilities using Bleve.

## 📋 Table of Contents

- [Features](#-features)
- [Quick Start](#-quick-start)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [Configuration](#-configuration)
- [API Documentation](#-api-documentation)
- [Development](#-development)
- [Testing](#-testing)
- [Deployment](#-deployment)
- [Security](#-security)
- [Performance](#-performance)
- [Troubleshooting](#-troubleshooting)
- [Contributing](#-contributing)
- [License](#-license)

## ✨ Features

### Core Functionality
- **Surah Management**: Retrieve list of all 114 surahs with metadata
- **Surah Details**: Get detailed surah information with verses, pagination support
- **Advanced Quran Search**: Full-text search across:
  - Translations
  - **Tafsir** (Interpretations)
  - **Thematic Topics** (e.g., "Faith", "Law")
  - Arabic Text and Transliteration
- **Prayer Times**: Get accurate prayer times for any location worldwide

### Production Features
- **Health Checks**: Comprehensive health, liveness, and readiness endpoints
- **Rate Limiting**: IP-based rate limiting (5 requests per 5 minutes per IP)
- **Response Caching**: Intelligent caching for external API responses
- **Security Headers**: Production-ready security headers (XSS, clickjacking protection, etc.)
- **Request Timeouts**: Configurable request timeouts (30 seconds default)
- **Graceful Shutdown**: Clean shutdown handling for zero-downtime deployments
- **Structured Logging**: Comprehensive logging with Zap logger

## 🚀 Quick Start

### Prerequisites

- **Go 1.25+** - [Download](https://golang.org/dl/)
- **Make** (optional) - For using Makefile commands
- **Docker & Docker Compose** (optional) - For containerized deployment

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/anugrahsputra/go-quran-api.git
   cd go-quran-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   # For development
   cp .env.dev.example .env
   
   # For production
   cp .env.example .env
   ```

4. **Index Quran data** (First time setup)
   ```bash
   go run main.go -reindex
   ```
   > ⚠️ **Note**: This process fetches all 114 surahs and their detailed Tafsir. It may take **10-15 minutes** to complete depending on your internet connection.

5. **Run the application**
   ```bash
   go run main.go
   # or
   make run
   ```

   The server will start on `http://localhost:8080` (or the port specified in your `.env` file).

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port number | `8080` | No |
| `ENV` | Environment mode (`development`/`production`) | - | No |
| `GIN_MODE` | Gin framework mode (`debug`/`release`/`test`) | `debug` (dev) / `release` (prod) | No |
| `SEARCH_INDEX_PATH` | Path to Bleve search index directory | `quran.bleve` | No |
| `KEMENAG_API` | Kemenag API base URL | `https://web-api.qurankemenag.net` | No |
| `PRAYER_TIME_API` | Prayer time API base URL | `https://api.aladhan.com/v1` | No |

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

All API endpoints return JSON responses with a consistent structure.

### Health Check Endpoints

#### Health Check
```http
GET /health
```

Returns comprehensive health status including search index status.

#### Liveness Probe
```http
GET /health/live
```

Simple check to verify the service is running.

#### Readiness Probe
```http
GET /health/ready
```

Checks if the service is ready to accept traffic.

### Surah Endpoints

#### Get List of Surahs
```http
GET /api/v1/surah/
```

Returns a list of all 114 surahs with basic information.

#### Get Surah Detail
```http
GET /api/v1/surah/detail/:surah_id?page=1&limit=10
```

Retrieves detailed information about a specific surah including all verses with pagination.

**Path Parameters:**
- `surah_id` (required): Surah ID (1-114)

**Query Parameters:**
- `page` (optional): Page number (default: `1`)
- `limit` (optional): Items per page (default: `10`, max: `100`)

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/surah/detail/1?page=1&limit=10"
```

### Quran Search Endpoint

#### Search Quran
```http
GET /api/v1/search?q=faith&page=1&limit=10
```

Performs full-text search across Quran translations, **Tafsir**, and **Topics** using the Bleve search engine.

**Query Parameters:**
- `q` (required): Search query (searches in Translation, Tafsir, and Topic)
- `page` (optional): Page number (default: `1`)
- `limit` (optional): Items per page (default: `10`, max: `100`)

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/search?q=faith&page=1&limit=10"
```

### Prayer Time Endpoints

#### Get Prayer Times
```http
GET /api/v1/prayer-time?city=Jakarta&country=Indonesia&date=2025-01-15
```

Retrieves prayer times for a specific location and date.

**Query Parameters:**
- `city` (required): City name
- `country` (required): Country name
- `date` (optional): Date in `YYYY-MM-DD` format (default: today)

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/prayer-time?city=Jakarta&country=Indonesia&date=2025-01-15"
```

## 🛠️ Development

### Architecture

The application follows **Clean Architecture** principles with clear separation of concerns:

```
Request Flow:
HTTP Request → Router → Handler → Service → Repository → External API
                                                              ↓
Response Flow:
External API → Repository → Service → Handler → Router → HTTP Response
```

### Adding New Features

See [FEATURE_GUIDE.md](./FEATURE_GUIDE.md) for a detailed guide on adding new features following the existing architecture.

### Makefile Commands

```bash
make build    # Build the application binary
make test     # Run all unit tests
make run      # Run the application
make lint     # Run the linter (requires golangci-lint)
make clean    # Clean build artifacts
make help     # Display help message
```

### Re-indexing Data

To rebuild the search index (e.g., after updating Quran data):

```bash
go run main.go -reindex
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## 🚢 Deployment

### Docker Deployment

```bash
# Build the image
docker build -t quran-api .

# Run the container
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/quran.bleve:/data/quran.bleve \
  -e ENV=production \
  -e GIN_MODE=release \
  -e SEARCH_INDEX_PATH=/data/quran.bleve \
  quran-api
```

## 🔒 Security

- **Security Headers**: X-Frame-Options, X-Content-Type-Options, XSS-Protection, CSP, Referrer-Policy.
- **Rate Limiting**: IP-based rate limiting (5 requests per 5 minutes per IP).
- **Request Timeouts**: 30-second timeout for all requests.

## ⚡ Performance

- **Response Caching**: Intelligent caching for external API responses.
- **Efficient Search**: Bleve full-text search engine.
- **Pagination**: All list endpoints support pagination.
- **Connection Pooling**: HTTP client connection pooling.

## 🐛 Troubleshooting

- **Search returns no results**: Verify `quran.bleve/` exists. Run `go run main.go -reindex`.
- **Server won't start**: Check port availability and `.env` file.
- **Rate limit errors**: Wait a few minutes or adjust limits.

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new features
5. Ensure all tests pass
6. Submit a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- **Quran Data**: Provided by [Kemenag API](https://web-api.qurankemenag.net)
- **Prayer Times**: Provided by [Aladhan API](https://api.aladhan.com)
- **Web Framework**: Built with [Gin](https://gin-gonic.com/)
- **Search Engine**: Full-text search powered by [Bleve](https://blevesearch.com/)
- **Logging**: Structured logging with [Zap](https://github.com/uber-go/zap)

