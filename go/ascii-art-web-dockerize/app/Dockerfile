# Use the official Golang image for building the app
FROM golang:1.23.2 AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory inside the container
WORKDIR /app

# Copy the Go modules files first to leverage caching
COPY go.mod ./


# Download the module dependencies
RUN go mod download

# Copy the rest of the app source code
COPY . .

# Build the Go application
RUN go build -ldflags="-w -s" -o webapp main.go

# Use a lightweight image to run the app
FROM alpine:3.18.2

# Set metadata
LABEL maintainer="Marinos Kouvaras, Ypatios Chaniotakos"
LABEL version="1.0"
LABEL description="Optimized Dockerized Go Web Application"


# Install essential certificates and bash
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Set working directory for the runtime container
WORKDIR /app


# Copy the binary from the builder stage
COPY --from=builder /app /app

# Copy the built binary and resources from the builder stage
# COPY --from=builder /app/webapp .
# COPY --from=builder /app/templates ./templates
# COPY --from=builder /app/static ./static
# COPY --from=builder /app/banners ./banners
# COPY --from=builder /app/tools ./tools


# Expose the port that the application listens on
EXPOSE 8080

# Healthcheck for monitoring
HEALTHCHECK --interval=30s --timeout=5s \
    CMD curl -f http://localhost:8080/health || exit 1

# Command to run the application
CMD ["./webapp"]
