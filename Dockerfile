FROM golang:1.23-alpine AS builder

# Install git for dependency fetching
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o go-insight cmd/main.go

# Use minimal alpine image for production
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/go-insight .
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations

# Copy config files
COPY --from=builder /app/config ./config
COPY --from=builder /app/.env.example ./.env

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Health check (using default port, will be overridden by actual PORT env)
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Expose default port (actual port set via ENV at runtime)
EXPOSE 8080

# Run the application
CMD ["./go-insight"]