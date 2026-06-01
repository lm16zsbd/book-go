#!/bin/bash
set -e

CONTAINER="serica-go"
IMAGE="404904371652.dkr.ecr.ap-southeast-1.amazonaws.com/serica-go-dev:latest"

echo "=== Pulling latest image ==="
docker pull $IMAGE

echo "=== Stopping old container ==="
docker stop $CONTAINER 2>/dev/null || true
docker rm $CONTAINER 2>/dev/null || true

echo "=== Starting new container ==="
docker run -d \
  --name $CONTAINER \
  --restart unless-stopped \
  -p 8080:8080 \
  --env-file .env \
  $IMAGE

echo "=== Deploy done ==="
docker ps --filter name=$CONTAINER
