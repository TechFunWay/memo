#!/bin/bash

set -e

PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$PROJECT_DIR"

VERSION_RAW=$(grep -o 'AppVersion = "[^"]*"' backend/config/config.go | cut -d'"' -f2)
[ -z "$VERSION_RAW" ] && echo "❌ 无法获取版本号" && exit 1
VERSION="${VERSION_RAW#v}"

echo "============================================"
echo "  Memo v${VERSION} 一键打包"
echo "============================================"
echo ""

echo "📝 [1/4] 生成发布说明..."
./scripts/generate-release-readmes.sh

echo ""
echo "📦 [2/4] 打包所有平台..."
./scripts/build-all.sh

echo ""
echo "📦 [3/4] 打包飞牛应用..."
./scripts/build-fnpack.sh

echo ""
echo "📦 [4/4] 构建 Docker 镜像..."
./scripts/build-docker.sh

echo ""
echo "============================================"
echo "  ✅ 全部完成!"
echo "============================================"
echo ""
echo "📁 release/v${VERSION}/"
ls -lh "release/v${VERSION}/"*.tar.gz 2>/dev/null | awk '{print "  平台包: " $NF " (" $5 ")"}'
ls -lh "release/v${VERSION}/"*.fpk 2>/dev/null | awk '{print "  飞牛包: " $NF " (" $5 ")"}'
echo ""
