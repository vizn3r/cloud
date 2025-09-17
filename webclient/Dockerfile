FROM golang:1.25-alpine3.21 AS builder

WORKDIR /app

# Copy dependency files first for better caching
COPY webclient/go.mod webclient/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY webclient/ ./

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -trimpath \
    -o main .

# Use distroless for runtime
FROM gcr.io/distroless/static-debian12:nonroot

# Copy binary and static files
COPY --from=builder --chown=nonroot:nonroot /app/main /app/main
COPY --from=builder --chown=nonroot:nonroot /app/public /app/public

# Use nonroot user
USER nonroot

WORKDIR /app

# Expose port
EXPOSE 8008

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app/main", "--health-check"]

# Start the application
ENTRYPOINT ["/app/main"]
