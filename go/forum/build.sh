#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
IMAGE_NAME="forum-app"
CONTAINER_NAME="forum"
PORT="8080"

echo -e "${YELLOW}Building Forum Application Docker Image...${NC}"

# Clean up existing containers and images
echo -e "${YELLOW}Cleaning up existing containers...${NC}"
docker stop $CONTAINER_NAME 2>/dev/null || true
docker rm $CONTAINER_NAME 2>/dev/null || true

# Build the Docker image
echo -e "${YELLOW}Building Docker image: $IMAGE_NAME${NC}"
if docker build -t $IMAGE_NAME .; then
    echo -e "${GREEN}✓ Docker image built successfully${NC}"
else
    echo -e "${RED}✗ Failed to build Docker image${NC}"
    exit 1
fi

# Show built images
echo -e "${YELLOW}Available Docker images:${NC}"
docker images | grep -E "(REPOSITORY|$IMAGE_NAME)"

# Run the container
echo -e "${YELLOW}Starting container: $CONTAINER_NAME${NC}"
if docker run -d --name $CONTAINER_NAME -p $PORT:$PORT $IMAGE_NAME; then
    echo -e "${GREEN}✓ Container started successfully${NC}"
else
    echo -e "${RED}✗ Failed to start container${NC}"
    exit 1
fi

# Show running containers
echo -e "${YELLOW}Running containers:${NC}"
docker ps -a | grep -E "(CONTAINER ID|$CONTAINER_NAME)"

# Show application status
echo -e "${GREEN}✓ Forum application is now running!${NC}"
echo -e "${YELLOW}Access the application at: http://localhost:$PORT${NC}"
echo -e "${YELLOW}Container logs: docker logs $CONTAINER_NAME${NC}"
echo -e "${YELLOW}Stop container: docker stop $CONTAINER_NAME${NC}"