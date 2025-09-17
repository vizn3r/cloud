#!/bin/bash

docker build -f server/Dockerfile -t cloud-server:latest .
docker build -f webclient/Dockerfile -t cloud-webclient:latest .
docker compose down
docker compose up -d
