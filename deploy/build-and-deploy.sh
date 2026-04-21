#!/usr/bin/env bash
#
# build-and-deploy.sh
#
# 一键编译前端 + 嵌入前端的后端二进制，并热更新到运行中的容器。
#
# 工作原理：
#   1. 本地编译前端 → 编译 Linux amd64 后端二进制
#   2. scp 上传二进制到服务器
#   3. docker cp 替换容器内二进制
#   4. docker restart 重启容器
#
# 注意：
#   - 此方式为热更新，替换的二进制仅在当前容器生命周期内有效
#   - 如果执行 docker compose down && up，会恢复为镜像内版本
#   - 如需持久化更新，请使用 docker compose pull 拉取新镜像
#
# 支持平台：
#   - macOS / Linux：直接运行
#   - Windows：通过 Git Bash 或 WSL 运行
#
# 用法：
#   cd <项目根目录>
#   bash deploy/build-and-deploy.sh
#
# 可通过环境变量覆盖默认配置：
#   SERVER_HOST=192.168.160.145 SERVER_USER=ggh bash deploy/build-and-deploy.sh
#
# 生产服务器（192.168.160.145）请使用 docker compose pull 更新，不走此脚本。
#

set -euo pipefail

# 默认配置：测试服务器（ggh@10.10.200.21）
SERVER_HOST="${SERVER_HOST:-10.10.200.21}"
SERVER_USER="${SERVER_USER:-ggh}"
REMOTE_DIR="${REMOTE_DIR:-/home/ggh/dockers/sub2api-plus}"
CONTAINER_NAME="${CONTAINER_NAME:-sub2api-plus}"
BINARY_NAME="server"

echo "==> Build & Deploy (Hot Update)"
echo "    Server: ${SERVER_USER}@${SERVER_HOST}"
echo "    Remote: ${REMOTE_DIR}"
echo ""

# 检查前置依赖
check_cmd() {
  if ! command -v "$1" &>/dev/null; then
    echo "Error: '$1' not found in PATH."
    exit 1
  fi
}

check_cmd pnpm
check_cmd go
check_cmd ssh
check_cmd scp

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "${PROJECT_ROOT}"

# 1. 编译前端
echo "==> Building frontend..."
cd frontend
pnpm install
pnpm build
cd "${PROJECT_ROOT}"

# 2. 编译后端（嵌入前端，交叉编译为 Linux amd64）
echo "==> Building backend (embedded frontend, linux/amd64)..."
cd backend
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
  -tags embed \
  -ldflags="-s -w" \
  -trimpath \
  -o "${PROJECT_ROOT}/${BINARY_NAME}" \
  ./cmd/server
cd "${PROJECT_ROOT}"

echo "==> Binary built: ${PROJECT_ROOT}/${BINARY_NAME}"

# 3. 上传到服务器
echo "==> Uploading binary to server..."
scp "${BINARY_NAME}" "${SERVER_USER}@${SERVER_HOST}:${REMOTE_DIR}/${BINARY_NAME}"

# 4. 替换容器内二进制并重启
echo "==> Hot updating container binary..."
ssh "${SERVER_USER}@${SERVER_HOST}" << EOF
  set -e
  cd "${REMOTE_DIR}"
  docker cp ${BINARY_NAME} ${CONTAINER_NAME}:/app/sub2api-plus
  docker restart ${CONTAINER_NAME}
EOF

# 5. 等待并检查健康状态
echo "==> Waiting for container to become healthy..."
sleep 8

HEALTH_STATUS=$(ssh "${SERVER_USER}@${SERVER_HOST}" "docker inspect --format='{{.State.Health.Status}}' ${CONTAINER_NAME} 2>/dev/null || echo 'unknown'")
UPTIME=$(ssh "${SERVER_USER}@${SERVER_HOST}" "docker ps --filter name=${CONTAINER_NAME} --format 'table {{.Names}}\t{{.Status}}'")

echo ""
echo "==> Container status:"
echo "${UPTIME}"
echo ""
echo "==> Health status: ${HEALTH_STATUS}"

if [ "${HEALTH_STATUS}" = "healthy" ]; then
  echo "==> Hot update succeeded."
  echo ""
  echo "NOTE: This binary is injected into the running container."
  echo "      If you run 'docker compose down && docker compose up -d',"
  echo "      the container will revert to the image's built-in binary."
  echo "      To persist the update, push a new Docker image and pull it."
  exit 0
elif [ "${HEALTH_STATUS}" = "starting" ]; then
  echo "==> Container is still starting. Please check manually in a few seconds."
  exit 0
else
  echo "==> WARNING: Container health status is '${HEALTH_STATUS}'. Please check logs."
  exit 1
fi
