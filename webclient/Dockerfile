FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY webclient/go.mod webclient/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY webclient/ ./

# Build the application
RUN go build -o main .

FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/public ./public

# Create non-root user
RUN addgroup -g 1001 -S app && \
    adduser -u 1001 -S app -G app

# Set proper ownership
RUN chown -R app:app /app

# Switch to non-root user
USER app

# Expose port (default Fiber port is 3000)
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

# Start the application
CMD ["./main"]
