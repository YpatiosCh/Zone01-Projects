#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Cleaning up Docker objects...${NC}"

# Stop and remove forum containers
echo -e "${YELLOW}Stopping and removing forum containers...${NC}"
docker stop forum 2>/dev/null || true
docker rm forum 2>/dev/null || true

# Remove forum-specific images
echo -e "${YELLOW}Removing forum images...${NC}"
docker rmi forum-app:latest 2>/dev/null || true
docker rmi forum-app 2>/dev/null || true

# Remove unused objects
echo -e "${YELLOW}Removing unused containers...${NC}"
docker container prune -f

echo -e "${YELLOW}Removing unused images...${NC}"
docker image prune -f

echo -e "${YELLOW}Removing unused volumes...${NC}"
docker volume prune -f

echo -e "${YELLOW}Removing unused networks...${NC}"
docker network prune -f

echo -e "${GREEN}✓ Cleanup completed!${NC}"

# Show current state
echo -e "${YELLOW}Current Docker state:${NC}"
echo "Images:"
docker images
echo -e "\nContainers:"
docker ps -a