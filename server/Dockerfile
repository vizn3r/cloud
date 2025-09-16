FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY server/go.mod server/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY server/ ./

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o cloud-server -ldflags="-s -w" main.go

# Final stage
FROM alpine:3.21

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite

# Create app user
RUN addgroup -g 1000 -S app && \
    adduser -u 1000 -S app -G app

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=app:app /app/cloud-server .

# Copy SQLite database schema
COPY --chown=app:app server/db/sql/tables.sql ./db/sql/

# Create data directory for SQLite
RUN mkdir -p /data && chown app:app /data

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Switch to non-root user
USER app

# Start the application
CMD ["./cloud-server"]