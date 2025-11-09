# Quran API

A RESTful API service for accessing Quranic data including surahs (chapters), verses, search functionality, and prayer times. Built with Go and Gin framework, featuring full-text search capabilities using Bleve.

## 🚀 Features

- **Surah Management**: Get list of all surahs and detailed surah information
- **Verse Search**: Full-text search across Quran translations with pagination
- **Prayer Times**: Get prayer times for any location
- **Health Checks**: Production-ready health, liveness, and readiness endpoints
- **Rate Limiting**: IP-based rate limiting to prevent abuse
- **Caching**: Response caching for improved performance
- **Production Ready**: Graceful shutdown, security headers, request timeouts

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Project Structure](#project-structure)
- [Development](#development)
- [Deployment](#deployment)
- [Contributing](#contributing)

## Prerequisites

- Go 1.25 or higher
- Make (optional, for using Makefile commands)
- Docker & Docker Compose (optional, for containerized deployment)

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/anugrahsputra/go-quran-api.git
cd go-quran-api
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Set up environment variables

For development:
```bash
cp .env.dev.example .env
```

For production:
```bash
cp .env.example .env
```

Edit `.env` file with your configuration (see [Configuration](#configuration) section).

### 4. Index Quran data (First time setup)

Before using the search functionality, you need to index the Quran data:

```bash
go run main.go -reindex
```

This process may take several minutes as it fetches and indexes all 114 surahs (~6,236 verses).

### 5. Run the application

```bash
go run main.go
```

The server will start on `http://localhost:8080` (or the port specified in your `.env` file).

## Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | No |
| `ENV` | Environment (development/production) | - | No |
| `GIN_MODE` | Gin framework mode (debug/release) | `debug` (dev) | No |
| `SEARCH_INDEX_PATH` | Path to Bleve search index | `quran.bleve` | No |
| `KEMENAG_API` | Kemenag API base URL | `https://web-api.qurankemenag.net` | No |
| `PRAYER_TIME_API` | Prayer time API base URL | `https://api.aladhan.com/v1` | No |

### Example `.env` files

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

## Usage

### Basic Usage

Start the server:
```bash
go run main.go
```

### Re-indexing Data

To rebuild the search index:
```bash
go run main.go -reindex
```

### Using Makefile (if available)

```bash
make run          # Run the application
make build        # Build the binary
make test         # Run tests
make clean        # Clean build artifacts
```

## API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### Health Check Endpoints

#### Health Check
```http
GET /health
```
Returns overall health status including search index status.

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
Returns a list of all 114 surahs.

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
    }
  ]
}
```

#### Get Surah Detail
```http
GET /api/v1/surah/detail/:surah_id?page=1&limit=10
```

**Parameters:**
- `surah_id` (path, required): Surah ID (1-114)
- `page` (query, optional): Page number (default: 1)
- `limit` (query, optional): Items per page (default: 10, max: 100)

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

### Search Endpoint

#### Search Verses
```http
GET /api/v1/search?q=allah&page=1&limit=10
```

**Parameters:**
- `q` (query, required): Search query (searches in translation text)
- `page` (query, optional): Page number (default: 1)
- `limit` (query, optional): Items per page (default: 10, max: 100)

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

### Prayer Time Endpoints

#### Get Prayer Times
```http
GET /api/v1/prayer-time?city=Jakarta&country=Indonesia&date=2025-01-15
```

**Parameters:**
- `city` (query, required): City name
- `country` (query, required): Country name
- `date` (query, optional): Date in YYYY-MM-DD format (default: today)

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

## Project Structure

```
quran-api/
├── common/              # Common utilities and test helpers
├── config/              # Configuration and logging setup
├── domain/               # Domain layer
│   ├── dto/             # Data Transfer Objects (API responses)
│   └── model/           # Domain models
├── handler/             # HTTP handlers (request/response)
├── repository/          # Data access layer
├── router/              # Route definitions
├── service/             # Business logic layer
├── utils/               # Utilities and middleware
│   ├── helper/          # Helper functions
│   └── middleware/     # HTTP middlewares
├── main.go              # Application entry point
├── go.mod               # Go module definition
├── Dockerfile           # Docker image definition
├── docker-compose.yml   # Docker Compose configuration
├── .env.example         # Production environment template
├── .env.dev.example     # Development environment template
├── FEATURE_GUIDE.md     # Guide for adding new features
└── README.md            # This file
```

## Development

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

See [FEATURE_GUIDE.md](./FEATURE_GUIDE.md) for a detailed guide on adding new features.

Quick steps:
1. Define domain models in `domain/model/`
2. Create DTOs in `domain/dto/`
3. Implement repository in `repository/`
4. Add business logic in `service/`
5. Create handler in `handler/`
6. Define routes in `router/`
7. Register routes in `router/main_route.go`

### Running Tests

```bash
go test ./...
```

### Code Style

- Follow Go standard formatting: `gofmt -w .`
- Use `golint` or `golangci-lint` for linting
- Follow existing code patterns and conventions

## Deployment

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
  quran-api
```

#### Using Docker Compose

```bash
docker-compose up -d
```

### Production Considerations

1. **Environment Variables**: Set all required environment variables
2. **Index Path**: Use a persistent volume for the search index
3. **Health Checks**: Configure health checks in your orchestration platform
4. **Rate Limiting**: Adjust rate limits based on your needs
5. **Monitoring**: Set up logging and monitoring
6. **SSL/TLS**: Use a reverse proxy (nginx, Traefik) for HTTPS

### Health Checks

The application provides three health check endpoints:
- `/health` - Full health check
- `/health/live` - Liveness probe (Kubernetes)
- `/health/ready` - Readiness probe (Kubernetes)

Configure these in your deployment platform for automatic health monitoring.

## Rate Limiting

The API implements IP-based rate limiting:
- **Rate**: 5 requests per 5 minutes per IP
- **Burst**: 5 requests
- Health check endpoints are not rate-limited

## Security Features

- Security headers (X-Frame-Options, X-Content-Type-Options, etc.)
- Request timeout (30 seconds)
- Rate limiting
- Input validation
- Graceful error handling

## Performance

- Response caching for external API calls
- Efficient Bleve full-text search
- Pagination support
- Connection pooling
- Graceful shutdown

## Troubleshooting

### Search returns no results

1. Ensure the index exists: Check if `quran.bleve/` directory exists
2. Re-index the data: `go run main.go -reindex`
3. Check index health: `GET /health`

### Server won't start

1. Check if port is available: `lsof -i :8080`
2. Verify environment variables are set correctly
3. Check logs for error messages

### Indexing fails

1. Check internet connection (requires external API access)
2. Verify `KEMENAG_API` is accessible
3. Check available disk space
4. Review logs for specific error messages

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow the existing code structure and patterns
- Add tests for new features
- Update documentation as needed
- Ensure all tests pass before submitting

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Quran data provided by [Kemenag API](https://web-api.qurankemenag.net)
- Prayer times provided by [Aladhan API](https://api.aladhan.com)
- Built with [Gin](https://gin-gonic.com/) web framework
- Full-text search powered by [Bleve](https://blevesearch.com/)

## Support

For issues, questions, or contributions, please open an issue on GitHub.

---

**Made with ❤️ for the Muslim community**
