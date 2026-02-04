# Build stage
FROM golang:1.24-alpine3.21 AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/web/main.go

# Production stage
FROM alpine:3.21

# Install runtime dependencies and security updates
RUN apk --no-cache add ca-certificates sqlite && \
    apk upgrade --no-cache

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

# Create app directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

COPY config.json /app/config.json

# Copy static files
COPY --from=builder /app/static ./static

# Copy templates
COPY --from=builder /app/templates ./templates

# Copy database schema (if needed for initialization)
COPY --from=builder /app/internal/database/schema.sql ./internal/database/

# Create directory for database with proper permissions
RUN mkdir -p /app/data && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080

# Run the application
CMD ["./main"]