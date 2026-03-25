#!/usr/bin/env bash
# 多架构构建并推送到 Docker Hub：docker.io/${REGISTRY}/sales-pilot-{backend,frontend,nginx}:${VERSION}
#
# 使用前：
#   1. docker login docker.io
#   2. 本地只看到 sales-pilot-* 无 liguoqiang 前缀属正常：compose 默认本地标签；
#      本脚本推送后镜像在远端；多架构 --push 不会在本地生成同名列表（可 docker pull 验证）。
#
# 可选环境变量：VERSION、REGISTRY（Docker Hub 用户名，默认 liguoqiang）、PLATFORMS、NEXT_PUBLIC_API_URL
set -euo pipefail

VERSION="${VERSION:-0.1.0}"
REGISTRY="${REGISTRY:-liguoqiang}"
# Docker Hub 完整仓库路径（省略 registry 前缀时 CLI 默认 docker.io）
IMAGE_PREFIX="docker.io/${REGISTRY}"
PLATFORMS="${PLATFORMS:-linux/amd64,linux/arm64}"

if ! docker buildx inspect multiarch-builder >/dev/null 2>&1; then
  echo ">>> 创建 buildx builder: multiarch-builder"
  docker buildx create --name multiarch-builder --driver docker-container --use
else
  docker buildx use multiarch-builder
fi

echo ">>> 构建并推送 ${PLATFORMS} -> ${IMAGE_PREFIX}/sales-pilot-backend:${VERSION}"
docker buildx build --platform "${PLATFORMS}" \
  -t "${IMAGE_PREFIX}/sales-pilot-backend:${VERSION}" \
  -t "${IMAGE_PREFIX}/sales-pilot-backend:latest" \
  -f backend/Dockerfile backend \
  --push

echo ">>> 构建并推送 ${PLATFORMS} -> ${IMAGE_PREFIX}/sales-pilot-frontend:${VERSION}"
docker buildx build --platform "${PLATFORMS}" \
  --build-arg NEXT_PUBLIC_API_URL="${NEXT_PUBLIC_API_URL:-}" \
  -t "${IMAGE_PREFIX}/sales-pilot-frontend:${VERSION}" \
  -t "${IMAGE_PREFIX}/sales-pilot-frontend:latest" \
  -f frontend/Dockerfile frontend \
  --push

echo ">>> 构建并推送 ${PLATFORMS} -> ${IMAGE_PREFIX}/sales-pilot-nginx:${VERSION}"
docker buildx build --platform "${PLATFORMS}" \
  -t "${IMAGE_PREFIX}/sales-pilot-nginx:${VERSION}" \
  -t "${IMAGE_PREFIX}/sales-pilot-nginx:latest" \
  -f nginx/Dockerfile nginx \
  --push

echo "已推送："
echo "  ${IMAGE_PREFIX}/sales-pilot-backend:${VERSION} (及 :latest)"
echo "  ${IMAGE_PREFIX}/sales-pilot-frontend:${VERSION} (及 :latest)"
echo "  ${IMAGE_PREFIX}/sales-pilot-nginx:${VERSION} (及 :latest)"
