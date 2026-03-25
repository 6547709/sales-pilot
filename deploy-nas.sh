#!/usr/bin/env bash
# 一键在 NAS（或任意远端主机）部署：需已安装 Docker / Compose plugin
# 用法：
#   export NAS_HOST=192.168.1.10
#   export NAS_USER=admin
#   export NAS_PATH=/volume1/docker/sales-pilot
#   ./deploy-nas.sh
set -euo pipefail

: "${NAS_HOST:?请设置 NAS_HOST}"
NAS_USER="${NAS_USER:-root}"
NAS_PATH="${NAS_PATH:-/share/sales-pilot}"
SSH_TARGET="${NAS_USER}@${NAS_HOST}"

echo ">>> 同步项目到 ${SSH_TARGET}:${NAS_PATH}"
ssh "${SSH_TARGET}" "mkdir -p '${NAS_PATH}'"
rsync -avz --delete \
  --exclude '.git' --exclude 'frontend/node_modules' --exclude 'frontend/.next' \
  --exclude '.venv' --exclude 'Request.pdf' \
  ./ "${SSH_TARGET}:${NAS_PATH}/"

echo ">>> 远端启动（数据卷持久在 Docker volume pgdata）"
ssh "${SSH_TARGET}" "cd '${NAS_PATH}' && \
  JWT_SECRET=\${JWT_SECRET:-$(openssl rand -hex 32)} \
  APP_BASE_URL=\${APP_BASE_URL:-http://${NAS_HOST}} \
  FRONTEND_ORIGIN=\${FRONTEND_ORIGIN:-http://${NAS_HOST}} \
  docker compose pull 2>/dev/null || true && \
  docker compose build && \
  docker compose up -d"

echo ">>> 完成。PostgreSQL 数据位于远端 Docker volume: sales-pilot_pgdata（名称因目录而异，可用 docker volume ls 查看）"
