#!/bin/bash

set -e

PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
cd "$PROJECT_DIR"

VERSION_RAW=$(grep -o 'AppVersion = "[^"]*"' backend/config/config.go | cut -d'"' -f2)
[ -z "$VERSION_RAW" ] && echo "❌ 无法获取版本号" && exit 1
VERSION="${VERSION_RAW#v}"

APP_NAME="techfunway-memo"
RELEASE_DIR="release/v${VERSION}"

echo "============================================"
echo "  多平台打包 v${VERSION}"
echo "============================================"

mkdir -p "${RELEASE_DIR}"

echo ""
echo "📦 编译前端..."
cd frontend && npm run build && cd ..

compile() {
    GOOS=$1
    GOARCH=$2
    LABEL=$3
    echo -n "  📦 ${LABEL}... "
    
    DIR="${RELEASE_DIR}/${APP_NAME}-${VERSION}-${LABEL}"
    rm -rf "${DIR}"
    mkdir -p "${DIR}/www"
    
    cd backend
    CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o "../${DIR}/memo" . 2>/dev/null
    cd ..
    
    if [ "${GOOS}" = "windows" ]; then
        mv "${DIR}/memo" "${DIR}/memo.exe"
        echo '@echo off
cd /d "%~dp0"
memo.exe %*' > "${DIR}/start.bat"
    else
        echo '#!/bin/bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"
chmod +x memo 2>/dev/null || true
./memo "$@"' > "${DIR}/start.sh"
        chmod +x "${DIR}/start.sh"
    fi
    
    cp -r frontend/dist/* "${DIR}/www/"

    # 复制发布说明文档（如果存在）
    for doc in release/v${VERSION}/README-*.md; do
        if [ -f "$doc" ]; then
            cp "$doc" "${DIR}/"
        fi
    done
	    # 复制 HTML 使用说明
	    if [ -f "release/v${VERSION}/usage.html" ]; then
	        cp "release/v${VERSION}/usage.html" "${DIR}/"
	    fi
    find "${DIR}" -name ".DS_Store" -delete 2>/dev/null || true
    
    echo "✅"
}

echo ""
echo "🔨 编译平台..."
echo ""
compile "linux" "amd64" "linux-amd64"
compile "linux" "arm64" "linux-arm64"
compile "darwin" "amd64" "macos-amd64"
compile "darwin" "arm64" "macos-arm64"
compile "windows" "amd64" "windows-amd64"

echo ""
echo "📦 创建压缩包..."
cd "${RELEASE_DIR}"
for dir in ${APP_NAME}-${VERSION}-*/; do
    name=$(basename "$dir")
    COPYFILE_DISABLE=1 tar -czf "${name}.tar.gz" --exclude='.DS_Store' "$dir"
done
cd "$PROJECT_DIR"

echo ""
echo "✅ 完成!"
ls -lh "${RELEASE_DIR}"/*.tar.gz | awk '{print "  " $NF " (" $5 ")"}'
