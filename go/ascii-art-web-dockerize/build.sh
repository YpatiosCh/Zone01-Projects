#!/bin/bash

# Enable strict error handling
set -euo pipefail

# Configuration variables
IMAGE_NAME="go-webapp"
IMAGE_TAG="1.0"
CONTAINER_NAME="go-webapp-container"
PORT_MAPPING="8080:8080"

# Set the path to the app directory
APP_DIR="./app"

# Color codes for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print success message
print_success() {
    echo -e "${GREEN}✔ $1${NC}"
}

# Function to print warning message
print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Build Docker image
build_image() {
    echo "Building Docker image..."
    if docker build -t "${IMAGE_NAME}:${IMAGE_TAG}" "${APP_DIR}"; then
        print_success "Docker image built successfully"
    else
        echo "Failed to build Docker image"
        exit 1
    fi
}

# Create and run container
run_container() {
    echo "Creating and running container..."
    
    # Remove existing container if it exists
    if docker ps -a | grep -q "${CONTAINER_NAME}"; then
        print_warning "Removing existing container..."
        docker rm -f "${CONTAINER_NAME}"
    fi

    # Run new container
    if docker run -d \
        --name "${CONTAINER_NAME}" \
        -p "${PORT_MAPPING}" \
        "${IMAGE_NAME}:${IMAGE_TAG}"; then
        print_success "Container started successfully"
        print_success "Application is now accessible at http://localhost:8080"
    else
        echo "Failed to start container"
        exit 1
    fi
}

# Verify container health
verify_container() {
    echo "Verifying container health..."
    sleep 3
    if docker ps | grep -q "${CONTAINER_NAME}"; then
        print_success "Container is running"
        docker ps | grep "${CONTAINER_NAME}"
    else
        echo "Container failed to start"
        exit 1
    fi
}

# Main build process
main() {
    echo "Starting build process for Go Web Application..."
    
    # Build image
    build_image
    
    # Run container
    run_container
    
    # Verify container
    verify_container
}

# Run the main process
main