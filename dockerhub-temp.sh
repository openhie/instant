#!/usr/bin/env bash
set -ex

# automate tagging with the short commit hash
docker build --no-cache -t openhie/instant-temp:$(git rev-parse --short HEAD) .
docker tag openhie/instant-temp:$(git rev-parse --short HEAD) openhie/instant-temp
docker push openhie/instant-temp:$(git rev-parse --short HEAD)
docker push openhie/instant-temp:latest