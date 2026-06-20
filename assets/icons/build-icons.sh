#!/bin/bash
# 备忘录图标生成脚本
# 从矢量母版 memo-icon.svg 生成所有 PNG 尺寸，并分发到前端 / 飞牛 fnpack 目录。
#
# 工作流：
#   1. 编辑 assets/icons/memo-icon.svg（矢量源，可任意修改）
#   2. 重新光栅化母版 master-1024.png：
#        - 若装有 rsvg-convert（brew install librsvg）→ 本脚本自动重生成
#        - 否则用浏览器打开 assets/icons/render.html 截图 1024×1024（透明背景）
#          覆盖 assets/icons/master-1024.png
#   3. 运行本脚本：bash assets/icons/build-icons.sh
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$( cd "$DIR/../.." && pwd )"
SVG="$DIR/memo-icon.svg"
MASTER="$DIR/master-1024.png"

# 1) 若有 rsvg-convert，则从 SVG 重新光栅化母版（保留透明角与阴影滤镜）
if command -v rsvg-convert >/dev/null 2>&1; then
  echo "🎨 rsvg-convert 光栅化母版..."
  rsvg-convert -w 1024 -h 1024 "$SVG" -o "$MASTER"
else
  echo "ℹ️  未检测到 rsvg-convert，使用已存在的 master-1024.png（如改过 SVG，请先用 render.html 重新截图）"
fi
[ -f "$MASTER" ] || { echo "❌ 缺少 $MASTER"; exit 1; }

# 2) 用 sips 生成各尺寸（保留 alpha 透明）
TMP="$(mktemp -d)"
gen() { cp "$MASTER" "$TMP/$1.png"; sips -Z "$2" "$TMP/$1.png" >/dev/null 2>&1; }
gen 256 256; gen 128 128; gen 64 64; gen 32 32

# 3) 分发到目标位置
echo "📦 分发图标..."
# 站点 title 图标（前端源；构建时同步到 dist/static）
cp "$SVG"        "$ROOT/frontend/public/favicon.svg"
cp "$TMP/64.png" "$ROOT/frontend/public/favicon.png"
# 飞牛 fnOS 桌面图标 images/icon_{0}.png
cp "$TMP/256.png" "$ROOT/fnpack/app/ui/images/icon_256.png"
cp "$TMP/128.png" "$ROOT/fnpack/app/ui/images/icon_128.png"
cp "$TMP/64.png"  "$ROOT/fnpack/app/ui/images/icon_64.png"
# 飞牛安装包图标
cp "$TMP/256.png" "$ROOT/fnpack/ICON_256.PNG"
cp "$TMP/128.png" "$ROOT/fnpack/ICON.PNG"
# 同步到构建产物（若存在）
for d in "$ROOT/backend/static" "$ROOT/frontend/dist"; do
  [ -d "$d" ] && cp "$SVG" "$d/favicon.svg" && cp "$TMP/64.png" "$d/favicon.png"
done

rm -rf "$TMP"
echo "✅ 图标已更新。如需进发布包，重跑 ./scripts/build-all.sh 与 ./scripts/build-fnpack.sh"
