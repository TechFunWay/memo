#!/bin/bash

set -e

echo "🚀 启动 Memo 服务..."

# 停止旧服务
if lsof -Pi :8905 -sTCP:LISTEN -t >/dev/null ; then
    echo "⚠️  端口 8905 已被占用，正在停止旧进程..."
    lsof -ti:8905 | xargs kill -9 2>/dev/null
    sleep 1
fi

# 编译前端
echo "📦 编译前端..."
cd "$(dirname "$0")/frontend"
npm run build 2>/dev/null

# 复制到后端
echo "📋 部署前端到后端..."
rm -rf ../backend/static/*
mkdir -p ../backend/static
cp -r dist/* ../backend/static/

# 启动后端
echo "🖥️  启动后端服务..."
cd ../backend
go run . -port 8905 -data-dir ./data &

BACKEND_PID=$!

echo ""
echo "✨ 服务已启动!"
echo "📌 访问地址: http://localhost:8905"
echo ""
echo "按 Ctrl+C 停止服务"

trap "kill $BACKEND_PID 2>/dev/null; echo '服务已停止'; exit" INT
wait
