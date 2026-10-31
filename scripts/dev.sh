#!/usr/bin/env bash

set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APP_PORT="${PORT:-8905}"
DEV_DIR="${PROJECT_DIR}/dev"
WEB_DIR="${DEV_DIR}/www"
DATA_DIR="${DEV_DIR}/data"
BINARY_PATH="${DEV_DIR}/memo"
NEXT_WEB_DIR="${DEV_DIR}/.www-next"
NEXT_BINARY_PATH="${DEV_DIR}/.memo-next"

if ! [[ "${APP_PORT}" =~ ^[0-9]+$ ]] || [ "${APP_PORT}" -lt 1 ] || [ "${APP_PORT}" -gt 65535 ]; then
    echo "❌ 无效端口：${APP_PORT}，请输入 1-65535 之间的数字"
    exit 1
fi

if ! command -v lsof >/dev/null 2>&1; then
    echo "❌ 未找到 lsof，无法检测并停止端口占用进程"
    exit 1
fi

cleanup_next() {
    rm -rf "${NEXT_WEB_DIR}"
    rm -f "${NEXT_BINARY_PATH}"
}

trap cleanup_next EXIT

echo "============================================"
echo "  Memo 本地正式环境构建"
echo "============================================"
echo ""

mkdir -p "${DEV_DIR}" "${DATA_DIR}"
cleanup_next

echo "📦 [1/4] 编译前端..."
(
    cd "${PROJECT_DIR}/frontend"
    npm run build
)

echo ""
echo "📋 [2/4] 组装前端资源..."
mkdir -p "${NEXT_WEB_DIR}"
cp -R "${PROJECT_DIR}/frontend/dist/." "${NEXT_WEB_DIR}/"

echo ""
echo "🔨 [3/4] 编译后端..."
(
    cd "${PROJECT_DIR}/backend"
    CGO_ENABLED=0 go build -o "${NEXT_BINARY_PATH}" .
)
chmod +x "${NEXT_BINARY_PATH}"

echo ""
echo "📁 [4/4] 更新 dev 运行目录..."
rm -rf "${WEB_DIR}"
mv "${NEXT_WEB_DIR}" "${WEB_DIR}"
mv "${NEXT_BINARY_PATH}" "${BINARY_PATH}"
cp "${PROJECT_DIR}/scripts/dev-start.sh" "${DEV_DIR}/start.sh"
chmod +x "${DEV_DIR}/start.sh"

trap - EXIT

echo ""
echo "✅ 构建完成"
echo ""
echo "dev/"
echo "├── memo       后端可执行文件"
echo "├── start.sh   启动脚本"
echo "├── www/       前端静态资源"
echo "└── data/      数据库、上传文件和日志（重新构建时保留）"
echo ""

LISTENER_PIDS="$(lsof -tiTCP:"${APP_PORT}" -sTCP:LISTEN 2>/dev/null || true)"
if [ -n "${LISTENER_PIDS}" ]; then
    echo "⚠️  端口 ${APP_PORT} 已被占用，正在停止旧服务..."
    while IFS= read -r listener_pid; do
        if [ -n "${listener_pid}" ]; then
            kill -TERM "${listener_pid}" 2>/dev/null || true
        fi
    done <<< "${LISTENER_PIDS}"

    for _ in 1 2 3 4 5 6 7 8 9 10; do
        if ! lsof -tiTCP:"${APP_PORT}" -sTCP:LISTEN >/dev/null 2>&1; then
            break
        fi
        sleep 0.5
    done

    REMAINING_PIDS="$(lsof -tiTCP:"${APP_PORT}" -sTCP:LISTEN 2>/dev/null || true)"
    if [ -n "${REMAINING_PIDS}" ]; then
        echo "⚠️  旧服务未及时退出，正在强制停止..."
        while IFS= read -r listener_pid; do
            if [ -n "${listener_pid}" ]; then
                kill -KILL "${listener_pid}" 2>/dev/null || true
            fi
        done <<< "${REMAINING_PIDS}"
        sleep 0.5
    fi

    if lsof -tiTCP:"${APP_PORT}" -sTCP:LISTEN >/dev/null 2>&1; then
        echo "❌ 无法释放端口 ${APP_PORT}"
        exit 1
    fi
    echo "✅ 端口 ${APP_PORT} 已释放"
fi

echo ""
echo "🚀 从 dev/ 启动 Memo..."
echo "📌 访问地址：http://localhost:${APP_PORT}"
echo "📁 数据目录：${DATA_DIR}"
echo "按 Ctrl+C 停止服务"
echo ""

cd "${DEV_DIR}"
PORT="${APP_PORT}" exec ./start.sh
