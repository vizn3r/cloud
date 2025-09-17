# Use a specific version for better reproducibility
FROM golang:1.25-alpine3.21 AS builder

# Install build dependencies in one layer
RUN apk add --no-cache git gcc musl-dev ca-certificates

# Set working directory
WORKDIR /app

# Copy dependency files first for better layer caching
COPY server/go.mod server/go.sum ./

# Download dependencies (cached if go.mod/go.sum unchanged)
RUN go mod download

# Copy source code
COPY server/ ./

# Build with optimizations: static linking, smaller binary, no debug info
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-s -w -extldflags=-static" \
    -trimpath \
    -o cloud-server \
    main.go

# Use distroless for minimal runtime image
FROM gcr.io/distroless/static-debian12:nonroot

# Copy only the essential files
COPY --from=builder --chown=nonroot:nonroot /app/cloud-server /app/cloud-server
COPY --from=builder --chown=nonroot:nonroot /app/db/sql/tables.sql /app/db/sql/tables.sql

# Use distroless nonroot user (already configured)
USER nonroot

# Set working directory
WORKDIR /app

# Expose port
EXPOSE 8080

# Health check using the built-in binary
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app/cloud-server", "--health-check"]

# Start the application
ENTRYPOINT ["/app/cloud-server"]
