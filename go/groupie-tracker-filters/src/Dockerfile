# Build stage
FROM golang:1.23.2 AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory inside the container
WORKDIR /website

# Copy the Go modules files first to leverage caching
COPY go.mod ./

# Download the module dependencies
RUN go mod download

# Copy the rest of the app source code
COPY . ./

# Build the Go application
RUN go build -ldflags="-w -s" -o website main.go

# Runtime stage
FROM alpine:3.18.2

# Set metadata
LABEL maintainer="Marinos Kouvaras, Ypatios Chaniotakos"
LABEL version="1.0"
LABEL description="Groupie tracker website"

# Install essential certificates and bash
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Set working directory for the runtime container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /website /app
EXPOSE 8080

# Command to run the application
CMD ["/app/website"]
