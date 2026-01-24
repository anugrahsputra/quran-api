# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git build-base

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# We use CGO_ENABLED=0 for a static binary, but we have build-base for transitive dependencies if needed
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o quran-api cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata curl

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/quran-api .

# Create directory for search index
RUN mkdir -p /data

# Expose port (default)
EXPOSE 8080

# Default environment variables
ENV PORT=8080

# Optimize Go memory management for container environments
ENV GOMEMLIMIT=512MiB
ENV GOGC=100

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
  CMD curl -f http://localhost:${PORT}/health/live || exit 1

# Run the application
CMD ["./quran-api"]
