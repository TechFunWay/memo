#!/bin/bash

# 生成发布说明文件（从模板生成 .md 和 .html）
# 在 build-all.sh 之前运行

set -e

PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
cd "$PROJECT_DIR"

VERSION_RAW=$(grep -o 'AppVersion = "[^"]*"' backend/config/config.go | cut -d'"' -f2)
[ -z "$VERSION_RAW" ] && echo "❌ 无法获取版本号" && exit 1
VERSION="${VERSION_RAW#v}"

APP_NAME="techfunway-memo"
RELEASE_DIR="release/v${VERSION}"
TEMPLATE_DIR="scripts/release-templates"

mkdir -p "${RELEASE_DIR}"

echo "📝 生成发布说明 (v${VERSION})..."

for type in binary docker fnos; do
	input="${TEMPLATE_DIR}/README-${type}.md"
	if [ ! -f "$input" ]; then
		echo "  ⚠️  跳过 ${type}：模板文件不存在"
		continue
	fi

	# 生成 .md（替换模板变量）
	sed \
		-e "s/{{VERSION}}/${VERSION}/g" \
		-e "s/{{APP_NAME}}/${APP_NAME}/g" \
		"$input" > "${RELEASE_DIR}/README-${type}.md"
	echo "  ✅ README-${type}.md"
done

# 生成 HTML 使用说明
input="${TEMPLATE_DIR}/usage.html"
if [ -f "$input" ]; then
	sed \
		-e "s/{{VERSION}}/${VERSION}/g" \
		-e "s/{{APP_NAME}}/${APP_NAME}/g" \
		"$input" > "${RELEASE_DIR}/usage.html"
	echo "  ✅ usage.html"
fi

echo ""
echo "✅ 发布说明生成完成"
ls -lh "${RELEASE_DIR}"/README-*.md "${RELEASE_DIR}"/usage.html 2>/dev/null | awk '{print "  " $NF}'
