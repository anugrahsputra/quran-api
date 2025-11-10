# Quran API

A production-ready RESTful API service for accessing Quranic data including surahs (chapters), verses, full-text search functionality, and prayer times. Built with Go and the Gin framework, featuring efficient full-text search capabilities using Bleve.

## 📋 Table of Contents

- [Features](#-features)
- [Quick Start](#-quick-start)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [Configuration](#-configuration)
- [API Documentation](#-api-documentation)
- [Project Structure](#-project-structure)
- [Development](#-development)
- [Testing](#-testing)
- [Deployment](#-deployment)
- [Architecture](#-architecture)
- [Security](#-security)
- [Performance](#-performance)
- [Troubleshooting](#-troubleshooting)
- [Contributing](#-contributing)
- [License](#-license)

## ✨ Features

### Core Functionality
- **Surah Management**: Retrieve list of all 114 surahs with metadata
- **Surah Details**: Get detailed surah information with verses, pagination support
- **Full-Text Search**: Search across Quran translations with advanced querying
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
   > ⚠️ **Note**: This process may take several minutes as it fetches and indexes all 114 surahs (~6,236 verses).

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

### Example Configuration Files

**Development** (`.env.dev.example`):
```env
PORT=8080
ENV=development
GIN_MODE=debug
SEARCH_INDEX_PATH=quran.bleve
KEMENAG_API=https://web-api.qurankemenag.net
PRAYER_TIME_API=https://api.aladhan.com/v1
```

**Production** (`.env.example`):
```env
PORT=8080
ENV=production
GIN_MODE=release
SEARCH_INDEX_PATH=/data/quran.bleve
KEMENAG_API=https://web-api.qurankemenag.net
PRAYER_TIME_API=https://api.aladhan.com/v1
```

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

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-01-15T10:30:00Z",
  "checks": {
    "search_index": {
      "status": "healthy",
      "message": "Search index is accessible (indexed documents available)",
      "response_time": "1.2ms"
    }
  }
}
```

**Status Codes:**
- `200 OK` - Service is healthy
- `503 Service Unavailable` - Service is degraded

#### Liveness Probe
```http
GET /health/live
```

Simple check to verify the service is running. Always returns `200 OK` if the service is alive.

**Response:**
```json
{
  "status": "alive",
  "timestamp": "2025-01-15T10:30:00Z",
  "checks": {}
}
```

#### Readiness Probe
```http
GET /health/ready
```

Checks if the service is ready to accept traffic. Verifies all dependencies are available.

**Response:**
```json
{
  "status": 200,
  "message": "ready",
  "data": {
    "timestamp": "2025-01-15T10:30:00Z",
    "checks": {
      "search_index": {
        "status": "healthy",
        "message": "Search index is accessible (indexed documents available)",
        "response_time": "1.2ms"
      }
    }
  }
}
```

> 💡 **Tip**: These endpoints are not rate-limited and are perfect for Kubernetes health checks.

### Surah Endpoints

#### Get List of Surahs
```http
GET /api/v1/surah/
```

Returns a list of all 114 surahs with basic information.

**Response:**
```json
{
  "status": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "arabic": "الفاتحة",
      "latin": "Al-Fatihah",
      "translation": "Pembukaan",
      "transliteration": "Al-Fatihah",
      "location": "Makkah",
      "num_ayah": 7
    },
    {
      "id": 2,
      "arabic": "البقرة",
      "latin": "Al-Baqarah",
      "translation": "Sapi Betina",
      "transliteration": "Al-Baqarah",
      "location": "Madinah",
      "num_ayah": 286
    }
  ]
}
```

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

**Response:**
```json
{
  "status": 200,
  "message": "success",
  "meta": {
    "total": 7,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  },
  "data": {
    "surah_id": 1,
    "arabic": "الفاتحة",
    "latin": "Al-Fatihah",
    "translation": "Pembukaan",
    "transliteration": "Al-Fatihah",
    "location": "Makkah",
    "audio": "https://api.qurankemenag.net/audio/1",
    "verses": [
      {
        "id": 1,
        "ayah": 1,
        "page": 1,
        "quarter_hizb": 0.5,
        "juz": 1,
        "manzil": 1,
        "arabic": "بِسْمِ اللّٰهِ الرَّحْمٰنِ الرَّحِيْمِ",
        "kitabah": "...",
        "latin": "Bismillahirrahmanirrahim",
        "translation": "Dengan nama Allah Yang Maha Pengasih, Maha Penyayang.",
        "audio": "..."
      }
    ]
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid surah_id or pagination parameters
- `404 Not Found` - Surah not found
- `500 Internal Server Error` - Server error

### Search Endpoint

#### Search Verses
```http
GET /api/v1/search?q=allah&page=1&limit=10
```

Performs full-text search across Quran translations using the Bleve search engine.

**Query Parameters:**
- `q` (required): Search query (searches in translation text)
- `page` (optional): Page number (default: `1`)
- `limit` (optional): Items per page (default: `10`, max: `100`)

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/search?q=allah&page=1&limit=10"
```

**Response:**
```json
{
  "code": 200,
  "status": "OK",
  "message": "Success",
  "meta": {
    "total": 286,
    "page": 1,
    "limit": 10,
    "total_pages": 29
  },
  "data": [
    {
      "surah_number": 2,
      "ayah_number": 255,
      "text": "اللّٰهُ لَآ إِلٰهَ إِلَّا هُوَ الْحَيُّ الْقَيُّومُ",
      "latin": "Allahu la ilaha illa Huwa, al-Hayyul Qayyum",
      "translation": "Allah! There is no god but He, the Living, the Self-Subsisting, Eternal."
    }
  ]
}
```

**Error Responses:**
- `400 Bad Request` - Missing or invalid query parameter
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Search index error

> ⚠️ **Note**: Search functionality requires the index to be built. Run `go run main.go -reindex` if you haven't already.

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

**Response:**
```json
{
  "status": 200,
  "message": "success",
  "data": {
    "city": "Jakarta",
    "country": "Indonesia",
    "date": "2025-01-15",
    "times": {
      "fajr": "04:30",
      "sunrise": "05:45",
      "dhuhr": "12:00",
      "asr": "15:15",
      "maghrib": "18:20",
      "isha": "19:35"
    }
  }
}
```

**Error Responses:**
- `400 Bad Request` - Missing required parameters
- `404 Not Found` - Location not found
- `500 Internal Server Error` - External API error

## 📁 Project Structure

```
quran-api/
├── common/                    # Common utilities and test helpers
│   └── test_helper.go
├── config/                    # Configuration and logging setup
│   ├── config.go             # Configuration loader
│   └── logger.go             # Logger initialization
├── domain/                    # Domain layer (business entities)
│   ├── dto/                  # Data Transfer Objects (API responses)
│   │   ├── detail_surah_resp.go
│   │   ├── prayer_time_resp.go
│   │   ├── response.go
│   │   ├── search_resp.go
│   │   └── surah_list_resp.go
│   └── model/                 # Domain models
│       ├── ayah.go
│       ├── detail_surah.go
│       ├── prayer_time.go
│       └── surah.go
├── handler/                   # HTTP handlers (request/response handling)
│   ├── health_handler.go
│   ├── prayer_time_handler.go
│   ├── search_handler.go
│   ├── surah_handler.go
│   └── surah_handler_test.go
├── repository/                # Data access layer
│   ├── prayer_time_repository.go
│   ├── quran_repository.go
│   └── search_repository.go
├── router/                    # Route definitions
│   ├── detail_surah_route.go
│   ├── health_route.go
│   ├── main_route.go
│   ├── prayer_time_route.go
│   ├── search_route.go
│   └── surah_route.go
├── service/                   # Business logic layer
│   ├── prayer_time_service.go
│   ├── search_service.go
│   ├── surah_service.go
│   └── surah_service_test.go
├── utils/                     # Utilities and middleware
│   ├── helper/               # Helper functions
│   │   ├── cache_helper.go
│   │   └── helper.go
│   └── middleware/           # HTTP middlewares
│       ├── cors_middleware.go
│       ├── ip_rate_limiter_middleware.go
│       ├── security_middleware.go
│       └── timeout_middleware.go
├── main.go                    # Application entry point
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── Dockerfile                 # Docker image definition
├── docker-compose.yml         # Docker Compose configuration
├── Makefile                   # Build automation
├── .env.example               # Production environment template
├── .env.dev.example           # Development environment template
├── FEATURE_GUIDE.md           # Guide for adding new features
├── TEMPLATE_FEATURE.md        # Feature template
└── README.md                  # This file
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

**Layer Responsibilities:**
- **Router**: Route definitions and middleware registration
- **Handler**: HTTP request/response handling, parameter validation
- **Service**: Business logic, data transformation, pagination
- **Repository**: Data access, external API calls, caching

### Adding New Features

See [FEATURE_GUIDE.md](./FEATURE_GUIDE.md) for a detailed guide on adding new features following the existing architecture.

**Quick Steps:**
1. Define domain models in `domain/model/`
2. Create DTOs in `domain/dto/`
3. Implement repository in `repository/`
4. Add business logic in `service/`
5. Create handler in `handler/`
6. Define routes in `router/`
7. Register routes in `router/main_route.go`

### Makefile Commands

```bash
make build    # Build the application binary
make test     # Run all unit tests
make run      # Run the application
make lint     # Run the linter (requires golangci-lint)
make clean    # Clean build artifacts
make help     # Display help message
```

### Code Style

- Follow Go standard formatting: `gofmt -w .`
- Use `golangci-lint` for linting
- Follow existing code patterns and conventions
- Write tests for new features

### Re-indexing Data

To rebuild the search index (e.g., after updating Quran data):

```bash
go run main.go -reindex
```

This will:
1. Fetch all surahs from the Kemenag API
2. Index all verses in the Bleve search engine
3. Store the index in the path specified by `SEARCH_INDEX_PATH`

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests in verbose mode
go test -v ./...

# Run specific test
go test -v ./handler -run TestSurahHandler
```

### Test Structure

Tests are located alongside the source files with the `_test.go` suffix. The project uses the standard Go testing package and `testify` for assertions.

## 🚢 Deployment

### Docker Deployment

#### Build and Run

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

#### Using Docker Compose

```bash
# Start the service
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the service
docker-compose down
```

### Production Considerations

1. **Environment Variables**: Set all required environment variables
2. **Index Path**: Use a persistent volume for the search index (`/data/quran.bleve`)
3. **Health Checks**: Configure health checks in your orchestration platform
   - Liveness: `GET /health/live`
   - Readiness: `GET /health/ready`
4. **Rate Limiting**: Adjust rate limits based on your needs (modify `router/main_route.go`)
5. **Monitoring**: Set up logging and monitoring (Prometheus, Grafana, etc.)
6. **SSL/TLS**: Use a reverse proxy (nginx, Traefik, Caddy) for HTTPS
7. **Graceful Shutdown**: The application supports graceful shutdown (30-second timeout)

### Kubernetes Deployment

Example Kubernetes configuration:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: quran-api
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: quran-api
        image: quran-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: ENV
          value: "production"
        - name: GIN_MODE
          value: "release"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        volumeMounts:
        - name: search-index
          mountPath: /data/quran.bleve
      volumes:
      - name: search-index
        persistentVolumeClaim:
          claimName: quran-index-pvc
```

## 🔒 Security

### Security Features

- **Security Headers**: 
  - `X-Frame-Options: DENY` - Prevents clickjacking
  - `X-Content-Type-Options: nosniff` - Prevents MIME type sniffing
  - `X-XSS-Protection: 1; mode=block` - XSS protection
  - `Content-Security-Policy: default-src 'self'` - CSP
  - `Referrer-Policy: strict-origin-when-cross-origin` - Referrer policy

- **Rate Limiting**: IP-based rate limiting (5 requests per 5 minutes per IP)
- **Request Timeouts**: 30-second timeout for all requests
- **Input Validation**: Parameter validation in handlers
- **Error Handling**: Graceful error handling without exposing internals

### Security Best Practices

1. **HTTPS**: Always use HTTPS in production (use a reverse proxy)
2. **Environment Variables**: Never commit `.env` files
3. **Rate Limiting**: Adjust rate limits based on your threat model
4. **Monitoring**: Monitor for suspicious activity
5. **Updates**: Keep dependencies updated

## ⚡ Performance

### Performance Features

- **Response Caching**: Intelligent caching for external API responses
- **Efficient Search**: Bleve full-text search engine
- **Pagination**: All list endpoints support pagination
- **Connection Pooling**: HTTP client connection pooling
- **Graceful Shutdown**: Zero-downtime deployments

### Performance Tips

1. **Index Location**: Store the search index on fast storage (SSD)
2. **Caching**: External API responses are cached to reduce load
3. **Pagination**: Use appropriate page sizes (10-50 items)
4. **Connection Pooling**: HTTP clients reuse connections

## 🐛 Troubleshooting

### Search returns no results

1. **Check if index exists**: Verify `quran.bleve/` directory exists
2. **Re-index the data**: Run `go run main.go -reindex`
3. **Check index health**: `GET /health` endpoint
4. **Verify permissions**: Ensure the application has read/write access to the index directory

### Server won't start

1. **Check port availability**: `lsof -i :8080` (macOS/Linux) or `netstat -ano | findstr :8080` (Windows)
2. **Verify environment variables**: Check `.env` file is present and correctly formatted
3. **Check logs**: Review application logs for error messages
4. **Verify Go version**: Ensure Go 1.25+ is installed (`go version`)

### Indexing fails

1. **Check internet connection**: Requires external API access to `KEMENAG_API`
2. **Verify API accessibility**: Test `https://web-api.qurankemenag.net` is reachable
3. **Check disk space**: Ensure sufficient disk space for the index
4. **Review logs**: Check for specific error messages in the logs
5. **Check permissions**: Ensure write permissions for the index directory

### Rate limit errors

- **429 Too Many Requests**: You've exceeded the rate limit (5 requests per 5 minutes)
- **Solution**: Wait a few minutes or adjust rate limits in `router/main_route.go`

### Health check failures

- **503 Service Unavailable**: One or more dependencies are unhealthy
- **Check**: Review the `/health` endpoint response for specific issues
- **Common causes**: Search index not accessible, external API down

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Make your changes**
4. **Add tests** for new features
5. **Ensure all tests pass** (`make test`)
6. **Commit your changes** (`git commit -m 'Add some amazing feature'`)
7. **Push to the branch** (`git push origin feature/amazing-feature`)
8. **Open a Pull Request**

### Development Guidelines

- Follow the existing code structure and patterns
- Add tests for new features
- Update documentation as needed
- Ensure all tests pass before submitting
- Follow Go code style guidelines
- Write clear commit messages

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- **Quran Data**: Provided by [Kemenag API](https://web-api.qurankemenag.net)
- **Prayer Times**: Provided by [Aladhan API](https://api.aladhan.com)
- **Web Framework**: Built with [Gin](https://gin-gonic.com/)
- **Search Engine**: Full-text search powered by [Bleve](https://blevesearch.com/)
- **Logging**: Structured logging with [Zap](https://github.com/uber-go/zap)

## 📞 Support

For issues, questions, or contributions:
- **GitHub Issues**: [Open an issue](https://github.com/anugrahsputra/go-quran-api/issues)
- **Documentation**: See [FEATURE_GUIDE.md](./FEATURE_GUIDE.md) for development guide

---

**Made with ❤️ for the Muslim community**
