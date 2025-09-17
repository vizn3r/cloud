#!/bin/bash

# Enable BuildKit for faster builds
export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}  Building optimized Docker images...${NC}"

# Build with cache optimization
echo -e "${YELLOW} Building server image...${NC}"
docker build \
  --build-arg BUILDKIT_INLINE_CACHE=1 \
  --cache-from cloud-server:latest \
  -f server/Dockerfile \
  -t cloud-server:latest \
  . || {
  echo -e "${RED} Server build failed${NC}"
  exit 1
}

echo -e "${YELLOW} Building web client image...${NC}"
docker build \
  --build-arg BUILDKIT_INLINE_CACHE=1 \
  --cache-from cloud-webclient:latest \
  -f webclient/Dockerfile \
  -t cloud-webclient:latest \
  . || {
  echo -e "${RED} Web client build failed${NC}"
  exit 1
}

# Show image sizes
echo -e "${GREEN} Image sizes:${NC}"
docker images | grep -E "(cloud-server|cloud-webclient)" | awk '{print $1":"$2" - "$7}'

# Stop existing containers
echo -e "${YELLOW} Stopping existing containers...${NC}"
docker compose down

# Start services with optimized compose
echo -e "${GREEN} Starting optimized services...${NC}"
docker compose up -d

# Wait for health checks
echo -e "${YELLOW} Waiting for services to be healthy...${NC}"
timeout=60
while [ $timeout -gt 0 ]; do
  if docker compose ps | grep -q "healthy"; then
    echo -e "${GREEN} Services are healthy!${NC}"
    break
  fi
  echo -n "."
  sleep 2
  ((timeout -= 2))
done

if [ $timeout -le 0 ]; then
  echo -e "${RED} Services failed to become healthy${NC}"
  docker compose logs
  exit 1
fi

echo -e "${GREEN} Deployment complete!${NC}"
echo -e "${GREEN} Web client: http://localhost:8818${NC}"
echo -e "${GREEN} API server: http://localhost:8808${NC}"
