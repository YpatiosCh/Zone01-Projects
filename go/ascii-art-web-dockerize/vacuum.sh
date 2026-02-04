#!/bin/bash

# Enable strict error handling
set -euo pipefail

echo "Starting Docker cleanup..."

# Prune stopped containers
docker container prune -f

# Prune unused images
docker image prune -a -f

# Prune unused volumes
docker volume prune -f

# Prune unused networks
docker network prune -f

# Stop and remove the container associated with go-webapp:1.0
container_id=$(docker ps -q --filter "ancestor=go-webapp:1.0")

if [[ -n "$container_id" ]]; then
    echo "Stopping container running image go-webapp:1.0..."
    docker stop "$container_id"
    docker rm "$container_id"
fi

# Remove the image go-webapp:1.0
if docker image inspect go-webapp:1.0 > /dev/null 2>&1; then
    echo "Removing image go-webapp:1.0..."
    docker image rm -f go-webapp:1.0
fi

echo "Docker cleanup completed successfully."
