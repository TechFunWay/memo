#!/bin/bash

set -e

PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
cd "$PROJECT_DIR"

VERSION_RAW=$(grep -o 'AppVersion = "[^"]*"' backend/config/config.go | cut -d'"' -f2)
[ -z "$VERSION_RAW" ] && echo "❌ 无法获取版本号" && exit 1
VERSION="${VERSION_RAW#v}"

APP_NAME="techfunway-memo"
IMAGE_NAME="techfunways/memo"
RELEASE_DIR="release/v${VERSION}"

[ ! -d "${RELEASE_DIR}/${APP_NAME}-${VERSION}-linux-amd64" ] && echo "❌ 请先运行 ./scripts/build-all.sh" && exit 1

echo "============================================"
echo "  Docker 多平台镜像构建 v${VERSION}"
echo "============================================"

# 清理旧镜像
docker rmi $(docker images "${IMAGE_NAME}" -q) 2>/dev/null || true

# 删除 .DS_Store
find "${RELEASE_DIR}" -name ".DS_Store" -delete 2>/dev/null || true

# 切换到 default context
docker context use default 2>/dev/null || true

echo ""
echo "🔨 构建多平台镜像..."
echo ""

cd "${RELEASE_DIR}"

# 创建 Dockerfile
cat > Dockerfile << EOF
FROM scratch

ARG TARGETARCH
ARG VERSION

COPY ${APP_NAME}-\${VERSION}-linux-\${TARGETARCH}/memo /app/memo
COPY ${APP_NAME}-\${VERSION}-linux-\${TARGETARCH}/www /app/www

VOLUME ["/app/data"]
EXPOSE 8905

ENTRYPOINT ["/app/memo", "-port", "8905", "-data-dir", "/app/data", "-web-dir", "/app/www"]
EOF

# 构建多平台镜像
docker buildx build \
    --builder default \
    --platform linux/amd64,linux/arm64 \
    --output type=docker \
    --build-arg VERSION=${VERSION} \
    -t "${IMAGE_NAME}:v${VERSION}" \
    -t "${IMAGE_NAME}:latest" \
    .

rm -f Dockerfile
cd "${PROJECT_DIR}"

# 生成 docker-compose 文件（每次构建 docker 后自动生成）
# - docker-compose-v${VERSION}.yml：固定版本，不自动更新
# - docker-compose-latest.yml：跟随 latest，并用 Watchtower 自动更新
COMPOSE_VER="${RELEASE_DIR}/docker-compose-v${VERSION}.yml"
COMPOSE_LATEST="${RELEASE_DIR}/docker-compose-latest.yml"

cat > "${COMPOSE_VER}" << EOF
# 备忘录 v${VERSION} — 固定版本部署（不自动更新）
# 用法: docker compose -f docker-compose-v${VERSION}.yml up -d
services:
  memo:
    image: ${IMAGE_NAME}:v${VERSION}
    container_name: memo
    restart: unless-stopped
    ports:
      - "8905:8905"
    volumes:
      - ./data:/app/data
EOF

cat > "${COMPOSE_LATEST}" << EOF
# 备忘录 latest — 跟随最新版，并通过 Watchtower 自动更新
# 用法: docker compose -f docker-compose-latest.yml up -d
# Watchtower 每小时检查一次，仅更新带 enable 标签的容器，并清理旧镜像
services:
  memo:
    image: ${IMAGE_NAME}:latest
    container_name: memo
    restart: unless-stopped
    ports:
      - "8905:8905"
    volumes:
      - ./data:/app/data
    labels:
      - "com.centurylinklabs.watchtower.enable=true"

  watchtower:
    image: containrrr/watchtower
    container_name: memo-watchtower
    restart: unless-stopped
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    command: --label-enable --cleanup --interval 3600
EOF

echo "  📄 docker-compose 已生成: ${COMPOSE_VER}"
echo "  📄 docker-compose 已生成: ${COMPOSE_LATEST}"

echo ""
echo "✅ 完成!"
echo ""
echo "📦 本地镜像:"
docker images "${IMAGE_NAME}" --format "  {{.Repository}}:{{.Tag}}" | sort -u
echo ""
echo "💡 运行: docker run -p 8905:8905 -v ./data:/app/data ${IMAGE_NAME}:v${VERSION}"
