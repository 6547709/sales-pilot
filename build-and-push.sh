#!/usr/bin/env bash
# 单架构（默认 linux/amd64）构建并推送到 Docker Hub：docker.io/${REGISTRY}/sales-pilot-{backend,frontend,nginx}:${VERSION}
#
# 使用前：docker login docker.io
# 可选环境变量：VERSION、REGISTRY（Docker Hub 用户名，默认 liguoqiang）、PLATFORM、NEXT_PUBLIC_API_URL
set -euo pipefail

VERSION="${VERSION:-0.1.0}"
REGISTRY="${REGISTRY:-liguoqiang}"
IMAGE_PREFIX="docker.io/${REGISTRY}"
PLATFORM="${PLATFORM:-linux/amd64}"

if ! docker buildx inspect sales-pilot-builder >/dev/null 2>&1; then
  echo ">>> 创建 buildx builder: sales-pilot-builder"
  docker buildx create --name sales-pilot-builder --driver docker-container --use
else
  docker buildx use sales-pilot-builder
fi

echo ">>> 构建并推送 ${PLATFORM} -> ${IMAGE_PREFIX}/sales-pilot-backend:${VERSION}"
docker buildx build --platform "${PLATFORM}" \
  -t "${IMAGE_PREFIX}/sales-pilot-backend:${VERSION}" \
  -t "${IMAGE_PREFIX}/sales-pilot-backend:latest" \
  -f backend/Dockerfile backend \
  --push

echo ">>> 构建并推送 ${PLATFORM} -> ${IMAGE_PREFIX}/sales-pilot-frontend:${VERSION}"
docker buildx build --platform "${PLATFORM}" \
  --build-arg NEXT_PUBLIC_API_URL="${NEXT_PUBLIC_API_URL:-}" \
  -t "${IMAGE_PREFIX}/sales-pilot-frontend:${VERSION}" \
  -t "${IMAGE_PREFIX}/sales-pilot-frontend:latest" \
  -f frontend/Dockerfile frontend \
  --push

echo ">>> 构建并推送 ${PLATFORM} -> ${IMAGE_PREFIX}/sales-pilot-nginx:${VERSION}"
docker buildx build --platform "${PLATFORM}" \
  -t "${IMAGE_PREFIX}/sales-pilot-nginx:${VERSION}" \
  -t "${IMAGE_PREFIX}/sales-pilot-nginx:latest" \
  -f nginx/Dockerfile nginx \
  --push

echo "已推送："
echo "  ${IMAGE_PREFIX}/sales-pilot-backend:${VERSION}（及 :latest）"
echo "  ${IMAGE_PREFIX}/sales-pilot-frontend:${VERSION}（及 :latest）"
echo "  ${IMAGE_PREFIX}/sales-pilot-nginx:${VERSION}（及 :latest）"
