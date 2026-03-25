#!/usr/bin/env bash
# 构建并推送镜像，标签：docker.io/liguoqiang/[IMAGE]:[VERSION]
set -euo pipefail

VERSION="${VERSION:-0.1.0}"
REGISTRY="${REGISTRY:-docker.io/liguoqiang}"

docker buildx build --platform linux/amd64 -t "${REGISTRY}/sales-pilot-backend:${VERSION}" -f backend/Dockerfile backend --push
docker buildx build --platform linux/amd64 \
  --build-arg NEXT_PUBLIC_API_URL="${NEXT_PUBLIC_API_URL:-}" \
  -t "${REGISTRY}/sales-pilot-frontend:${VERSION}" -f frontend/Dockerfile frontend --push

echo "已推送 ${REGISTRY}/sales-pilot-backend:${VERSION}"
echo "已推送 ${REGISTRY}/sales-pilot-frontend:${VERSION}"
